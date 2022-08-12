// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-05
// Based on aurservd by liasica, magicrolan@qq.com.

package ec

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
)

// DoorStatus 柜门状态(处理换电用)
type DoorStatus uint8

const (
    DoorStatusUnknown            DoorStatus = iota // 未知
    DoorStatusClose                                // 关闭
    DoorStatusOpen                                 // 开启
    DoorStatusFail                                 // 故障
    DoorStatusBatteryFull                          // 电池未取出
    DoorStatusBatteryEmpty                         // 电池未放入
    DoorStatusCabinetStatusError                   // 电柜状态获取失败
)

func (bds DoorStatus) String() string {
    switch bds {
    case DoorStatusClose:
        return "关闭"
    case DoorStatusOpen:
        return "开启"
    case DoorStatusFail:
        return "故障"
    case DoorStatusBatteryFull:
        return "电池未取出"
    case DoorStatusBatteryEmpty:
        return "电池未放入"
    case DoorStatusCabinetStatusError:
        return "电柜状态获取失败"
    }
    return "未知"
}

// DoorError 换电故障
var DoorError = map[DoorStatus]string{
    DoorStatusUnknown:      "仓门状态未知",
    DoorStatusFail:         "仓门故障",
    DoorStatusBatteryFull:  "电池未取出",
    DoorStatusBatteryEmpty: "电池未放入",
}

// BinInfo 任务电柜仓位信息
type BinInfo struct {
    Index       int                      `json:"index" bson:"index"`             // 仓位index
    Electricity model.BatteryElectricity `json:"electricity" bson:"electricity"` // 电量
    Voltage     float64                  `json:"voltage" bson:"voltage"`         // 电压(V)
}

func (b *BinInfo) String() string {
    return fmt.Sprintf(
        "%d号仓, 电压: %.2fV, 电流: %2.fA",
        b.Index+1,
        b.Voltage,
        b.Electricity,
    )
}

// Cabinet 任务电柜设备信息
type Cabinet struct {
    Health         uint8             `json:"health" bson:"health"`                 // 电柜健康状态 0离线 1正常 2故障
    Doors          int               `json:"doors" bson:"doors"`                   // 总仓位
    BatteryNum     int               `json:"batteryNum" bson:"batteryNum"`         // 总电池数
    BatteryFullNum int               `json:"batteryFullNum" bson:"batteryFullNum"` // 总满电电池数
    Bins           model.CabinetBins `json:"bins" bson:"bins"`                     // 仓位信息
}
