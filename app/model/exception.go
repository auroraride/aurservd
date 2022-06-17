// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-17
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type ExceptionEmployeeSetting struct {
    Items   []InventoryItem `json:"items"`   // 物资列表
    Reasons []interface{}   `json:"reasons"` // 异常项
}
