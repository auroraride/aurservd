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

type ExchangeOverview struct {
    Times int `json:"times"` // 换电次数
    Days  int `json:"days"`  // 换电天数
}

type ExchangeLogReq struct {
    PaginationReq
}

type ExchangeLogBinInfo struct {
    EmptyIndex int `json:"emptyIndex"` // 空电仓位
    FullIndex  int `json:"fullIndex"`  // 满电仓位
}

type ExchangeLogRes struct {
    UUID    string             `json:"uuid"`    // 换电编号
    Name    string             `json:"name"`    // 门店或电柜名称
    Type    string             `json:"type"`    // 门店或电柜
    Time    string             `json:"time"`    // 换电时间
    Success bool               `json:"success"` // 是否成功
    City    City               `json:"city"`    // 城市
    BinInfo ExchangeLogBinInfo `json:"binInfo"` // 仓位信息
}
