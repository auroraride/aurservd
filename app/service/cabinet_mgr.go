// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-03
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"
	"time"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/defs/cabdef"
	"github.com/golang-module/carbon/v2"
	"github.com/lithammer/shortuuid/v4"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/ec"
	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type cabinetMgrService struct {
	ctx      context.Context
	modifier *model.Modifier
}

func NewCabinetMgr() *cabinetMgrService {
	return &cabinetMgrService{
		ctx: context.Background(),
	}
}

func NewCabinetMgrWithModifier(m *model.Modifier) *cabinetMgrService {
	s := NewCabinetMgr()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

// Maintain 设置电柜操作维护
func (s *cabinetMgrService) Maintain(operator *logging.Operator, req *model.CabinetMaintainReq) (detail *model.CabinetDetailRes) {
	if req.Maintain == nil {
		snag.Panic("参数请求错误")
	}
	cab := NewCabinet().QueryOne(req.ID)

	if model.CabinetStatus(cab.Status) == model.CabinetStatusPending {
		snag.Panic("未投放电柜无法操作")
	}

	status := model.CabinetStatusNormal
	if *req.Maintain {
		status = model.CabinetStatusMaintenance
	}

	detail = NewCabinet().Detail(cab)

	_, err := cab.Update().SetStatus(status.Value()).Save(s.ctx)
	if err != nil {
		snag.Panic(err)
	}

	detail.Status = status

	// 记录日志
	go logging.NewOperateLog().
		SetRef(cab).
		SetOperator(operator).
		SetOperate(model.OperateCabinetMaintain).
		SetDiff(model.CabinetStatus(cab.Status).String(), status.String()).
		SetRemark(cab.Remark).
		Send()

	return
}

// BinOperate 仓位操作
func (s *cabinetMgrService) BinOperate(operator *logging.Operator, id uint64, data any) bool {
	if operator.Type != model.OperatorTypeManager && operator.Type != model.OperatorTypeMaintainer {
		snag.Panic("无权限操作")
	}

	cab := NewCabinet().QueryOne(id)

	if model.CabinetStatus(cab.Status) != model.CabinetStatusMaintenance {
		snag.Panic("非操作维护中不可操作")
	}

	switch req := data.(type) {
	case *model.CabinetDoorOperateReq:
		var op cabdef.Operate
		switch *req.Operation {
		case model.CabinetDoorOperateOpen:
			op = cabdef.OperateDoorOpen
		case model.CabinetDoorOperateLock:
			op = cabdef.OperateBinDisable
		case model.CabinetDoorOperateUnlock:
			op = cabdef.OperateBinEnable
		}
		success, results := NewIntelligentCabinet().Operate(operator, cab, op, req)
		if op == cabdef.OperateDoorOpen {
			s.BinOperateAssetTransfer(operator, cab, results)
		}
		return success
	case *model.CabinetBinDeactivateReq:
		return NewIntelligentCabinet().Deactivate(operator, cab, &cabdef.BinDeactivateRequest{
			Serial:     cab.Serial,
			Ordinal:    *req.Index + 1,
			Deactivate: silk.Bool(req.Operation == 2),
			Reason:     silk.String(req.Remark),
		})
	default:
		return false
	}
}

// BinOperateAssetTransfer 仓位操作处理资产调拨
func (s *cabinetMgrService) BinOperateAssetTransfer(operator *logging.Operator, cab *ent.Cabinet, data []*cabdef.BinOperateResult) {
	// 操作人判定 只有运维和门店管理员才会记录
	if operator.Type == model.OperatorTypeMaintainer || operator.Type == model.OperatorTypeEmployee || operator.Type == model.OperatorTypeAgent {
		// 根据返回数据更新仓位状态
		var before, after *cabdef.BinInfo
		for _, result := range data {
			if before == nil {
				before = result.Before
			}
			after = result.After
		}

		if after != nil {
			var fromLocationType, toLocationType model.AssetLocationsType
			var fromLocationID, toLocationID uint64
			var sn string
			var assetType model.AssetType
			var modelID uint64
			var in bool
			// 非智能柜根据电池在位判定电池资产变化
			if !cab.Intelligent && before.BatteryExists != after.BatteryExists {
				assetType = model.AssetTypeNonSmartBattery

				models, _ := cab.QueryModels().All(s.ctx)
				if len(models) > 0 {
					modelID = models[0].ID
				}
				// 电池取出
				if before.BatteryExists {
					toLocationType, toLocationID = s.OperatorAssetLocation(operator, cab)
					fromLocationType = model.AssetLocationsTypeCabinet
					fromLocationID = cab.ID
					in = false
				}

				// 电池放入
				if !before.BatteryExists {
					// 电池放入 查询运维/站点/门店 如果有电池调拨到电柜
					fromLocationType, fromLocationID = s.OperatorAssetLocation(operator, cab)
					toLocationType = model.AssetLocationsTypeCabinet
					toLocationID = cab.ID
					in = true
					// 查找电池数量是否满足
					bat, _ := NewAsset().QueryNonSmartBattery(&model.QueryAssetReq{
						LocationsID:   silk.UInt64(fromLocationID),
						LocationsType: &fromLocationType,
						ModelID:       modelID,
					})
					// todo 这里应该上报异常 因为非智能电池不知道插入哪一个 所以没办法根据原来的位置找到电池
					if bat == nil {
						snag.Panic("电池不存在")
					}
				}
				err := s.BatteryTransfer(operator, model.BatteryTransferReq{
					FromLocationType: fromLocationType,
					FromLocationID:   fromLocationID,
					ToLocationType:   toLocationType,
					ToLocationID:     toLocationID,
					Details: []model.AssetTransferCreateDetail{
						{
							AssetType: assetType,
							ModelID:   silk.UInt64(modelID),
							Num:       silk.UInt(1),
						},
					},
					In: in,
				})
				if err != nil {
					zap.L().Error("电池调拨失败", zap.Error(err))
				}
			}

			// 智能电柜根据电池编码判定电池资产变化
			if cab.Intelligent && before.BatterySN != after.BatterySN {
				assetType = model.AssetTypeSmartBattery
				// 记录取出电池状态变动，产生调拨单
				if before.BatterySN != "" {
					// 电池取出
					toLocationType, toLocationID = s.OperatorAssetLocation(operator, cab)
					fromLocationType = model.AssetLocationsTypeCabinet
					fromLocationID = cab.ID
					sn = before.BatterySN
					in = false
				}
				if after.BatterySN != "" {
					// 电池放入 查询运维/站点/门店 如果有电池 从运维调拨到电柜
					// 如果没有电池 从电池原本位置调拨到电柜
					fromLocationType, fromLocationID, _ = s.GetBattery(operator, cab, after.BatterySN)
					toLocationType = model.AssetLocationsTypeCabinet
					toLocationID = cab.ID
					sn = after.BatterySN
					in = true
				}
				err := s.BatteryTransfer(operator, model.BatteryTransferReq{
					FromLocationType: fromLocationType,
					FromLocationID:   fromLocationID,
					ToLocationType:   toLocationType,
					ToLocationID:     toLocationID,
					Details: []model.AssetTransferCreateDetail{
						{
							AssetType: assetType,
							SN:        silk.String(sn),
						},
					},
					In: in,
				})
				if err != nil {
					zap.L().Error("电池调拨失败", zap.Error(err))
				}
			}
		}
	}
}

// OperatorAssetLocation 获取操作人资产位置
func (s *cabinetMgrService) OperatorAssetLocation(operator *logging.Operator, cab *ent.Cabinet) (locationType model.AssetLocationsType, locationID uint64) {
	switch operator.Type {
	case model.OperatorTypeEmployee:
		locationType = model.AssetLocationsTypeStore
		st, _ := cab.QueryStore().First(s.ctx)
		if st == nil {
			snag.Panic("改电池未绑定门店")
			return
		}
		locationID = st.ID
	case model.OperatorTypeMaintainer:
		locationType = model.AssetLocationsTypeOperation
		locationID = operator.ID
	case model.OperatorTypeAgent:
		locationType = model.AssetLocationsTypeStation
		sta, _ := cab.QueryStation().First(s.ctx)
		if sta == nil {
			snag.Panic("改电池未绑定站点")
			return
		}
		locationID = sta.ID
	default:
		snag.Panic("无权限操作")
	}
	return
}

// GetBattery 获取电池
func (s *cabinetMgrService) GetBattery(operator *logging.Operator, cab *ent.Cabinet, sn string) (locationType model.AssetLocationsType, locationID uint64, at *ent.Asset) {
	// 如果没有电池 从电池原本位置调拨到电柜
	locationType, locationID = s.OperatorAssetLocation(operator, cab)
	q := ent.Database.Asset.QueryNotDeleted().Where(
		asset.LocationsType(locationType.Value()),
		asset.LocationsID(locationID),
		asset.Sn(sn),
	)

	at, _ = q.First(s.ctx)
	if at == nil {
		// 查询电池原本位置
		at, _ = ent.Database.Asset.QueryNotDeleted().Where(asset.Sn(sn)).First(s.ctx)
		if at == nil {
			snag.Panic("电池不存在")
			return
		}
		locationType = model.AssetLocationsType(at.LocationsType)
		locationID = at.LocationsID
	}
	return locationType, locationID, at
}

// BatteryTransfer 电池调拨
func (s *cabinetMgrService) BatteryTransfer(operator *logging.Operator, req model.BatteryTransferReq) error {
	var reason string
	reason = "操作取出电池"
	if req.In {
		reason = "操作放入电池"
	}
	_, failed, err := NewAssetTransfer().Transfer(s.ctx, &model.AssetTransferCreateReq{
		FromLocationType:  &req.FromLocationType,
		FromLocationID:    &req.FromLocationID,
		ToLocationType:    req.ToLocationType,
		ToLocationID:      req.ToLocationID,
		Details:           req.Details,
		Reason:            reason,
		AssetTransferType: model.AssetTransferTypeTransfer,
		OperatorID:        operator.ID,
		OperatorType:      operator.Type,
		AutoIn:            true,
		SkipLimit:         true,
	}, &model.Modifier{
		ID:    operator.ID,
		Name:  operator.Name,
		Phone: operator.Phone,
	})
	if err != nil {
		return err
	}
	if len(failed) > 0 {
		return errors.New(failed[0])
	}
	return nil
}

// Reboot 重启电柜
func (s *cabinetMgrService) Reboot(req *model.IDPostReq) bool {
	if s.modifier == nil {
		snag.Panic("权限校验失败")
	}

	ec.BusyFromIDX(req.ID)

	now := time.Now()
	opId := shortuuid.New()

	cab := NewCabinetWithModifier(s.modifier).QueryOne(req.ID)

	if model.CabinetStatus(cab.Status) != model.CabinetStatusMaintenance {
		snag.Panic("非操作维护中不可操作")
	}

	if cab.Brand == adapter.CabinetBrandKaixin {
		snag.Panic("凯信电柜不支持该操作")
	}

	var status bool

	// 创建并开始任务
	task := &ec.Task{
		CabinetID: req.ID,
		Serial:    cab.Serial,
		Job:       model.JobManagerReboot,
		Cabinet:   cab.GetTaskInfo(),
	}

	task.Create().Start()

	// 结束回调
	defer func() {
		ts := model.TaskStatusSuccess
		if !status {
			ts = model.TaskStatusFail
			task.Message = "重启失败"
		}
		task.Stop(ts)
	}()

	// 请求云动重启
	// TODO 云动 - 重启

	brand := cab.Brand
	go func() {
		// 上传日志
		dlog := &logging.DoorOperateLog{
			ID:            opId,
			Brand:         brand.String(),
			OperatorName:  s.modifier.Name,
			OperatorID:    s.modifier.ID,
			OperatorPhone: s.modifier.Phone,
			OperatorRole:  model.CabinetDoorOperatorRoleManager, // operator.OperatorRole()
			Serial:        cab.Serial,
			Operation:     "重启",
			Success:       status,
			Time:          now.Format(carbon.DateTimeLayout),
		}
		dlog.Send()
	}()

	return status
}
