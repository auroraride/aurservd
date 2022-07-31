// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/provider"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/lithammer/shortuuid/v4"
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

    model            string
    step             model.RiderCabinetOperateStep
    cabinet          *ent.Cabinet
    putInElectricity model.BatteryElectricity
    alternative      bool
    info             *model.RiderCabinetInfo
    operating        *model.RiderCabinetOperating
    batteryNum       uint // 换电开始时电池数量, 业务时应该监听业务发生时的电池数量，当业务流程中电池数量变动大于1的时候视为异常

    subscribe *ent.Subscribe

    start time.Time

    ex *ent.Exchange

    emptyCloseTime time.Time
    fullCloseTime  time.Time
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
    cab := cs.QueryWithSerial(req.Serial)

    // 检查可用电池型号
    if !cs.ModelInclude(cab, subd.Model) {
        snag.Panic("电池型号不兼容")
    }

    // 查询电柜
    if !cs.Health(cab) {
        snag.Panic("电柜目前不可用")
    }
    // 更新一次电柜状态
    NewCabinet().UpdateStatus(cab)
    info := NewCabinet().Usable(cab)
    if info.EmptyBin == nil {
        snag.Panic("电柜仓位不可用")
    }

    uid := shortuuid.New()
    res := &model.RiderCabinetInfo{
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
        Model:                      subd.Model,
        CityID:                     subd.City.ID,
        Brand:                      model.CabinetBrand(cab.Brand),
    }

    tools.NewLog().Infof("[换电信息:%s]\n%s\n", uid, res)

    err := cache.Set(s.ctx, uid, res, 10*time.Second).Err()
    if err != nil {
        log.Error(err)
        snag.Panic("信息获取失败")
    }
    return res
}

// updateCabinetExchangeProcess 更新换电步骤信息
func (s *riderCabinetService) updateCabinetExchangeProcess() {
    cache.Set(s.ctx, s.info.Serial, &model.CabinetExchangeProcess{
        Info:       s.operating,
        Step:       s.step,
        BatteryNum: s.batteryNum,
        Rider: &model.RiderBasic{
            ID:    s.rider.ID,
            Phone: s.rider.Phone,
            Name:  s.rider.Edges.Person.Name,
        },
    }, s.maxTime)
}

// Process 处理换电
func (s *riderCabinetService) Process(req *model.RiderCabinetOperateReq) {
    // 检查用户是否可以办理业务
    NewRiderPermissionWithRider(s.rider).BusinessX()

    sm, _ := NewSetting().GetSetting(model.SettingMaintain).(bool)
    if sm {
        snag.Panic("系统维护中, 请稍后重试")
    }

    // 校验换电信息
    iv := cache.Int(model.SettingExchangeInterval)
    if exist, _ := ent.Database.Exchange.QueryNotDeleted().Where(
        exchange.RiderID(s.rider.ID),
        exchange.Success(true),
        exchange.CreatedAtGTE(time.Now().Add(-time.Duration(cache.Int(model.SettingExchangeInterval))*time.Minute)),
    ).Exist(s.ctx); exist {
        snag.Panic(fmt.Sprintf("换电过于频繁, %d分钟可再次换电", iv))
    }

    info := new(model.RiderCabinetInfo)
    uid := *req.UUID
    err := cache.Get(s.ctx, uid).Scan(info)
    if err != nil || info == nil || info.EmptyBin == nil {
        snag.Panic("未找到信息, 请重新扫码")
    }

    // 是否有生效中套餐
    subd, sub := NewSubscribe().RecentDetail(s.rider.ID)
    if subd == nil || subd.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效中的骑行卡")
    }

    s.subscribe = sub

    cab := NewCabinet().QueryOne(info.ID)
    index := -1
    var be model.BatteryElectricity
    if info.FullBin == nil {
        if req.Alternative != nil && *req.Alternative {
            index = info.Alternative.Index
            be = info.Alternative.Electricity
            s.alternative = true
        } else {
            snag.Panic("用户取消")
        }
    } else {
        index = info.FullBin.Index
        be = info.FullBin.Electricity
    }

    // 缓存步骤
    if model.CabinetBusying(cab.Serial) {
        snag.Panic("电柜忙")
    }

    s.logger = logging.NewExchangeLog(s.rider.ID, uid, info.Serial, s.rider.Phone, be.IsBatteryFull())
    s.step = model.RiderCabinetOperateStepOpenEmpty
    s.cabinet = cab
    s.batteryNum = cab.BatteryNum
    s.info = info
    s.operating = &model.RiderCabinetOperating{
        UUID:        uid,
        ID:          info.ID,
        Name:        info.Name,
        EmptyIndex:  info.EmptyBin.Index,
        FullIndex:   index,
        Serial:      info.Serial,
        Electricity: be,
        Model:       info.Model,
    }
    s.model = sub.Model

    // 更换UUID缓存内容为步骤状态
    res := &model.RiderCabinetOperateRes{
        Step:   model.RiderCabinetOperateStepOpenEmpty,
        Status: model.RiderCabinetOperateStatusProcessing,
    }
    cache.Del(s.ctx, uid)
    cache.Set(s.ctx, uid, res, s.maxTime)

    s.updateCabinetExchangeProcess()

    s.start = time.Now()

    // 记录换电人
    // TODO 超时处理
    s.ex, _ = ent.Database.Exchange.
        Create().
        SetRiderID(s.rider.ID).
        SetCityID(s.info.CityID).
        SetDetail(&model.ExchangeCabinet{
            Alternative: s.alternative,
            Info:        s.operating,
            Result:      res,
        }).
        SetUUID(s.operating.UUID).
        SetCabinetID(s.operating.ID).
        SetSuccess(false).
        SetModel(s.subscribe.Model).
        SetNillableEnterpriseID(s.subscribe.EnterpriseID).
        SetNillableStationID(s.subscribe.StationID).
        SetSubscribeID(s.subscribe.ID).
        SetAlternative(s.alternative).
        SetStartAt(s.start).
        Save(s.ctx)

    if s.ex == nil {
        snag.Panic("换电失败")
    }

    // 处理换电流程
    go s.ProcessByStep()
}

// ProcessStepStart 步骤开始
func (s *riderCabinetService) ProcessStepStart() {
    s.updateCabinetExchangeProcess()
    cache.Set(s.ctx, s.operating.UUID, &model.RiderCabinetOperateRes{
        Step:   s.step,
        Status: model.RiderCabinetOperateStatusProcessing,
        Stop:   false,
    }, s.maxTime)
}

// ProcessStepEnd 结束换电流程
func (s *riderCabinetService) ProcessStepEnd() {
    var panicErr string
    if r := recover(); r != nil {
        log.Errorf("换电异常结束 -> %s {%s}: %v", s.operating.UUID, s.step, r)
        panicErr = fmt.Sprintf("%v", r)
    }

    res := new(model.RiderCabinetOperateRes)
    _ = cache.Get(s.ctx, s.operating.UUID).Scan(res)

    // 释放占用
    cache.Del(s.ctx, s.operating.Serial)
    cache.Del(s.ctx, s.operating.UUID)

    if panicErr != "" {
        res.Message = "换电故障: " + panicErr
        res.Stop = true
        res.Status = model.RiderCabinetOperateStatusFail

        index, _ := s.currentDoorInfo()
        log.Errorf("换电异常结束 -> %s {%s}: %v", s.operating.UUID, res.Step, s.logProcessRes(index, res))
    }

    now := time.Now()

    // 保存数据库
    _, _ = s.ex.Update().
        SetRiderID(s.rider.ID).
        SetCityID(s.info.CityID).
        SetDetail(&model.ExchangeCabinet{
            Alternative: s.alternative,
            Info:        s.operating,
            Result:      res,
        }).
        SetUUID(s.operating.UUID).
        SetCabinetID(s.operating.ID).
        SetSuccess(res.Status == model.RiderCabinetOperateStatusSuccess && res.Step == model.RiderCabinetOperateStepPutOut).
        SetModel(s.subscribe.Model).
        SetNillableEnterpriseID(s.subscribe.EnterpriseID).
        SetNillableStationID(s.subscribe.StationID).
        SetSubscribeID(s.subscribe.ID).
        SetAlternative(s.alternative).
        SetStartAt(s.start).
        SetFinishAt(now).
        SetDuration(int(now.Sub(s.start).Seconds())).
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
    s.step += 1
    if !s.ProcessLog(s.ProcessDoorStatus()) {
        return
    }

    // 第三步: 开启满电仓
    s.step += 1
    if !s.ProcessLog(s.ProcessOpenBin()) {
        return
    }
    // 手动延时处理
    time.Sleep(5 * time.Second)

    // 第四步: 长轮询判断仓门是否关闭
    s.step += 1
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

        // 判断是否 有电池 并且 (电压大于0 或 电量大于0)
        if bin.Battery && (bin.Voltage > 0 || bin.Electricity > 0) {
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

    s.ProcessStepStart()
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
            res.Status = model.RiderCabinetOperateStatusSuccess
            return
        case model.CabinetBinDoorStatusOpen:
            break
        default:
            res.Message = model.CabinetBinDoorError[ds]
            res.Stop = true
            res.Status = model.RiderCabinetOperateStatusFail
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
func (s *riderCabinetService) ProcessOpenBin() (res *model.RiderCabinetOperateRes) {
    res = &model.RiderCabinetOperateRes{}

    index := s.operating.EmptyIndex
    id := s.operating.ID
    step := s.step
    if step == model.RiderCabinetOperateStepOpenFull {
        index = s.operating.FullIndex
    }
    s.ProcessStepStart()

    var status bool
    var err error

    r := s.rider
    operator := model.CabinetDoorOperator{
        ID:    r.ID,
        Role:  model.CabinetDoorOperatorRoleRider,
        Name:  r.Edges.Person.Name,
        Phone: r.Phone,
    }

    // 开始处理
    reason := model.RiderCabinetOperateReasonEmpty
    if step == model.RiderCabinetOperateStepOpenFull {
        reason = model.RiderCabinetOperateReasonFull
    }
    operation := model.CabinetDoorOperateOpen
    status, err = NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
        ID:        &id,
        Index:     &index,
        Remark:    fmt.Sprintf("骑手换电 - %s", reason),
        Operation: &operation,
    }, operator, true)
    if err != nil {
        snag.Panic(err)
    }

    log.Infof(`[换电步骤 - 仓门检测]: {step:%s} %d, 仓门Index: %d, 操作反馈: %t, 用户电话: %s`,
        s.step,
        id,
        index,
        status,
        s.rider.Phone,
    )

    provider.AutoBinFault(operator, s.cabinet, index, status, func() {
        operation := model.CabinetDoorOperateLock
        _, _ = NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
            ID:        &id,
            Index:     &index,
            Remark:    fmt.Sprintf("换电仓门处理失败自动锁仓 - %s", s.rider.Phone),
            Operation: &operation,
        }, operator, true)
    })

    if status {
        res.Status = model.RiderCabinetOperateStatusSuccess
    } else {
        log.Errorf("[ProcessOpenBin] 处理失败: %t -> %#v", status, err)
        // 发生故障处理
        res.Status = model.RiderCabinetOperateStatusFail
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
        if res.Status == model.RiderCabinetOperateStatusSuccess || res.Stop || time.Now().Sub(start).Seconds() > 30 {
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

func (s *riderCabinetService) logProcessRes(index int, res *model.RiderCabinetOperateRes) string {
    return fmt.Sprintf(`[换电步骤 - 结果]: {step:%s} %s - %d, 用户电话: %s, 状态: %s, 消息: %s, 终止: %t`,
        s.step,
        s.cabinet.Serial,
        index,
        s.rider.Phone,
        res.Status,
        res.Message,
        res.Stop,
    )
}

// ProcessLog 处理步骤日志
func (s *riderCabinetService) ProcessLog(res *model.RiderCabinetOperateRes) bool {
    res.Step = s.step
    res.Stop = s.step == model.RiderCabinetOperateStepPutOut || res.Status == model.RiderCabinetOperateStatusFail

    index, be := s.currentDoorInfo()

    log.Info(s.logProcessRes(index, res))

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

    return res.Status == model.RiderCabinetOperateStatusSuccess
}
