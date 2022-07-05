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
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/lithammer/shortuuid/v4"
    log "github.com/sirupsen/logrus"
    "time"
)

// TODO 服务器崩溃后自动启动继续换电进程

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
    debug     bool
}

func NewRiderCabinet() *riderCabinetService {
    return &riderCabinetService{
        ctx:     context.Background(),
        maxTime: 180 * time.Second,
        debug:   ar.Config.Cabinet.Debug,
    }
}

func NewRiderCabinetWithRider(rider *ent.Rider) *riderCabinetService {
    s := NewRiderCabinet()
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

func NewRiderCabinetWithModifier(m *model.Modifier) *riderCabinetService {
    s := NewRiderCabinet()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// GetProcess 获取待换电信息
func (s *riderCabinetService) GetProcess(req *model.RiderCabinetOperateInfoReq) *model.RiderCabinetInfo {
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

    tools.NewLog().Infof("[换电信息] %s, %s", uid, res)

    err := cache.Set(s.ctx, uid, res, s.maxTime).Err()
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

    info := new(model.RiderCabinetInfo)
    uid := *req.UUID
    err := cache.Get(s.ctx, uid).Scan(info)
    if err != nil || info == nil || info.EmptyBin == nil {
        snag.Panic("未找到信息, 请重新操作")
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
    if model.CabinetBusying(info.Serial) {
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
        EmptyIndex:  info.EmptyBin.Index,
        FullIndex:   index,
        Serial:      info.Serial,
        Electricity: be,
        Model:       info.Model,
    }
    s.model = sub.Model

    // 更换UUID缓存内容为步骤状态
    cache.Del(s.ctx, uid)
    cache.Set(s.ctx, uid, &model.RiderCabinetOperateRes{
        Step:   model.RiderCabinetOperateStepOpenEmpty,
        Status: model.RiderCabinetOperateStatusProcessing,
    }, s.maxTime)

    s.updateCabinetExchangeProcess()

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
    // 释放占用
    cache.Del(s.ctx, s.operating.Serial)

    res := new(model.RiderCabinetOperateRes)
    _ = cache.Get(s.ctx, s.operating.UUID).Scan(res)

    // 保存数据库
    _, _ = ent.Database.Exchange.Create().
        SetRiderID(s.rider.ID).
        SetCityID(s.info.CityID).
        SetDetail(&model.ExchangeCabinet{
            Alternative: s.alternative,
            Info:        s.operating,
            Result:      res,
        }).
        SetUUID(s.operating.UUID).
        SetCabinetID(s.operating.ID).
        SetSuccess(res.Status == model.RiderCabinetOperateStatusSuccess && res.Step == model.RiderCabinetOperateStepClose).
        SetModel(s.subscribe.Model).
        SetNillableEnterpriseID(s.subscribe.EnterpriseID).
        SetNillableStationID(s.subscribe.StationID).
        SetSubscribeID(s.subscribe.ID).
        SetAlternative(s.alternative).
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
    if !s.ProcessLog(s.ProcessDoor()) {
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
    if !s.ProcessLog(s.ProcessDoor()) {
        return
    }
}

// ProcessDoorBatteryStatus 格式化仓门状态, 电池放入取出检测
func (s *riderCabinetService) ProcessDoorBatteryStatus() (ds model.CabinetBinDoorStatus) {
    // DEBUG START
    if s.debug {
        start := time.Now()
        ds = model.CabinetBinDoorStatusUnknown
        for {
            if time.Now().Sub(start).Seconds() > 2 {
                return model.CabinetBinDoorStatusClose
            }
            time.Sleep(1 * time.Second)
        }
    }
    // DEBUG END

    // 获取仓位index
    index := s.operating.EmptyIndex
    step := s.step
    if step == model.RiderCabinetOperateStepClose {
        index = s.operating.FullIndex
    }

    // 获取仓门状态
    ds = NewCabinet().DoorOpenStatus(s.cabinet, index, true)
    bin := s.cabinet.Bin[index]

    log.Infof(`[换电步骤 - 仓门检测]: {step:%s} %s - %d, 仓门Index: %d, 仓门状态: %s, 是否有电池: %t, 当前电压: %.2fV`,
        s.step,
        s.cabinet.Serial,
        index,
        bin.Index,
        ds,
        bin.Battery,
        bin.Voltage,
    )

    // 当仓门未关闭时跳过
    if ds != model.CabinetBinDoorStatusClose {
        return ds
    }

    // 验证是否放入旧电池
    if step == model.RiderCabinetOperateStepPutInto {
        // 凯信电柜需要绑定电池
        // bind := true
        // if s.info.Brand == model.CabinetBrandKaixin {
        //     bind = provider.NewKaixin().BatteryBind(s.rider.Edges.Person.Name+"-"+shortuuid.New(), s.info.Serial, s.model, s.operating.EmptyIndex)
        // }
        // 判断是否有电池并且电压大于0
        if bin.Battery && bin.Voltage > 0 {
            s.putInElectricity = bin.Electricity
            return model.CabinetBinDoorStatusClose
        }
        return model.CabinetBinDoorStatusBatteryEmpty
    }

    // 验证满电电池是否取走
    if step == model.RiderCabinetOperateStepClose {
        if bin.Battery {
            return model.CabinetBinDoorStatusBatteryFull
        }
        return model.CabinetBinDoorStatusClose
    }

    return ds
}

// ProcessDoor 操作换电中检查柜门并处理状态
func (s *riderCabinetService) ProcessDoor() (res *model.RiderCabinetOperateRes) {
    res = &model.RiderCabinetOperateRes{}

    s.ProcessStepStart()
    start := time.Now()

    for {
        // 检测仓门/电池
        ds := s.ProcessDoorBatteryStatus()
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

    if s.debug {
        // DEBUG START
        log.Println(id, index)
        debug := time.Now()
        for {
            if time.Now().Sub(debug).Seconds() > 2 {
                status = true
                err = nil
                break
            }
            time.Sleep(1 * time.Second)
        }
        // DEBUG END
    } else {
        // 开始处理
        reason := model.RiderCabinetOperateReasonEmpty
        if step == model.RiderCabinetOperateStepOpenFull {
            reason = model.RiderCabinetOperateReasonFull
        }
        r := s.rider
        operation := model.CabinetDoorOperateOpen
        status, err = NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
            ID:        &id,
            Index:     &index,
            Remark:    fmt.Sprintf("骑手换电 - %s", reason),
            Operation: &operation,
        }, model.CabinetDoorOperator{
            ID:    r.ID,
            Role:  model.CabinetDoorOperatorRoleRider,
            Name:  r.Edges.Person.Name,
            Phone: r.Phone,
        }, true)
        if err != nil {
            snag.Panic(err)
        }
    }

    log.Infof(`[换电步骤 - 仓门检测]: {step:%s} %d, 仓门Index: %d, 操作反馈: %t`,
        s.step,
        id,
        index,
        status,
    )

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

// ProcessLog 处理步骤日志
func (s *riderCabinetService) ProcessLog(res *model.RiderCabinetOperateRes) bool {
    res.Step = s.step
    res.Stop = s.step == model.RiderCabinetOperateStepClose || res.Status == model.RiderCabinetOperateStatusFail

    // 前两步是空仓, 后两步是满电仓位
    index := s.operating.EmptyIndex
    be := s.putInElectricity
    if s.step >= model.RiderCabinetOperateStepOpenFull {
        index = s.operating.FullIndex
        be = s.operating.Electricity
    }

    tools.NewLog().Infof("[换电步骤] {step:%s} %s", s.step, res)

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
