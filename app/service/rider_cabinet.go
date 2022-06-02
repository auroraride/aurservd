// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
    "time"
)

type riderCabinetService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    maxTime  time.Duration // 单步骤最大处理时长
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
    // 是否有生效中套餐
    o := NewSubscribe().Recent(s.rider.ID)
    if o == nil || o.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效中的骑行卡")
    }
    // 查询电柜
    cs := NewCabinet()
    cab := cs.QueryWithSerial(req.Serial)
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

    index := -1
    if info.FullBin == nil {
        if *req.Alternative {
            index = info.Alternative.Index
        } else {
            snag.Panic("用户取消")
        }
    } else {
        index = info.FullBin.Index
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

    // 处理换电流程
    go s.ProcessByStep(&model.RiderCabinetOperating{
        UUID:       uid,
        ID:         info.ID,
        EmptyIndex: info.EmptyBin.Index,
        FullIndex:  index,
        Serial:     info.Serial,
    })
}

// ProcessStepStart 步骤开始
func (s *riderCabinetService) ProcessStepStart(uid, serial string, step model.RiderCabinetOperateStep) {
    cache.Set(s.ctx, serial, step, s.maxTime)
    cache.Set(s.ctx, uid, &model.RiderCabinetOperateRes{
        Step:   step,
        Status: model.RiderCabinetOperateStatusProcessing,
        Stop:   false,
    }, s.maxTime)
}

// ProcessByStep 按步骤换电操作
func (s *riderCabinetService) ProcessByStep(req *model.RiderCabinetOperating) {
    cab := NewCabinet().QueryOne(req.ID)
    // 第一步: 开启空电仓
    if !s.ProcessOpenBin(req.UUID, req.Serial, req.ID, req.EmptyIndex, model.RiderCabinetOperateStepOpenEmpty) {
        return
    }
    // 第二步: 长轮询判断仓门是否关闭
    if !s.ProcessDoorStatus(model.RiderCabinetOperateStepPutInto, req.EmptyIndex, req.UUID, cab) {
        return
    }
    // 第三步: 开启满电仓
    if !s.ProcessOpenBin(req.UUID, req.Serial, req.ID, req.FullIndex, model.RiderCabinetOperateStepOpenFull) {
        return
    }
    // 第四步: 长轮询判断仓门是否关闭
    if !s.ProcessDoorStatus(model.RiderCabinetOperateStepClose, req.FullIndex, req.UUID, cab) {
        return
    }
}

// ProcessDoorStatus 操作换电中检查柜门并处理状态
// TODO 检测电池拿出和放入
func (s *riderCabinetService) ProcessDoorStatus(step model.RiderCabinetOperateStep, index int, uid string, cab *ent.Cabinet) bool {
    s.ProcessStepStart(uid, cab.Serial, step)
    start := time.Now()

    // TODO DEBUG START
    ds := model.CabinetBinDoorStatusUnknown
    for {
        if time.Now().Sub(start).Seconds() > 5 {
            ds = model.CabinetBinDoorStatusClose
            break
        }
        time.Sleep(1 * time.Second)
    }
    // TODO DEBUG END

    for {
        // ds := NewCabinet().DoorOpenStatus(cab, index)
        switch ds {
        case model.CabinetBinDoorStatusClose:
            stop := false
            // 当步骤为最后一步的时候终止步骤
            if step == model.RiderCabinetOperateStepClose {
                stop = true
                cache.Del(s.ctx, cab.Serial)
            }
            cache.Set(s.ctx, uid, &model.RiderCabinetOperateRes{
                Step:   step,
                Status: model.RiderCabinetOperateStatusSuccess,
                Stop:   stop,
            }, s.maxTime)
            return true
        case model.CabinetBinDoorStatusOpen:
            break
        case model.CabinetBinDoorStatusFail, model.CabinetBinDoorStatusUnknown:
            s.ProcessDoorFail(uid, cab.Serial, step, ds)
            return false
        }

        // 超时标记为任务失败
        if time.Now().Sub(start).Seconds() > s.maxTime.Seconds() {
            return false
        }
        time.Sleep(1 * time.Second)
    }
}

// ProcessDoorFail 操作换电柜门关闭检查失败
func (s *riderCabinetService) ProcessDoorFail(uid, serial string, step model.RiderCabinetOperateStep, ds model.CabinetBinDoorStatus) {
    // 失败
    msg := "柜门状态未知"
    if ds == model.CabinetBinDoorStatusFail {
        msg = "柜门故障"
    }
    cache.Set(s.ctx, uid, &model.RiderCabinetOperateRes{
        Step:    step,
        Status:  model.RiderCabinetOperateStatusFail,
        Message: msg,
        Stop:    true,
    }, s.maxTime)
    cache.Del(s.ctx, serial)
}

// ProcessOpenBin 开仓门
func (s *riderCabinetService) ProcessOpenBin(uid, serial string, id uint64, index int, step model.RiderCabinetOperateStep) bool {
    s.ProcessStepStart(uid, serial, step)

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
        if time.Now().Sub(debug).Seconds() > 5 {
            status = true
            err = nil
            break
        }
        time.Sleep(1 * time.Second)
    }
    // TODO DEBUG END
    res := &model.RiderCabinetOperateRes{
        Step: step,
    }
    if status {
        res.Status = model.RiderCabinetOperateStatusSuccess
    } else {
        // 发生故障处理
        res.Status = model.RiderCabinetOperateStatusFail
        res.Message = err.Error()
        res.Stop = true
        cache.Del(s.ctx, serial)
    }

    cache.Set(s.ctx, uid, res, s.maxTime)
    return status
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
