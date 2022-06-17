// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type Inventory struct {
    Name     string `json:"name" validate:"required" trans:"物资名称"`
    Count    bool   `json:"count" trans:"是否需要盘点"`
    Transfer bool   `json:"transfer" trans:"是否可调拨"`
    Purchase bool   `json:"purchase" trans:"是否可采购"`
}

type InventoryDelete struct {
    Name *string `json:"name" validate:"required" trans:"物资名称"` // *POST参数*
}

type InventoryListReq struct {
    Count    bool `json:"count"`
    Transfer bool `json:"transfer"`
    Purchase bool `json:"purchase"`
}

type InventoryItem struct {
    Name    string   `json:"name"`              // 物资名称
    Battery bool     `json:"battery"`           // 是否电池
    Voltage *float64 `json:"voltage,omitempty"` // 电压型号
}
