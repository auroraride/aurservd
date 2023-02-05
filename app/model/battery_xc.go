// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-05
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "github.com/auroraride/adapter/defs/xcdef"
    "github.com/auroraride/adapter/defs/xcdef/proto/entpb"
    "github.com/auroraride/aurservd/pkg/silk"
)

type XcBmsBattery struct {
    Status BatteryStatus `json:"status"` // 状态, 0:静置 1:充电 2:放电 3:异常(此时faults字段存在)
    Soc    uint32        `json:"soc"`    // 电量, 单位1%

    Charge    *bool         `json:"charge"`           // 充电是否开启
    DisCharge *bool         `json:"disCharge"`        // 放电是否开启
    Faults    *xcdef.Faults `json:"faults,omitempty"` // 故障列表, 0:总压低, 1:总压高, 2:单体低, 3:单体高, 6:放电过流, 7:充电过流, 8:SOC低, 11:充电高温, 12:充电低温, 13:放电高温, 14:放电低温, 15:短路, 16:MOS高温
}

func NewXcBmsBattery(hb *entpb.Heartbeat) (item *XcBmsBattery) {
    item = new(XcBmsBattery)
    item.Soc = hb.Soc

    if hb.MosStatus != nil {
        ms := new(xcdef.MosStatus)
        ms.FromBytes(hb.MosStatus)
        item.Charge = silk.Pointer(ms.CanCharge())
        item.DisCharge = silk.Pointer(ms.CanDisCharge())
    }

    switch {
    default:
        item.Status = BatteryStatusIdle
    case hb.Current > 0:
        item.Status = BatteryStatusCharging
    case hb.Current < 0:
        item.Status = BatteryStatusDisCharging
    }

    if hb.Faults != nil {
        faults := new(xcdef.Faults)
        faults.FromBytes(hb.Faults)
        item.Faults = faults
        item.Status = BatteryStatusFault
    }

    return
}
