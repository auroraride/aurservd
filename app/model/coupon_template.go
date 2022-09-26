// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-26
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "golang.org/x/exp/slices"
)

// CouponRule 优惠券规则
type CouponRule uint8

const (
    CouponRuleExclusive CouponRule = iota + 1 // 互斥
    CouponRuleStackable                       // 可叠加
)

func (c CouponRule) IsValid() bool {
    return slices.Contains(CouponRules, c)
}

// CouponExpiration 优惠券过期时间
type CouponExpiration uint8

const (
    CouponExpirationFixed    CouponExpiration = iota + 1 // 固定时间 (固定自然时间)
    CouponExpirationRelative                             // 相对时间 (自领取日算起)
)

func (c CouponExpiration) IsValid() bool {
    return slices.Contains(CouponExpirations, c)
}

var (
    CouponRules       = []CouponRule{CouponRuleExclusive, CouponRuleStackable}
    CouponExpirations = []CouponExpiration{CouponExpirationFixed, CouponExpirationRelative}
)

type CouponTemplate struct {
    Name       string           `json:"name" validate:"required" trans:"名称"`
    Rule       CouponRule       `json:"rule" validate:"required,enum" trans:"规则" enums:"1,2"`           // 1:互斥 2:可叠加
    Expiration CouponExpiration `json:"expiration" validate:"required,enum" trans:"过期时间" enums:"1,2"` // 1:固定时间(固定自然时间) 2:相对时间(自领取日算起)

    Multiple bool `json:"multiple"` // 同类券是否可多选

    Plans   *[]uint64 `json:"plans,omitempty" validate:"required_without=Cities" trans:"绑定骑士卡"` // 和`cities`不能同时为空也不能同时存在
    Cities  *[]uint64 `json:"cities,omitempty" validate:"required_without=Plans" trans:"可用城市"`   // 和`plans`不能同时为空也不能同时存在
    Coupons *[]uint64 `json:"coupons,omitempty" validate:"required_if=Rule 2" trans:"叠加优惠券"`    // `rule = 1`时必填
}

type CouponTemplateDetail struct {
    ID          uint64   `json:"id"`
    UpdatedAt   string   `json:"updatedAt"`             // 更新时间
    Total       int      `json:"total"`                 // 总数
    Remaining   int      `json:"remaining"`             // 剩余
    Used        int      `json:"used"`                  // 已使用
    Unused      int      `json:"unused"`                // 未使用
    Expired     int      `json:"expired"`               // 已过期
    PlanNames   []string `json:"planNames,omitempty"`   // 绑定骑士卡名称
    CityNames   []string `json:"cityNames,omitempty"`   // 可用城市名称
    CouponNames []string `json:"couponNames,omitempty"` // 可叠加优惠券名称
    CouponTemplate
}
