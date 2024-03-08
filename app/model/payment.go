// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	PaymentCacheTypePlan              uint = iota + 1 // 购买骑士卡订单
	PaymentCacheTypeRefund                            // 退款订单
	PaymentCacheTypeOverdueFee                        // 欠费订单
	PaymentCacheTypeAssistance                        // 救援订单
	PaymentCacheTypeAgentPrepay                       // 代理商预储值
	PaymentCacheTypeAliDepositFree                    // 芝麻信用免押
	PaymentCacheTypeWechatDepositFree                 // 微信免押
)

// PaymentSubscribe 购买骑士卡订单
type PaymentSubscribe struct {
	CityID      uint64  `json:"cityId"`            // 逾期时城市ID
	OrderType   uint    `json:"orderType"`         // 订单类型
	OutTradeNo  string  `json:"outTradeNo"`        // 订单号
	RiderID     uint64  `json:"riderId"`           // 骑手ID
	Name        string  `json:"name"`              // 订单名称
	Amount      float64 `json:"amount"`            // 总金额 = 套餐 + 押金
	Payway      uint8   `json:"payway"`            // 支付方式
	TradeNo     string  `json:"tradeNo,omitempty"` // 平台单号
	Plan        *Plan   `json:"plan"`              // 骑士卡
	Deposit     float64 `json:"deposit"`           // 附带押金
	PastDays    *int    `json:"pastDays"`          // 距离上次退订天数
	Commission  float64 `json:"commission"`        // 提成金额
	Model       string  `json:"model"`             // 可用电池型号
	Days        uint    `json:"days"`              // 骑士卡天数
	OrderID     *uint64 `json:"orderId"`           // 子订单ID
	SubscribeID *uint64 `json:"subscribeId"`       // 续费订单携带订阅ID

	Points        int64    `json:"points"`                 // 使用积分
	PointRatio    float64  `json:"pointRatio"`             // 积分兑换比例
	CouponAmount  float64  `json:"couponAmount"`           // 优惠券金额
	Coupons       []uint64 `json:"coupons"`                // 使用优惠券列表
	DiscountNewly float64  `json:"discountNewly"`          // 新签优惠
	EbikeBrandID  *uint64  `json:"ebikeBrandId,omitempty"` // 电车型号
}

// PaymentRefund 退款详情
type PaymentRefund struct {
	OrderID      uint64  `json:"orderId"`      // 原始订单ID
	TradeNo      string  `json:"tradeNo"`      // 原订单平台ID
	Total        float64 `json:"total"`        // 原支付订单总支付金额
	RefundAmount float64 `json:"refundAmount"` // 退款金额
	Reason       string  `json:"reason"`       // 退款缘由
	OutRefundNo  string  `json:"outRefundNo"`  // 退款单号

	// 退款申请同步结果
	Request bool      `json:"request"` // 是否申请成功
	Success bool      `json:"success"` // 是否退款成功
	Time    time.Time `json:"time"`    // 退款时间
}

type PaymentOverdueFee struct {
	Subject string `json:"subject"` // 支付信息

	OutTradeNo string  `json:"outTradeNo"` // 订单号
	OrderType  uint    `json:"orderType"`  // 订单类型
	Days       int     `json:"days"`       // 逾期天数
	Amount     float64 `json:"amount"`     // 逾期费用
	Payway     uint8   `json:"payway"`     // 支付方式

	RiderID     uint64 `json:"riderId"`     // 骑手ID
	PlanID      uint64 `json:"planId"`      // 逾期时套餐ID
	OrderID     uint64 `json:"orderId"`     // 逾期时订单ID
	SubscribeID uint64 `json:"subscribeId"` // 逾期时订阅ID
	CityID      uint64 `json:"cityId"`      // 逾期时城市ID

	TradeNo string `json:"tradeNo,omitempty"` // 平台单号
}

type PaymentAssistance struct {
	ID         uint64  `json:"id"`                // 救援ID
	Payway     uint8   `json:"payway"`            // 支付方式
	Subject    string  `json:"subject"`           // 支付信息
	Cost       float64 `json:"cost"`              // 费用
	OutTradeNo string  `json:"outTradeNo"`        // 订单号
	TradeNo    string  `json:"tradeNo,omitempty"` // 支付单号
}

// PaymentAgentPrepay 代理商预储值
type PaymentAgentPrepay struct {
	*AgentPrepay
	Payway     Payway `json:"payway"`            // 支付方式
	OutTradeNo string `json:"outTradeNo"`        // 订单号
	TradeNo    string `json:"tradeNo,omitempty"` // 支付单号
	Attach     string `json:"attach,omitempty"`  // 订单备注
}

// DepositFree 押金免押订单
type DepositFree struct {
	PlanID     uint64  `json:"planId"` // 套餐ID
	Payway     uint8   `json:"payway"`
	OutTradeNo string  `json:"outTradeNo"`
	TradeNo    string  `json:"tradeNo,omitempty"` // 支付单号
	Amount     float64 `json:"amount"`            // 押金金额
	Plan       *Plan   `json:"plan"`              // 骑士卡
	RiderID    uint64  `json:"riderId"`           // 骑手ID
}

// PaymentCache 支付缓存
type PaymentCache struct {
	CacheType   uint                `json:"cacheType"`             // 订单类型
	Subscribe   *PaymentSubscribe   `json:"create,omitempty"`      // 购买骑士卡订单
	Refund      *PaymentRefund      `json:"refund,omitempty"`      // 退款订单
	OverDueFee  *PaymentOverdueFee  `json:"overDueFee,omitempty"`  // 逾期费用订单
	Assistance  *PaymentAssistance  `json:"assistance,omitempty"`  // 救援订单
	AgentPrepay *PaymentAgentPrepay `json:"agentPrepay,omitempty"` // 代理商预充值
	DepositFree *DepositFree        `json:"depositFree,omitempty"` // 押金免押订单
}

func (pc *PaymentCache) MarshalBinary() ([]byte, error) {
	return jsoniter.Marshal(pc)
}

func (pc *PaymentCache) UnmarshalBinary(data []byte) error {
	return jsoniter.Unmarshal(data, pc)
}

func (pc *PaymentCache) GetPaymentArgs() (amount float64, desc string, outTradeNo string, attach string) {
	switch pc.CacheType {
	case PaymentCacheTypePlan:
		return pc.Subscribe.Amount, pc.Subscribe.Name, pc.Subscribe.OutTradeNo, ""
	case PaymentCacheTypeOverdueFee:
		return pc.OverDueFee.Amount, pc.OverDueFee.Subject, pc.OverDueFee.OutTradeNo, ""
	case PaymentCacheTypeAssistance:
		return pc.Assistance.Cost, pc.Assistance.Subject, pc.Assistance.OutTradeNo, ""
	case PaymentCacheTypeAgentPrepay:
		return pc.AgentPrepay.Amount, "代理商自主储值", pc.AgentPrepay.OutTradeNo, pc.AgentPrepay.Attach
	}
	return 0, "", "", ""
}

type PaymentConfigure struct {
	Points      int64    `json:"points"`      // 可用积分
	Ratio       float64  `json:"ratio"`       // 兑换比例(实际抵扣金额 = 兑换比例 × 积分数量)
	Coupons     []Coupon `json:"coupons"`     // 用户优惠券
	MaxDiscount float64  `json:"maxDiscount"` // 最大优惠金额
}
