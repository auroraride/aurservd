// Copyright (C) liasica. 2022-present.
//
// Created at 2022-12-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/defs/cabdef"
	"github.com/auroraride/adapter/log"
	"github.com/auroraride/adapter/rpc/pb"
	"github.com/auroraride/adapter/rpc/pb/timestamppb"
	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/rpc"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type intelligentCabinetService struct {
	*BaseService
}

func NewIntelligentCabinet(params ...any) *intelligentCabinetService {
	return &intelligentCabinetService{
		BaseService: newService(params...),
	}
}

// ExchangeUsable 获取换电信息
func (s *intelligentCabinetService) ExchangeUsable(bm string, cab *ent.Cabinet) (uid string, info *model.RiderCabinetOperateProcess) {
	payload := &cabdef.ExchangeUsableRequest{
		Serial: cab.Serial,
		Minsoc: cache.Float64(model.SettingExchangeMinBatteryKey),
		Lock:   10,
		Model:  bm,
	}

	v, err := adapter.Post[cabdef.CabinetBinUsableResponse](s.GetCabinetAdapterUrlX(cab, "/exchange/usable"), s.GetAdapterUserX(), payload)

	if err != nil {
		snag.Panic(err)
	}

	info = &model.RiderCabinetOperateProcess{
		EmptyBin: &model.BinInfo{
			Index: v.Empty.Ordinal - 1,
		},
	}

	fully := &model.BinInfo{
		Index:       v.Fully.Ordinal - 1,
		Voltage:     v.Fully.Voltage,
		Electricity: model.BatterySoc(v.Fully.Soc),
	}

	if v.Fully.Soc >= cache.Float64(model.SettingBatteryFullKey) {
		info.FullBin = fully
	} else {
		info.Alternative = fully
	}

	uid = v.UUID

	return
}

func (s *intelligentCabinetService) exchangeCacheKey(uid string) string {
	return "EXCHANGE:" + uid
}

// Exchange 请求换电
func (s *intelligentCabinetService) Exchange(uid string, ex *ent.Exchange, sub *ent.Subscribe, old *ent.Asset, cab *ent.Cabinet) {
	if cab.Intelligent && old == nil {
		snag.Panic("请求参数错误")
	}

	user := s.GetAdapterUserX()

	var (
		err      error
		stopAt   *time.Time
		duration float64
		putout   string
		putin    string
		empty    *model.BinInfo
		bs       = NewBattery()
		key      = s.exchangeCacheKey(uid)
	)

	success := atomic.NewBool(false)
	message := atomic.NewString("")

	defer func() {
		updater := ex.Update()

		if err != nil {
			updater.SetRemark(err.Error())
		}

		if stopAt == nil {
			stopAt = silk.Pointer(time.Now())
			duration = stopAt.Sub(ex.StartAt).Seconds()
		}

		if empty != nil {
			updater.SetEmpty(empty)
			// ex.Info.Empty = empty
		}

		var msg *string

		if message.Load() != "" {
			msg = silk.Pointer(message.Load())
		}

		// 保存数据库
		_ = updater.
			SetSuccess(success.Load()).
			SetFinishAt(*stopAt).
			SetDuration(int(duration)).
			SetPutoutBattery(putout).
			SetPutinBattery(putin).
			SetNillableMessage(msg).
			// SetInfo(ex.Info).
			Exec(s.ctx)
	}()

	// 缓存第一步
	ar.Redis.RPush(s.ctx, key, &pb.CabinetExchangeResponse{
		Uuid:     uid,
		Step:     model.ExchangeStepOpenEmpty.Uint32(),
		Business: adapter.BusinessExchange.String(),
		StartAt:  timestamppb.New(time.Now()),
	})

	// 设置10分钟超时
	_ = ar.Redis.Expire(s.ctx, key, 10*time.Minute)

	// 使用gRPC请求换电
	var batSN string
	if old != nil {
		batSN = old.Sn
	}
	err = rpc.CabinetExchange(
		rpc.CabinetKey(cab.Brand, cab.Intelligent),
		user,
		&pb.CabinetExchangeRequest{
			Uuid:    uid,
			Serial:  cab.Serial,
			Battery: batSN,
			Expires: model.CabinetBusinessScanExpires,
			Timeout: model.CabinetBusinessStepTimeout,
			Minsoc:  cache.Float64(model.SettingExchangeMinBatteryKey),
		}, func(result *pb.CabinetExchangeResponse, stop bool) {
			zap.L().Info("换电步骤记录回调", log.Payload(result))
			if result == nil {
				return
			}

			message.Store(result.Message)
			success.Store(result.Success)

			duration += result.Duration
			stopAt = silk.Pointer(result.StopAt.AsTime())

			// todo 失败的第二步骤逻辑处理

			// 如果成功并且是智能柜, 记录电池编码
			if result.Success {
				if cab.Intelligent {
					after := result.After
					before := result.Before

					// 记录用户放入的电池
					if result.Step == model.ExchangeStepPutInto.Uint32() && after != nil {
						putin = after.BatterySn
						empty = &model.BinInfo{
							Index:       int(after.Ordinal) - 1,
							Electricity: model.BatterySoc(after.Current),
							Voltage:     after.Voltage,
						}

						// 清除旧电池分配信息
						err = NewBattery(s.operator).Unallocate(old, model.AssetLocationsTypeCabinet, cab.ID, model.AssetTransferTypeExchange)
						if err != nil {
							return
						}

						go bs.RiderBusiness(true, putin, s.rider, cab, int(after.Ordinal))
					}

					// 记录用户取走的电池
					// 判定第三步是否成功, 只要柜门开启就把电池绑定到骑手 - BY: 曹博文
					if result.Step == model.ExchangeStepOpenFull.Uint32() && before != nil {
						putout = before.BatterySn

						go bs.RiderBusiness(false, putout, s.rider, cab, int(before.Ordinal))

						// 更新新电池信息
						bat, _ := bs.QuerySn(putout)
						if bat != nil {
							_ = ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
								return NewBattery(s.operator).Allocate(bat, sub, model.AssetTransferTypeExchange)
							})
						}
					}
				} else {
					if result.Step == model.ExchangeStepPutInto.Uint32() {
						at, _ := ent.Database.Asset.QueryNotDeleted().Where(
							asset.Type(model.AssetTypeNonSmartBattery.Value()),
							asset.LocationsType(model.AssetLocationsTypeRider.Value()),
							asset.LocationsID(s.rider.ID),
						).WithModel().First(s.ctx)
						if at == nil {
							zap.L().Error(fmt.Sprintf("换电-放电 非智能柜未找到电池信息 %s", cab.Serial))
							return
						}
						// 清除旧电池分配信息
						err = NewBattery(s.operator).Unallocate(at, model.AssetLocationsTypeCabinet, cab.ID, model.AssetTransferTypeExchange)
						if err != nil {
							zap.L().Error("换电-放电 非智能柜清除电池分配信息失败", zap.Error(err))
							return
						}
					}

					if result.Step == model.ExchangeStepOpenFull.Uint32() {
						var modelID uint64
						models, _ := cab.QueryModels().All(s.ctx)
						if len(models) > 0 {
							modelID = models[0].ID
						}
						locationsType := model.AssetLocationsTypeCabinet
						newBattery, _ := NewAsset().QueryNonSmartBattery(&model.QueryAssetBatteryReq{
							LocationsType: &locationsType,
							LocationsID:   silk.UInt64(cab.ID),
							ModelID:       modelID,
						})
						if newBattery == nil {
							zap.L().Error(fmt.Sprintf("换电-取电 非智能柜未找到电池信息 %s", cab.Serial))
							return
						}
						_ = ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
							return NewBattery(s.operator).Allocate(newBattery, sub, model.AssetTransferTypeExchange)
						})
					}
				}
			}
			// 缓存结果
			ar.Redis.RPush(s.ctx, key, result)
		},
	)

	if err != nil {
		zap.L().Error("换电请求失败", zap.Error(err), user.ZapField(), zap.String("uuid", uid))
		message.Store(err.Error())
		success.Store(false)
	}
}

// ExchangeResult 查询换电结果
func (s *intelligentCabinetService) ExchangeResult(uid string) (res *model.RiderExchangeProcessRes) {
	key := s.exchangeCacheKey(uid)

	start := time.Now()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		if time.Since(start).Seconds() > 30 {
			return
		}

		// 获取缓存信息
		var data pb.CabinetExchangeResponse
		err := ar.Redis.LPop(s.ctx, key).Scan(&data)
		// 未获取到缓存信息
		// && errors.Is(err, redis.Nil)
		if err != nil {
			continue
		}
		status := model.TaskStatusProcessing
		if data.StopAt != nil {
			if data.Success {
				status = model.TaskStatusSuccess
			} else {
				status = model.TaskStatusFail
			}
		}

		return &model.RiderExchangeProcessRes{
			Step:    uint8(data.Step),
			Status:  status,
			Message: data.Message,
			Stop:    data.Step == model.ExchangeStepPutOut.Uint32() || status == model.TaskStatusFail || data.Message != "",
		}
	}

	return nil
}

// BusinessCensorX 校验用户是否可以办理电柜业务
func (s *intelligentCabinetService) BusinessCensorX(bus adapter.Business, sub *ent.Subscribe, cab *ent.Cabinet) (bat *ent.Asset) {
	// if !cab.Intelligent {
	// 	return
	// }

	// 判定电柜状态
	if cab.Status == model.CabinetStatusMaintenance.Value() {
		snag.Panic("电柜维护中，请联系客服")
	}

	// // 判定是否智能电柜套餐
	// if !sub.Intelligent {
	// 	snag.Panic("套餐不匹配")
	// }

	// 获取电池
	bat, _ = sub.QueryBattery().WithModel().Where(asset.TypeIn(model.AssetTypeSmartBattery.Value(), model.AssetTypeNonSmartBattery.Value())).First(s.ctx)

	// 业务如果需要电池, 查找电池信息
	if bus.BatteryNeed() {
		// 未找到当前绑定的电池信息
		if bat == nil {
			snag.Panic(adapter.ErrorBatteryNotFound)
			return
		}
		if bat.Edges.Model != nil {
			// 检查电池型号与电柜型号兼容
			if !NewCabinet().ModelInclude(cab, bat.Edges.Model.Model) {
				snag.Panic("电池型号不匹配，请更换电柜重试")
			}
		}
	}

	return
}

// BusinessUsable 获取可用的业务仓位信息
func (s *intelligentCabinetService) BusinessUsable(cab *ent.Cabinet, bus adapter.Business, bm string) (uid string, index int, err error) {
	payload := &cabdef.BusinuessUsableRequest{
		Minsoc:   cache.Float64(model.SettingExchangeMinBatteryKey),
		Business: bus,
		Serial:   cab.Serial,
		Model:    bm,
	}

	var v cabdef.CabinetBinUsableResponse
	v, err = adapter.Post[cabdef.CabinetBinUsableResponse](s.GetCabinetAdapterUrlX(cab, "/business/usable"), s.GetAdapterUserX(), payload)
	if err != nil {
		return
	}

	uid = v.UUID
	index = v.BusinessBin.Ordinal - 1
	return
}

// DoBusiness 请求办理业务
func (s *intelligentCabinetService) DoBusiness(uidstr string, bus adapter.Business, sub *ent.Subscribe, riderBat *ent.Asset, cab *ent.Cabinet) (info *model.BinInfo, batinfo *model.Battery, err error) {
	defer func() {
		// 缓存任务返回
		data := &model.BusinessCabinetStatusRes{
			Success: err == nil,
			Stop:    true,
		}
		if err != nil {
			data.Message = err.Error()
		}
		cache.Set(s.ctx, uidstr, data, 10*time.Minute)
	}()

	var uid uuid.UUID
	uid, err = uuid.Parse(uidstr)
	if err != nil {
		return
	}

	var batterySN string
	if riderBat != nil {
		batterySN = riderBat.Sn
	}

	payload := &cabdef.BusinessRequest{
		UUID:     uid,
		Business: bus,
		Serial:   cab.Serial,
		Timeout:  model.CabinetBusinessStepTimeout,
		Battery:  batterySN,
		Model:    sub.Model,
	}

	var v cabdef.BusinessResponse
	v, err = adapter.Post[cabdef.BusinessResponse](s.GetCabinetAdapterUrlX(cab, "/business/do"), s.GetAdapterUserX(), payload)

	if err != nil {
		return
	}

	// TODO 失败后电池信息是否更新
	if v.Error != "" {
		err = errors.New(v.Error)
		return
	}

	var sn string
	var putin bool
	results := v.Results

	switch bus {
	case adapter.BusinessActive, adapter.BusinessContinue:
		sn = results[0].Before.BatterySN
	case adapter.BusinessPause, adapter.BusinessUnsubscribe:
		sn = results[1].After.BatterySN
		putin = true
	}

	b := results[1].After
	info = &model.BinInfo{
		Index:       b.Ordinal - 1,
		Electricity: model.BatterySoc(b.Soc),
		Voltage:     b.Voltage,
	}

	// 业务类型对应调拨类型
	var at model.AssetTransferType
	switch bus {
	case adapter.BusinessActive:
		at = model.AssetTransferTypeActive
	case adapter.BusinessPause:
		at = model.AssetTransferTypePause
	case adapter.BusinessContinue:
		at = model.AssetTransferTypeContinue
	case adapter.BusinessUnsubscribe:
		at = model.AssetTransferTypeUnSubscribe
	}

	var bat *ent.Asset
	var m string
	// 若智能电柜, 需记录电池信息
	if cab.Intelligent {
		// 获取电池
		bat, err = NewBattery().QuerySn(sn)
		if err != nil && bat == nil {
			zap.L().Error("业务记录失败", zap.Error(err))
			return
		}

		if bat.Edges.Model != nil {
			m = bat.Edges.Model.Model
		}

	} else {
		// 获取非智能电池
		bat, _ = ent.Database.Asset.QueryNotDeleted().Where(
			asset.LocationsType(model.AssetLocationsTypeCabinet.Value()),
			asset.HasCabinetWith(cabinet.Serial(cab.Serial)),
		).WithModel().First(s.ctx)
		if bat == nil {
			zap.L().Error("未找到电池信息")
			return
		}
		if bat.Edges.Model != nil {
			m = bat.Edges.Model.Model
		}
		// 非智能电柜, 没有电池编码
		sn = ""
	}
	batinfo = &model.Battery{
		ID:    bat.ID,
		SN:    sn,
		Model: m,
	}

	// 取走电池 激活,取消寄存会走这里
	if !putin {
		_ = ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
			return NewBattery(s.operator).Allocate(bat, sub, at)
		})
	}

	return
}

// Operate 操作电柜
// waitClose 是否等待关闭仓门（仅开仓动作有效）
func (s *intelligentCabinetService) Operate(operator *logging.Operator, cab *ent.Cabinet, op cabdef.Operate, req *model.CabinetDoorOperateReq, waitClose bool) (success bool, data []*cabdef.BinOperateResult) {
	now := time.Now()
	br := cab.Brand
	ordinal := *req.Index + 1

	var err error

	defer func() {
		go func() {
			// 上传日志
			dlog := &logging.DoorOperateLog{
				ID:            shortuuid.New(),
				Brand:         br.String(),
				OperatorName:  operator.Name,
				OperatorID:    operator.ID,
				OperatorPhone: operator.Phone,
				Serial:        cab.Serial,
				Name:          strconv.Itoa(ordinal) + "号仓",
				Operation:     req.Operation.String(),
				OperatorRole:  operator.OperatorRole(),
				Success:       success,
				Remark:        req.Remark,
				Time:          now.Format(carbon.DateTimeLayout),
			}
			dlog.Send()
		}()
	}()

	payload := &cabdef.OperateBinRequest{
		Operate: op,
		Ordinal: silk.Int(ordinal),
		Serial:  cab.Serial,
		Remark:  req.Remark,
	}

	apiUrl := "/operate/bin"
	if op == cabdef.OperateDoorOpen && waitClose {
		apiUrl = "/operate/bin/open-close"
	}

	data, err = adapter.Post[[]*cabdef.BinOperateResult](s.GetCabinetAdapterUrlX(cab, apiUrl), operator.GetAdapterUserX(), payload)
	zap.L().Info("电柜操作", zap.Bool("success", success), log.Payload(data), zap.Error(err))

	success = err == nil
	return success, data
}

func (s *intelligentCabinetService) Deactivate(operator *logging.Operator, cab *ent.Cabinet, payload *cabdef.BinDeactivateRequest) (success bool) {
	now := time.Now()
	br := cab.Brand

	operation := "启用仓位"
	if *payload.Deactivate {
		operation = "禁用仓位"
	}

	res, _ := adapter.Post[model.StatusResponse](s.GetCabinetAdapterUrlX(cab, "/bin/deactivate"), operator.GetAdapterUserX(), payload)

	success = res.Status

	if success {
		// 上传日志
		dlog := &logging.DoorOperateLog{
			ID:            shortuuid.New(),
			Brand:         br.String(),
			OperatorName:  operator.Name,
			OperatorID:    operator.ID,
			OperatorPhone: operator.Phone,
			Serial:        cab.Serial,
			Name:          strconv.Itoa(payload.Ordinal) + "号仓",
			Operation:     operation,
			OperatorRole:  operator.OperatorRole(),
			Success:       true,
			Remark:        *payload.Reason,
			Time:          now.Format(carbon.DateTimeLayout),
		}
		go dlog.Send()
	}
	return
}

// OpenBind 开电池仓并绑定骑手
func (s *intelligentCabinetService) OpenBind(req *model.CabinetOpenBindReq) {
	bs := NewAsset(s.modifier)
	// 查找骑手
	rd := NewRider().QueryPhoneX(req.Phone)
	// 查找订阅
	sub := NewSubscribe().QueryEffectiveX(rd.ID, ent.SubscribeQueryWithBattery, ent.SubscribeQueryWithRider)

	// 判定
	if exists, _ := sub.QueryBattery().Where(
		asset.TypeIn(model.AssetTypeNonSmartBattery.Value(), model.AssetTypeSmartBattery.Value()),
		asset.LocationsType(model.AssetLocationsTypeRider.Value()),
		asset.LocationsID(rd.ID),
	).Exist(s.ctx); exists {
		snag.Panic("该骑手当前有绑定的电池")
	}
	var bat *ent.Asset
	var err error
	// 查询电柜
	cab := NewCabinet().QueryOne(req.ID)
	if cab.Intelligent {
		// 查询电柜最新信息
		info, _ := s.Bininfo(cab, *req.Index+1)
		if info == nil {
			snag.Panic("获取最新仓位信息失败")
			return
		}
		if info.BatterySN != req.BatterySN {
			snag.Panic("电池编码有变动, 请刷新后重试")
		}
		// 查找电池
		bat, err = bs.QuerySn(req.BatterySN)
		if err != nil {
			snag.Panic("未找到电池信息")
		}
	} else {
		locationsType := model.AssetLocationsTypeCabinet
		locationsID := req.ID
		md := ent.Database.BatteryModel.Query().Where(batterymodel.Model(sub.Model)).FirstX(s.ctx)
		if md == nil {
			snag.Panic("未找到电池型号")
			return
		}
		bat, _ = NewAsset().QueryNonSmartBattery(&model.QueryAssetBatteryReq{
			LocationsType: &locationsType,
			LocationsID:   &locationsID,
			ModelID:       md.ID,
		})
	}
	// 开门
	success, _ := s.Operate(logging.GetOperatorX(s.modifier), cab, cabdef.OperateDoorOpen, &model.CabinetDoorOperateReq{
		ID:        req.ID,
		Index:     req.Index,
		Remark:    req.Remark,
		Operation: silk.Pointer(model.CabinetDoorOperateOpen),
	}, false)
	if !success {
		snag.Panic("仓门开启失败")
	}
	// 绑定
	err = NewBattery(s.operator).Bind(bat, sub, rd)
	if err != nil {
		snag.Panic(err)
	}
}

func (s *intelligentCabinetService) Bininfo(cab *ent.Cabinet, ordinal int) (*cabdef.BinInfo, error) {
	return adapter.Post[*cabdef.BinInfo](s.GetCabinetAdapterUrlX(cab, "/device/bininfo"), s.GetAdapterUserX(), &cabdef.BinInfoRequest{
		Serial:  cab.Serial,
		Ordinal: silk.Int(ordinal),
	})
}
