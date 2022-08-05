// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-03
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/alibabacloud-go/tea/tea"
    sls "github.com/aliyun/aliyun-log-go-sdk"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/provider"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "github.com/lithammer/shortuuid/v4"
    log "github.com/sirupsen/logrus"
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
func (s *cabinetMgrService) Maintain(req *model.CabinetMaintainReq) {
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

    _, _ = cab.Update().SetStatus(status.Raw()).Save(s.ctx)

    // 记录日志
    go logging.NewOperateLog().
        SetRef(cab).
        SetModifier(s.modifier).
        SetOperate(model.OperateAssistanceAllocate).
        SetDiff(model.CabinetStatus(cab.Status).String(), status.String()).
        Send()
}

// BinOperate 仓门操控
func (s *cabinetMgrService) BinOperate(req *model.CabinetDoorOperateReq) bool {
    if s.modifier == nil {
        snag.Panic("权限校验失败")
    }

    ec.BusyFromIDX(*req.ID)

    cs := NewCabinetWithModifier(s.modifier)

    cab := cs.QueryOne(*req.ID)

    if model.CabinetStatus(cab.Status) != model.CabinetStatusMaintenance {
        snag.Panic("非操作维护中不可操作")
    }

    task := &ec.Task{
        CabinetID: *req.ID,
        Serial:    cab.Serial,
        Cabinet: ec.Cabinet{
            Health:         cab.Health,
            Doors:          cab.Doors,
            BatteryNum:     cab.BatteryNum,
            BatteryFullNum: cab.BatteryFullNum,
        },
    }

    switch *req.Operation {
    case model.CabinetDoorOperateOpen:
        task.Job = ec.JobManagerOpen
        break
    case model.CabinetDoorOperateLock:
        task.Job = ec.JobManagerLock
        break
    case model.CabinetDoorOperateUnlock:
        task.Job = ec.JobManagerUnLock
        break
    }

    var status bool
    var err error

    // 创建并开始任务
    task.CreateX().Start()

    // 结束回调
    defer func() {
        ts := ec.TaskStatusSuccess
        if !status {
            ts = ec.TaskStatusFail
            task.Message = err.Error()
        }
        task.Stop(ts)
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
        Job:       ec.JobManagerReboot,
        Cabinet: ec.Cabinet{
            Health:         cab.Health,
            Doors:          cab.Doors,
            BatteryNum:     cab.BatteryNum,
            BatteryFullNum: cab.BatteryFullNum,
        },
    }

    task.CreateX().Start()

    // 结束回调
    defer func() {
        ts := ec.TaskStatusSuccess
        if !status {
            ts = ec.TaskStatusFail
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
        slsCfg := ar.Config.Aliyun.Sls
        lg := &sls.LogGroup{
            Logs: []*sls.Log{{
                Time: tea.Uint32(uint32(now.Unix())),
                Contents: logging.GenerateLogContent(&logging.DoorOperateLog{
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
                }),
            }},
        }
        err := ali.NewSls().PutLogs(slsCfg.Project, slsCfg.DoorLog, lg)
        if err != nil {
            log.Error(err)
            return
        }
    }()

    return status
}
