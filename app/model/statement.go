// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-06
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type StatementBillReq struct {
    End string `json:"end" validate:"required" trans:"账单截止日期"`
    ID  uint64 `json:"id" validate:"required" trans:"企业ID"`
}

type StatementBillRes struct {
    Voltage     float64 `json:"voltage"`     // 电压型号
    RiderNumber int     `json:"riderNumber"` // 使用骑手数量
    Price       float64 `json:"price"`       // 单价
    Days        int     `json:"days"`        // 天数
    Cost        float64 `json:"cost"`        // 账单金额
}
