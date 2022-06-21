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

    Step             model.RiderCabinetOperateStep
    Cabinet          *ent.Cabinet
    PutInElectricity model.BatteryElectricity
    Alternative      bool
    Info             *model.RiderCabinetInfo

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

    // 是否企业用户
    if s.rider.Edges.Enterprise != nil {
        // TODO 企业状态是否可以换电
    }

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
        snag.Panic("电柜目前不可用")
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
        Model:       info.Model,
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
        SetModel(s.subscribe.Model).
        SetNillableEnterpriseID(s.subscribe.EnterpriseID).
        SetNillableStationID(s.subscribe.StationID).
        SetSubscribeID(s.subscribe.ID).
        Save(s.ctx)
}

// ProcessByStep 按步骤换电操作
func (s *riderCabinetService) ProcessByStep(req *model.RiderCabinetOperating) {
    defer s.ProcessStepEnd(req)

    // 第一步: 开启空电仓
    if !s.ProcessLog(req, s.ProcessOpenBin(req)) {
        return
    }
    // 手动延时处理
    time.Sleep(5 * time.Second)

    // 第二步: 长轮询判断仓门是否关闭
    s.Step += 1
    if !s.ProcessLog(req, s.ProcessDoor(req)) {
        return
    }

    // 第三步: 开启满电仓
    s.Step += 1
    if !s.ProcessLog(req, s.ProcessOpenBin(req)) {
        return
    }
    // 手动延时处理
    time.Sleep(5 * time.Second)

    // 第四步: 长轮询判断仓门是否关闭
    s.Step += 1
    if !s.ProcessLog(req, s.ProcessDoor(req)) {
        return
    }
}

// ProcessDoorBatteryStatus 格式化仓门状态, 电池放入取出检测
func (s *riderCabinetService) ProcessDoorBatteryStatus(req *model.RiderCabinetOperating) (ds model.CabinetBinDoorStatus) {
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
    index := req.EmptyIndex
    step := s.Step
    if step == model.RiderCabinetOperateStepClose {
        index = req.FullIndex
    }

    // 获取仓门状态
    ds = NewCabinet().DoorOpenStatus(s.Cabinet, index)
    bin := s.Cabinet.Bin[index]

    log.Infof(`[换电步骤 - 仓门检测]: {step:%s} %s - %d, 仓门Index: %d, 仓门状态: %s, 是否有电池: %t`,
        s.Step,
        s.Cabinet.Serial,
        index,
        bin.Index,
        ds,
        bin.Battery,
    )

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

    for {
        // 检测仓门/电池
        ds := s.ProcessDoorBatteryStatus(req)
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
        })
        if err != nil {
            snag.Panic(err)
        }
    }

    log.Infof(`[换电步骤 - 仓门检测]: {step:%s} %d, 仓门Index: %d, 操作反馈: %t`,
        s.Step,
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

    tools.NewLog().Infof("[换电步骤] {step:%s} %s", s.Step, res)

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
