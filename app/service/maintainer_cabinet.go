// Copyright (C) liasica. 2023-present.
//
// Created at 2023-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"

	"github.com/LucaTheHacker/go-haversine"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type maintainerCabinetService struct {
	*BaseService
	orm *ent.CabinetClient
}

func NewMaintainerCabinet(params ...any) *maintainerCabinetService {
	return &maintainerCabinetService{
		BaseService: newService(params...),
		orm:         ent.Database.Cabinet,
	}
}

// List 运维归属电柜列表
func (s *maintainerCabinetService) List(cityIDs []uint64, req *model.MaintainerCabinetListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().Where(cabinet.CityIDIn(cityIDs...))

	if req.Status != nil {
		q.Where(cabinet.Status(req.Status.Value()))
	}

	if req.ModelID != nil {
		q.Where(cabinet.HasModelsWith(batterymodel.ID(*req.ModelID)))
	}

	if req.Keyword != nil {
		q.Where(
			cabinet.Or(
				cabinet.NameContainsFold(*req.Keyword),
				cabinet.SerialContains(*req.Keyword),
			),
		)
	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(cab *ent.Cabinet) model.CabinetListByDistanceRes {
		return model.CabinetListByDistanceRes{
			CabinetBasicInfo: model.CabinetBasicInfo{
				ID:     cab.ID,
				Brand:  cab.Brand,
				Serial: cab.Serial,
				Name:   cab.Name,
			},
			Status: cab.Status,
			Health: cab.Health,
		}
	}, NewCabinet().SyncCabinets)
}

// Detail 获取电柜详情
func (s *maintainerCabinetService) Detail(req *model.MaintainerCabinetDetailReq) (res model.MaintainerCabinetDetailRes) {
	cab := NewCabinet().DetailFromSerial(req.Serial)
	res.CabinetDetailRes = cab
	// 查找所属网点
	if cab.BranchID == nil {
		snag.Panic("电柜网点未找到")
	}
	b := NewBranch().Query(*cab.BranchID)
	res.Branch = model.Branch{
		ID:   b.ID,
		Name: b.Name,
	}
	// 查找当前最新维护信息
	res.Maintenance = NewAssetMaintenance().QueryByID(cab.ID)

	return
}

// Operatable 校验权限并获取操作人
func (s *maintainerCabinetService) Operatable(m any, cities []uint64, serial string, lng, lat float64, maintenance bool) (cab *ent.Cabinet, operator *logging.Operator, err error) {
	operator, err = logging.GetOperator(m)
	if err != nil {
		return
	}

	// 查找维护中的电柜
	cab, _ = s.orm.QueryNotDeleted().Where(
		cabinet.CityIDIn(cities...),
		cabinet.Serial(serial),
	).First(s.ctx)
	if cab == nil {
		err = errors.New("未找到电柜")
		return
	}

	// 判定距离
	distance := haversine.Distance(haversine.NewCoordinates(lat, lng), haversine.NewCoordinates(cab.Lat, cab.Lng)).Kilometers() * 1000.0
	if distance > 100 {
		err = errors.New("距离过远")
		return
	}

	// 判定维护
	if maintenance && cab.Status != model.CabinetStatusMaintenance.Value() {
		err = errors.New("电柜必须维护")
		return
	}

	return
}

// OperatableX 校验权限并获取操作人
// 如有错误则panic
func (s *maintainerCabinetService) OperatableX(m any, cities []uint64, serial string, lng, lat float64, maintenance bool) (*ent.Cabinet, *logging.Operator) {
	cab, operator, err := s.Operatable(m, cities, serial, lng, lat, maintenance)
	if err != nil {
		snag.Panic(err)
	}
	return cab, operator
}

func (s *maintainerCabinetService) Operate(m any, cities []uint64, req *model.MaintainerCabinetOperateReq) {
	// 校验权限并获取操作人
	cab, operator := s.OperatableX(m, cities, req.Serial, req.Lng, req.Lat, req.Operate.NeedMaintenance())

	// 获取Modifier 创建维保记录使用
	md := s.getModifier(m)
	if md == nil {
		snag.Panic("错误操作人数据")
	}

	switch req.Operate {
	case model.MaintainerCabinetOperateInterrupt:
		if req.Reason == "" {
			snag.Panic("中断原因必填")
		}
		NewCabinet().Interrupt(operator, &model.CabinetInterruptRequest{
			Serial:  req.Serial,
			Message: "运维中断业务:" + req.Reason,
		})
	case model.MaintainerCabinetOperateMaintenance:
		// 门店、代理不参与维护记录
		if !req.Mini {
			// 创建维护记录
			err := NewAssetMaintenance().Create(s.ctx, &model.AssetMaintenanceCreateReq{
				CabinetID:       cab.ID,
				OperatorID:      operator.ID,
				OperateRoleType: model.OperatorTypeMaintainer.Value(),
			}, md)
			if err != nil {
				return
			}
		}
		NewCabinetMgr().Maintain(operator, &model.CabinetMaintainReq{
			ID:       cab.ID,
			Maintain: silk.Bool(true),
		})
	case model.MaintainerCabinetOperateMaintenanceCancel:
		// 门店、代理不参与维护记录
		if !req.Mini {
			if req.Content == "" {
				snag.Panic("维保内容必填")
			}
			switch req.Status {
			case model.AssetMaintenanceStatusSuccess:
			case model.AssetMaintenanceStatusFail:
				if req.Reason == "" {
					snag.Panic("维保失败原因必填")
				}
			default:
				snag.Panic("无效维保状态")
			}

			// 查询维保单
			mt := NewAssetMaintenance().QueryMaintenanceByCabinetID(cab.ID)
			if mt == nil {
				snag.Panic("维保数据不存在")
				return
			}
			// 填写维保结果
			err := NewAssetMaintenance().Modify(s.ctx, &model.AssetMaintenanceModifyReq{
				ID:      mt.ID,
				Reason:  req.Reason,
				Content: req.Content,
				Status:  req.Status,
				Details: req.Details,
			}, md)
			if err != nil {
				return
			}
		}
		NewCabinetMgr().Maintain(operator, &model.CabinetMaintainReq{
			ID:       cab.ID,
			Maintain: silk.Bool(false),
		})
	}
}

// BinOperate 仓位操作
func (s *maintainerCabinetService) BinOperate(m any, cities []uint64, req *model.MaintainerBinOperateReq, waitClose bool) bool {
	// 校验权限并获取操作人
	cab, operator := s.OperatableX(m, cities, req.Serial, req.Lng, req.Lat, true)

	var op any
	switch req.Operate {
	default:
		return false
	case model.MaintainerBinOperateOpen:
		op = &model.CabinetDoorOperateReq{
			ID:        cab.ID,
			Index:     silk.Int(req.Ordinal - 1),
			Remark:    req.Reason,
			Operation: silk.Pointer(model.CabinetDoorOperateOpen),
		}
	case model.MaintainerBinOperateLock:
		op = &model.CabinetDoorOperateReq{
			ID:        cab.ID,
			Index:     silk.Int(req.Ordinal - 1),
			Remark:    req.Reason,
			Operation: silk.Pointer(model.CabinetDoorOperateLock),
		}
	case model.MaintainerBinOperateUnlock:
		op = &model.CabinetDoorOperateReq{
			ID:        cab.ID,
			Index:     silk.Int(req.Ordinal - 1),
			Remark:    req.Reason,
			Operation: silk.Pointer(model.CabinetDoorOperateUnlock),
		}
	case model.MaintainerBinOperateDisable:
		op = &model.CabinetBinDeactivateReq{
			ID:        cab.ID,
			Index:     silk.Int(req.Ordinal - 1),
			Remark:    req.Reason,
			Operation: model.CabinetBinDeactiveOperationDisable,
		}
	case model.MaintainerBinOperateEnable:
		op = &model.CabinetBinDeactivateReq{
			ID:        cab.ID,
			Index:     silk.Int(req.Ordinal - 1),
			Remark:    req.Reason,
			Operation: model.CabinetBinDeactiveOperationEnable,
		}
	}

	return NewCabinetMgr().BinOperate(operator, cab.ID, op, waitClose)
}

// Pause 暂停维护
func (s *maintainerCabinetService) Pause(m any, cities []uint64, req *model.MaintainerCabinetPauseReq) {
	// 校验权限并获取操作人
	cab, operator := s.OperatableX(m, cities, req.Serial, req.Lng, req.Lat, false)

	// 查询维保单
	mt := NewAssetMaintenance().QueryMaintenanceByCabinetID(cab.ID)
	if mt == nil {
		snag.Panic("维保数据不存在")
		return
	}
	switch req.Status {
	case model.AssetMaintenanceStatusPause:
		// 更新维保数据
		err := mt.Update().SetStatus(model.AssetMaintenanceStatusPause.Value()).Exec(context.Background())
		if err != nil {
			snag.Panic(err)
		}
		// 暂停维护
		NewCabinetMgr().Maintain(operator, &model.CabinetMaintainReq{
			ID:       cab.ID,
			Maintain: silk.Bool(false),
		})

	case model.AssetMaintenanceStatusUnder:
		// 更新维保数据
		err := mt.Update().SetStatus(model.AssetMaintenanceStatusUnder.Value()).Exec(context.Background())
		if err != nil {
			snag.Panic(err)
		}
		// 继续维护
		NewCabinetMgr().Maintain(operator, &model.CabinetMaintainReq{
			ID:       cab.ID,
			Maintain: silk.Bool(true),
		})
	}
}

func (s *maintainerCabinetService) getModifier(data any) (md *model.Modifier) {
	switch v := data.(type) {
	default:
		return nil
	case *ent.Employee:
		return &model.Modifier{
			ID:    v.ID,
			Phone: v.Phone,
			Name:  v.Name,
		}
	case *ent.Agent:
		return &model.Modifier{
			ID:    v.ID,
			Phone: v.Phone,
			Name:  v.Name,
		}
	case *ent.Maintainer:
		return &model.Modifier{
			ID:    v.ID,
			Phone: v.Phone,
			Name:  v.Name,
		}
	}
}
