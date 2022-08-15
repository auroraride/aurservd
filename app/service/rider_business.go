// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/business"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "time"
)

// TODO 服务器崩溃后自动启动继续换电进程
// TODO 电柜缓存优化

type riderBusinessService struct {
    ctx     context.Context
    rider   *ent.Rider
    maxTime time.Duration // 单步骤最大处理时长
    logger  *logging.ExchangeLog

    cabinet   *ent.Cabinet
    subscribe *ent.Subscribe

    task  *ec.Task
    max   *ec.BinInfo
    empty *ec.BinInfo

    bt business.Type
}

func NewRiderBusiness(rider *ent.Rider) *riderBusinessService {
    s := &riderBusinessService{
        ctx:     context.Background(),
        maxTime: 180 * time.Second,
    }
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

// preprocess 预处理业务
func (s *riderBusinessService) preprocess(serial string, bt business.Type) {
    NewSetting().SystemMaintainX()

    cs := NewCabinet()

    cab := cs.QueryOneSerialX(serial)
    if !cab.Transferred {
        snag.Panic("电柜资产异常")
    }

    // 是否有生效中套餐
    sub := NewSubscribe().Recent(s.rider.ID)
    if sub == nil {
        snag.Panic("无生效中的骑行卡")
    }

    s.subscribe = sub

    // 检查可用电池型号
    if !cs.ModelInclude(cab, sub.Model) {
        snag.Panic("电池型号不兼容")
    }

    // 查询电柜是否可使用
    NewCabinet().BusinessableX(cab)

    ec.BusyX(serial)

    err := NewCabinet().UpdateStatus(cab)
    if err != nil {
        snag.Panic(err)
    }

    max, empty := cab.Bin.MaxEmpty()

    var bn, en int

    for _, bin := range cab.Bin {
        // 仓位锁仓跳过计算
        if !bin.DoorHealth {
            continue
        }
        if bin.Battery {
            // 有电池
            bn += 1
        } else {
            // 无电池
            en += 1
        }
    }

    switch bt {
    case business.TypePause, business.TypeUnsubscribe:
        if en < 2 {
            snag.Panic("仓位不足, 无法处理当前业务")
        }
        if empty == nil {
            snag.Panic("电柜异常")
        }
        s.empty = &ec.BinInfo{Index: empty.Index}
        break
    case business.TypeActive, business.TypeContinue:
        if bn < 2 {
            snag.Panic("电池不足, 无法处理当前业务")
        }
        if max == nil {
            snag.Panic("电柜异常")
        }
        s.max = &ec.BinInfo{
            Index:       max.Index,
            Electricity: max.Electricity,
            Voltage:     max.Voltage,
        }
        break
    }

    s.bt = bt
    s.cabinet = cab

    jobs := map[business.Type]ec.Job{
        business.TypeActive:      ec.JobRiderActive,
        business.TypeUnsubscribe: ec.JobRiderUnSubscribe,
        business.TypePause:       ec.JobPause,
        business.TypeContinue:    ec.JobContinue,
    }

    task := &ec.Task{
        Job:       jobs[bt],
        Serial:    cab.Serial,
        CabinetID: cab.ID,
        Cabinet:   cab.GetTaskInfo(),
        Rider: &ec.Rider{
            ID:    s.rider.ID,
            Name:  s.rider.Edges.Person.Name,
            Phone: s.rider.Phone,
        },
    }
    s.task = task.CreateX()
}

func (s *riderBusinessService) open(bin *ec.BinInfo, remark string) (status bool, err error) {
    s.task.Start()

    operation := model.CabinetDoorOperateOpen

    status, err = NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
        ID:        tools.NewPointer().UInt64(s.cabinet.ID),
        Index:     tools.NewPointerInterface(bin.Index),
        Remark:    fmt.Sprintf("%s - %s", s.task.Job.Label(), remark),
        Operation: &operation,
    }, model.CabinetDoorOperator{
        ID:    s.rider.ID,
        Role:  model.CabinetDoorOperatorRoleRider,
        Name:  s.rider.Edges.Person.Name,
        Phone: s.rider.Phone,
    })

    return
}

func (s *riderBusinessService) putin() *ec.BinInfo {
    status, err := s.open(s.empty, model.RiderCabinetOperateReasonEmpty)
    ts := s.task.Status

    defer func() {
        s.task.Stop(ts)
    }()

    if !status {
        snag.Panic("仓门开启失败")
    }

    if err != nil {
        snag.Panic(err)
    }

    for {
        ds := s.battery()
        if ds == ec.DoorStatusClose {
            // 强制睡眠两秒: 原因是有可能柜门会晃动导致似关非关, 延时来获取正确状态
            time.Sleep(2 * time.Second)
            ds = s.battery()
        }
        switch ds {
        case ec.DoorStatusClose:
            ts = ec.TaskStatusSuccess
            return s.empty
        case ec.DoorStatusOpen:
            ts = ec.TaskStatusProcessing
            break
        default:
            s.task.Message = ec.DoorError[ds]
            ts = ec.TaskStatusFail
            break
        }

        // 超时标记为任务失败
        if time.Now().Sub(*s.task.StartAt).Seconds() > s.maxTime.Seconds() && s.task.Message == "" {
            s.task.Message = "超时"
            ts = ec.TaskStatusFail
        }

        if ts != ec.TaskStatusProcessing {
            if s.task.Message == "" {
                s.task.Message = "操作失败"
            }
            snag.Panic(s.task.Message)
        }

        time.Sleep(1 * time.Second)
    }
}

// battery 电池检测
func (s *riderBusinessService) battery() (ds ec.DoorStatus) {
    ds = NewCabinet().DoorOpenStatus(s.cabinet, s.empty.Index)
    cbin := s.cabinet.Bin[s.empty.Index]
    pe := cbin.Electricity
    pv := cbin.Voltage

    // 当仓门未关闭时跳过
    if ds != ec.DoorStatusClose {
        return
    }

    // 获取骑手放入电池信息, 验证是否放入旧电池
    if s.empty.Electricity == 0 {
        s.empty.Electricity = pe
    }

    if s.empty.Voltage < 40 {
        s.empty.Voltage = pv
    }

    // 判断是否 有电池 并且 (电压大于40 或 电量大于0)
    if cbin.Battery && (pv > 40 || pe > 0) {
        return ec.DoorStatusClose
    }

    // 仓门关闭但是检测不到电池的情况下, 继续检测30s
    if time.Now().Sub(*s.task.StartAt).Seconds() > 30 {
        return ec.DoorStatusBatteryEmpty
    } else {
        time.Sleep(1 * time.Second)
        return s.battery()
    }
}

func (s *riderBusinessService) putout() *ec.BinInfo {
    status, err := s.open(s.max, model.RiderCabinetOperateReasonFull)

    defer func() {
        ts := ec.TaskStatusSuccess
        if !status {
            ts = ec.TaskStatusFail
        }
        s.task.Stop(ts)
    }()

    if err != nil {
        snag.Panic(err)
    }

    return s.max
}

func (s *riderBusinessService) Active(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
    s.preprocess(req.Serial, business.TypeActive)
    if s.subscribe.Status != model.SubscribeStatusInactive {
        snag.Panic("骑士卡状态错误")
    }

    srv := NewBusinessRider(s.rider)
    srv.SetCabinet(s.cabinet).
        SetTask(func() *ec.BinInfo {
            return s.putout()
        }).
        Active(srv.Inactive(req.ID))
    return model.BusinessCabinetStatus{
        UUID:  s.task.ID.Hex(),
        Index: s.max.Index,
    }
}

func (s *riderBusinessService) Continue(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
    s.preprocess(req.Serial, business.TypeContinue)
    if s.subscribe.Status != model.SubscribeStatusPaused {
        snag.Panic("骑士卡状态错误")
    }

    NewBusinessRider(s.rider).SetTask(func() *ec.BinInfo {
        return s.putout()
    }).SetCabinet(s.cabinet).Continue(req.ID)
    return model.BusinessCabinetStatus{
        UUID:  s.task.ID.Hex(),
        Index: s.max.Index,
    }
}

func (s *riderBusinessService) Unsubscribe(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
    s.preprocess(req.Serial, business.TypeUnsubscribe)
    if s.subscribe.Status != model.SubscribeStatusUsing {
        snag.Panic("骑士卡未在计费中")
    }

    go func() {
        err := snag.WithPanic(func() {
            NewBusinessRider(s.rider).SetTask(func() *ec.BinInfo {
                return s.putin()
            }).SetCabinet(s.cabinet).UnSubscribe(req.ID)
        })

        if err != nil {
            log.Error(err)
        }
    }()

    return model.BusinessCabinetStatus{
        UUID:  s.task.ID.Hex(),
        Index: s.empty.Index,
    }
}

func (s *riderBusinessService) Pause(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
    s.preprocess(req.Serial, business.TypePause)
    if s.subscribe.Status != model.SubscribeStatusUsing {
        snag.Panic("骑士卡未在计费中")
    }

    go func() {
        err := snag.WithPanic(func() {
            NewBusinessRider(s.rider).
                SetTask(func() *ec.BinInfo {
                    return s.putin()
                }).
                SetCabinet(s.cabinet).
                Pause(req.ID)
        })
        if err != nil {
            log.Error(err)
        }
    }()

    return model.BusinessCabinetStatus{
        UUID:  s.task.ID.Hex(),
        Index: s.empty.Index,
    }
}

// Status 业务操作状态
func (s *riderBusinessService) Status(req *model.BusinessCabinetStatusReq) (res model.BusinessCabinetStatusRes) {
    start := time.Now()
    for {
        task := ec.QueryID(req.UUID)
        if task == nil {
            snag.Panic("未找到业务操作")
        }
        if task.Status == ec.TaskStatusFail || task.Status == ec.TaskStatusSuccess {
            res.Success = task.Status == ec.TaskStatusSuccess
            res.Stop = true
            res.Message = task.Message
        }
        if res.Stop || time.Now().Sub(start).Seconds() > 30 {
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
    sub, _ := ent.Database.Subscribe.QueryNotDeleted().WithPauses(func(query *ent.SubscribePauseQuery) {
        query.Where(subscribepause.EndAtIsNil())
    }).First(s.ctx)
    if sub == nil {
        snag.Panic("未找到骑士卡信息")
    }
    ps := sub.Edges.Pauses
    if len(ps) == 0 {
        snag.Panic("未找到寄存信息")
    }

    p := ps[0]
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

    return
}
