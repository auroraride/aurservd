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

type ExchangeStoreReq struct {
    Code string `json:"code" validate:"required,startswith=STORE:"` // 二维码
}

type ExchangeStoreRes struct {
    Voltage   float64 `json:"voltage"`   // 电池电压
    StoreName string  `json:"storeName"` // 门店名称
    Time      int64   `json:"time"`      // 时间戳
    UUID      string  `json:"uuid"`      // 编码
}
