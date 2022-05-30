// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package model

//
// const (
//     RiderOrderStatusPending  uint8 = iota // 未激活
//     RiderOrderStatusNormal                // 计费中
//     RiderOrderStatusPaused                // 暂停中
//     RiderOrderStatusOverdue               // 已逾期
//     RiderOrderStatusRemanded              // 已过期, 已归还电池
//     RiderOrderStatusRefunded              // 已退款
// )
// 
// type OrderPlan struct {
//     ID   uint64 `json:"id"`   // 骑士卡ID
//     Name string `json:"name"` // 骑士卡名称
// }
//
// type RiderOrderDays struct {
//     Days       int `json:"days"`       // 总天数
//     Remaining  int `json:"remaining"`  // 剩余天数
//     PausedDays int `json:"pausedDays"` // 暂停天数
//     AlterDays  int `json:"alterDays"`  // 改动天数
// }
//
// // RiderRecentOrder 骑手最近骑士卡
// type RiderRecentOrder struct {
//     ID      uint64         `json:"id"` // 订单ID
//     Days    RiderOrderDays `json:"days"`
//     Status  uint8          `json:"status" enums:"0,1,2,3,4"` // 状态 0未激活 1计费中 2暂停中 3已逾期 4已归还(已过期)
//     StartAt string         `json:"startAt"`                  // 开始时间
//     EndAt   string         `json:"endAt"`                    // 结束时间 / 预计套餐结束时间
//     Voltage float64        `json:"voltage"`                  // 可用电压型号
//
//     PayAt   string  `json:"payAt"`   // 支付时间
//     Payway  uint8   `json:"payway"`  // 支付方式
//     Amount  float64 `json:"amount"`  // 骑士卡金额
//     Deposit float64 `json:"deposit"` // 押金(只在未启用骑士卡中显示), 若押金为0则押金一行不显示
//     Total   float64 `json:"total"`   // 总金额, 总金额为 amount + deposit
//
//     City   City           `json:"city"`   // 所属城市
//     Models []BatteryModel `json:"models"` // 可用电池型号, 显示为`72V30AH`即Voltage(V)+Capacity(AH), 逗号分隔
//     Plan   OrderPlan      `json:"plan"`   // 骑士卡信息
// }
