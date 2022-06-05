// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/google/uuid"
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

    Step             model.RiderCabinetOperateStep
    Cabinet          *ent.Cabinet
    PutInElectricity model.BatteryElectricity
    Alternative      bool
    Info             *model.RiderCabinetInfo
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

func NewRiderCabinetWithModifier(m *model.Modifier) *riderCabinetService {
    s := NewRiderCabinet()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// GetProcess 获取待换电信息
func (s *riderCabinetService) GetProcess(req *model.RiderCabinetOperateInfoReq) *model.RiderCabinetInfo {
    // 是否企业用户
    if s.rider.Edges.Enterprise != nil {
        // TODO 企业状态是否可以换电
    }

    // 是否有生效中套餐
    o := NewSubscribe().Recent(s.rider.ID)
    if o == nil || o.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效中的骑行卡")
    }

    // 查询电柜
    cs := NewCabinet()
    cab := cs.QueryWithSerial(req.Serial)

    // 检查可用电池型号
    if !cs.VoltageInclude(cab, o.Voltage) {
        snag.Panic("电池型号不兼容")
    }

    // 查询套餐
    if !cs.Health(cab) {
        snag.Panic("电柜目前不可用")
    }
    info := NewCabinet().Usable(cab)
    if info.EmptyBin == nil {
        snag.Panic("电柜目前不可用")
    }

    uid := uuid.New().String()
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
        Voltage:                    o.Voltage,
        CityID:                     o.City.ID,
    }
    err := cache.Set(s.ctx, uid, res, s.maxTime).Err()
    if err != nil {
        log.Error(err)
        snag.Panic("信息获取失败")
    }
    return res
}

// Process 处理换电
func (s *riderCabinetService) Process(req *model.RiderCabinetOperateReq) {
    info := new(model.RiderCabinetInfo)
    uid := *req.UUID
    err := cache.Get(s.ctx, uid).Scan(info)
    if err != nil || info == nil || info.EmptyBin == nil {
        snag.Panic("未找到信息, 请重新操作")
    }

    cab := NewCabinet().QueryOne(info.ID)
    index := -1
    var be model.BatteryElectricity
    if info.FullBin == nil {
        if req.Alternative != nil && *req.Alternative {
            index = info.Alternative.Index
            be = info.Alternative.Electricity
            s.Alternative = true
        } else {
            snag.Panic("用户取消")
        }
    } else {
        index = info.FullBin.Index
        be = info.FullBin.Electricity
    }

    // 缓存步骤
    n, _ := cache.Get(s.ctx, info.Serial).Int()
    if n > 0 {
        snag.Panic("电柜忙")
    }

    // 更换UUID缓存内容为步骤状态
    cache.Del(s.ctx, uid)
    cache.Set(s.ctx, uid, &model.RiderCabinetOperateRes{
        Step:   model.RiderCabinetOperateStepOpenEmpty,
        Status: model.RiderCabinetOperateStatusProcessing,
    }, s.maxTime)
    cache.Set(s.ctx, info.Serial, model.RiderCabinetOperateStepOpenEmpty, s.maxTime)

    s.logger = logging.NewExchangeLog(s.rider.ID, uid, info.Serial, s.rider.Phone, be.IsBatteryFull())
    s.Step = model.RiderCabinetOperateStepOpenEmpty
    s.Cabinet = cab
    s.Info = info

    // 处理换电流程
    go s.ProcessByStep(&model.RiderCabinetOperating{
        UUID:        uid,
        ID:          info.ID,
        EmptyIndex:  info.EmptyBin.Index,
        FullIndex:   index,
        Serial:      info.Serial,
        Electricity: be,
    })
}

// ProcessStepStart 步骤开始
func (s *riderCabinetService) ProcessStepStart(req *model.RiderCabinetOperating) {
    cache.Set(s.ctx, req.Serial, s.Step, s.maxTime)
    cache.Set(s.ctx, req.UUID, &model.RiderCabinetOperateRes{
        Step:   s.Step,
        Status: model.RiderCabinetOperateStatusProcessing,
        Stop:   false,
    }, s.maxTime)
}

// ProcessStepEnd 结束换电流程
func (s *riderCabinetService) ProcessStepEnd(req *model.RiderCabinetOperating) {
    // 释放占用
    cache.Del(s.ctx, req.Serial)

    res := new(model.RiderCabinetOperateRes)
    _ = cache.Get(s.ctx, req.UUID).Scan(res)

    // 保存数据库
    _, _ = ar.Ent.Exchange.Create().
        SetRiderID(s.rider.ID).
        SetCityID(s.Info.CityID).
        SetDetail(&model.ExchangeCabinet{
            Alternative: s.Alternative,
            Info:        req,
            Result:      res,
        }).
        SetUUID(req.UUID).
        SetCabinetID(req.ID).
        SetSuccess(res.Status == model.RiderCabinetOperateStatusSuccess && res.Step == model.RiderCabinetOperateStepClose).
        Save(s.ctx)
}

// ProcessByStep 按步骤换电操作
func (s *riderCabinetService) ProcessByStep(req *model.RiderCabinetOperating) {
    defer s.ProcessStepEnd(req)

    // 第一步: 开启空电仓
    if !s.ProcessLog(req, s.ProcessOpenBin(req)) {
        return
    }
    s.Step += 1

    // 第二步: 长轮询判断仓门是否关闭
    if !s.ProcessLog(req, s.ProcessDoor(req)) {
        return
    }
    s.Step += 1

    // 第三步: 开启满电仓
    if !s.ProcessLog(req, s.ProcessOpenBin(req)) {
        return
    }
    s.Step += 1

    // 第四步: 长轮询判断仓门是否关闭
    if !s.ProcessLog(req, s.ProcessDoor(req)) {
        return
    }
}

// ProcessDoorBatteryStatus 格式化仓门状态, 电池放入取出检测
func (s *riderCabinetService) ProcessDoorBatteryStatus(req *model.RiderCabinetOperating) (ds model.CabinetBinDoorStatus) {
    // 获取仓位index
    index := req.EmptyIndex
    step := s.Step
    if step == model.RiderCabinetOperateStepClose {
        index = req.FullIndex
    }

    bin := s.Cabinet.Bin[index]
    // 获取仓门状态
    ds = NewCabinet().DoorOpenStatus(s.Cabinet, index)

    // 当仓门未关闭时跳过
    if ds != model.CabinetBinDoorStatusClose {
        return ds
    }

    // 验证是否放入旧电池
    if step == model.RiderCabinetOperateStepPutInto {
        if bin.Battery {
            s.PutInElectricity = bin.Electricity
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
func (s *riderCabinetService) ProcessDoor(req *model.RiderCabinetOperating) (res *model.RiderCabinetOperateRes) {
    res = &model.RiderCabinetOperateRes{}

    s.ProcessStepStart(req)
    start := time.Now()

    // TODO DEBUG START
    ds := model.CabinetBinDoorStatusUnknown
    for {
        if time.Now().Sub(start).Seconds() > 2 {
            ds = model.CabinetBinDoorStatusClose
            break
        }
        time.Sleep(1 * time.Second)
    }
    // TODO DEBUG END

    for {
        // TODO 上线
        // ds := s.ProcessDoorBatteryStatus(req)
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
func (s *riderCabinetService) ProcessOpenBin(req *model.RiderCabinetOperating) (res *model.RiderCabinetOperateRes) {
    res = &model.RiderCabinetOperateRes{}

    index := req.EmptyIndex
    id := req.ID
    step := s.Step
    if step == model.RiderCabinetOperateStepOpenFull {
        index = req.FullIndex
    }
    s.ProcessStepStart(req)

    // TODO 上线
    // 开始处理
    // reason := model.RiderCabinetOperateReasonEmpty
    // if step == model.RiderCabinetOperateStepOpenFull {
    //     reason = model.RiderCabinetOperateReasonFull
    // }
    // r := s.rider
    // operation := model.CabinetDoorOperateOpen
    // status, err = NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
    //     ID:        &id,
    //     Index:     &index,
    //     Remark:    fmt.Sprintf("骑手换电 - %s", reason),
    //     Operation: &operation,
    // }, model.CabinetDoorOperator{
    //     ID:    r.ID,
    //     Role:  model.CabinetDoorOperatorRoleRider,
    //     Name:  r.Edges.Person.Name,
    //     Phone: r.Phone,
    // })
    // if err != nil {
    //     snag.Panic(err)
    // }

    // TODO DEBUG START
    log.Println(id, index)
    debug := time.Now()
    var status bool
    var err error
    for {
        if time.Now().Sub(debug).Seconds() > 2 {
            status = true
            err = nil
            break
        }
        time.Sleep(1 * time.Second)
    }
    // TODO DEBUG END

    if status {
        res.Status = model.RiderCabinetOperateStatusSuccess
    } else {
        // 发生故障处理
        res.Status = model.RiderCabinetOperateStatusFail
        res.Message = err.Error()
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
func (s *riderCabinetService) ProcessLog(req *model.RiderCabinetOperating, res *model.RiderCabinetOperateRes) bool {
    res.Step = s.Step
    res.Stop = s.Step == model.RiderCabinetOperateStepClose || res.Status == model.RiderCabinetOperateStatusFail

    // 前两步是空仓, 后两步是满电仓位
    index := req.EmptyIndex
    be := s.PutInElectricity
    if s.Step >= model.RiderCabinetOperateStepOpenFull {
        index = req.FullIndex
        be = req.Electricity
    }
    s.logger.Clone().
        SetBin(index).
        SetStatus(res.Status).
        SetMessage(res.Message).
        SetStep(s.Step).
        SetElectricity(be).
        Send()

    // 存储缓存
    err := cache.Set(s.ctx, req.UUID, res, s.maxTime).Err()
    if err != nil {
        log.Error(err)
        return false
    }

    return res.Status == model.RiderCabinetOperateStatusSuccess
}
