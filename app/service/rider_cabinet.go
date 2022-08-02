// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/actuator"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/provider"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    log "github.com/sirupsen/logrus"
    "time"
)

// TODO 服务器崩溃后自动启动继续换电进程
// TODO 电柜缓存优化

type riderCabinetService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    maxTime  time.Duration // 单步骤最大处理时长
    logger   *logging.ExchangeLog

    cabinet   *ent.Cabinet
    subscribe *ent.Subscribe
    exchange  *ent.Exchange
    task      *actuator.Task
}

func NewRiderCabinet() *riderCabinetService {
    return &riderCabinetService{
        ctx:     context.Background(),
        maxTime: 180 * time.Second,
    }
}

func NewRiderCabinetWithRider(rider *ent.Rider) *riderCabinetService {
    s := NewRiderCabinet()
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

// GetProcess 获取待换电信息
func (s *riderCabinetService) GetProcess(req *model.RiderCabinetOperateInfoReq) *model.RiderCabinetInfo {
    sm, _ := NewSetting().GetSetting(model.SettingMaintain).(bool)
    if sm {
        snag.Panic("系统维护中, 请稍后重试")
    }

    // 检查用户换电间隔
    iv := cache.Int(model.SettingExchangeInterval)
    if exist, _ := ent.Database.Exchange.QueryNotDeleted().Where(
        exchange.RiderID(s.rider.ID),
        exchange.Success(true),
        exchange.CreatedAtGTE(time.Now().Add(-time.Duration(cache.Int(model.SettingExchangeInterval))*time.Minute)),
    ).Exist(s.ctx); exist {
        snag.Panic(fmt.Sprintf("换电过于频繁, %d分钟可再次换电", iv))
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
    if !cs.Health(cab) {
        snag.Panic("电柜目前不可用")
    }

    // 是否忙
    if actuator.Busy(cab.Serial) {
        snag.Panic("电柜忙, 请稍后")
    }

    // 更新一次电柜状态
    NewCabinet().UpdateStatus(cab)
    info := NewCabinet().Usable(cab)
    if info.EmptyBin == nil || (info.FullBin == nil && info.Alternative == nil) {
        snag.Panic("电柜仓位不可用")
    }

    var fully *actuator.BinInfo

    if info.Alternative == nil {
        fully = &actuator.BinInfo{
            Index:       info.Alternative.Index,
            Electricity: info.Alternative.Electricity,
            Voltage:     info.Alternative.Voltage,
        }
    } else {
        fully = &actuator.BinInfo{
            Index:       info.FullBin.Index,
            Electricity: info.FullBin.Electricity,
            Voltage:     info.FullBin.Voltage,
        }
    }

    // TODO 存储换电信息到MongoDB
    task := &actuator.Task{
        Task: actuator.JobExchange,
        Cabinet: actuator.Cabinet{
            Serial:         cab.Serial,
            Health:         cab.Health,
            Doors:          cab.Doors,
            BatteryNum:     cab.BatteryNum,
            BatteryFullNum: cab.BatteryFullNum,
        },
        Exchange: &actuator.Exchange{
            Model:       subd.Model,
            Alternative: info.Alternative != nil,
            Empty: &actuator.BinInfo{
                Index: info.EmptyBin.Index,
            },
            Fully: fully,
        },
    }
    id := task.CreateX()

    // TODO 修改前端返回值
    res := &model.RiderCabinetInfo{
        ID:                         cab.ID,
        UUID:                       id,
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

    tools.NewLog().Infof("[换电信息:%s]\n%s\n", id, res)

    return res
}

// Start 开始换电
func (s *riderCabinetService) Start(req *model.RiderCabinetOperateReq) {
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
    task := actuator.Obtain(actuator.ObtainReq{ID: req.UUID})

    // TODO 存储骑手信息并比对骑手信息是否相符
    if task == nil || task.Task != actuator.JobExchange || task.Exchange == nil {
        snag.Panic("未找到信息, 请重新扫码")
    }

    // 是否有生效中套餐
    subd, sub := NewSubscribe().RecentDetail(s.rider.ID)
    if subd == nil || subd.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效中的骑行卡")
    }

    s.subscribe = sub

    cab := NewCabinet().QueryOneSerialX(task.Cabinet.Serial)
    var be model.BatteryElectricity
    if task.Exchange.Alternative && !req.Alternative {
        snag.Panic("非满电换电取消")
    }

    // 检查电柜是否繁忙
    if x := actuator.Obtain(actuator.ObtainReq{Serial: cab.Serial}); x.Id != req.UUID {
        snag.Panic("电柜忙, 请稍后重试")
    }

    s.logger = logging.NewExchangeLog(s.rider.ID, task.Id.Hex(), cab.Serial, s.rider.Phone, be.IsBatteryFull())
    s.cabinet = cab

    // 开始任务
    task.Start(func(task *actuator.Task) {
        task.Exchange.Step = actuator.ExchangeStepOpenEmpty
    })

    // 记录换电人
    // TODO 超时处理
    s.exchange, _ = ent.Database.Exchange.
        Create().
        SetRiderID(s.rider.ID).
        SetCityID(*cab.CityID).
        SetInfo(&actuator.ExchangeInfo{
            Cabinet:  task.Cabinet,
            Exchange: task.Exchange,
        }).
        SetUUID(task.Id.Hex()).
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
    go s.ProcessByStep()
}

// ProcessNextStep 下一个步骤
func (s *riderCabinetService) ProcessNextStep() {
    s.task.Update(func(task *actuator.Task) {
        task.Exchange.SetNextStep()
    })
}

// ProcessStepEnd 结束换电流程
func (s *riderCabinetService) ProcessStepEnd() {
    status := actuator.TaskStatusFail
    if s.task.Exchange.CurrentBin().Success {
        status = actuator.TaskStatusSuccess
    }

    if r := recover(); r != nil {
        log.Errorf("换电异常结束 -> %s {%s}: %v", s.task.Id.Hex(), s.task.Exchange.Step, r)
        s.task.Message = fmt.Sprintf("%v", r)
        status = actuator.TaskStatusFail
    }

    s.task.Stop(status)

    now := time.Now()

    // 保存数据库
    _, _ = s.exchange.Update().
        SetRiderID(s.rider.ID).
        SetCityID(*s.cabinet.CityID).
        SetInfo(&actuator.ExchangeInfo{
            Cabinet:  s.task.Cabinet,
            Exchange: s.task.Exchange,
        }).
        SetUUID(s.task.Id.Hex()).
        SetCabinetID(s.cabinet.ID).
        SetSuccess(status == actuator.TaskStatusSuccess).
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
func (s *riderCabinetService) ProcessByStep() {
    defer s.ProcessStepEnd()

    // 第一步: 开启空电仓
    if !s.ProcessLog(s.ProcessOpenBin()) {
        return
    }
    // 手动延时处理
    time.Sleep(5 * time.Second)

    // 第二步: 长轮询判断仓门是否关闭
    s.ProcessNextStep()
    if !s.ProcessLog(s.ProcessDoorStatus()) {
        return
    }

    // 第三步: 开启满电仓
    s.ProcessNextStep()
    if !s.ProcessLog(s.ProcessOpenBin()) {
        return
    }
    // 手动延时处理
    time.Sleep(5 * time.Second)

    // 第四步: 长轮询判断仓门是否关闭
    s.ProcessNextStep()
    if !s.ProcessLog(s.ProcessDoorStatus()) {
        return
    }
}

// ProcessDoorBatteryStatus 格式化仓门状态, 电池放入取出检测
func (s *riderCabinetService) ProcessDoorBatteryStatus() (ds model.CabinetBinDoorStatus) {
    // 获取仓位index
    index := s.operating.EmptyIndex
    step := s.step
    if step == model.RiderCabinetOperateStepPutOut {
        index = s.operating.FullIndex
    }

    // 获取仓门状态
    ds = NewCabinet().DoorOpenStatus(s.cabinet, index, true)
    bin := s.cabinet.Bin[index]

    // 放入电池电量检测
    ebin := s.cabinet.Bin[s.operating.EmptyIndex]
    ee := ebin.Electricity
    if s.operating.RiderElectricity != 0 {
        ee = s.operating.RiderElectricity
    } else {
        s.operating.RiderElectricity = ee
    }

    log.Infof(`[换电步骤 - 仓门检测]: {step:%s} %s - %d, 用户电话: %s, 仓门Index: %d, 仓门状态: %s, 是否有电池: %t, 当前电压: %.2fV, 当前电量: %.2f%%, 放入电池电量: %.2f%%`,
        s.step,
        s.cabinet.Serial,
        index,
        s.rider.Phone,
        bin.Index,
        ds,
        bin.Battery,
        bin.Voltage,
        bin.Electricity,
        ee,
    )

    // 当仓门未关闭时跳过
    if ds != model.CabinetBinDoorStatusClose {
        return ds
    }

    // 验证是否放入旧电池
    if step == model.RiderCabinetOperateStepPutInto {
        // 获取骑手放入电池电量
        s.operating.RiderElectricity = bin.Electricity

        // 放入时间
        if s.emptyCloseTime.IsZero() {
            s.emptyCloseTime = time.Now()
        }

        // 凯信电柜需要绑定电池<已废弃>
        // bind := true
        // if s.info.Brand == model.CabinetBrandKaixin {
        //     bind = provider.NewKaixin().BatteryBind(s.rider.Edges.Person.Name+"-"+shortuuid.New(), s.info.Serial, s.model, s.operating.EmptyIndex)
        // }

        // 判断是否 有电池 并且 (电压大于40 或 电量大于0)
        if bin.Battery && (bin.Voltage > 40 || bin.Electricity > 0) {
            s.putInElectricity = bin.Electricity
            return model.CabinetBinDoorStatusClose
        }

        // 检测不到电池的情况下, 继续检测60s
        if time.Now().Sub(s.emptyCloseTime).Seconds() > 60 {
            return model.CabinetBinDoorStatusBatteryEmpty
        } else {
            time.Sleep(1 * time.Second)
            return s.ProcessDoorBatteryStatus()
        }
    }

    // 验证满电电池是否取走
    if step == model.RiderCabinetOperateStepPutOut {
        if s.fullCloseTime.IsZero() {
            s.fullCloseTime = time.Now()
        }

        // 如果已取走直接返回
        if !bin.Battery {
            return model.CabinetBinDoorStatusClose
        }

        // 如果未取走则继续检测60s
        if time.Now().Sub(s.fullCloseTime).Seconds() > 60 {
            return model.CabinetBinDoorStatusBatteryFull
        } else {
            time.Sleep(1 * time.Second)
            return s.ProcessDoorBatteryStatus()
        }
    }

    return ds
}

// ProcessDoorStatus 操作换电中检查柜门并处理状态
func (s *riderCabinetService) ProcessDoorStatus() (res *model.RiderCabinetOperateRes) {
    res = &model.RiderCabinetOperateRes{}

    start := time.Now()

    for {
        // 检测仓门/电池
        ds := s.ProcessDoorBatteryStatus()
        if ds == model.CabinetBinDoorStatusClose {
            // 强制睡眠两秒: 原因是有可能柜门会晃动导致似关非关, 延时来获取正确状态
            time.Sleep(2 * time.Second)
            ds = s.ProcessDoorBatteryStatus()
        }

        if s.step == model.RiderCabinetOperateStepPutInto {
            s.operating.PutInDoor = ds
        }

        if s.step == model.RiderCabinetOperateStepPutOut {
            s.operating.PutOutDoor = ds
        }

        switch ds {
        case model.CabinetBinDoorStatusClose:
            res.Status = model.TaskStatusSuccess
            return
        case model.CabinetBinDoorStatusOpen:
            break
        default:
            res.Message = model.CabinetBinDoorError[ds]
            res.Stop = true
            res.Status = model.TaskStatusFailFail
            return
        }

        // 超时标记为任务失败
        if time.Now().Sub(start).Seconds() > s.maxTime.Seconds() {
            res.Message = "超时"
            return
        }
        time.Sleep(1 * time.Second)
    }
}

// ProcessOpenBin 开仓门
func (s *riderCabinetService) ProcessOpenBin() (err error) {
    bin := s.task.Exchange.CurrentBin()
    step := s.task.Exchange.Step

    var status bool

    r := s.rider
    operator := model.CabinetDoorOperator{
        ID:    r.ID,
        Role:  model.CabinetDoorOperatorRoleRider,
        Name:  r.Edges.Person.Name,
        Phone: r.Phone,
    }

    // 开始处理
    reason := model.RiderCabinetOperateReasonEmpty
    if step == actuator.ExchangeStepOpenFull {
        reason = model.RiderCabinetOperateReasonFull
    }
    operation := model.CabinetDoorOperateOpen
    id := s.cabinet.ID
    index := tools.NewPointerInterface(bin.Index)
    status, err = NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
        ID:        &id,
        Index:     index,
        Remark:    fmt.Sprintf("骑手换电 - %s", reason),
        Operation: &operation,
    }, operator, true)
    if err != nil {
        snag.Panic(err)
    }

    log.Infof(`[换电步骤 - 仓门检测]: {step:%s} %d, 仓门Index: %d, 操作反馈: %t, 用户电话: %s`,
        step,
        id,
        bin.Index,
        status,
        s.rider.Phone,
    )

    provider.AutoBinFault(operator, s.cabinet, bin.Index, status, func() {
        _, _ = NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
            ID:        &id,
            Index:     index,
            Remark:    fmt.Sprintf("换电仓门处理失败自动锁仓 - %s", s.rider.Phone),
            Operation: tools.NewPointerInterface(model.CabinetDoorOperateLock),
        }, operator, true)
    })

    if status {
        res.Status = model.TaskStatusSuccess
    } else {
        log.Errorf("[ProcessOpenBin] 处理失败: %t -> %#v", status, err)
        // 发生故障处理
        res.Status = model.TaskStatusFailFail
        if err != nil {
            res.Message = err.Error()
        } else {
            res.Message = "柜门处理失败"
        }
        res.Stop = true
    }

    return
}

// ProcessStatus 长轮询获取状态
func (s *riderCabinetService) ProcessStatus(req *model.RiderCabinetOperateStatusReq) (res *model.RiderCabinetOperateRes) {
    start := time.Now()
    for {
        res = new(model.RiderCabinetOperateRes)
        err := cache.Get(s.ctx, *req.UUID).Scan(res)
        if err != nil {
            snag.Panic("未找到换电操作")
            return
        }
        if res.Status == model.TaskStatusSuccess || res.Stop || time.Now().Sub(start).Seconds() > 30 {
            return res
        }
    }
}

func (s *riderCabinetService) currentDoorInfo() (index int, be model.BatteryElectricity) {
    // 前两步是空仓, 后两步是满电仓位
    index = s.operating.EmptyIndex
    be = s.putInElectricity
    if s.step >= model.RiderCabinetOperateStepOpenFull {
        index = s.operating.FullIndex
        be = s.operating.Electricity
    }
    return
}

func (s *riderCabinetService) processLogText() string {
    ex := s.task.Exchange
    return fmt.Sprintf(`[换电步骤 - 结果]: {step:%s} %s -> %s, 用户电话: %s, 状态: %s, 消息: %s, 终止: %t`,
        ex.Step,
        s.cabinet.Serial,
        s.rider.Phone,
        ex.CurrentBin(),
        s.task.Status,
        s.task.Message,
        s.task.StopAt != nil,
    )
}

// ProcessLog 处理步骤日志
func (s *riderCabinetService) ProcessLog(res *model.RiderCabinetOperateRes) bool {
    res.Step = s.step
    res.Stop = s.step == model.RiderCabinetOperateStepPutOut || res.Status == model.TaskStatusFailFail

    index, be := s.currentDoorInfo()

    log.Info(s.processLogText(index, res))

    s.logger.Clone().
        SetBin(index).
        SetStatus(res.Status).
        SetMessage(res.Message).
        SetStep(s.step).
        SetElectricity(be).
        Send()

    // 存储缓存
    err := cache.Set(s.ctx, s.operating.UUID, res, s.maxTime).Err()
    if err != nil {
        log.Error(err)
        return false
    }

    return res.Status == model.TaskStatusSuccess
}
