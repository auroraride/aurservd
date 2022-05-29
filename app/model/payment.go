// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    jsoniter "github.com/json-iterator/go"
    "time"
)

const (
    PaymentCacheTypePlan   uint = iota + 1 // 购买骑士卡订单
    PaymentCacheTypeRefund                 // 退款订单
)

// PaymentCache 支付缓存
type PaymentCache struct {
    CacheType uint           `json:"cacheType"`        // 订单类型
    Plan      *PaymentPlan   `json:"create,omitempty"` // 购买骑士卡订单
    Refund    *PaymentRefund `json:"refund,omitempty"` // 退款订单
}

// PaymentPlan 购买骑士卡订单
type PaymentPlan struct {
    CityID     uint64    `json:"cityID"`            // 城市ID
    OrderType  uint      `json:"orderType"`         // 订单类型
    OutTradeNo string    `json:"outTradeNo"`        // 订单号
    RiderID    uint64    `json:"riderId"`           // 骑手ID
    Name       string    `json:"name"`              // 订单名称
    Amount     float64   `json:"amount"`            // 总金额 = 套餐 + 押金
    Payway     uint8     `json:"payway"`            // 支付方式
    Expire     time.Time `json:"expire"`            // 过期时间
    TradeNo    string    `json:"tradeNo,omitempty"` // 平台单号
    Plan       *PlanItem `json:"plan,omitempty"`    // 骑士卡
    Deposit    float64   `json:"deposit"`           // 附带押金
}

func (oc *PaymentCache) MarshalBinary() ([]byte, error) {
    return jsoniter.Marshal(oc)
}

func (oc *PaymentCache) UnmarshalBinary(data []byte) error {
    return jsoniter.Unmarshal(data, oc)
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
