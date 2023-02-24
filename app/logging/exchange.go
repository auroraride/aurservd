// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-03
// Based on aurservd by liasica, magicrolan@qq.com.

package logging

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
)

type ExchangeLog struct {
    UUID       string `json:"uuid" index:"doc"`       // 操作ID
    Serial     string `json:"serial" index:"doc"`     // 电柜编码
    RiderID    uint64 `json:"riderId" index:"doc"`    // 骑手ID
    RiderPhone string `json:"riderPhone" index:"doc"` // 骑手电话
    Full       bool   `json:"full" index:"doc"`       // 是否满电操作

    Step        uint8  `json:"step" index:"doc"`        // 步骤序号
    Description string `json:"description" index:"doc"` // 描述

    Status uint8  `json:"status" index:"doc"` // 状态
    Result string `json:"result" index:"doc"` // 结果

    BinIndex int    `json:"binIndex" index:"doc"`          // 仓门index
    Bin      string `json:"bin" index:"doc" string:"true"` // 操作仓门

    Message     string  `json:"message" index:"doc"`     // 消息
    Electricity float64 `json:"electricity" index:"doc"` // 电池电量

    Time string `json:"time" index:"doc"` // 时间
}

func (e *ExchangeLog) GetLogstoreName() string {
    return ar.Config.Aliyun.Sls.ExchangeLog
}

func NewExchangeLog(riderID uint64, uid, serial, phone string, full bool) *ExchangeLog {
    return &ExchangeLog{
        RiderID:    riderID,
        UUID:       uid,
        RiderPhone: phone,
        Serial:     serial,
        Full:       full,
    }
}

func (e *ExchangeLog) Clone() *ExchangeLog {
    return &ExchangeLog{
        RiderID:    e.RiderID,
        UUID:       e.UUID,
        RiderPhone: e.RiderPhone,
        Serial:     e.Serial,
    }
}

func (e *ExchangeLog) SetBin(index int) *ExchangeLog {
    e.BinIndex = index
    e.Bin = fmt.Sprintf("%d号仓", index+1)
    return e
}

func (e *ExchangeLog) SetStep(step model.ExchangeStep) *ExchangeLog {
    e.Step = uint8(step)
    e.Description = step.String()
    return e
}

func (e *ExchangeLog) SetStatus(status model.TaskStatus) *ExchangeLog {
    e.Status = uint8(status)
    if status == model.TaskStatusSuccess {
        e.Result = "成功"
    } else {
        e.Result = "失败"
    }
    return e
}

func (e *ExchangeLog) SetElectricity(electricity model.BatterySoc) *ExchangeLog {
    e.Electricity = float64(electricity)
    return e
}

func (e *ExchangeLog) SetMessage(message string) *ExchangeLog {
    e.Message = message
    return e
}

func (e *ExchangeLog) Send() {
    PutLog(e)
}
