// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-08
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type CtxRiderKey struct{}
type CtxModifierKey struct{}
type CtxEmployeeKey struct{}

type StoreCabiletGoal uint8

const (
	StockGoalAll StoreCabiletGoal = iota
	StockGoalStore
	StockGoalCabinet
	StockGoalStation
)

func (sg StoreCabiletGoal) String() string {
	switch sg {
	case StockGoalStore:
		return "门店"
	case StockGoalCabinet:
		return "电柜"
	case StockGoalStation:
		return "站点"
	default:
		return ""
	}
}

func (sg StoreCabiletGoal) SQLString() string {
	return map[StoreCabiletGoal]string{
		StockGoalStore:   "store",
		StockGoalCabinet: "cabinet",
	}[sg]
}
