// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    jsoniter "github.com/json-iterator/go"
    "time"
)

const (
    OrderTypeNew           uint8 = iota + 1 // 骑士卡
    OrderTypeRenewal                        // 续签
    OrderTypeReUse                          // 重签
    OrderTypeChangeBattery                  // 更换电池
    OrderTypeRescue                         // 救援
    OrderTypeFee                            // 滞纳金
)

const (
    OrderPaywayAlipay uint8 = iota + 1 // 支付宝支付
    OrderPaywayWechat                  // 微信支付
)

const (
    OrderStatusPending       uint8 = iota // 未支付
    OrderStatusPaid                       // 已支付
    OrderStatusRefundPending              // 申请退款
    OrderStatusRefundSuccess              // 已退款
)

// OrderCreateReq 订单创建请求
type OrderCreateReq struct {
    PlanID    uint64 `json:"planId" validate:"required" trans:"套餐ID"`
    Payway    uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"1,2"`            // 1支付宝 2微信
    OrderType uint8  `json:"orderType" validate:"required" trans:"订单类型" enums:"1,2,3,4,5,6"` // 1新签 2续签 3重签 4更改电池 5救援 6滞纳金
}

// OrderCreateRes 订单创建返回
type OrderCreateRes struct {
    Prepay string `json:"prepay"` // 预支付字符串
}

// OrderCache 订单缓存
type OrderCache struct {
    OrderType  uint8     `json:"orderType"` // 订单类型
    OutTradeNo string    `json:"outTradeNo"`
    RiderID    uint64    `json:"riderId"`
    Name       string    `json:"name"`
    Amount     float64   `json:"amount"`
    Payway     uint8     `json:"payway"`
    Expire     time.Time `json:"expire"`            // 过期时间
    TradeNo    string    `json:"tradeNo,omitempty"` // 平台单号
    Plan       *PlanItem `json:"plan,omitempty"`    // 骑士卡
}

func (oc *OrderCache) MarshalBinary() ([]byte, error) {
    return jsoniter.Marshal(oc)
}

func (oc *OrderCache) UnmarshalBinary(data []byte) error {
    return jsoniter.Unmarshal(data, oc)
}

// OrderSubordinate 从属订单信息
type OrderSubordinate struct {
}
