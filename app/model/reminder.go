// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-20
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type ReminderFilter struct {
    Keyword string `json:"keyword" query:"keyword"`           // 骑手关键词 (姓名或手机号)
    Start   string `json:"start" query:"start"`               // 开始时间
    End     string `json:"end" query:"end"`                   // 结束时间
    Days    *int   `json:"days" query:"days"`                 // 剩余天数 (不携带为不筛选)
    Type    string `json:"type" query:"type" enums:"vms,sms"` // 催费方式, vms:语音 sms:短信
    RiderID uint64 `json:"riderId" query:"riderId"`           // 骑手ID
}

type ReminderListReq struct {
    PaginationReq
    ReminderFilter
}

type ReminderListRes struct {
    Phone      string  `json:"phone"`      // 电话
    Name       string  `json:"name"`       // 姓名
    Success    bool    `json:"success"`    // 是否成功
    Time       string  `json:"time"`       // 发送时间
    PlanName   string  `json:"planName"`   // 骑士卡
    Fee        float64 `json:"fee"`        // 逾期费用
    FeeFormula string  `json:"feeFormula"` // 费用计算公式
}
