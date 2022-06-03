// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-03
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// ExchangeCabinet 电柜换电
type ExchangeCabinet struct {
    Alternative bool                    `json:"alternative"` // 是否使用备选方案
    Info        *RiderCabinetOperating  `json:"info"`        // 换电信息
    Result      *RiderCabinetOperateRes `json:"result"`      // 换电结果
}
