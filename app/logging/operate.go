// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import (
    "github.com/alibabacloud-go/tea/tea"
    sls "github.com/aliyun/aliyun-log-go-sdk"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/golang-module/carbon/v2"
    "github.com/lithammer/shortuuid/v4"
    log "github.com/sirupsen/logrus"
    "time"
)

type Operate uint

const (
    OperatePersonBan    = iota // 封禁身份
    OperatePersonUnBan         // 解封身份
    OperateRiderBLock          // 封禁账户
    OperateRiderUnBLock        // 解封账户
)

func (o Operate) String() string {
    switch o {
    case OperatePersonBan:
        return "封禁用户"
    case OperatePersonUnBan:
        return "解封用户"
    case OperateRiderBLock:
        return "封禁账户"
    case OperateRiderUnBLock:
        return "解封账户"
    default:
        return "未知操作"
    }
}

// OperateLog 系统操作日志
type OperateLog struct {
    ID string `json:"id" sls:"操作ID"`

    RefID    uint64 `json:"refId" sls:"关联ID" index:""`
    RefTable string `json:"refTable" sls:"关联表" index:""`

    Operate Operate `json:"operate" sls:"类别" index:"doc"`
    Remark  string  `json:"remark" sls:"备注"`
    Before  string  `json:"before" sls:"操作前"`
    After   string  `json:"after" sls:"操作后"`

    ManagerID    uint64 `json:"managerID" sls:"操作人ID" index:""`
    ManagerPhone string `json:"phone" sls:"操作人电话" index:""`
    ManagerName  string `json:"managerName" sls:"操作人" index:"doc"`

    Time string `json:"time" sls:"时间" index:""`
}

func CreateOperateLog() *OperateLog {
    return &OperateLog{
        ID: shortuuid.New(),
    }
}

func (o *OperateLog) SetRef(ref model.Table) *OperateLog {
    o.RefTable = ref.GetTableName()
    o.RefID = ref.GetID()
    return o
}

func (o *OperateLog) SetOperate(operate Operate) *OperateLog {
    o.Operate = operate
    return o
}

func (o *OperateLog) SetRemark(remark string) *OperateLog {
    o.Remark = remark
    return o
}

func (o *OperateLog) SetDiff(before, after string) *OperateLog {
    o.Before = before
    o.After = after
    return o
}

func (o *OperateLog) SetModifier(m *model.Modifier) *OperateLog {
    o.ManagerID = m.ID
    o.ManagerPhone = m.Phone
    o.ManagerName = m.Name
    return o
}

// PutOperateLog 提交操作日志
func (o *OperateLog) PutOperateLog() {
    go func() {
        now := time.Now()
        if o.Time == "" {
            o.Time = now.Format(carbon.DateTimeLayout)
        }
        cfg := ar.Config.Aliyun.Sls
        err := ali.NewSls().PutLogs(cfg.Project, cfg.OperateLog, &sls.LogGroup{
            Logs: []*sls.Log{{
                Time:     tea.Uint32(uint32(now.Unix())),
                Contents: GenerateLogContent(o),
            }},
        })
        if err != nil {
            log.Error(err)
            return
        }
    }()
}
