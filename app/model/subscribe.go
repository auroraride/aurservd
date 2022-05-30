// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-30
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    SubscribeStatusInactive     uint8 = iota // 未激活
    SubscribeStatusUsing                     // 计费中
    SubscribeStatusPaused                    // 寄存中
    SubscribeStatusOverdue                   // 已逾期
    SubscribeStatusUnSubscribed              // 已退订, 已过期, 已归还电池
    SubscribeStatusCanceled                  // 已取消, 已退款
)

var (
    // OrderSubscribeTypes 骑手骑士卡订单类型
    OrderSubscribeTypes = []uint{OrderTypeNewly, OrderTypeAgain, OrderTypeRenewal, OrderTypeTransform}
)

// SubscribeOrderInfo 订阅订单信息
type SubscribeOrderInfo struct {
    ID      uint64  `json:"id"`      // 订单ID
    Status  uint    `json:"status"`  // 订单状态 0未支付 1已支付 2申请退款 3已退款 4退款被拒绝
    PayAt   string  `json:"payAt"`   // 支付时间
    Payway  uint8   `json:"payway"`  // 支付方式
    Amount  float64 `json:"amount"`  // 骑士卡金额
    Deposit float64 `json:"deposit"` // 押金(只在未启用骑士卡中显示), 若押金为0则押金一行不显示
    Total   float64 `json:"total"`   // 总金额, 总金额为 amount + deposit
}

type Subscribe struct {
    ID        uint64  `json:"id"`                       // 订阅ID
    Status    uint8   `json:"status" enums:"0,1,2,3,4"` // 状态 0未激活 1计费中 2暂停中 3已逾期 4已归还(已过期)
    Voltage   float64 `json:"voltage"`                  // 可用电压型号
    Days      int     `json:"days"`                     // 总天数
    Remaining int     `json:"remaining"`                // 剩余天数
    PauseDays int     `json:"pauseDays"`                // 暂停天数
    AlterDays int     `json:"alterDays"`                // 改动天数
    StartAt   string  `json:"startAt"`                  // 开始时间
    EndAt     string  `json:"endAt"`                    // 结束时间 / 预计套餐结束时间

    City   *City               `json:"city,omitempty"`   // 所属城市
    Models []BatteryModel      `json:"models,omitempty"` // 可用电池型号, 显示为`72V30AH`即Voltage(V)+Capacity(AH), 逗号分隔
    Plan   *Plan               `json:"plan,omitempty"`   // 骑士卡信息
    Order  *SubscribeOrderInfo `json:"order,omitempty"`  // 订单信息
}

// UnsubscribedDaysToNewly 退订多少天后视为有效重签并计算佣金
func UnsubscribedDaysToNewly() int {
    return 15
}
