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
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

// TODO 服务器崩溃后自动启动继续换电进程
// TODO 电柜缓存优化

type riderCabinetService struct {
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

func NewRiderCabinet(rider *ent.Rider) *riderCabinetService {
    s := &riderCabinetService{
        ctx:     context.Background(),
        maxTime: 180 * time.Second,
    }
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

// preprocess 预处理业务
func (s *riderCabinetService) preprocess(serial string, bt business.Type) {
    ec.BusyX(serial)

    cab := NewCabinet().QueryOneSerialX(serial)
    if !cab.Transferred {
        snag.Panic("电柜资产异常")
    }

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
    case business.TypePause:
    case business.TypeUnsubscribe:
        if en < 2 {
            snag.Panic("仓位不足, 无法处理当前业务")
        }
        if empty == nil {
            snag.Panic("电柜异常")
        }
        s.empty = &ec.BinInfo{Index: empty.Index}
        break
    case business.TypeActive:
    case business.TypeContinue:
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

    s.task = &ec.Task{
        ID:  primitive.NewObjectID(),
        Job: ec.JobRiderActive,
        Cabinet: ec.Cabinet{
            Health:         cab.Health,
            Doors:          cab.Doors,
            BatteryNum:     cab.BatteryNum,
            BatteryFullNum: cab.BatteryFullNum,
        },
    }
}

func (s *riderCabinetService) open(bin *ec.BinInfo) (status bool, err error) {
    s.task.Start()

    operation := model.CabinetDoorOperateOpen

    status, err = NewCabinet().DoorOperate(&model.CabinetDoorOperateReq{
        ID:        tools.NewPointer().UInt64(s.cabinet.ID),
        Index:     tools.NewPointerInterface(bin.Index),
        Remark:    fmt.Sprintf("%s - %s", s.task.Job.Label(), model.RiderCabinetOperateReasonFull),
        Operation: &operation,
    }, model.CabinetDoorOperator{
        ID:    s.rider.ID,
        Role:  model.CabinetDoorOperatorRoleRider,
        Name:  s.rider.Edges.Person.Name,
        Phone: s.rider.Phone,
    }, true)

    return
}

func (s *riderCabinetService) putin() *ec.BinInfo {
    status, err := s.open(s.empty)
    ts := ec.TaskStatusFail

    defer s.task.Stop(ts)

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
            snag.Panic(s.task.Message)
        }

        time.Sleep(1 * time.Second)
    }
}

// battery 电池检测
func (s *riderCabinetService) battery() (ds ec.DoorStatus) {
    ds = NewCabinet().DoorOpenStatus(s.cabinet, s.empty.Index, true)
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

func (s *riderCabinetService) putout() *ec.BinInfo {
    status, err := s.open(s.max)

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

func (s *riderCabinetService) Active(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
    s.preprocess(req.Serial, business.TypeActive)
    srv := NewBusinessRiderWithRider(s.rider)
    srv.SetCabinet(s.cabinet).
        SetTask(s.putout).
        Active(srv.Inactive(req.ID))
    return model.BusinessCabinetStatus{
        UUID:  s.task.ID.Hex(),
        Index: s.empty.Index,
    }
}

func (s *riderCabinetService) Continue(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
    s.preprocess(req.Serial, business.TypeContinue)
    NewBusinessRiderWithRider(s.rider).SetTask(s.putout).Continue(req.ID)
    return model.BusinessCabinetStatus{
        UUID:  s.task.ID.Hex(),
        Index: s.empty.Index,
    }
}

func (s *riderCabinetService) Unsubscribe(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
    s.preprocess(req.Serial, business.TypeUnsubscribe)
    NewBusinessRiderWithRider(s.rider).SetTask(s.putin).UnSubscribe(req.ID)
    return model.BusinessCabinetStatus{
        UUID:  s.task.ID.Hex(),
        Index: s.empty.Index,
    }
}

func (s *riderCabinetService) Pause(req *model.BusinessCabinetReq) model.BusinessCabinetStatus {
    s.preprocess(req.Serial, business.TypePause)
    NewBusinessRiderWithRider(s.rider).SetTask(s.putin).Pause(req.ID)
    return model.BusinessCabinetStatus{
        UUID:  s.task.ID.Hex(),
        Index: s.empty.Index,
    }
}

// Status 业务操作状态
func (s *riderCabinetService) Status(req *model.BusinessCabinetStatusReq) (res model.BusinessCabinetStatusRes) {
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
        if res.Stop || time.Now().Sub(start) > 30 {
            return
        }
        time.Sleep(1 * time.Second)
    }
}
