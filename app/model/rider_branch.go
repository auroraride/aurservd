// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-11
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type RiderBranchDetailReq struct {
	BranchID uint64 `json:"branchId" query:"branchId" validate:"required" trans:"网点ID"`
	Model    string `json:"model" query:"model" validate:"required" trans:"电柜型号"`
	Type     uint8  `json:"type" query:"type" validate:"required" enums:"1,2" trans:"类别"` // 1:门店 2:电柜
}

type RiderBranchDetailRes struct {
	Name     string               `json:"name"`     // 门店或电柜名称
	Address  string               `json:"address"`  // 详细地址
	Lng      float64              `json:"lng"`      // 经度
	Lat      float64              `json:"lat"`      // 纬度
	Cabinets []Cabinet            `json:"cabinets"` // 电柜信息
	Batterys []RiderBranchBattery `json:"batterys"` // 可用电池信息
}

type RiderBranchBattery struct {
	Model    string `json:"model,omitempty"`    // 电池型号
	Fully    int    `json:"fully,omitempty"`    // 满电数量, 门店无此字段
	Charging int    `json:"charging,omitempty"` // 充电中数量, 门店无此字段
}

type RiderBranchCabinet struct {
	ID     uint64 `json:"id"`     // 电柜ID
	Name   string `json:"name"`   // 电柜名称
	Serial string `json:"serial"` // 电柜编码
}
