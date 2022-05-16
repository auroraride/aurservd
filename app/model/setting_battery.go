// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-16
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type BatteryElectricity float64

var (
    // BatteryElectricityFull 满电设置
    BatteryElectricityFull BatteryElectricity = 80
)

// BatteryElectricityBootstrap 从数据库中加载满电配置
// TODO 当满电设置更新的时候自动更新该配置
func BatteryElectricityBootstrap() {
    BatteryElectricityFull = 80
}

func NewBatteryElectricity(pos float64) BatteryElectricity {
    if pos < 0 {
        pos = 0
    }
    return BatteryElectricity(pos)
}

// IsBatteryFull 电池是否满电
func (be BatteryElectricity) IsBatteryFull() bool {
    return be >= BatteryElectricityFull
}
