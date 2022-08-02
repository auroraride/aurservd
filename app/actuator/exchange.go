// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-31
// Based on aurservd by liasica, magicrolan@qq.com.

package actuator

import (
    "github.com/auroraride/aurservd/app/model"
    "time"
)

// ExchangeStep 换电步骤
// RiderCabinetOperateStep
type ExchangeStep uint8

const (
    ExchangeStepOpenEmpty ExchangeStep = iota + 1 // 第一步, 开启空电仓
    ExchangeStepPutInto                           // 第二步, 放入旧电池并关闭仓门
    ExchangeStepOpenFull                          // 第三步, 开启满电仓
    ExchangeStepPutOut                            // 第四步, 取出新电池并关闭仓门
)

func (ros ExchangeStep) String() string {
    switch ros {
    case ExchangeStepOpenEmpty:
        return "第一步, 开启空电仓"
    case ExchangeStepPutInto:
        return "第二步, 放入旧电池并关闭仓门"
    case ExchangeStepOpenFull:
        return "第三步, 开启满电仓"
    case ExchangeStepPutOut:
        return "第四步, 取出新电池并关闭仓门"
    }
    return "未知"
}

// ExchangeDoorStatus 柜门状态(处理换电用)
// CabinetBinDoorStatus
type ExchangeDoorStatus uint8

const (
    ExchangeDoorStatusUnknown      ExchangeDoorStatus = iota // 未知
    ExchangeDoorStatusClose                                  // 关闭
    ExchangeDoorStatusOpen                                   // 开启
    ExchangeDoorStatusFail                                   // 故障
    ExchangeDoorStatusBatteryFull                            // 电池未取出
    ExchangeDoorStatusBatteryEmpty                           // 电池未放入
)

func (bds ExchangeDoorStatus) String() string {
    switch bds {
    case ExchangeDoorStatusClose:
        return "关闭"
    case ExchangeDoorStatusOpen:
        return "开启"
    case ExchangeDoorStatusFail:
        return "故障"
    case ExchangeDoorStatusBatteryFull:
        return "电池未取出"
    case ExchangeDoorStatusBatteryEmpty:
        return "电池未放入"
    }
    return "未知"
}

// ExchangeDoorError 换电故障
// CabinetBinDoorError
var ExchangeDoorError = map[ExchangeDoorStatus]string{
    ExchangeDoorStatusUnknown:      "仓门状态未知",
    ExchangeDoorStatusFail:         "仓门故障",
    ExchangeDoorStatusBatteryFull:  "电池未取出",
    ExchangeDoorStatusBatteryEmpty: "电池未放入",
}

// Exchange 换电信息
type Exchange struct {
    Alternative bool         `json:"alternative" bson:"alternative"` // 是否备选方案
    Success     bool         `json:"success" bson:"success"`         // 是否成功
    Step        ExchangeStep `json:"step" bson:"step"`               // 当前步骤
    Empty       BinInfo      `json:"empty" bson:"empty"`             // 空仓位
    Fully       BinInfo      `json:"fully" bson:"fully"`             // 满电仓位
    Model       string       `json:"model" bson:"model"`             // 电池型号
}

// BinInfo 任务电柜仓位信息
type BinInfo struct {
    Index       int                      `json:"index" bson:"index"`               // 仓位index
    Electricity model.BatteryElectricity `json:"electricity" bson:"electricity"`   // 电量
    Voltage     float64                  `json:"voltage" bson:"voltage"`           // 电压(V)
    OpenAt      *time.Time               `json:"openAt,omitempty" bson:"openAt"`   // 开门时间
    CloseAt     *time.Time               `json:"closeAt,omitempty" bson:"closeAt"` // 关门时间
}
