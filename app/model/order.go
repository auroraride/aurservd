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
    OrderTypeNewPlan       uint = iota + 1 // 新签
    OrderTypeRePlan                        // 续签
    OrderTypeRenewal                       // 重签
    OrderTypeChangeBattery                 // 更换电池
    OrderTypeRescue                        // 救援
    OrderTypeFee                           // 滞纳金
    OrderTypeDeposit                       // 押金
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
    CityID    uint64 `json:"cityId" validate:"required" trans:"城市ID"`
    PlanID    uint64 `json:"planId" validate:"required" trans:"套餐ID"`
    Payway    uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"1,2"`              // 1支付宝 2微信
    OrderType uint   `json:"orderType" validate:"required" trans:"订单类型" enums:"1,2,3,4,5,6,7"` // 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金
}

// OrderCreateRes 订单创建返回
type OrderCreateRes struct {
    Prepay     string `json:"prepay"`     // 预支付字符串
    OutTradeNo string `json:"outTradeNo"` // 交易编码
}

// OrderCache 订单缓存
type OrderCache struct {
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

func (oc *OrderCache) MarshalBinary() ([]byte, error) {
    return jsoniter.Marshal(oc)
}

func (oc *OrderCache) UnmarshalBinary(data []byte) error {
    return jsoniter.Unmarshal(data, oc)
}

// OrderSubordinate 从属订单信息
type OrderSubordinate struct {
}

// OrderNotActived 骑手未激活订单信息
type OrderNotActived struct {
    ID      uint64         `json:"id"`                 // 订单编号
    Amount  float64        `json:"amount"`             // 骑士卡金额
    Deposit float64        `json:"deposit"`            // 押金, 若押金为0则押金一行不显示
    Total   float64        `json:"total"`              // 总金额, 总金额为 amount + deposit
    Payway  uint8          `json:"payway" enums:"1,2"` // 支付方式 1支付宝 2微信
    Plan    PlanItem       `json:"plan"`               // 骑行卡详情
    City    City           `json:"city"`               // 所属城市
    Models  []BatteryModel `json:"models"`             // 可用电池型号, 显示为`72V30AH`即Voltage(V)+Capacity(AH), 逗号分隔
    Time    string         `json:"time"`               // 支付时间
}
