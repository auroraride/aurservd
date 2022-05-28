// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    RiderOrderStatusPending  uint8 = iota // 未激活
    RiderOrderStatusNormal                // 计费中
    RiderOrderStatusPaused                // 暂停中
    RiderOrderStatusOverdue               // 已逾期
    RiderOrderStatusRemanded              // 已过期, 已归还电池
)

// RiderRecentOrder 骑手最近骑士卡
type RiderRecentOrder struct {
    PlanID     uint64  `json:"planId"`                   // 骑士卡ID
    PlanName   string  `json:"planName"`                 // 骑士卡名称
    Days       int     `json:"days"`                     // 总天数
    Remaining  int     `json:"remaining"`                // 剩余天数
    PausedDays int     `json:"pausedDays"`               // 暂停天数
    AlterDays  int     `json:"alterDays"`                // 改动天数
    Status     uint8   `json:"status" enums:"0,1,2,3,4"` // 状态 0未激活 1计费中 2暂停中 3已逾期 4已归还(已过期)
    StartAt    string  `json:"startAt"`                  // 开始时间
    EndAt      string  `json:"endAt"`                    // 结束时间
    Voltage    float64 `json:"voltage"`                  // 可用电压型号
}
