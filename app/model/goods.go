// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-10-15, by liasica

package model

import (
	"errors"
	"slices"
	"time"

	"github.com/golang-module/carbon/v2"
)

var (
	ErrorGoodsPaymentDuplicate = errors.New("商品付款阶段重复")
	ErrorGoodsPaymentInvalid   = errors.New("商品付款阶段无效")
	ErrorGoodsPaymentEmpty     = errors.New("商品付款阶段为空")
	ErrorGoodsPaymentAmount    = errors.New("商品付款金额无效")
)

type GoodsPaymentPeriod int

const (
	GoodsPaymentPeriodOnce      GoodsPaymentPeriod = iota // 一次性
	GoodsPaymentPeriodDaily                               // 按天
	GoodsPaymentPeriodMonthly                             // 按月
	GoodsPaymentPeriodQuarterly                           // 按季度
	GoodsPaymentPeriodYearly                              // 按年
)

// GoodsPaymentPlanStage 商品付款阶段
type GoodsPaymentPlanStage struct {
	Period GoodsPaymentPeriod `json:"period"` // 付款周期
	Unit   int                `json:"unit"`   // 周期单位
	Amount float64            `json:"amount"` // 付款金额
}

// Equal 判断两个商品付款阶段是否相等
// 两个商品付款阶段相等，当且仅当周期和单位相等
func (stage GoodsPaymentPlanStage) Equal(other GoodsPaymentPlanStage) bool {
	return stage.Period == other.Period && stage.Unit == other.Unit
}

// GoodsPaymentPlan 商品付款阶段列表
type GoodsPaymentPlan []GoodsPaymentPlanStage

// Equal 判断两个商品付款阶段列表是否相等
func (plan GoodsPaymentPlan) Equal(other GoodsPaymentPlan) bool {
	if len(plan) != len(other) {
		return false
	}

	for i, option := range plan {
		if !option.Equal(other[i]) {
			return false
		}
	}

	return true
}

// BillingDates 获取商品付款阶段的账单日期
func (plan GoodsPaymentPlan) BillingDates(t time.Time) (dates []time.Time) {
	c := carbon.CreateFromStdTime(t).StartOfDay()
	dates = make([]time.Time, len(plan))
	for i, option := range plan {
		if i == 0 {
			dates[i] = c.StdTime()
			continue
		}
		switch option.Period {
		case GoodsPaymentPeriodOnce:
			if dates[i].IsZero() {
				dates[i] = c.StdTime()
			}
		case GoodsPaymentPeriodDaily:
			c = c.AddDays(option.Unit)
		case GoodsPaymentPeriodMonthly:
			c = c.AddMonthsNoOverflow(option.Unit)
		case GoodsPaymentPeriodQuarterly:
			c = c.AddQuartersNoOverflow(option.Unit)
		case GoodsPaymentPeriodYearly:
			c = c.AddYearsNoOverflow(option.Unit)
		}
		if option.Period != GoodsPaymentPeriodOnce {
			dates[i] = c.StdTime()
		}
	}
	return
}

// Stage 获取商品付款阶段的某个阶段
func (plan GoodsPaymentPlan) Stage(stage int) GoodsPaymentPlanStage {
	if stage >= len(plan) {
		return GoodsPaymentPlanStage{}
	}
	return plan[stage]
}

// GoodsPaymentPlans 商品付款方案列表
type GoodsPaymentPlans []GoodsPaymentPlan

// Valid 对商品付款方案表进行校验和排序
func (plans GoodsPaymentPlans) Valid() error {
	// 先进行一次排序
	slices.SortFunc(plans, func(i, j GoodsPaymentPlan) int {
		if len(i) < len(j) {
			return -1
		}
		return 1
	})

	var last GoodsPaymentPlan
	for _, options := range plans {
		// 判断是否重复
		if last != nil && last.Equal(options) {
			return ErrorGoodsPaymentDuplicate
		}

		// 判断是否为空
		if len(options) == 0 {
			return ErrorGoodsPaymentEmpty
		}

		for _, option := range options {
			// 判定付款周期是否有效
			if option.Period == GoodsPaymentPeriodOnce && len(options) != 1 {
				return ErrorGoodsPaymentInvalid
			}

			// 判定付款金额是否有效
			if option.Amount <= 0 {
				return ErrorGoodsPaymentAmount
			}
		}

		last = options
	}
	return nil
}

// BillingDates 获取商品付款方案表的账单日期
func (plans GoodsPaymentPlans) BillingDates(index int, t time.Time) (dates []time.Time) {
	if index >= len(plans) {
		return
	}
	return plans[index].BillingDates(t)
}

// Plan 获取商品付款方案表的某个方案
func (plans GoodsPaymentPlans) Plan(index int) GoodsPaymentPlan {
	if index >= len(plans) {
		return nil
	}
	return plans[index]
}

// PlanStage 获取商品付款方案表的某个阶段
// stage从1开始
func (plans GoodsPaymentPlans) PlanStage(index, stage int) (o GoodsPaymentPlanStage) {
	if index >= len(plans) {
		return
	}
	if stage >= len(plans[index]) {
		return
	}
	return plans[index][stage]
}

// PlanIndex 获取商品付款方案索引
func (plans GoodsPaymentPlans) PlanIndex(plan GoodsPaymentPlan) int {
	for i, p := range plans {
		if p.Equal(plan) {
			return i
		}
	}
	return -1
}
