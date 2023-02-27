// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-12
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
    "github.com/auroraride/adapter"
    "github.com/auroraride/aurservd/app/ec"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent/business"
)

func (c *Cabinet) GetTaskInfo() *ec.Cabinet {
    return &ec.Cabinet{
        Health:         c.Health,
        Doors:          c.Doors,
        BatteryNum:     c.BatteryNum,
        BatteryFullNum: c.BatteryFullNum,
        Bins:           c.Bin,
    }
}

// ReserveAble 是否可预约
// num 当前已有预约数量
func (c *Cabinet) ReserveAble(typ business.Type, num int) bool {
    if c.Status != model.CabinetStatusNormal.Value() || c.Health != model.CabinetHealthStatusOnline {
        return false
    }
    switch typ {
    // 取电
    case business.TypeActive, business.TypeContinue:
        return c.BatteryNum-num >= 2
    // 放电
    case business.TypePause, business.TypeUnsubscribe:
        return c.EmptyBinNum-num >= 2
    }
    return false
}

func (c *Cabinet) UsingMicroService() bool {
    return c.Intelligent || c.Brand == adapter.CabinetYundong.String()
}
