// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-16
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "github.com/auroraride/aurservd/pkg/cache"
)

type BatteryElectricity float64

func NewBatteryElectricity(pos float64) BatteryElectricity {
    if pos < 0 {
        pos = 0
    }
    return BatteryElectricity(pos)
}

// IsBatteryFull 电池是否满电
func (be BatteryElectricity) IsBatteryFull() bool {
    return be >= BatteryElectricity(cache.BatteryFull(SettingBatteryFull))
}
