// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-05
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/provider"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    log "github.com/sirupsen/logrus"
    "math"
    "time"
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
    s.ctx = context.WithValue(context.Background(), "rider", r)
    return s
}

// GetProcess 获取待换电信息
func (s *riderExchangeService) GetProcess(req *model.RiderCabinetOperateInfoReq) *model.RiderExchangeInfo {
    sm, _ := NewSetting().GetSetting(model.SettingMaintain).(bool)
    if sm {
        snag.Panic("系统维护中, 请稍后重试")
    }

    // 检查用户换电间隔
    iv := cache.Int(model.SettingExchangeInterval)
    if exist, _ := ent.Database.Exchange.QueryNotDeleted().Where(
        exchange.RiderID(s.rider.ID),
        exchange.Success(true),
    ).First(s.ctx); exist != nil {
        m := int(math.Ceil(time.Now().Sub(exist.FinishAt).Minutes()))
        if iv-m > 0 {
            snag.Panic(fmt.Sprintf("换电过于频繁, %d分钟可再次换电", iv-m))
        }
    }
    // 检查用户是否可以办理业务
    NewRiderPermissionWithRider(s.rider).BusinessX()

    // 是否有生效中套餐
    subd, _ := NewSubscribe().RecentDetail(s.rider.ID)
    if subd == nil || subd.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效中的骑行卡")
    }

    // 查询电柜
    cs := NewCabinet()
    cab := cs.QueryOneSerialX(req.Serial)

    // 检查可用电池型号
    if !cs.ModelInclude(cab, subd.Model) {
        snag.Panic("电池型号不兼容")
    }

    // 查询电柜
    if !cs.Businessable(cab) {
        snag.Panic("电柜目前不可用")
    }

    ec.BusyX(cab.Serial)

    // 更新一次电柜状态
    err := cs.UpdateStatus(cab)
    if err != nil {
        log.Error(err)
        snag.Panic("电柜状态获取失败")
    }

    info := cs.Usable(cab)
    if info.EmptyBin == nil || (info.FullBin == nil && info.Alternative == nil) {
        snag.Panic("电柜仓位不可用")
    }

    var fully *ec.BinInfo

    if info.Alternative != nil {
        fully = &ec.BinInfo{
            Index:       info.Alternative.Index,
            Electricity: info.Alternative.Electricity,
            Voltage:     info.Alternative.Voltage,
        }
    } else {
        fully = &ec.BinInfo{
            Index:       info.FullBin.Index,
            Electricity: info.FullBin.Electricity,
            Voltage:     info.FullBin.Voltage,
        }
    }

    task := &ec.Task{
        Serial:    cab.Serial,
        CabinetID: cab.ID,
        Job:       ec.JobExchange,
        Rider: &ec.Rider{
            ID:    s.rider.ID,
            Name:  s.rider.Edges.Person.Name,
            Phone: s.rider.Phone,
        },
        Cabinet: ec.Cabinet{
            Health:         cab.Health,
            Doors:          cab.Doors,
            BatteryNum:     cab.BatteryNum,
            BatteryFullNum: cab.BatteryFullNum,
        },
        Exchange: &ec.Exchange{
            Model:       subd.Model,
            Alternative: info.Alternative != nil,
            Empty: &ec.BinInfo{
                Index: info.EmptyBin.Index,
            },
            Fully: fully,
        },
    }
    task = task.CreateX()

    // TODO 修改前端返回值
    res := &model.RiderExchangeInfo{
        ID:                         cab.ID,
        UUID:                       task.ID.Hex(),
        Full:                       info.FullBin != nil,
        Name:                       cab.Name,
        Health:                     cab.Health,
        Serial:                     cab.Serial,
        Doors:                      cab.Doors,
        BatteryNum:                 cab.BatteryNum,
        BatteryFullNum:             cab.BatteryFullNum,
        RiderCabinetOperateProcess: info,
        Model:                      subd.Model,
        CityID:                     subd.City.ID,
        Brand:                      model.CabinetBrand(cab.Brand),
    }

    tools.NewLog().Infof("[换电信息:%s]\n%s\n", task.ID.Hex(), res)

    return res
}

// Start 开始换电
func (s *riderExchangeService) Start(req *model.RiderExchangeProcessReq) {
    // 检查是否维护中
    sm, _ := NewSetting().GetSetting(model.SettingMaintain).(bool)
    if sm {
        snag.Panic("系统维护中, 请稍后重试")
    }

    // 检查用户是否可以办理业务
    NewRiderPermissionWithRider(s.rider).BusinessX()

    // 校验换电信息
    iv := cache.Int(model.SettingExchangeInterval)
    if exist, _ := ent.Database.Exchange.QueryNotDeleted().Where(
        exchange.RiderID(s.rider.ID),
        exchange.Success(true),
        exchange.CreatedAtGTE(time.Now().Add(-time.Duration(cache.Int(model.SettingExchangeInterval))*time.Minute)),
    ).Exist(s.ctx); exist {
        snag.Panic(fmt.Sprintf("换电过于频繁, %d分钟可再次换电", iv))
    }

    // 查找任务
    task := ec.QueryID(req.UUID)

    // TODO 存储骑手信息并比对骑手信息是否相符
    if task == nil || task.Status > 0 || task.StartAt != nil || task.Job != ec.JobExchange || task.Exchange == nil {
        snag.Panic("未找到信息, 请重新扫码")
    }

    // 是否有生效中套餐
    subd, sub := NewSubscribe().RecentDetail(s.rider.ID)
    if subd == nil || subd.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效中的骑行卡")
    }

    s.subscribe = sub

    cab := NewCabinet().QueryOneSerialX(task.Serial)
    var be model.BatteryElectricity
    if task.Exchange.Alternative && !req.Alternative {
        snag.Panic("非满电换电取消")
    }

    // 检查电柜是否繁忙
    if x := ec.Obtain(ec.ObtainReq{Serial: cab.Serial}); x != nil && x.ID != req.UUID {
        snag.Panic("电柜忙, 请稍后重试")
    }

    // 查询电柜
    if !NewCabinet().Businessable(cab) {
        snag.Panic("电柜目前不可用")
    }

    s.logger = logging.NewExchangeLog(s.rider.ID, task.ID.Hex(), cab.Serial, s.rider.Phone, be.IsBatteryFull())
    s.cabinet = cab

    // 开始任务
    task.Start(func(task *ec.Task) {
        task.Exchange.Steps = []*ec.ExchangeStepInfo{
            {Step: ec.ExchangeStepOpenEmpty, Time: time.Now()},
        }
    })

    // 记录换电人
    // TODO 超时处理
    s.exchange, _ = ent.Database.Exchange.
        Create().
        SetRiderID(s.rider.ID).
        SetCityID(*cab.CityID).
        SetInfo(&ec.ExchangeInfo{
            Cabinet:  task.Cabinet,
            Exchange: task.Exchange,
        }).
        SetUUID(task.ID.Hex()).
        SetCabinetID(cab.ID).
        SetSuccess(false).
        SetModel(s.subscribe.Model).
        SetNillableEnterpriseID(s.subscribe.EnterpriseID).
        SetNillableStationID(s.subscribe.StationID).
        SetSubscribeID(s.subscribe.ID).
        SetAlternative(task.Exchange.Alternative).
        SetStartAt(*task.StartAt).
        Save(s.ctx)

    if s.exchange == nil {
        snag.Panic("换电失败")
    }

    // 处理换电流程
    s.task = task
    go s.ProcessByStep()
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
    status := ec.TaskStatusFail
    if s.task.Exchange.IsSuccess() {
        status = ec.TaskStatusSuccess
    }

    if r := recover(); r != nil {
        log.Errorf("换电异常结束 -> [%s: %s] %s: %v", s.task.ID.Hex(), s.cabinet.Serial, s.task.Exchange.CurrentStep(), r)
        s.task.Message = fmt.Sprintf("%v", r)
        status = ec.TaskStatusFail
    }

    s.task.Update(func(task *ec.Task) {
        task.Stop(status)
    })

    now := time.Now()

    // 保存数据库
    _, _ = s.exchange.Update().
        SetRiderID(s.rider.ID).
        SetCityID(*s.cabinet.CityID).
        SetInfo(&ec.ExchangeInfo{
            Cabinet:  s.task.Cabinet,
            Exchange: s.task.Exchange,
            Message:  s.task.Message,
        }).
        SetUUID(s.task.ID.Hex()).
        SetCabinetID(s.cabinet.ID).
        SetSuccess(status == ec.TaskStatusSuccess).
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
    ds = NewCabinet().DoorOpenStatus(s.cabinet, bin.Index, true)

    // 当前仓位信息
    cbin := s.cabinet.Bin[bin.Index]
    pe := cbin.Electricity
    pv := cbin.Voltage

    log.Infof(`[电柜操作 - 仓门检测]: [ %s ] %s, 用户电话: %s, 仓位: %d号仓, 仓门状态: %s, 是否有电池: %t, 电池信息: %.2f%%[%.2fV]`,
        s.cabinet.Serial,
        step,
        s.rider.Phone,
        bin.Index+1,
        ds,
        cbin.Battery,
        pe,
        pv,
    )

    // 当仓门未关闭时跳过
    if ds != ec.DoorStatusClose {
        return
    }

    // 关门时间
    if step.Time.IsZero() {
        step.Time = time.Now()
    }

    // 验证是否放入旧电池
    if step.Step == ec.ExchangeStepPutInto {
        // 获取骑手放入电池信息
        if s.task.Exchange.Empty.Electricity == 0 {
            s.task.Exchange.Empty.Electricity = pe
        }

        if s.task.Exchange.Empty.Voltage < 40 {
            s.task.Exchange.Empty.Voltage = pv
        }

        // 判断是否 有电池 并且 (电压大于40 或 电量大于0)
        if cbin.Battery && (pv > 40 || pe > 0) {
            return ec.DoorStatusClose
        }

        // 仓门关闭但是检测不到电池的情况下, 继续检测30s
        if time.Now().Sub(step.Time).Seconds() > 30 {
            return ec.DoorStatusBatteryEmpty
        } else {
            time.Sleep(1 * time.Second)
            return s.ProcessDoorBatteryStatus()
        }
    }

    // 验证满电电池是否取走
    if step.Step == ec.ExchangeStepPutOut {
        // 如果已取走直接返回
        if !cbin.Battery {
            return ec.DoorStatusClose
        }

        // 仓门关闭, 如果未取走则继续检测10s
        if time.Now().Sub(step.Time).Seconds() > 10 {
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
            step.Status = ec.TaskStatusSuccess
            break
        case ec.DoorStatusOpen:
            break
        default:
            message = ec.DoorError[ds]
            step.Status = ec.TaskStatusFail
            break
        }

        // 超时标记为任务失败
        if time.Now().Sub(start).Seconds() > s.maxTime.Seconds() && message == "" {
            message = "超时"
            step.Status = ec.TaskStatusFail
            step.Time = time.Now()
        }

        if step.Status != ec.TaskStatusProcessing {
            if !step.IsSuccess() {
                s.task.Message = message
                s.task.Stop(ec.TaskStatusFail)
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
        Name:  r.Edges.Person.Name,
        Phone: r.Phone,
    }

    // 开始处理
    reason := model.RiderCabinetOperateReasonEmpty
    if step.Step == ec.ExchangeStepOpenFull {
        reason = model.RiderCabinetOperateReasonFull
    }

    operation := model.CabinetDoorOperateOpen
    id := s.cabinet.ID
    index := tools.NewPointerInterface(bin.Index)

    status, err := NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
        ID:        &id,
        Index:     index,
        Remark:    fmt.Sprintf("骑手换电 - %s", reason),
        Operation: &operation,
    }, operator, true)
    if err != nil {
        log.Error(err)
    }

    s.task.Update(func(t *ec.Task) {
        step.Time = time.Now()
        if status {
            step.Status = ec.TaskStatusSuccess
        } else {
            step.Status = ec.TaskStatusFail
            t.Message = err.Error()
            t.Stop(ec.TaskStatusFail)
        }
    })

    log.Infof(`[电柜操作 - 开启仓门]: [ %s ] %s, 用户电话: %s, 仓位: %d号仓, 操作反馈: %t`,
        s.cabinet.Serial,
        step,
        s.rider.Phone,
        bin.Index+1,
        status,
    )

    provider.AutoBinFault(operator, s.cabinet, bin.Index, status, func() {
        _, _ = NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
            ID:        &id,
            Index:     index,
            Remark:    fmt.Sprintf("换电仓门处理失败自动锁仓 - %s", s.rider.Phone),
            Operation: tools.NewPointerInterface(model.CabinetDoorOperateLock),
        }, operator, true)
    })

    return s
}

// ProcessLog 处理步骤日志
func (s *riderExchangeService) ProcessLog() bool {
    ex := s.task.Exchange
    log.Printf(`[电柜操作 - 步骤结果]: [ %s ] %s, 用户电话: %s, 状态: %s, 终止: %t %s`,
        s.cabinet.Serial,
        s.task.Exchange.CurrentStep(),
        s.rider.Phone,
        ex.CurrentStep().Status,
        ex.IsLastStep() || s.task.StopAt != nil,
        s.task.Message,
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
    start := time.Now()
    for {
        task := ec.QueryID(req.UUID)
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
        if cs.IsSuccess() || res.Stop || time.Now().Sub(start).Seconds() > 30 {
            return
        }
        time.Sleep(1 * time.Second)
    }
}
