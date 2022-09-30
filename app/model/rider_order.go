// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// RiderOrder 骑手订单
type RiderOrder struct {
    ID         uint64    `json:"id"`                 // 订单ID
    Type       uint      `json:"type"`               // 订单类型 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金
    Status     uint8     `json:"status"`             // 订单状态 0未支付 1已支付 2申请退款 3已退款 4退款被拒绝
    Payway     uint8     `json:"payway"`             // 支付方式 1支付宝 2微信
    PayAt      string    `json:"payAt"`              // 支付时间
    Amount     float64   `json:"amount"`             // 支付金额
    OutTradeNo string    `json:"outTradeNo"`         // 订单编号
    TradeNo    string    `json:"tradeNo"`            // 订单编号 (支付平台)
    City       City      `json:"city"`               // 城市
    Rider      Rider     `json:"rider"`              // 骑手
    Plan       *Plan     `json:"plan,omitempty"`     // 骑士卡, 非骑士卡订阅订单无此字段 (可为空)
    Model      string    `json:"model,omitempty"`    // 电池型号 (可为空)
    Store      *Store    `json:"store,omitempty"`    // 门店 (可为空)
    Employee   *Employee `json:"employee,omitempty"` // 店员 (可为空)
    Refund     *Refund   `json:"refund,omitempty"`   // 退款详情 (可为空)
}
