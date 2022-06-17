// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-17
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    ExceptionStatusPending uint8 = iota // 异常未处理
    ExceptionStatusDone                 // 异常已处理
)

type ExceptionEmployeeSetting struct {
    Items   []InventoryItem `json:"items"`   // 物资列表
    Reasons []interface{}   `json:"reasons"` // 异常项
}

type ExceptionEmployeeReq struct {
    Name    *string  `json:"name"`    // 物资名称 (非电池物资, 和`voltage`不能同时存在, 也不能同时为空), 物资列表查看接口 `E3001 物资异常配置`
    Voltage *float64 `json:"voltage"` // 电压型号 (电池, 和`name`不能同时存在, 也不能同时为空)

    Reason      string `json:"reason" validate:"required" trans:"异常原因"` // 异常原因查看接口 `E3001 物资异常配置`
    Description string `json:"description" validate:"required" trans:"描述"`
    Num         int    `json:"num" validate:"required,gte=1" trans:"异常数量"`

    Attachments []string `json:"attachments" validate:"max=3"` // 附件, *注意, 门店端的附件需要以employee/开头*
}
