// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-28
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "time"

// BusinessCabinetType 业务类型
type BusinessCabinetType string

const (
    BusinessCabinetTypeRiderExchange    BusinessCabinetType = "RDR_EXCHANGE"    // 骑手-换电
    BusinessCabinetTypeRiderActive                          = "RDR_ACTIVE"      // 骑手-激活
    BusinessCabinetTypeRiderUnSubscribe                     = "RDR_UNSUBSCRIBE" // 骑手-退租
    BusinessCabinetTypePause                                = "RDR_PAUSE"       // 骑手-寄存
    BusinessCabinetTypeContinue                             = "RDR_CONTINUE"    // 骑手-取消寄存
    BusinessCabinetTypeManagerOpen                          = "MGR_OPEN"        // 管理-开门
)

type BusinessCabinetSerial struct {
    Serial string `json:"serial"` // 电柜编码
}

// BusinessCabinet 电柜业务详情
type BusinessCabinet struct {
    Type  BusinessCabinetType // 业务类别
    Start time.Time           // 开始时间
}
