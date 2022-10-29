// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    SettingDeposit                = "DEPOSIT"                  // 押金
    SettingRenewal                = "RENEWAL"                  // 退订多久后重签计算佣金
    SettingCabinetFault           = "CABINET_FAULT"            // 电柜故障
    SettingRescueReason           = "RESCUE_REASON"            // 救援原因
    SettingRescueFee              = "RESCUE_FEE"               // 救援费用
    SettingReminder               = "REMINDER"                 // 催费通知
    SettingBatteryFull            = "BATTERY_FULL"             // 满电电量
    SettingException              = "EXCEPTION"                // 物资异常
    SettingPauseMaxDays           = "PAUSE_MAX_DAYS"           // 最大寄存时间
    SettingExchangeInterval       = "EXCHANGE_INTERVAL"        // 限制换电间隔
    SettingMaintain               = "MAINTAIN"                 // 维护中
    SettingReserveDuration        = "RESERVE_MAX_DURATION"     // 最长预约时间
    SettingExchangeMinBattery     = "EXCHANGE_MIN_BATTERY"     // 换电最低电量
    SettingPlanBatteryDescription = "PLAN_BATTERY_DESCRIPTION" // 单电订阅介绍
    SettingPlanEbikeDescription   = "PLAN_EBIKE_DESCRIPTION"   // 车电订阅介绍
    SettingQuestions              = "QUESTION"                 // 常见问题
    SettingAppVersion             = "APP_VERSION"              // App版本
    SettingConsumePoints          = "CONSUME_POINTS"           // 消费赠送积分
)

type SettingValueConvert func(content string) any

type SettingItem struct {
    Desc    string // 描述
    Default any    // 默认值
}

type SettingReminderNotice struct {
    App []int `json:"app"` // APP提醒
    Sms []int `json:"sms"` // 短信通知
    Vms []int `json:"vms"` // 电话通知
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
    SettingException: {
        Desc:    "物资异常",
        Default: []string{"丢失", "故障"},
    },
    SettingRescueReason: {
        Desc:    "救援原因",
        Default: []string{},
    },
    SettingRescueFee: {
        Desc:    "救援费用(元/公里)",
        Default: "0",
    },
    SettingReminder: {
        Desc: "催费通知",
        Default: SettingReminderNotice{
            App: []int{5},
            Sms: []int{3},
            Vms: []int{1, -3},
        },
    },
    SettingBatteryFull: {
        Desc:    "满电电量百分比",
        Default: "80",
    },
    SettingPauseMaxDays: {
        Desc:    "最大寄存时间",
        Default: "31",
    },
    SettingExchangeInterval: {
        Desc:    "限制换电间隔",
        Default: "20",
    },
    SettingMaintain: {
        Desc:    "是否维护中",
        Default: false,
    },
    SettingReserveDuration: {
        Desc:    "最长预约时间",
        Default: "60",
    },
    SettingExchangeMinBattery: {
        Desc:    "换电最低电量(%)",
        Default: "50",
    },
    SettingPlanBatteryDescription: {
        Desc:    "单电订阅介绍",
        Default: SettingPlanDescription{},
    },
    SettingPlanEbikeDescription: {
        Desc:    "车电订阅介绍",
        Default: SettingPlanDescription{},
    },
    SettingQuestions: {
        Desc:    "常见问题",
        Default: []SettingQuestion{},
    },
    SettingAppVersion: {
        Desc:    "App版本",
        Default: map[string]SettingAppVersionValue{},
    },
    SettingConsumePoints: {
        Desc:    "消费赠送积分",
        Default: map[uint64]float64{},
    },
}

type SettingRiderApp struct {
    AssistanceFee   float64 `json:"assistanceFee"`   // 救援费用
    ReserveDuration int     `json:"reserveDuration"` // 预约最长时间(分钟)
}

type SettingPlanDescription struct {
    Banner    string `json:"banner"`    // banner图
    Product   string `json:"product"`   // 商品介绍
    Pickup    string `json:"pickup"`    // 提货方式
    Attention string `json:"attention"` // 注意事项
}

type SettingQuestion struct {
    Question string `json:"question"` // 问题
    Answer   string `json:"answer"`   // 解答
}

type SettingAppVersionValue struct {
    Version     string `json:"version"`
    Description string `json:"description"`
    Link        string `json:"link"`
    Enable      bool   `json:"enable"`
}
