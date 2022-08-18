// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    CommissionStatusPending uint8 = iota // 未发放
    CommissionStatusIssued               // 已发放
)

type ComissionExportItem struct {
    Business string  `json:"business"` // 业务
    Rider    string  `json:"rider"`    // 骑手
    Amount   float64 `json:"amount"`   // 提成金额
    Plan     string  `json:"plan"`     // 骑士卡
    Cost     float64 `json:"cost"`     // 订单总额
    Time     string  `json:"time"`     // 时间
}
