// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// Plan 骑士卡基础信息
type Plan struct {
    ID   uint64 `json:"id"`   // 骑士卡ID
    Name string `json:"name"` // 骑士卡名称
    Days uint   `json:"days"` // 骑士卡天数
}

type PlanComplex struct {
    ID         uint64  `json:"id,omitempty"` // 子项ID (可为空, 编辑的时候需要携带此字段)
    Price      float64 `json:"price" validate:"required" trans:"价格"`
    Days       uint    `json:"days" validate:"required,min=1" trans:"有效天数"`
    Original   float64 `json:"original"`   // 原价
    Desc       string  `json:"desc"`       // 优惠信息
    Commission float64 `json:"commission"` // 提成
}

type PlanCreateReq struct {
    Name      string        `json:"name" validate:"required" trans:"骑士卡名称"`
    Enable    bool          `json:"enable"` // 是否启用
    Start     string        `json:"start" validate:"required,datetime=2006-01-02" trans:"开始日期"`
    End       string        `json:"end" validate:"required,datetime=2006-01-02" trans:"结束日期"`
    Cities    []uint64      `json:"cities" validate:"required,min=1" trans:"启用城市"`
    Models    []string      `json:"models" validate:"required,min=1" trans:"电池型号"`
    Complexes []PlanComplex `json:"complexes" validate:"required,min=1" trans:"骑士卡详细信息"`
}

type PlanModifyReq struct {
    ID uint64 `json:"id" param:"id" validate:"required" trans:"骑士卡ID"`
    PlanCreateReq
}

// PlanEnableModifyReq 骑士卡状态修改请求
type PlanEnableModifyReq struct {
    ID     uint64 `json:"id" validate:"required" param:"id"` // 骑士卡ID
    Enable *bool  `json:"enable" validate:"required"`        // 启用或禁用
}

// PlanListReq 列表请求
type PlanListReq struct {
    PaginationReq

    CityID *uint64 `json:"cityId" query:"cityId"` // 城市ID
    Name   *string `json:"name" query:"name"`     // 骑士卡名称
    Enable *bool   `json:"enable" query:"enable"` // 启用状态
}

type PlanWithComplexes struct {
    ID        uint64         `json:"id"`        // 骑士卡ID
    Name      string         `json:"name"`      // 骑士卡名称
    Enable    bool           `json:"enable"`    // 是否启用
    Start     string         `json:"start"`     // 开始日期
    End       string         `json:"end"`       // 结束日期
    Cities    []City         `json:"cities"`    // 可用城市
    Models    []BatteryModel `json:"models"`    // 可用型号
    Complexes []PlanComplex  `json:"complexes"` // 详情集合
}

// PlanListRiderReq 骑士卡列表请求
type PlanListRiderReq struct {
    CityID uint64 `json:"cityId" query:"cityId" validate:"required" trans:"城市ID"`
    Min    uint   `json:"min" swaggerignore:"true"` // 最小天数
}

type RiderPlanListRes struct {
    Model   string          `json:"model"`   // 电池型号
    Plans   []RiderPlanItem `json:"plans"`   // 套餐列表
    Deposit float64         `json:"deposit"` // 需缴纳押金
}

// RiderPlanItem 骑士返回数据
type RiderPlanItem struct {
    ID       uint64  `json:"id"`
    Name     string  `json:"name"`     // 骑士卡名称
    Price    float64 `json:"price"`    // 价格
    Days     uint    `json:"days"`     // 天数
    Original float64 `json:"original"` // 原价
    Desc     string  `json:"desc"`     // 优惠信息
}

type RiderPlanRenewalRes struct {
    Items   []RiderPlanItem `json:"items"`             // 骑士卡列表
    Overdue bool            `json:"overdue"`           // 是否需要支付逾期费用
    Days    uint            `json:"days,omitempty"`    // 逾期天数, 可能为空
    Fee     float64         `json:"fee,omitempty"`     // 逾期费用, 可能为空
    Formula string          `json:"formula,omitempty"` // 逾期费用计算公式, 可能为空
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
