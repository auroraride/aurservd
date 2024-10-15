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
	ErrorGoodsPaymentDuplicate = errors.New("商品付款选项重复")
	ErrorGoodsPaymentInvalid   = errors.New("商品付款选项无效")
	ErrorGoodsPaymentEmpty     = errors.New("商品付款选项为空")
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

// GoodsPaymentPlanOption 商品付款选项
type GoodsPaymentPlanOption struct {
	Period GoodsPaymentPeriod `json:"period"` // 付款周期
	Unit   int                `json:"unit"`   // 周期单位
	Amount float64            `json:"amount"` // 付款金额
}

// Equal 判断两个商品付款选项是否相等
// 两个商品付款选项相等，当且仅当周期和单位相等
func (option GoodsPaymentPlanOption) Equal(other GoodsPaymentPlanOption) bool {
	return option.Period == other.Period && option.Unit == other.Unit
}

// GoodsPaymentPlan 商品付款选项列表
type GoodsPaymentPlan []GoodsPaymentPlanOption

// Equal 判断两个商品付款选项列表是否相等
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

// BillingDates 获取商品付款选项的账单日期
func (plan GoodsPaymentPlan) BillingDates(t time.Time) (dates []time.Time) {
	c := carbon.CreateFromStdTime(t).StartOfDay()
	dates = make([]time.Time, len(plan))
	for i, option := range plan {
		if option.Period == GoodsPaymentPeriodOnce {
			dates[i] = c.StdTime()
			continue
		}
		switch option.Period {
		case GoodsPaymentPeriodOnce:
			dates[i] = c.StdTime()
			return
		case GoodsPaymentPeriodDaily:
			dates[i] = c.AddDays(option.Unit).StdTime()
		case GoodsPaymentPeriodMonthly:
			dates[i] = c.AddMonthsNoOverflow(option.Unit).StdTime()
		case GoodsPaymentPeriodQuarterly:
			dates[i] = c.AddQuartersNoOverflow(option.Unit).StdTime()
		case GoodsPaymentPeriodYearly:
			dates[i] = c.AddYearsNoOverflow(option.Unit).StdTime()
		}
	}
	return
}

// GoodsPaymentPlans 商品付款方案表
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
