// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"time"

	"github.com/auroraride/adapter"
	"github.com/golang-module/carbon/v2"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/ec"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/allocate"
	"github.com/auroraride/aurservd/internal/ent/business"
	"github.com/auroraride/aurservd/internal/ent/subscribepause"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
)

// TODO 服务器崩溃后自动启动继续换电进程
// TODO 电柜缓存优化

type riderBusinessService struct {
	ctx     context.Context
	rider   *ent.Rider
	maxTime time.Duration // 单步骤最大处理时长

	cabinet   *ent.Cabinet
	subscribe *ent.Subscribe

	bt business.Type

	battery *ent.Battery

	response model.BusinessCabinetStatus
}

func NewRiderBusiness(rider *ent.Rider) *riderBusinessService {
	s := &riderBusinessService{
		ctx:     context.Background(),
		maxTime: 180 * time.Second,
	}
	s.ctx = context.WithValue(s.ctx, model.CtxRiderKey{}, rider)
	s.rider = rider
	return s
}

// preprocess 预处理业务
func (s *riderBusinessService) preprocess(serial string, bt business.Type) {
	NewSetting().SystemMaintainX()

	cs := NewCabinet()

	cab := cs.QueryOneSerialX(serial)
	// TODO 智能电池是否需要调拨
	if !cab.Transferred {
		snag.Panic("电柜资产异常")
	}

	// 是否有生效中套餐
	sub := NewSubscribe().Recent(s.rider.ID)
	if sub == nil {
		snag.Panic("无生效中的骑行卡")
	}

	if sub.BrandID != nil {
		snag.Panic("车电订阅无法自主办理业务")
	}

	s.bt = bt
	s.cabinet = cab
	s.subscribe = sub

	// 检查可用电池型号
	if !cs.ModelInclude(cab, sub.Model) {
		snag.Panic("电池型号不兼容")
	}

	// 查询电柜是否可使用
	NewCabinet().BusinessableX(cab)

	// 转换业务
	bus, _ := NewBusiness().Convert(bt)

	// 查询电柜业务是否可被允许
	// 代理商管理的电柜只允许非本代理商骑手进行换电业务
	if cab.EnterpriseID != nil && cab.EnterpriseID != sub.EnterpriseID && bus != adapter.BusinessExchange {
		snag.Panic("该电柜无法办理业务, 请前往其他电柜进行")
	}

	// 验证是否可以办理业务
	s.battery = NewIntelligentCabinet(s.rider).BusinessCensorX(bus, sub, cab)

	// 获取仓位信息
	var err error
	s.response.UUID, s.response.Index, err = NewIntelligentCabinet(s.rider).BusinessUsable(cab, bus, sub.Model)
	if err != nil {
		snag.Panic(err)
	}
}

// Active 骑手自主激活
// TODO 分配信息是否需要记录电池编号
func (s *riderBusinessService) Active(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
	// 预处理
	s.preprocess(req.Serial, business.TypeActive)

	// 检查骑士卡状态
	if s.subscribe.Status != model.SubscribeStatusInactive {
		snag.Panic("骑士卡状态错误")
	}

	// 检查是否需要签约
	if NewSubscribe().NeedContract(s.subscribe) {
		// 查询分配信息是否存在, 如果存在则删除
		NewAllocate().SubscribeDeleteIfExists(s.subscribe.ID)

		// 存储分配信息
		err := ent.Database.Allocate.Create().
			SetType(allocate.TypeBattery).
			SetSubscribe(s.subscribe).
			SetRider(s.rider).
			SetStatus(model.AllocateStatusPending.Value()).
			SetTime(time.Now()).
			SetModel(s.subscribe.Model).
			SetCabinetID(s.cabinet.ID).
			SetRemark("用户自主扫码").
			Exec(s.ctx)
		if err != nil {
			snag.Panic("请求失败")
		}

		// 返回签约URL
		snag.Panic(snag.StatusRequireSign, NewContractWithRider(s.rider).Sign(&model.ContractSignReq{
			SubscribeID: s.subscribe.ID,
		}))
	}

	// 查找分配信息
	allo, _ := ent.Database.Allocate.Query().Where(
		allocate.SubscribeID(s.subscribe.ID),
		allocate.RiderID(s.subscribe.RiderID),
		// allocate.Status(model.AllocateStatusSigned.Value()),
		allocate.CabinetIDNotNil(),
	).First(s.ctx)
	if allo == nil {
		snag.Panic("未找到分配信息")
	}

	NewBusinessRider(s.rider).
		SetCabinet(s.cabinet).
		SetCabinetTask(func() (*model.BinInfo, *model.Battery, error) {
			// 更新分配信息
			_ = allo.Update().SetStatus(model.AllocateStatusSigned.Value()).SetCabinetID(s.cabinet.ID).Exec(s.ctx)
			return NewIntelligentCabinet(s.rider).DoBusiness(s.response.UUID, adapter.BusinessActive, s.subscribe, nil, s.cabinet)
		}).
		Active(s.subscribe, allo)

	return s.response
}

func (s *riderBusinessService) Continue(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
	s.preprocess(req.Serial, business.TypeContinue)
	if s.subscribe.Status != model.SubscribeStatusPaused {
		snag.Panic("骑士卡状态错误")
	}

	// ↓ 2023-01-02 添加了异步操作
	go func() {
		err := snag.WithPanic(func() {
			// ↑ 2023-01-02 添加了异步操作
			NewBusinessRider(s.rider).
				SetCabinet(s.cabinet).
				SetCabinetTask(func() (*model.BinInfo, *model.Battery, error) {
					return NewIntelligentCabinet(s.rider).DoBusiness(s.response.UUID, adapter.BusinessContinue, s.subscribe, nil, s.cabinet)
				}).
				Continue(req.ID)
			// ↓ 2023-01-02 添加了异步操作
		})

		if err != nil {
			zap.L().Error("骑手取消寄存业务更新失败", zap.Error(err))
		}
	}()
	// ↑ 2023-01-02 添加了异步操作

	return s.response
}

func (s *riderBusinessService) Unsubscribe(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
	s.preprocess(req.Serial, business.TypeUnsubscribe)

	if s.subscribe.Status != model.SubscribeStatusUsing {
		snag.Panic("骑士卡状态异常")
	}

	go func() {
		err := snag.WithPanic(func() {
			NewBusinessRider(s.rider).
				SetCabinet(s.cabinet).
				SetCabinetTask(func() (*model.BinInfo, *model.Battery, error) {
					return NewIntelligentCabinet(s.rider).DoBusiness(s.response.UUID, adapter.BusinessUnsubscribe, s.subscribe, s.battery, s.cabinet)
				}).
				UnSubscribe(req.ID)
		})

		if err != nil {
			zap.L().Error("骑手取退租业务更新失败", zap.Error(err))
		}
	}()

	return s.response
}

func (s *riderBusinessService) Pause(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
	s.preprocess(req.Serial, business.TypePause)
	if s.subscribe.Status != model.SubscribeStatusUsing {
		snag.Panic("骑士卡未在计费中")
	}

	go func() {
		err := snag.WithPanic(func() {
			NewBusinessRider(s.rider).
				SetCabinetTask(func() (*model.BinInfo, *model.Battery, error) {
					return NewIntelligentCabinet(s.rider).DoBusiness(s.response.UUID, adapter.BusinessPause, s.subscribe, s.battery, s.cabinet)
				}).
				SetCabinet(s.cabinet).
				Pause(req.ID)
		})
		if err != nil {
			zap.L().Error("骑手取寄存业务更新失败", zap.Error(err))
		}
	}()

	return s.response
}

// Status 业务操作状态
func (s *riderBusinessService) Status(req *model.BusinessCabinetStatusReq) (res model.BusinessCabinetStatusRes) {
	start := time.Now()
	for {
		// 通过ID类型判断是否智能柜业务
		// 尝试解析 xid, 若成功解析, 则是非智能柜业务
		taskID, err := xid.FromString(req.UUID)
		if err == nil {
			t := ec.QueryID(taskID)
			if t == nil {
				snag.Panic("未找到业务操作")
			}
			if t.Status == model.TaskStatusFail || t.Status == model.TaskStatusSuccess {
				res.Success = t.Status == model.TaskStatusSuccess
				res.Stop = true
				res.Message = t.Message
			}
		} else {
			_ = cache.Get(s.ctx, req.UUID).Scan(&res)
		}

		if res.Stop || time.Since(start).Seconds() > 30 {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// Executable 业务是否可执行
func (s *riderBusinessService) Executable(sub *ent.Subscribe, typ business.Type) bool {
	if sub == nil {
		return false
	}
	switch typ {
	case business.TypePause, business.TypeUnsubscribe:
		return sub.Status == model.SubscribeStatusUsing
	case business.TypeActive:
		return sub.Status == model.SubscribeStatusInactive
	case business.TypeContinue:
		return sub.Status == model.SubscribeStatusPaused
	}
	return false
}

// PauseInfo 寄存信息
func (s *riderBusinessService) PauseInfo() (res model.BusinessPauseInfoRes) {
	// 查找订阅信息
	sub := NewSubscribe().Recent(s.rider.ID)
	if sub == nil {
		snag.Panic("无生效中的骑行卡")
	}

	// 查找寄存信息
	p, _ := ent.Database.SubscribePause.QueryNotDeleted().
		Where(subscribepause.SubscribeID(sub.ID), subscribepause.EndAtIsNil()).
		Order(ent.Desc(subscribepause.FieldCreatedAt)).
		First(s.ctx)

	if p == nil {
		snag.Panic("未找到寄存信息")
	}

	res = model.BusinessPauseInfoRes{
		Days:      p.Days,
		Overdue:   p.OverdueDays,
		Remaining: sub.Remaining,
	}
	start := p.StartAt
	// 判断寄存开始日期
	if carbon.Time2Carbon(start).Timestamp() != carbon.Time2Carbon(start).StartOfDay().Timestamp() {
		start = carbon.Time2Carbon(start).Tomorrow().StartOfDay().Carbon2Time()
	}
	res.Start = start.Format(carbon.DateLayout)
	now := carbon.Now()
	if now.Timestamp() != now.StartOfDay().Timestamp() {
		now = now.Yesterday()
	}
	res.End = now.Carbon2Time().Format(carbon.DateLayout)

	if p.Days < 1 {
		res.Start = ""
		res.End = ""
		p.Days = 0
	}

	return
}
