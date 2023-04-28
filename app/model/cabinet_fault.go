// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "time"

const (
	CabinetFaultStatusPending uint8 = iota // 未处理
	CabinetFaultStatusDone                 // 已处理
	CabinetFaultStatusRejected
)

// CabinetFaultReportReq 故障上报请求体
type CabinetFaultReportReq struct {
	CabinetID   uint64   `json:"cabinetId" validate:"required" trans:"电柜ID"`
	Fault       string   `json:"fault" validate:"required" trans:"故障"`
	Description string   `json:"description" validate:"required" trans:"故障描述"`
	Attachments []string `json:"attachments" validate:"max=3"` // 附件
}

// CabinetFaultListReq 故障列表
type CabinetFaultListReq struct {
	PaginationReq

	CityID      *uint64 `json:"cityId" query:"cityID"`
	CabinetName *string `json:"cabinetName" query:"cabinetName"`
	Serial      *string `json:"serial" query:"serial"`
	Status      *uint8  `json:"status" query:"status"`
	Start       *string `json:"start" query:"start"`
	End         *string `json:"end" query:"end"`
}

type CabinetFaultCabinet struct {
	ID     uint64         `json:"id"`     // 电柜ID
	Name   string         `json:"name"`   // 电柜名称
	Brand  string         `json:"brand"`  // 电柜品牌
	Serial string         `json:"serial"` // 电柜编号
	Models []BatteryModel `json:"models"` // 电池型号
}

// CabinetFaultItem 故障信息
type CabinetFaultItem struct {
	ID          uint64              `json:"id"`          // 故障ID
	Status      uint8               `json:"status"`      // 故障状态 0未处理 1已处理
	City        City                `json:"city"`        // 城市信息
	Cabinet     CabinetFaultCabinet `json:"cabinet"`     // 电柜信息
	Fault       string              `json:"fault"`       // 故障原因
	Rider       RiderSampleInfo     `json:"rider"`       // 骑手信息
	Attachments []string            `json:"attachments"` // 故障附件
	Description string              `json:"description"` // 故障描述
	CreatedAt   time.Time           `json:"createdAt"`   // 提交时间
}

// CabinetFaultDealReq 故障处理请求
type CabinetFaultDealReq struct {
	ID     *uint64 `json:"id" validate:"required" trans:"故障ID" param:"id"`
	Status *uint8  `json:"status" validate:"required,gte=0,lte=1" enums:"0,1" trans:"故障状态"` // 0未处理 1已处理
	Remark *string `json:"remark" trans:"备注"`
}

type CabinetFaultNotice struct {
	City        string `json:"city"`
	Branch      string `json:"branch"`
	Name        string `json:"name"`
	Serial      string `json:"serial"`
	Phone       string `json:"phone"`
	Fault       string `json:"fault"`
	Description string `json:"description"`
}
