// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import (
    sls "github.com/aliyun/aliyun-log-go-sdk"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ali"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/lithammer/shortuuid/v4"
    "time"
)

type Operate uint

const (
    OperatePersonBan      = iota // 封禁身份
    OperatePersonUnBan           // 解封身份
    OperateRiderBLock            // 封禁账户
    OperateRiderUnBLock          // 解封账户
    OperateSubscribeAlter        // 修改订阅时间
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
    case OperateSubscribeAlter:
        return "修改时间"
    default:
        return "未知操作"
    }
}

// OperateLog 系统操作日志
type OperateLog struct {
    ID string `json:"id" sls:"操作ID"`

    RefID    uint64 `json:"refId" sls:"关联ID" index:"doc"`
    RefTable string `json:"refTable" sls:"关联表" index:"doc"`

    Operate Operate `json:"operate" sls:"类别" string:"true" index:"doc"`
    Remark  string  `json:"remark" sls:"备注" index:"doc"`
    Before  string  `json:"before" sls:"操作前" index:"doc"`
    After   string  `json:"after" sls:"操作后" index:"doc"`

    ManagerID    uint64 `json:"managerID" sls:"操作人ID" index:"doc"`
    ManagerPhone string `json:"phone" sls:"操作人电话" index:"doc"`
    ManagerName  string `json:"managerName" sls:"操作人" index:"doc"`

    Time string `json:"time" sls:"时间" index:"doc"`
}

func (o *OperateLog) GetLogstoreName() string {
    return ar.Config.Aliyun.Sls.OperateLog
}

func NewOperateLog() *OperateLog {
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

func (o *OperateLog) Send() {
    PutLog(o)
}

// GetLogs 获取日志
// 参考 https://help.aliyun.com/document_detail/29029.htm?spm=a2c4g.11186623.0.0.21752a73O1u0I4#t13238.html
func (o *OperateLog) GetLogs(from time.Time, query string, params ...int64) []map[string]string {
    cfg := ar.Config.Aliyun.Sls
    var offset int64
    if len(params) > 0 {
        offset = params[0]
    }
    response, err := ali.NewSls().GetLogsV2(cfg.Project, cfg.OperateLog, &sls.GetLogRequest{
        From:    from.Unix(),
        To:      time.Now().Unix(),
        Reverse: true,
        Query:   query,
        Lines:   50,
        Offset:  offset,
    })
    if err != nil {
        return make([]map[string]string, 0)
    }
    return response.Logs
}
