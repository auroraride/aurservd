// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type StockTransferReq struct {
    Voltage    float64 `json:"voltage,omitempty"` // 电池型号 (和`物资名称`不能同时存在, 也不能同时为空)
    Name       string  `json:"name,omitempty"`    // 物资名称 (和`电池型号`不能同时存在, 也不能同时为空)
    OutboundID uint64  `json:"outboundId"`        // 调出自 (0:平台)
    InboundID  uint64  `json:"inboundId"`         // 调入至 (0:平台)
    Num        int     `json:"num"`               // 调拨数量
}

type StockListReq struct {
    PaginationReq

    Name   *string `json:"name" query:"name"`     // 门店名称
    CityID *uint64 `json:"cityId" query:"cityId"` // 城市ID
    Start  *string `json:"start" query:"start"`   // 开始时间
    End    *string `json:"end" query:"end"`       // 结束时间
}

type StockMaterial struct {
    Name     string `json:"name"`     // 物资名称
    Outbound int    `json:"outbound"` // 出库数量
    Inbound  int    `json:"inbound"`  // 入库数量
    Surplus  int    `json:"surplus"`  // 剩余
}

type StockListRes struct {
    Store        Store            `json:"store"`        // 门店
    City         City             `json:"city"`         // 城市
    BatteryTotal int              `json:"batteryTotal"` // 电池总数
    Batteries    []*StockMaterial `json:"batteries"`    // 电池详情
    Materials    []*StockMaterial `json:"materials"`    // 非电池物资详情
}
