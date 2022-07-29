// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-28
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "time"

// CabinetTaskType 业务类型
type CabinetTaskType string

const (
    CabinetTaskTypeRiderExchange    CabinetTaskType = "RDR_EXCHANGE"    // 骑手-换电
    CabinetTaskTypeRiderActive                      = "RDR_ACTIVE"      // 骑手-激活
    CabinetTaskTypeRiderUnSubscribe                 = "RDR_UNSUBSCRIBE" // 骑手-退租
    CabinetTaskTypePause                            = "RDR_PAUSE"       // 骑手-寄存
    CabinetTaskTypeContinue                         = "RDR_CONTINUE"    // 骑手-取消寄存
    CabinetTaskTypeManagerOpen                      = "MGR_OPEN"        // 管理-开门
)

type CabinetTaskSerial struct {
    Serial string `json:"serial"` // 电柜编码
}

// CabinetTask 电柜任务详情
type CabinetTask struct {
    Type  CabinetTaskType // 业务类别
    Start time.Time       // 开始时间
}
