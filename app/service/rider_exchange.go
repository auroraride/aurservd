// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-05
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"time"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/async"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/ec"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type riderExchangeService struct {
	ctx   context.Context
	rider *ent.Rider

	maxTime time.Duration // 单步骤最大处理时长
	logger  *logging.ExchangeLog

	cabinet   *ent.Cabinet
	subscribe *ent.Subscribe
	exchange  *ent.Exchange
}

func NewRiderExchange(r *ent.Rider) *riderExchangeService {
	s := &riderExchangeService{
		maxTime: 180 * time.Second,
		rider:   r,
	}
	s.ctx = context.WithValue(context.Background(), model.CtxRiderKey{}, r)
	return s
}

// GetProcess 获取待换电信息
func (s *riderExchangeService) GetProcess(req *model.RiderCabinetOperateInfoReq) (res *model.RiderExchangeInfo) {
	NewSetting().SystemMaintainX()

	// 是否有生效中套餐
	sub := NewSubscribe().RecentX(s.rider.ID)

	// 检查用户是否可以办理业务
	NewRiderPermissionWithRider(s.rider).BusinessX().SubscribeX(model.RiderPermissionTypeExchange, sub)

	NewExchange().RiderInterval(s.rider, sub.CityID)

	// 查询电柜
	cs := NewCabinet()
	cab := cs.QueryOneSerialX(req.Serial)

	// 检查可用电池型号
	if !cs.ModelInclude(cab, sub.Model) {
		snag.Panic("电池型号不兼容")
	}

	// 查询电柜是否可使用
	NewCabinet().BusinessableX(cab)

	// 判定是否可以换电
	NewIntelligentCabinet(s.rider).BusinessCensorX(adapter.BusinessExchange, sub, cab)
	uid, info := NewIntelligentCabinet(s.rider).ExchangeUsable(sub.Model, cab)

	// TODO 修改前端返回值
	res = &model.RiderExchangeInfo{
		ID:                         cab.ID,
		UUID:                       uid,
		Full:                       info.FullBin != nil,
		Name:                       cab.Name,
		Health:                     cab.Health,
		Serial:                     cab.Serial,
		Doors:                      cab.Doors,
		BatteryNum:                 cab.BatteryNum,
		BatteryFullNum:             cab.BatteryFullNum,
		RiderCabinetOperateProcess: info,
		Model:                      sub.Model,
		CityID:                     sub.CityID,
		Brand:                      cab.Brand,
	}

	cache.Set(s.ctx, uid, res, 10*time.Minute)

	b, _ := jsoniter.Marshal(res)
	zap.L().Info("换电信息: uuid=" + uid + ", data=" + adapter.ConvertBytes2String(b))

	return res
}

// Start 开始换电
func (s *riderExchangeService) Start(req *model.RiderExchangeProcessReq) {
	// 是否有生效中套餐
	sub := NewSubscribe().RecentX(s.rider.ID)

	s.subscribe = sub

	// 检查用户是否可以办理业务
	NewRiderPermissionWithRider(s.rider).BusinessX().SubscribeX(model.RiderPermissionTypeExchange, sub)

	// 检查是否维护中
	NewSetting().SystemMaintainX()

	// 校验换电间隔
	NewExchange().RiderInterval(s.rider, sub.CityID)

	var (
		cab   *ent.Cabinet
		info  model.RiderExchangeInfo
		bat   *ent.Battery
		batSN *string
	)

	// 尝试从缓存获取智能电柜换电信息
	err := cache.Get(s.ctx, req.UUID).Scan(&info)

	if err != nil {
		snag.Panic("未获取到换电信息")
	}

	cab = NewCabinet().QueryOneSerialX(info.Serial)

	// 查询电柜是否可使用
	NewCabinet().BusinessableX(cab)

	bat = NewIntelligentCabinet(s.rider).BusinessCensorX(adapter.BusinessExchange, sub, cab)
	if bat != nil {
		batSN = silk.String(bat.Sn)
	}

	fully := info.FullBin
	if info.Alternative != nil {
		fully = info.Alternative
	}

	if cab == nil || cab.CityID == nil {
		snag.Panic("未找到电柜信息, 请重试")
	}

	// 记录换电人
	// TODO 超时处理
	s.exchange, _ = ent.Database.Exchange.
		Create().
		SetRiderID(s.rider.ID).
		SetCityID(*cab.CityID).
		SetCabinetInfo(&model.ExchangeCabinetInfo{
			Health:         cab.Health,
			Doors:          cab.Doors,
			BatteryNum:     cab.BatteryNum,
			BatteryFullNum: cab.BatteryFullNum,
		}).
		SetEmpty(&model.BinInfo{
			Index: info.EmptyBin.Index,
		}).
		SetFully(fully).
		SetUUID(req.UUID).
		SetCabinetID(cab.ID).
		SetSuccess(false).
		SetModel(s.subscribe.Model).
		SetNillableEnterpriseID(s.subscribe.EnterpriseID).
		SetNillableStationID(s.subscribe.StationID).
		SetSubscribeID(s.subscribe.ID).
		SetAlternative(info.Alternative != nil).
		SetStartAt(time.Now()).
		SetNillableRiderBattery(batSN).
		Save(s.ctx)

	if s.exchange == nil {
		snag.Panic("换电失败")
	}

	// 异步处理换电任务
	go async.WithTask(func() {
		NewIntelligentCabinet(s.rider).Exchange(req.UUID, s.exchange, sub, bat, cab)
	})
}

// GetProcessStatus 长轮询获取状态
func (s *riderExchangeService) GetProcessStatus(req *model.RiderExchangeProcessStatusReq) (res *model.RiderExchangeProcessRes) {
	info := new(model.RiderExchangeInfo)
	// 尝试从缓存获取智能电柜换电信息
	err := cache.Get(s.ctx, req.UUID).Scan(info)
	if err == nil {
		return NewIntelligentCabinet(s.rider).ExchangeResult(req.UUID)
	}

	start := time.Now()
	var uid xid.ID
	uid, err = xid.FromString(req.UUID)
	if err != nil {
		snag.Panic("未找到换电操作")
	}
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		task := ec.QueryID(uid)
		if task == nil {
			snag.Panic("未找到换电操作")
		}
		cs := task.Exchange.CurrentStep()
		res = &model.RiderExchangeProcessRes{
			Step:    uint8(cs.Step),
			Status:  uint8(cs.Status),
			Message: task.Message,
			Stop:    task.StopAt != nil,
		}
		if cs.IsSuccess() || res.Stop || time.Since(start).Seconds() > 30 {
			return
		}
	}
}
