// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-26
// Based on aurservd by liasica, magicrolan@qq.com.

// 优惠券码最佳使用算法, 贪心算法参考
// https://blog.csdn.net/qq_44112474/article/details/123616038
// https://segmentfault.com/q/1010000015942183/a-1020000015946180

package model

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	CouponShuffle  = []int{9, 0, 1, 4, 3, 5, 6, 8, 7, 2}
	CouponStatuses = []CouponStatus{CouponStatusInStock, CouponStatusUnused, CouponStatusUsed, CouponStatusExpired}
	CouponTargets  = []CouponTarget{CouponTargetRider, CouponTargetStock}
)

type CouponTarget uint8

const (
	CouponTargetRider CouponTarget = iota + 1 // 骑手
	CouponTargetStock                         // 库存
)

func (t CouponTarget) IsValid() bool {
	return slices.Contains(CouponTargets, t)
}

func (t CouponTarget) String() string {
	switch t {
	case CouponTargetRider:
		return "骑手"
	case CouponTargetStock:
		return "库存"
	}
	return " - "
}

func (t CouponTarget) Value() uint8 {
	return uint8(t)
}

type CouponStatus uint8

const (
	CouponStatusInStock CouponStatus = iota // 库存中
	CouponStatusUnused                      // 未使用
	CouponStatusUsed                        // 已使用
	CouponStatusExpired                     // 已过期
)

func (s CouponStatus) IsValid() bool {
	return slices.Contains(CouponStatuses, s)
}

func (s CouponStatus) String() string {
	switch s {
	case CouponStatusInStock:
		return "库存中"
	case CouponStatusUnused:
		return "未使用"
	case CouponStatusUsed:
		return "已使用"
	case CouponStatusExpired:
		return "已过期"
	}
	return " - "
}

type CouponCode string

func (c CouponCode) String() string {
	return string(c.Original())
}

func (c CouponCode) IsValid() bool {
	return len(c.Original()) == 12
}

func (c CouponCode) Original() CouponCode {
	code := strings.ReplaceAll(string(c), "-", "")
	return CouponCode(strings.ToUpper(code))
}

func (c CouponCode) Humanity() string {
	str := string(c)
	return fmt.Sprintf("%s-%s-%s", str[:4], str[4:8], str[8:])
}

type CouponGenerateReq struct {
	TemplateID uint64   `json:"templateId" validate:"required" trans:"优惠券模板"`
	Phones     []string `json:"phones" validate:"required_without=Number,excluded_with=Number" trans:"骑手电话"`
	Number     int      `json:"number" validate:"required_without=Phones,excluded_with=Phones" trans:"生成数量"`
	Amount     float64  `json:"amount" validate:"required" trans:"金额"`
	Remark     string   `json:"remark" validate:"required" trans:"备注"`
}

type CouponListFilter struct {
	RiderID uint64        `json:"riderId" query:"riderId"`
	Target  *CouponTarget `json:"target" query:"target" enums:"1,2"`     // 发送对象, 不携带参数为全部, 1:骑手 2:库存
	Status  *CouponStatus `json:"status" query:"status" enums:"0,1,2,3"` // 状态, 不携带参数为全部, 0:库存中 1:未使用 2:已使用 3:已过期
	Keyword string        `json:"keyword" query:"keyword"`               // 骑手查询关键词
	Code    CouponCode    `json:"code" query:"code"`                     // 券码
}

type CouponListReq struct {
	PaginationReq
	CouponListFilter
}

type CouponListRes struct {
	ID         uint64       `json:"id"`
	Amount     float64      `json:"amount"`               // 金额
	Code       string       `json:"code"`                 // 券码
	Name       string       `json:"name"`                 // 名称
	Rider      string       `json:"rider,omitempty"`      // 骑手
	Phone      string       `json:"phone,omitempty"`      // 骑手电话
	Status     CouponStatus `json:"status"`               // 状态, 0:库存中 1:未使用 2:已使用 3:已过期
	Creator    string       `json:"creator"`              // 创建者
	Cities     []string     `json:"cities,omitempty"`     // 可用城市
	Plans      []string     `json:"plans,omitempty"`      // 可用骑士卡
	Time       string       `json:"time"`                 // 生成时间
	UsedAt     string       `json:"usedAt,omitempty"`     // 使用时间
	ExpiredAt  string       `json:"expiredAt,omitempty"`  // 过期时间
	TradeNo    string       `json:"outTradeNo,omitempty"` // 使用订单编号
	Plan       string       `json:"plan,omitempty"`       // 使用骑士卡
	TemplateID uint64       `json:"templateId"`           // 模板ID
	AssemblyID uint64       `json:"assemblyId"`           // 发券记录ID
}

type CouponAllocateReq struct {
	ID      uint64 `json:"id" validate:"required"`      // 优惠券
	RiderID uint64 `json:"riderId" validate:"required"` // 骑手
}

type Coupon struct {
	ID        uint64  `json:"id"`
	Cate      string  `json:"cate"`             // 类型标识
	Useable   bool    `json:"useable"`          // 是否可使用
	Amount    float64 `json:"amount"`           // 金额
	Name      string  `json:"name"`             // 名称
	ExpiredAt string  `json:"expiredAt"`        // 过期时间
	Code      string  `json:"code"`             // 券码
	Exclusive bool    `json:"exclusive"`        // 与其他类型券是否互斥
	Plans     []*Plan `json:"plans,omitempty"`  // 可用骑士卡, 不存在此字段则不限制
	Cities    []City  `json:"cities,omitempty"` // 可用城市, 不存在此字段则不限制
}

type CouponRiderListReq struct {
	Type uint8 `json:"type" enums:"0,1,2" query:"type"` // 查询类别 0:可使用 1:已使用 2:已过期
}

type CouponRider struct {
	Name      string  `json:"name"`             // 名称
	ExpiredAt string  `json:"expiredAt"`        // 到期时间
	Amount    float64 `json:"amount"`           // 金额
	Code      string  `json:"code"`             // 券码
	UsedAt    string  `json:"usedAt,omitempty"` // 使用时间
}
