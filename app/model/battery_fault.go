// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-06
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "database/sql/driver"
    "fmt"
)

type BatteryFault string

const (
    BatteryFaultVoltageLow                 BatteryFault = "TVL" // 总压低
    BatteryFaultVoltageHigh                BatteryFault = "TVH" // 总压高
    BatteryFaultMonVoltageLow              BatteryFault = "MVL" // 单体低
    BatteryFaultMonVoltageHigh             BatteryFault = "MVH" // 单体高
    BatteryFaultDischargeOvercurrent       BatteryFault = "DOC" // 放电过流
    BatteryFaultChargeOvercurrent          BatteryFault = "COC" // 充电过流
    BatteryFaultSocLow                     BatteryFault = "SCL" // SOC低
    BatteryFaultChargingTemperatureHigh    BatteryFault = "CTH" // 充电高温
    BatteryFaultChargingTemperatureLow     BatteryFault = "CTL" // 充电低温
    BatteryFaultDisChargingTemperatureHigh BatteryFault = "DTH" // 放电高温
    BatteryFaultDisChargingTemperatureLow  BatteryFault = "DTL" // 放电低温
    BatteryFaultShortCircuit               BatteryFault = "SCT" // 短路
    BatteryFaultMosTemperatureHigh         BatteryFault = "MTH" // MOS高温
)

func (b BatteryFault) Text() string {
    switch b {
    case BatteryFaultVoltageLow:
        return "总压低"
    case BatteryFaultVoltageHigh:
        return "总压高"
    case BatteryFaultMonVoltageLow:
        return "单体低"
    case BatteryFaultMonVoltageHigh:
        return "单体高"
    case BatteryFaultDischargeOvercurrent:
        return "放电过流"
    case BatteryFaultChargeOvercurrent:
        return "充电过流"
    case BatteryFaultSocLow:
        return "SOC低"
    case BatteryFaultChargingTemperatureHigh:
        return "充电高温"
    case BatteryFaultChargingTemperatureLow:
        return "充电低温"
    case BatteryFaultDisChargingTemperatureHigh:
        return "放电高温"
    case BatteryFaultDisChargingTemperatureLow:
        return "放电低温"
    case BatteryFaultShortCircuit:
        return "短路"
    case BatteryFaultMosTemperatureHigh:
        return "MOS高温"
    }
    return " - "
}

func (b BatteryFault) String() string {
    return string(b)
}

func (BatteryFault) Values() []string {
    return []string{
        BatteryFaultVoltageLow.String(),
        BatteryFaultVoltageHigh.String(),
        BatteryFaultMonVoltageLow.String(),
        BatteryFaultMonVoltageHigh.String(),
        BatteryFaultDischargeOvercurrent.String(),
        BatteryFaultChargeOvercurrent.String(),
        BatteryFaultSocLow.String(),
        BatteryFaultChargingTemperatureHigh.String(),
        BatteryFaultChargingTemperatureLow.String(),
        BatteryFaultDisChargingTemperatureHigh.String(),
        BatteryFaultDisChargingTemperatureLow.String(),
        BatteryFaultShortCircuit.String(),
        BatteryFaultMosTemperatureHigh.String(),
    }
}

func (b *BatteryFault) Scan(src interface{}) error {
    switch v := src.(type) {
    case nil:
        return nil
    case string:
        *b = BatteryFault(v)
    }
    return nil
}

func (b BatteryFault) Value() (driver.Value, error) {
    return b, nil
}

func BatteryFaultValidator(t BatteryFault) error {
    switch t {
    case BatteryFaultVoltageLow,
        BatteryFaultVoltageHigh,
        BatteryFaultMonVoltageLow,
        BatteryFaultMonVoltageHigh,
        BatteryFaultDischargeOvercurrent,
        BatteryFaultChargeOvercurrent,
        BatteryFaultSocLow,
        BatteryFaultChargingTemperatureHigh,
        BatteryFaultChargingTemperatureLow,
        BatteryFaultDisChargingTemperatureHigh,
        BatteryFaultDisChargingTemperatureLow,
        BatteryFaultShortCircuit,
        BatteryFaultMosTemperatureHigh:
        return nil
    default:
        return fmt.Errorf("未知的故障类别: %q", t)
    }
}
