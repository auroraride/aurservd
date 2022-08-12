// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-12
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import "github.com/auroraride/aurservd/app/ec"

func (c *Cabinet) GetTaskInfo() ec.Cabinet {
    return ec.Cabinet{
        Health:         c.Health,
        Doors:          c.Doors,
        BatteryNum:     c.BatteryNum,
        BatteryFullNum: c.BatteryFullNum,
        Bins:           c.Bin,
    }
}
