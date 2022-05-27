// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    RiderPlanStatusPending uint8 = iota // 未激活
    RiderPlanStatusNormal               // 计费中
)

type RiderPlan struct {
    PlanID   uint64  `json:"planId"`             // 骑士卡ID
    Name     string  `json:"name"`               // 骑士卡名称
    LastDays uint    `json:"lastDays"`           // 剩余天数
    Status   uint8   `json:"status" enums:"0,1"` // 状态 0未激活 1计费中
    Start    string  `json:"start"`              // 开始时间
    End      string  `json:"end"`                // 预计结束时间
    Voltage  float64 `json:"voltage"`            // 可用电压型号
}
