// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-16
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "github.com/auroraride/aurservd/pkg/cache"
)

type BatterySoc float64

func NewBatterySoc(pos float64) BatterySoc {
    if pos < 0 {
        pos = 0
    }
    return BatterySoc(pos)
}

// IsBatteryFull 电池是否满电
func (be BatterySoc) IsBatteryFull() bool {
    return be >= BatterySoc(cache.Float64(SettingBatteryFullKey))
}

func (be BatterySoc) Value() float64 {
    return float64(be)
}
