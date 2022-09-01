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
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/lithammer/shortuuid/v4"
    "time"
)

// OperateLog 系统操作日志
type OperateLog struct {
    ID string `json:"id" sls:"操作ID"`

    RefID    uint64 `json:"refId" sls:"关联ID" index:"doc"`
    RefTable string `json:"refTable" sls:"关联表" index:"doc"`

    Operate model.Operate `json:"operate" sls:"类别" string:"true" index:"doc"`
    Remark  string        `json:"remark" sls:"备注"`
    Before  string        `json:"before" sls:"操作前"`
    After   string        `json:"after" sls:"操作后"`

    OperatorType  model.OperatorType `json:"operatorType" sls:"操作人类别" index:"doc"` // 0管理员 1员工
    OperatorID    uint64             `json:"operatorId" sls:"操作人ID" index:"doc"`
    OperatorPhone string             `json:"operatorPhone" sls:"操作人电话" index:"doc"`
    OperatorName  string             `json:"operatorName" sls:"操作人" index:"doc"`

    Info string `json:"info" sls:"信息"`
    Time string `json:"time" sls:"时间"`
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
    switch ref.(type) {
    case model.TableSlsInfo:
        o.Info = ref.(model.TableSlsInfo).GetSLSLogInfo()
        break
    }
    return o
}

func (o *OperateLog) SetRefManually(table string, id uint64) *OperateLog {
    o.RefTable = table
    o.RefID = id
    return o
}

func (o *OperateLog) SetOperate(operate model.Operate) *OperateLog {
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
    if m != nil {
        o.OperatorID = m.ID
        o.OperatorPhone = m.Phone
        o.OperatorName = m.Name
        o.OperatorType = model.OperatorTypeManager
    }
    return o
}

func (o *OperateLog) SetAgent(ag *ent.Agent) *OperateLog {
    if ag != nil {
        o.OperatorID = ag.ID
        o.OperatorPhone = ag.Phone
        o.OperatorName = ag.Name
        o.OperatorType = model.OperatorTypeAgent
    }
    return o
}

func (o *OperateLog) SetEmployee(e *model.Employee) *OperateLog {
    if e != nil {
        o.OperatorID = e.ID
        o.OperatorPhone = e.Phone
        o.OperatorName = e.Name
        o.OperatorType = model.OperatorTypeEmployee
    }
    return o
}

func (o *OperateLog) SetCabinet(c *model.CabinetBasicInfo) *OperateLog {
    if c != nil {
        o.OperatorID = c.ID
        o.OperatorPhone = c.Serial
        o.OperatorName = c.Name
        o.OperatorType = model.OperatorTypeCabinet
    }
    return o
}

func (o *OperateLog) Send() {
    PutLog(o)
}

// GetLogs 获取日志
// 参考 https://help.aliyun.com/document_detail/29029.htm?spm=a2c4g.11186623.0.0.21752a73O1u0I4#t13238.html
func (o *OperateLog) GetLogs(from time.Time, query string, offset, limit int64) []map[string]string {
    cfg := ar.Config.Aliyun.Sls
    response, err := ali.NewSls().GetLogsV2(cfg.Project, cfg.OperateLog, &sls.GetLogRequest{
        From:    from.Unix(),
        To:      time.Now().Unix(),
        Reverse: true,
        Query:   query,
        Lines:   limit,
        Offset:  offset,
    })
    if err != nil {
        return make([]map[string]string, 0)
    }
    return response.Logs
}
