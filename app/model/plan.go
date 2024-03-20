// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "golang.org/x/exp/slices"

type PlanType uint8

const (
	PlanTypeBattery          PlanType = 1 + iota // 单电
	PlanTypeEbikeWithBattery                     // 车加电
)

func (t PlanType) Value() uint8 {
	return uint8(t)
}

func (t PlanType) IsValid() bool {
	return slices.Contains(PlanTypes, t)
}

var (
	PlanTypes = []PlanType{PlanTypeBattery, PlanTypeEbikeWithBattery}
)

// Plan 骑士卡基础信息
type Plan struct {
	ID          uint64  `json:"id"`          // 骑士卡ID
	Name        string  `json:"name"`        // 骑士卡名称
	Days        uint    `json:"days"`        // 骑士卡天数
	Intelligent bool    `json:"intelligent"` // 是否智能电柜套餐
	Price       float64 `json:"price"`       // 售价
}

type PlanComplexes []PlanComplex

type PlanComplex struct {
	ID uint64 `json:"id,omitempty"` // ID (可为空, 编辑的时候需要携带此字段)

	Price float64 `json:"price" validate:"required" trans:"价格"`
	Days  uint    `json:"days" validate:"required,min=1" trans:"有效天数"`

	Original      float64 `json:"original"`      // 原价
	Desc          string  `json:"desc"`          // 优惠信息
	Commission    float64 `json:"commission"`    // 提成
	DiscountNewly float64 `json:"discountNewly"` // 新签优惠

	CommissionName string `json:"commissionName,omitempty"` // 佣金方案名称
	CommissionID   uint64 `json:"commissionId,omitempty"`   // 佣金方案ID

	Model string `json:"model" validate:"required"` // 电池型号, 单电需要每一项都补充此字段
}

type PlanCreateReq struct {
	Type      PlanType       `json:"type" validate:"required,enum" trans:"骑士卡类别"`
	Name      string         `json:"name" validate:"required" trans:"骑士卡名称"`
	Start     string         `json:"start" validate:"required,datetime=2006-01-02" trans:"开始日期"`
	End       string         `json:"end" validate:"required,datetime=2006-01-02" trans:"结束日期"`
	Cities    []uint64       `json:"cities" validate:"required,min=1" trans:"启用城市"`
	Complexes []*PlanComplex `json:"complexes" validate:"required,min=1" trans:"骑士卡详细信息"`

	BrandID uint64 `json:"brandId" validate:"required_if=Type 2" trans:"电车型号"` // 车加电必填

	Enable bool     `json:"enable"` // 是否启用
	Notes  []string `json:"notes"`  // 购买须知

	Intelligent *bool `json:"intelligent" validate:"required"` // 是否智能柜套餐

	Deposit                 *bool    `json:"deposit" validate:"required"` // 是否开启押金
	DepositAmount           *float64 `json:"depositAmount"`               // 押金金额
	DepositPay              *bool    `json:"depositPay"`                  // 是否支持支付押金 true:支持 false:不支持
	DepositWechatPayscore   *bool    `json:"depositWechatPayscore"`       // 是否支持微信支付分免押金 true:支持 false:不支持
	DepositAlipayAuthFreeze *bool    `json:"depositAlipayAuthFreeze"`     // 是否支持预授权信用免押金 true:支持 false:不支持
	DepositContract         *bool    `json:"depositContract"`             // 是否支持合同免押金 true:支持 false:不支持
}

// PlanEnableModifyReq 骑士卡状态修改请求
type PlanEnableModifyReq struct {
	ID     uint64 `json:"id" validate:"required" param:"id"` // 骑士卡ID
	Enable *bool  `json:"enable" validate:"required"`        // 启用或禁用
}

// PlanListReq 列表筛选
type PlanListReq struct {
	PaginationReq

	Intelligent *uint8    `json:"intelligent" query:"intelligent"` // 智能柜套餐筛选 0:全部 1:是 2:否
	CityID      *uint64   `json:"cityId" query:"cityId"`           // 城市ID
	Name        *string   `json:"name" query:"name"`               // 骑士卡名称
	Enable      *bool     `json:"enable" query:"enable"`           // 启用状态
	Model       *string   `json:"model" query:"model"`             // 电池型号
	Type        *PlanType `json:"type" query:"type" enums:"1,2"`   // 骑士卡类别, 不携带字段为全部, 1:单电 2:车加电
	BrandID     *uint64   `json:"brandId" query:"brandId"`         // 电车型号
	Deposit     *bool     `json:"deposit" query:"deposit"`         // 是否开启押金 fales:未开启 true:开启
}

type PlanNotSettedDailyRent struct {
	City       City        `json:"city"`                 // 城市
	Model      string      `json:"model"`                // 电池型号
	EbikeBrand *EbikeBrand `json:"ebikeBrand,omitempty"` // 电车型号
}

type PlanListRes struct {
	ID        uint64           `json:"id"`
	Type      PlanType         `json:"type"`      // 类别
	Name      string           `json:"name"`      // 名称
	Enable    bool             `json:"enable"`    // 是否启用
	Start     string           `json:"start"`     // 开始日期
	End       string           `json:"end"`       // 结束日期
	Cities    []City           `json:"cities"`    // 可用城市
	Complexes []*PlanComplexes `json:"complexes"` // 详情集合(按电池型号分组)

	Brand       *EbikeBrand `json:"brand,omitempty"` // 电车型号
	Notes       []string    `json:"notes,omitempty"` // 购买须知
	Intelligent bool        `json:"intelligent"`     // 是否智能柜套餐
	Model       string      `json:"model"`           // 电池型号

	NotSettedDailyRent []*PlanNotSettedDailyRent `json:"notSettedDailyRent,omitempty"` // 未设定的日租金

	Deposit                 bool    `json:"deposit"`                 // 是否开启押金 fales:未开启 true:开启
	DepositAmount           float64 `json:"depositAmount"`           // 押金金额
	DepositPay              bool    `json:"depositPay"`              // 是否支持支付押金 true:支持 false:不支持
	DepositWechatPayscore   bool    `json:"depositWechatPayscore"`   // 是否支持微信支付分免押金 true:支持 false:不支持
	DepositAlipayAuthFreeze bool    `json:"depositAlipayAuthFreeze"` // 是否支持预授权信用免押金 true:支持 false:不支持
	DepositContract         bool    `json:"depositContract"`         // 是否支持合同免押金 true:支持 false:不支持
}

// PlanListRiderReq 骑士卡列表请求
type PlanListRiderReq struct {
	CityID uint64 `json:"cityId" query:"cityId" validate:"required" trans:"城市ID"`

	Min          uint    `json:"min" swaggerignore:"true"`          // 最小天数
	Model        string  `json:"model" swaggerignore:"true"`        // 电池型号
	EbikeBrandID *uint64 `json:"ebikeBrandId" swaggerignore:"true"` // 电车型号
	Intelligent  bool    `json:"intelligent" swaggerignore:"true"`  // 是否智能
}

// RiderPlanItem 骑士返回数据
type RiderPlanItem struct {
	ID            uint64  `json:"id"`
	Name          string  `json:"name"`          // 骑士卡名称
	Price         float64 `json:"price"`         // 价格
	Days          uint    `json:"days"`          // 天数
	Original      float64 `json:"original"`      // 原价
	DiscountNewly float64 `json:"discountNewly"` // 新签优惠
	Desc          string  `json:"desc"`          // 优惠信息
}

type RiderPlanListRes struct {
	Model   string          `json:"model"`   // 电池型号
	Plans   []RiderPlanItem `json:"plans"`   // 套餐列表
	Deposit float64         `json:"deposit"` // 需缴纳押金
}

// // PlanItem 单项骑士卡详情(用做订单备份)
// type PlanItem struct {
//     ID         uint64  `json:"id"`
//     Name       string  `json:"name"`       // 骑士卡名称
//     Enable     bool    `json:"enable"`     // 是否启用
//     Start      string  `json:"start"`      // 开始日期
//     End        string  `json:"end"`        // 结束日期
//     Price      float64 `json:"price"`      // 价格
//     Days       uint    `json:"days"`       // 有效天数
//     Original   float64 `json:"original"`   // 原价
//     Desc       string  `json:"desc"`       // 优惠信息
//     Commission float64 `json:"commission"` // 提成
// }

type PlanSelectionReq struct {
	Effect *uint8 `json:"effect" query:"effect" enums:"0,1,2"` // 筛选生效中 0:全部(默认) 1:生效中 2:未生效
	Status *uint8 `json:"status" query:"status" enums:"0,1,2"` // 筛选状态 0:全部(默认) 1:启用 2:禁用
}

type CommissionPlanSelectionReq struct {
	Keyword *string `json:"keyword" query:"keyword"` // 关键字
}

type PlanDaysPriceOptions []PlanDaysPriceOption

// PlanDaysPriceOption 骑士卡天数价格选项
type PlanDaysPriceOption struct {
	ID            uint64  `json:"id"`
	Name          string  `json:"name"`          // 骑士卡名称
	Price         float64 `json:"price"`         // 价格, 应支付价格 = 价格 - 新签优惠
	Days          uint    `json:"days"`          // 天数
	Original      float64 `json:"original"`      // 原价
	DiscountNewly float64 `json:"discountNewly"` // 新签优惠
	HasEbike      bool    `json:"hasEbike"`      // 是否包含电车

	Deposit                 bool    `json:"deposit"`                 // 是否启用押金
	DepositAmount           float64 `json:"depositAmount"`           // 押金金额
	DepositWechatPayscore   bool    `json:"depositWechatPayscore"`   // 是否支持微信支付分免押金
	DepositAlipayAuthFreeze bool    `json:"depositAlipayAuthFreeze"` // 是否支持预授权信用免押金
	DepositContract         bool    `json:"depositContract"`         // 是否支持合同免押金
	DepositPay              bool    `json:"depositPay"`              // 是否支持支付押金
}

type PlanModelOptions []*PlanModelOption

// PlanModelOption 新签电池型号选项
type PlanModelOption struct {
	Children *PlanDaysPriceOptions `json:"children"` // 天数和价格信息
	Model    string                `json:"model"`    // 型号
	Intro    string                `json:"intro"`    // 介绍图
	Notes    []string              `json:"notes"`    // 购买须知
}

type PlanEbikeBrandOptions []*PlanEbikeBrandOption

// PlanEbikeBrandOption 新签电车品牌选项
type PlanEbikeBrandOption struct {
	Children *PlanModelOptions `json:"children"`        // 子项
	Name     string            `json:"name"`            // 名称
	Cover    string            `json:"cover,omitempty"` // 封面图
}

type RiderPlanRenewalRes struct {
	Items     []RiderPlanItem   `json:"items"`               // 骑士卡列表
	Overdue   bool              `json:"overdue"`             // 是否需要支付逾期费用
	Days      uint              `json:"days,omitempty"`      // 逾期天数, 可能为空
	Fee       float64           `json:"fee,omitempty"`       // 逾期费用, 可能为空
	Formula   string            `json:"formula,omitempty"`   // 逾期费用计算公式, 可能为空
	Configure *PaymentConfigure `json:"configure,omitempty"` // 支付配置
}

type PlanNewlyRes struct {
	Brands    []*PlanEbikeBrandOption `json:"brands,omitempty"`    // 车电选项
	Models    []*PlanModelOption      `json:"models,omitempty"`    // 单电选项
	Deposit   float64                 `json:"deposit"`             // 需缴纳押金
	Configure *PaymentConfigure       `json:"configure,omitempty"` // 支付配置

	BatteryDescription SettingPlanDescription `json:"batteryDescription"` // 单电介绍
	EbikeDescription   SettingPlanDescription `json:"ebikeDescription"`   // 车电介绍
}

type PlanModifyTimeReq struct {
	ID    uint64 `json:"id" validate:"required" trans:"骑士卡ID"` // 使用items[n].id
	Start string `json:"start" validate:"required,datetime=2006-01-02" trans:"开始日期"`
	End   string `json:"end" validate:"required,datetime=2006-01-02" trans:"结束日期"`
}
