// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-05
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/async"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/ec"
	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/provider"
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
	task      *ec.Task
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
		snag.Panic("电池型号不匹配，请更换电柜重试")
	}

	// 限制换电城市
	if cab != nil && sub.PlanID != nil {
		citys, err := NewPlan().PlanCity(*sub.PlanID)
		if err != nil {
			snag.Panic("未找到套餐")
		}
		if !NewRiderBusiness(s.rider).IsCabinetCityInCities(citys, *cab.CityID) {
			snag.Panic("请在指定城市办理业务")
		}
	}

	var (
		info  *model.RiderCabinetOperateProcess
		fully *model.BinInfo
		uid   string
	)

	// 查询电柜是否可使用
	NewCabinet().BusinessableX(cab)

	// 判断设备是否智能设备
	if cab.UsingMicroService() {
		// 判定是否可以换电
		NewIntelligentCabinet(s.rider).BusinessCensorX(adapter.BusinessExchange, sub, cab)
		uid, info = NewIntelligentCabinet(s.rider).ExchangeUsable(sub.Model, cab)
	} else {
		// 更新一次电柜状态
		err := cs.UpdateStatus(cab)
		if err != nil {
			snag.Panic("电柜状态获取失败")
		}

		info = cs.Usable(cab)

		ec.BusyX(cab.Serial)

		if info.EmptyBin == nil || (info.FullBin == nil && info.Alternative == nil) {
			snag.Panic("电柜仓位不可用")
		}

		if info.Alternative != nil {
			fully = &model.BinInfo{
				Index:       info.Alternative.Index,
				Electricity: info.Alternative.Electricity,
				Voltage:     info.Alternative.Voltage,
			}
		} else {
			fully = &model.BinInfo{
				Index:       info.FullBin.Index,
				Electricity: info.FullBin.Electricity,
				Voltage:     info.FullBin.Voltage,
			}
		}

		t := &ec.Task{
			Serial:    cab.Serial,
			CabinetID: cab.ID,
			Job:       model.JobExchange,
			Rider: &ec.Rider{
				ID:    s.rider.ID,
				Name:  s.rider.Name,
				Phone: s.rider.Phone,
			},
			Cabinet: cab.GetTaskInfo(),
			Exchange: &ec.Exchange{
				Model:       sub.Model,
				Alternative: info.Alternative != nil,
				Empty: &model.BinInfo{
					Index: info.EmptyBin.Index,
				},
				Fully: fully,
			},
		}

		t = t.Create()

		uid = t.ID
	}

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

	if cab.UsingMicroService() {
		cache.Set(s.ctx, uid, res, 10*time.Minute)
	}

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
		tcab  *ec.Cabinet
		tex   *ec.Exchange
		t     *ec.Task
		bat   *ent.Asset
		batSN *string
	)

	// 尝试从缓存获取智能电柜换电信息
	err := cache.Get(s.ctx, req.UUID).Scan(&info)

	// 判断设备是否智能设备
	if err == nil {
		cab = NewCabinet().QueryOneSerialX(info.Serial)

		// 查询电柜是否可使用
		NewCabinet().BusinessableX(cab)

		bat = NewIntelligentCabinet(s.rider).BusinessCensorX(adapter.BusinessExchange, sub, cab)
		if bat != nil {
			batSN = silk.String(bat.Sn)
		}

		tex = &ec.Exchange{
			Model:       sub.Model,
			Alternative: info.Alternative != nil,
			Empty: &model.BinInfo{
				Index: info.EmptyBin.Index,
			},
		}

		tcab = cab.GetTaskInfo()
		tex.Fully = info.FullBin
		if info.Alternative != nil {
			tex.Fully = info.Alternative
		}
	} else {
		// 查找任务
		var uid xid.ID
		uid, err = xid.FromString(req.UUID)
		if err != nil {
			snag.Panic("换电任务获取失败, 请重新扫码")
		}

		t = ec.QueryID(uid)

		// 判断任务是否存在, 并且比对存储骑手信息是否相符
		if t == nil || t.Status > 0 || t.StartAt != nil || t.Job != model.JobExchange || t.Exchange == nil || t.Rider == nil || t.Rider.ID != s.rider.ID {
			snag.Panic("未找到信息, 请重新扫码")
		}

		cab = NewCabinet().QueryOneSerialX(t.Serial)
		var be model.BatterySoc
		if t.Exchange.Alternative && !req.Alternative {
			snag.Panic("非满电，换电取消")
		}

		// 更新一次电柜状态
		err = NewCabinet().UpdateStatus(cab)
		if err != nil {
			snag.Panic("电柜状态获取失败")
		}

		// 检查电柜是否繁忙
		if x := ec.Obtain(ec.ObtainReq{Serial: cab.Serial}); x != nil && x.ID != uid.String() {
			snag.Panic("电柜忙, 请稍后重试")
		}

		// 查询电柜是否可使用
		NewCabinet().BusinessableX(cab)

		s.logger = logging.NewExchangeLog(s.rider.ID, t.ID, cab.Serial, s.rider.Phone, be.IsBatteryFull())
		s.cabinet = cab
		s.task = t

		tcab = t.Cabinet
		tex = t.Exchange
	}

	if cab == nil || cab.CityID == nil {
		snag.Panic("未找到电柜信息, 请重试")
		return
	}

	// 记录换电人
	// TODO 超时处理
	s.exchange, _ = ent.Database.Exchange.
		Create().
		SetRiderID(s.rider.ID).
		SetCityID(*cab.CityID).
		// SetInfo(&model.ExchangeInfo{
		//     Cabinet: &model.ExchangeCabinetInfo{
		//         Health:         tcab.Health,
		//         Doors:          tcab.Doors,
		//         BatteryNum:     tcab.BatteryNum,
		//         BatteryFullNum: tcab.BatteryFullNum,
		//     },
		//     Empty: tex.Empty,
		//     Fully: tex.Fully,
		//     Steps: tex.Steps,
		// }).
		SetCabinetInfo(&model.ExchangeCabinetInfo{
			Health:         tcab.Health,
			Doors:          tcab.Doors,
			BatteryNum:     tcab.BatteryNum,
			BatteryFullNum: tcab.BatteryFullNum,
		}).
		SetEmpty(tex.Empty).
		SetFully(tex.Fully).
		SetUUID(req.UUID).
		SetCabinetID(cab.ID).
		SetSuccess(false).
		SetModel(s.subscribe.Model).
		SetNillableEnterpriseID(s.subscribe.EnterpriseID).
		SetNillableStationID(s.subscribe.StationID).
		SetSubscribeID(s.subscribe.ID).
		SetAlternative(tex.Alternative).
		SetStartAt(time.Now()).
		SetNillableRiderBattery(batSN).
		Save(s.ctx)

	if s.exchange == nil {
		snag.Panic("换电失败")
	}

	// 异步进行换电任务并存储换电任务
	go async.WithTask(func() {
		NewIntelligentCabinet(s.rider).Exchange(req.UUID, s.exchange, sub, bat, cab)
	})
}

// ProcessNextStep 开始下一个步骤
func (s *riderExchangeService) ProcessNextStep() *riderExchangeService {
	s.task.Update(func(task *ec.Task) {
		task.Exchange.StartNextStep()
	})
	return s
}

// ProcessStepEnd 结束换电流程
func (s *riderExchangeService) ProcessStepEnd() {
	status := model.TaskStatusFail
	if s.task.Exchange.IsSuccess() {
		status = model.TaskStatusSuccess
	}

	if r := recover(); r != nil {
		zap.L().Error("换电异常结束 -> ["+s.task.ID+
			": "+s.cabinet.Serial+
			"] "+s.task.Exchange.CurrentStep().String()+
			"", zap.Error(fmt.Errorf("%v", r)))
		s.task.Message = fmt.Sprintf("%v", r)
		status = model.TaskStatusFail
	}

	s.task.Update(func(task *ec.Task) {
		task.Stop(status)
	})

	now := time.Now()

	// 保存数据库
	_, _ = s.exchange.Update().
		SetRiderID(s.rider.ID).
		SetCityID(*s.cabinet.CityID).
		SetFully(s.task.Exchange.Fully).
		SetEmpty(s.task.Exchange.Empty).
		SetSteps(s.task.Exchange.Steps).
		SetMessage(s.task.Message).
		SetUUID(s.task.ID).
		SetCabinetID(s.cabinet.ID).
		SetSuccess(status == model.TaskStatusSuccess).
		SetModel(s.subscribe.Model).
		SetNillableEnterpriseID(s.subscribe.EnterpriseID).
		SetNillableStationID(s.subscribe.StationID).
		SetSubscribeID(s.subscribe.ID).
		SetAlternative(s.task.Exchange.Alternative).
		SetNillableStartAt(s.task.StartAt).
		SetFinishAt(now).
		SetDuration(int(s.task.StopAt.Sub(*s.task.StartAt).Seconds())).
		Save(s.ctx)
}

// ProcessByStep 按步骤换电操作
func (s *riderExchangeService) ProcessByStep() {
	defer s.ProcessStepEnd()

	// 第一步: 开启空电仓
	if !s.ProcessOpenBin().ProcessLog() {
		return
	}
	// 手动延时处理
	time.Sleep(5 * time.Second)

	// 第二步: 长轮询判断仓门是否关闭
	if !s.ProcessNextStep().ProcessDoorStatus().ProcessLog() {
		return
	}

	// 第三步: 开启满电仓
	if !s.ProcessNextStep().ProcessOpenBin().ProcessLog() {
		return
	}
	// 手动延时处理
	time.Sleep(5 * time.Second)

	// 第四步: 长轮询判断仓门是否关闭
	if !s.ProcessNextStep().ProcessDoorStatus().ProcessLog() {
		return
	}
}

// ProcessDoorBatteryStatus 格式化仓门状态, 电池放入取出检测
func (s *riderExchangeService) ProcessDoorBatteryStatus() (ds ec.DoorStatus) {
	// 获取仓位
	bin := s.task.Exchange.CurrentBin()

	// 获取步骤
	step := s.task.Exchange.CurrentStep()

	// 获取仓门状态
	ds = NewCabinet().DoorOpenStatus(s.cabinet, bin.Index)

	// 当前仓位信息
	cbin := s.cabinet.Bin[bin.Index]
	pe := cbin.Electricity
	pv := cbin.Voltage

	// zap.L().Info("[电柜操作 - 仓门检测]: [ " + s.cabinet.Serial +
	//     " ] " + step.String() +
	//     ", 用户电话: " + s.rider.Phone +
	//     ", 仓位: " + strconv.Itoa(bin.Index+1) +
	//     "号仓, 仓门状态: " + ds.String() +
	//     ", 是否有电池: " + adapter.Or(cbin.Battery, "是", "否") +
	//     ", 电池信息: " + fmt.Sprintf("%.2f%%[%.2fV]", pe, pv),
	// )

	// 当仓门未关闭时跳过
	if ds != ec.DoorStatusClose {
		return
	}

	// 关门时间
	if step.Time.IsZero() {
		step.Time = time.Now()
	}

	// 验证是否放入旧电池
	if step.Step == model.ExchangeStepPutInto {
		// 获取骑手放入电池信息
		if s.task.Exchange.Empty.Electricity == 0 {
			s.task.Exchange.Empty.Electricity = pe
		}

		if s.task.Exchange.Empty.Voltage < 40 {
			s.task.Exchange.Empty.Voltage = pv
		}

		// 曹博文说: 判断是否 有电池 并且 (电压大于40 或 电量大于0)
		if cbin.Battery && (pv > 45 || pe > 0) {
			return ec.DoorStatusClose
		}

		// 仓门关闭但是检测不到电池的情况下, 继续检测30s
		if time.Since(step.Time).Seconds() > 30 {
			return ec.DoorStatusBatteryEmpty
		} else {
			time.Sleep(1 * time.Second)
			return s.ProcessDoorBatteryStatus()
		}
	}

	// 验证满电电池是否取走
	if step.Step == model.ExchangeStepPutOut {
		// 如果已取走直接返回
		if !cbin.Battery {
			return ec.DoorStatusClose
		}

		// 仓门关闭, 如果未取走则继续检测10s
		if time.Since(step.Time).Seconds() > 10 {
			return ec.DoorStatusBatteryFull
		} else {
			time.Sleep(1 * time.Second)
			return s.ProcessDoorBatteryStatus()
		}
	}

	return ds
}

// ProcessDoorStatus 操作换电中检查柜门并处理状态
func (s *riderExchangeService) ProcessDoorStatus() *riderExchangeService {
	start := time.Now()
	step := s.task.Exchange.CurrentStep()

	for {
		// 检测仓门/电池
		ds := s.ProcessDoorBatteryStatus()
		if ds == ec.DoorStatusClose {
			// 强制睡眠两秒: 原因是有可能柜门会晃动导致似关非关, 延时来获取正确状态
			time.Sleep(2 * time.Second)
			ds = s.ProcessDoorBatteryStatus()
		}

		var message string

		switch ds {
		case ec.DoorStatusClose:
			step.Status = model.TaskStatusSuccess
		case ec.DoorStatusOpen:
			break
		default:
			message = ec.DoorError[ds]
			step.Status = model.TaskStatusFail
		}

		// 超时标记为任务失败
		if time.Since(start).Seconds() > s.maxTime.Seconds() && message == "" {
			message = "超时"
			step.Status = model.TaskStatusFail
			step.Time = time.Now()
		}

		if step.Status != model.TaskStatusProcessing {
			if !step.IsSuccess() {
				s.task.Message = message
				s.task.Stop(model.TaskStatusFail)
			}
			return s
		}

		time.Sleep(1 * time.Second)
	}
}

// ProcessOpenBin 开仓门
func (s *riderExchangeService) ProcessOpenBin() *riderExchangeService {
	bin := s.task.Exchange.CurrentBin()
	step := s.task.Exchange.CurrentStep()

	r := s.rider
	operator := model.CabinetDoorOperator{
		ID:    r.ID,
		Role:  model.CabinetDoorOperatorRoleRider,
		Name:  r.Name,
		Phone: r.Phone,
	}

	// 开始处理
	reason := model.RiderCabinetOperateReasonEmpty
	if step.Step == model.ExchangeStepOpenFull {
		reason = model.RiderCabinetOperateReasonFull
	}

	operation := model.CabinetDoorOperateOpen
	id := s.cabinet.ID
	index := silk.Pointer(bin.Index)

	status, err := NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
		ID:        id,
		Index:     index,
		Remark:    fmt.Sprintf("骑手换电 - %s", reason),
		Operation: &operation,
	}, operator)
	if err != nil {
		zap.L().Error("仓门开启失败", zap.Error(err))
	}

	s.task.Update(func(t *ec.Task) {
		step.Time = time.Now()
		if status {
			step.Status = model.TaskStatusSuccess
		} else {
			step.Status = model.TaskStatusFail
			t.Message = err.Error()
			t.Stop(model.TaskStatusFail)
		}
	})

	// log.Infof(`[电柜操作 - 开启仓门]: [ %s ] %s, 用户电话: %s, 仓位: %d号仓, 操作反馈: %t`,
	//     s.cabinet.Serial,
	//     step,
	//     s.rider.Phone,
	//     bin.Index+1,
	//     status,
	// )

	provider.AutoBinFault(operator, s.cabinet, bin.Index, status, func() {
		_, _ = NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
			ID:        id,
			Index:     index,
			Remark:    fmt.Sprintf("换电仓门处理失败自动锁仓 - %s", s.rider.Phone),
			Operation: silk.Pointer(model.CabinetDoorOperateLock),
		}, operator)
	})

	return s
}

// ProcessLog 处理步骤日志
func (s *riderExchangeService) ProcessLog() bool {
	ex := s.task.Exchange
	stop := "否"
	if ex.IsLastStep() || s.task.StopAt != nil {
		stop = "是"
	}
	zap.L().Info("[电柜操作 - 步骤结果]: [ " + s.cabinet.Serial +
		" ] " + s.task.Exchange.CurrentStep().String() +
		", 用户电话: " + s.rider.Phone +
		", 状态: " + ex.CurrentStep().Status.String() +
		", 终止: " + stop + " ->> " + s.task.Message,
	)

	step := ex.CurrentStep()

	s.logger.Clone().
		SetBin(ex.CurrentBin().Index).
		SetStatus(ex.CurrentStep().Status).
		SetMessage(s.task.Message).
		SetStep(ex.CurrentStep().Step).
		SetElectricity(ex.CurrentBin().Electricity).
		Send()

	return step.IsSuccess()
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
			Status:  cs.Status,
			Message: task.Message,
			Stop:    task.StopAt != nil,
		}
		if cs.IsSuccess() || res.Stop || time.Since(start).Seconds() > 30 {
			return
		}
	}
}
