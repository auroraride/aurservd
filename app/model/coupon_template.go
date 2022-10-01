// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-28
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "golang.org/x/exp/slices"
    "time"
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

func (c CouponRule) Value() uint8 {
    return uint8(c)
}

// CouponDurationRule 优惠券过期时间
type CouponDurationRule uint8

const (
    CouponDurationFixed    CouponDurationRule = iota + 1 // 固定时间 (固定自然时间)
    CouponDurationRelative                               // 相对时间 (自领取日算起)
)

func (c CouponDurationRule) IsValid() bool {
    return slices.Contains(CouponDurations, c)
}

var (
    CouponRules     = []CouponRule{CouponRuleExclusive, CouponRuleStackable}
    CouponDurations = []CouponDurationRule{CouponDurationFixed, CouponDurationRelative}
)

type CouponDuration struct {
    DurationRule CouponDurationRule `json:"durationRule" validate:"required,enum" trans:"有效期规则" enums:"1,2"` // 1:固定时间(固定自然时间) 2:相对时间(自领取日算起)
    DurationTime string             `json:"durationTime" validate:"required_if=DurationRule 1" trans:"过期时间"`  // `durationRule=1`必填 yyyy-mm-dd, 如 2022-09-25
    DurationDays int                `json:"durationDays" validate:"required_if=DurationRule 2" trans:"相对天数"`  // `durationRule=2`必填
}

// ExpiresAt 计算到期时间
// toRider 是否发放到骑手
func (d *CouponDuration) ExpiresAt(toRider bool) *time.Time {
    if d.DurationRule == CouponDurationFixed {
        return silk.Pointer(carbon.Time2Carbon(tools.NewTime().ParseDateStringX(d.DurationTime)).EndOfDay().Carbon2Time())
    }
    if toRider {
        return silk.Pointer(carbon.Now().AddDays(d.DurationDays).EndOfDay().Carbon2Time())
    }
    return nil
}

type CouponTemplate struct {
    *CouponDuration // 有效期规则

    Rule     CouponRule `json:"rule" validate:"required,enum" trans:"使用规则" enums:"1,2"` // 1:互斥 2:可叠加
    Multiple bool       `json:"multiple"`                                                   // 重复使用(该券是否可叠加使用)

    PlanIDs *[]uint64 `json:"planIds,omitempty" validate:"required_without=CityIDs,excluded_with=CityIDs" trans:"绑定骑士卡"` // 和`cities`不能同时为空也不能同时存在
    CityIDs *[]uint64 `json:"cityIds,omitempty" validate:"required_without=PlanIDs,excluded_with=PlanIDs" trans:"可用城市"`   // 和`plans`不能同时为空也不能同时存在
}

type CouponTemplateCreateReq struct {
    Name   string `json:"name" validate:"required" trans:"名称"`
    Remark string `json:"remark" validate:"required,max=10" trans:"备注信息"` // 10字以内, 需要显示在优惠券列表中
    CouponTemplate
}

type CouponTemplateMeta struct {
    Plans  []Plan `json:"plans,omitempty"`  // 绑定骑士卡
    Cities []City `json:"cities,omitempty"` // 可用城市
    CouponTemplate
}

type CouponTemplateListRes struct {
    ID uint64 `json:"id"`
    CouponTemplateMeta
    Name    string `json:"name"`
    Enable  bool   `json:"enable"`
    Total   int    `json:"total"`   // 总数
    InStock int    `json:"inStock"` // 库存
    Used    int    `json:"used"`    // 已使用
    Unused  int    `json:"unused"`  // 未使用
    Expired int    `json:"expired"` // 已过期
    Time    string `json:"time"`    // 更新时间
    Remark  string `json:"remark"`
}

type CouponTemplateListReq struct {
    Enable *bool `json:"enable" query:"enable"` // 是否启用, 默认`true`是
    PaginationReq
}

type CouponTemplateStatusReq struct {
    ID     uint64 `json:"id" validate:"required"`
    Enable bool   `json:"enable"`
}

type CouponTemplateSelection struct {
    Enabled  []SelectOption `json:"enabled"`  // 已启用
    Disabled []SelectOption `json:"disabled"` // 已禁用
}
