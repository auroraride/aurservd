// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
	CabinetBusinessScanExpires int64   = 30  // 电柜业务 - 扫码有效期(s)
	CabinetBusinessStepTimeout int64   = 200 // 电柜业务 - 操作超时(s)
	DailyRentDefault           float64 = 99999.0
)

const (
	SettingDepositKey                = "DEPOSIT"                  // 押金
	SettingRenewalKey                = "RENEWAL"                  // 退订多久后重签计算佣金
	SettingCabinetFaultKey           = "CABINET_FAULT"            // 电柜故障
	SettingRescueReasonKey           = "RESCUE_REASON"            // 救援原因
	SettingRescueFeeKey              = "RESCUE_FEE"               // 救援费用
	SettingReminderKey               = "REMINDER"                 // 催费通知
	SettingBatteryFullKey            = "BATTERY_FULL"             // 电池满电电量判定
	SettingExceptionKey              = "EXCEPTION"                // 物资异常
	SettingPauseMaxDaysKey           = "PAUSE_MAX_DAYS"           // 最大寄存时间
	SettingExchangeIntervalKey       = "EXCHANGE_INTERVAL"        // 限制换电间隔
	SettingExchangeLimitKey          = "EXCHANGE_LIMIT"           // 换电限制
	SettingExchangeFrequencyKey      = "EXCHANGE_FREQUENCY"       // 换电频次
	SettingMaintainKey               = "MAINTAIN"                 // 维护中
	SettingReserveDurationKey        = "RESERVE_MAX_DURATION"     // 最长预约时间
	SettingExchangeMinBatteryKey     = "EXCHANGE_MIN_BATTERY"     // 换电最低电量
	SettingPlanBatteryDescriptionKey = "PLAN_BATTERY_DESCRIPTION" // 单电订阅介绍
	SettingPlanEbikeDescriptionKey   = "PLAN_EBIKE_DESCRIPTION"   // 车电订阅介绍
	SettingQuestionKey               = "QUESTION"                 // 常见问题
	SettingAppVersionKey             = "APP_VERSION"              // App版本
	SettingConsumePointKey           = "CONSUME_POINTS"           // 消费赠送积分
	SettingDailyRent                 = "DAILY_RENT"               // 日租金, 格式为 key(cityid + batterymodal + bikebrandid): amount
	SettingBatteryFaultKey           = "BATTERY_FAULT"            // 电池故障原因
	SettingEbikeFaultKey             = "EBIKE_FAULT"              // 电车故障原因
	SetttingOtherFaultKey            = "OTHER_FAULT"              // 其他故障原因
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

type SettingConsumePoint struct {
	CityID     uint64  `json:"cityId"`     // 城市ID
	Proportion float64 `json:"proportion"` // 赠送比例
}

var Settings = map[string]SettingItem{
	SettingDepositKey: {
		Desc:    "平台押金",
		Default: "99",
	},
	SettingRenewalKey: {
		Desc:    "重签判定时间",
		Default: "7",
	},
	SettingCabinetFaultKey: {
		Desc:    "电柜故障",
		Default: []string{},
	},
	SettingBatteryFaultKey: {
		Desc:    "电池故障",
		Default: []string{},
	},
	SettingEbikeFaultKey: {
		Desc:    "电车故障",
		Default: []string{},
	},
	SetttingOtherFaultKey: {
		Desc:    "其他故障",
		Default: []string{},
	},
	SettingExceptionKey: {
		Desc:    "物资异常",
		Default: []string{"丢失", "故障"},
	},
	SettingRescueReasonKey: {
		Desc:    "救援原因",
		Default: []string{},
	},
	SettingRescueFeeKey: {
		Desc:    "救援费用(元/公里)",
		Default: "0",
	},
	SettingReminderKey: {
		Desc: "催费通知",
		Default: SettingReminderNotice{
			App: []int{5},
			Sms: []int{3},
			Vms: []int{1, -3},
		},
	},
	SettingBatteryFullKey: {
		Desc:    "满电电量百分比",
		Default: "90",
	},
	// SettingPauseMaxDaysKey: {
	//     Desc:    "最大寄存时间",
	//     Default: "31",
	// },
	SettingExchangeIntervalKey: {
		Desc:    "限制换电间隔",
		Default: "20",
	},
	SettingExchangeLimitKey: {
		Desc:    "城市换电间隔",
		Default: make(SettingExchangeLimits),
	},
	SettingExchangeFrequencyKey: {
		Desc:    "城市换电频次",
		Default: make(SettingExchangeFrequencies),
	},
	SettingMaintainKey: {
		Desc:    "是否维护中",
		Default: false,
	},
	SettingReserveDurationKey: {
		Desc:    "最长预约时间",
		Default: "60",
	},
	SettingExchangeMinBatteryKey: {
		Desc:    "换电最低电量(%)",
		Default: "50",
	},
	SettingPlanBatteryDescriptionKey: {
		Desc:    "单电订阅介绍",
		Default: SettingPlanDescription{},
	},
	SettingPlanEbikeDescriptionKey: {
		Desc:    "车电订阅介绍",
		Default: SettingPlanDescription{},
	},
	SettingQuestionKey: {
		Desc:    "常见问题",
		Default: []SettingQuestion{},
	},
	SettingAppVersionKey: {
		Desc:    "App版本",
		Default: map[string]SettingAppVersionValue{},
	},
	SettingConsumePointKey: {
		Desc:    "消费赠送积分",
		Default: []SettingConsumePoint{},
	},
	SettingDailyRent: {
		Desc:    "日租金",
		Default: make(map[string]float64),
	},
}
