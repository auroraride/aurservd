// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-12
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
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
func (c *Cabinet) ReserveAble(typ business.Type, num map[model.ReserveBusinessKey]int) bool {
	if c.Status != model.CabinetStatusNormal.Value() || c.Health != model.CabinetHealthStatusOnline {
		return false
	}
	switch typ {
	// 取电
	case business.TypeActive, business.TypeContinue:
		activeNum := num[model.NewReserveBusinessKey(c.ID, business.TypeActive.String())]
		continueNum := num[model.NewReserveBusinessKey(c.ID, business.TypeContinue.String())]
		return c.BatteryNum-(activeNum+continueNum) >= 2
	// 放电
	case business.TypePause, business.TypeUnsubscribe:
		pauseNum := num[model.NewReserveBusinessKey(c.ID, business.TypePause.String())]
		unsubscribeNum := num[model.NewReserveBusinessKey(c.ID, business.TypeUnsubscribe.String())]
		return c.EmptyBinNum-(pauseNum+unsubscribeNum) >= 2
	}
	return false
}

// UsingMicroService 是否使用微服务
// TODO 删除无用逻辑[是否使用微服务]
func (c *Cabinet) UsingMicroService() bool {
	// return c.Intelligent || c.Brand == adapter.CabinetBrandYundong || c.Brand == adapter.CabinetBrandTuobang
	return true
}
