// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-12
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
	"github.com/auroraride/aurservd/app/ec"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/pkg/cache"
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
func (c *Cabinet) ReserveAble(typ model.BusinessType, num map[model.ReserveBusinessKey]int) bool {
	if c.Status != model.CabinetStatusNormal.Value() || c.Health != model.CabinetHealthStatusOnline {
		return false
	}
	// 可用电池数量,可用空仓数量 锁仓不算可用
	var availableBatteryNum, availableEmptyBinNum int
	for _, bin := range c.Bin {
		if bin.DoorHealth && bin.Electricity.Value() >= cache.Float64(model.SettingExchangeMinBatteryKey) {
			availableBatteryNum += 1
		}
		if !bin.Battery && bin.DoorHealth {
			availableEmptyBinNum += 1
		}
	}
	switch typ {
	// 取电
	case model.BusinessTypeActive, model.BusinessTypeContinue:
		activeNum := num[model.NewReserveBusinessKey(c.ID, model.BusinessTypeActive)]
		continueNum := num[model.NewReserveBusinessKey(c.ID, model.BusinessTypeContinue)]
		return availableBatteryNum-(activeNum+continueNum) >= 2
	// 放电
	case model.BusinessTypePause, model.BusinessTypeUnsubscribe:
		pauseNum := num[model.NewReserveBusinessKey(c.ID, model.BusinessTypePause)]
		unsubscribeNum := num[model.NewReserveBusinessKey(c.ID, model.BusinessTypeUnsubscribe)]
		return availableEmptyBinNum-(pauseNum+unsubscribeNum) >= 2
	}
	return false
}

// UsingMicroService 是否使用微服务
// TODO 删除无用逻辑[是否使用微服务]
func (c *Cabinet) UsingMicroService() bool {
	// return c.Intelligent || c.Brand == adapter.CabinetBrandYundong || c.Brand == adapter.CabinetBrandTuobang
	return true
}
