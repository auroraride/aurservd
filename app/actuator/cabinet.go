// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-05
// Based on aurservd by liasica, magicrolan@qq.com.

package actuator

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
)

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
    Health         uint8 `json:"health" bson:"health"`                 // 电柜健康状态 0离线 1正常 2故障
    Doors          uint  `json:"doors" bson:"doors"`                   // 总仓位
    BatteryNum     uint  `json:"batteryNum" bson:"batteryNum"`         // 总电池数
    BatteryFullNum uint  `json:"batteryFullNum" bson:"batteryFullNum"` // 总满电电池数
}
