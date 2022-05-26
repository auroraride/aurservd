// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    SettingDeposit      = "DEPOSIT"       // 押金
    SettingRenewal      = "RENEWAL"       // 重签判定时间
    SettingCabinetFault = "CABINET_FAULT" // 电柜故障
    SettingRescueReason = "RESCUE_REASON" // 救援原因
    SettingRescueFee    = "RESCUE_FEE"    // 救援费用
    SettingOverdue      = "OVERDUE"       // 逾期通知
    SettingBatteryFull  = "BATTERY_FULL"  // 满电电量
)

type SettingValueConvert func(content string) any

type SettingItem struct {
    Desc  string // 描述
    Value string // 值

    Default any // 默认值
}

type SettingOverdueNotice struct {
    App  string `json:"app"`  // APP提醒
    Sms  string `json:"sms"`  // 短信通知
    Call string `json:"call"` // 电话通知
}

type SettingReq struct {
    Key     *string `json:"key" validate:"required" param:"key" trans:"设置项"`
    Content *string `json:"content" validate:"required" trans:"值"`
}

type SettingRes struct {
    Key     string `json:"key"`     // 设置项
    Content string `json:"content"` // 设置值
    Desc    string `json:"desc"`    // 描述
}

var Settings = map[string]SettingItem{
    SettingDeposit: {
        Desc:    "平台押金",
        Default: "99",
    },
    SettingRenewal: {
        Desc:    "重签判定时间",
        Default: "7",
    },
    SettingCabinetFault: {
        Desc:    "电柜故障",
        Default: []string{},
    },
    SettingRescueReason: {
        Desc:    "救援原因",
        Default: []string{},
    },
    SettingRescueFee: {
        Desc:    "救援费用(元/公里)",
        Default: "0",
    },
    SettingOverdue: {
        Desc: "逾期通知",
        Default: SettingOverdueNotice{
            App:  "3",
            Sms:  "1",
            Call: "0",
        },
    },
    SettingBatteryFull: {
        Desc:    "满电电量百分比",
        Default: "80",
    },
}
