// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-03
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/adapter/defs/cabdef"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/provider"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "github.com/lithammer/shortuuid/v4"
    "time"
)

type cabinetMgrService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewCabinetMgr() *cabinetMgrService {
    return &cabinetMgrService{
        ctx: context.Background(),
    }
}

func NewCabinetMgrWithModifier(m *model.Modifier) *cabinetMgrService {
    s := NewCabinetMgr()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// Maintain 设置电柜操作维护
func (s *cabinetMgrService) Maintain(req *model.CabinetMaintainReq) (detail *model.CabinetDetailRes) {
    if req.Maintain == nil {
        snag.Panic("参数请求错误")
    }
    cab := NewCabinet().QueryOne(req.ID)

    if model.CabinetStatus(cab.Status) == model.CabinetStatusPending {
        snag.Panic("未投放电柜无法操作")
    }

    status := model.CabinetStatusNormal
    if *req.Maintain {
        status = model.CabinetStatusMaintenance
    }

    detail = NewCabinetWithModifier(s.modifier).Detail(cab)

    _, err := cab.Update().SetStatus(status.Value()).Save(s.ctx)
    if err != nil {
        snag.Panic(err)
    }

    detail.Status = status

    // 记录日志
    go logging.NewOperateLog().
        SetRef(cab).
        SetModifier(s.modifier).
        SetOperate(model.OperateCabinetMaintain).
        SetDiff(model.CabinetStatus(cab.Status).String(), status.String()).
        SetRemark(cab.Remark).
        Send()

    return
}

// BinOperate 仓门操控
func (s *cabinetMgrService) BinOperate(req *model.CabinetDoorOperateReq) bool {
    if s.modifier == nil {
        snag.Panic("权限校验失败")
    }

    cs := NewCabinetWithModifier(s.modifier)

    cab := cs.QueryOne(req.ID)

    if model.CabinetStatus(cab.Status) != model.CabinetStatusMaintenance {
        snag.Panic("非操作维护中不可操作")
    }

    if cab.Intelligent {
        var op cabdef.Operate
        switch *req.Operation {
        case model.CabinetDoorOperateOpen:
            op = cabdef.OperateDoorOpen
        case model.CabinetDoorOperateLock:
            op = cabdef.OperateBinDisable
        case model.CabinetDoorOperateUnlock:
            op = cabdef.OperateBinEnable
        }
        return NewIntelligentCabinet(s.modifier).Operate(cab, op, req)
    }

    ec.BusyFromIDX(req.ID)

    task := &ec.Task{
        CabinetID: req.ID,
        Serial:    cab.Serial,
        Cabinet:   cab.GetTaskInfo(),
    }

    switch *req.Operation {
    case model.CabinetDoorOperateOpen:
        task.Job = model.JobManagerOpen
        break
    case model.CabinetDoorOperateLock:
        task.Job = model.JobManagerLock
        break
    case model.CabinetDoorOperateUnlock:
        task.Job = model.JobManagerUnLock
        break
    }

    var status bool
    var err error

    // 创建并开始任务
    task.Create().Start()

    // 结束回调
    defer func() {
        if v := recover(); v != nil {
            err = fmt.Errorf("%v", v)
        }

        ts := model.TaskStatusSuccess
        if !status {
            ts = model.TaskStatusFail
            task.Message = err.Error()
        }

        task.Stop(ts)

        if err != nil {
            snag.Panic(err)
        }
    }()

    // 柜门操作
    status, err = cs.DoorOperate(req, model.CabinetDoorOperator{
        ID:    s.modifier.ID,
        Role:  model.CabinetDoorOperatorRoleManager,
        Name:  s.modifier.Name,
        Phone: s.modifier.Phone,
    })

    if err != nil {
        snag.Panic(err)
    }

    return status
}

// Reboot 重启电柜
func (s *cabinetMgrService) Reboot(req *model.IDPostReq) bool {
    if s.modifier == nil {
        snag.Panic("权限校验失败")
    }

    ec.BusyFromIDX(req.ID)

    now := time.Now()
    opId := shortuuid.New()

    cab := NewCabinetWithModifier(s.modifier).QueryOne(req.ID)

    if model.CabinetStatus(cab.Status) != model.CabinetStatusMaintenance {
        snag.Panic("非操作维护中不可操作")
    }

    if cab.Brand == model.CabinetBrandKaixin.Value() {
        snag.Panic("凯信电柜不支持该操作")
    }

    var prov provider.Provider
    var status bool

    // 创建并开始任务
    task := &ec.Task{
        CabinetID: req.ID,
        Serial:    cab.Serial,
        Job:       model.JobManagerReboot,
        Cabinet:   cab.GetTaskInfo(),
    }

    task.Create().Start()

    // 结束回调
    defer func() {
        ts := model.TaskStatusSuccess
        if !status {
            ts = model.TaskStatusFail
            task.Message = "重启失败"
        }
        task.Stop(ts)
    }()

    // 请求云动重启
    prov = provider.NewYundong()
    status = prov.Reboot(s.modifier.Name+"-"+opId, cab.Serial)

    brand := model.CabinetBrand(cab.Brand)
    go func() {
        // 上传日志
        dlog := &logging.DoorOperateLog{
            ID:            opId,
            Brand:         brand.String(),
            OperatorName:  s.modifier.Name,
            OperatorID:    s.modifier.ID,
            OperatorPhone: s.modifier.Phone,
            OperatorRole:  model.CabinetDoorOperatorRoleManager,
            Serial:        cab.Serial,
            Operation:     "重启",
            Success:       status,
            Time:          now.Format(carbon.DateTimeLayout),
        }
        dlog.Send()
    }()

    return status
}
