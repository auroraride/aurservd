// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "time"

type CabinetBrand string

const (
    CabinetBrandKaixin  CabinetBrand = "KAIXIN"  // 凯信电柜
    CabinetBrandYundong CabinetBrand = "YUNDONG" // 云动电柜
)

func (cb CabinetBrand) String() string {
    return string(cb)
}

type CabinetStatus uint

const (
    CabinetStatusPending     CabinetStatus = iota // 未投放
    CabinetStatusOK                               // 运营中
    CabinetStatusMaintenance                      // 维护中
)

// Cabinet 电柜基础属性
type Cabinet struct {
    BranchID *uint64       `json:"branchId"`                                                    // 网点
    Status   CabinetStatus `json:"status" enums:"0,1,2"`                                        // 电柜状态 0未投放 1运营中 2维护中
    Brand    CabinetBrand  `json:"brand" validate:"required" trans:"品牌" enums:"KAIXIN,YUNDONG"` // KAIXIN(凯信) YUNDONG(云动)
    Serial   string        `json:"serial" validate:"required" trans:"电柜原始编码"`
    Name     string        `json:"name" validate:"required" trans:"电柜名称"`
    Doors    uint          `json:"doors" validate:"required" trans:"柜门数量"`
    Remark   *string       `json:"remark" trans:"备注"`
}

// CabinetCreateReq 电柜创建请求
type CabinetCreateReq struct {
    Cabinet
    Models []uint64 `json:"models" trans:"电池型号" validate:"required"`
}

// CabinetItem 电柜属性
type CabinetItem struct {
    ID uint64 `json:"id"` // 电柜ID
    Sn string `json:"sn"` // 平台编码
    Cabinet
    Models    []BatteryModel `json:"models"`    // 电池型号
    City      City           `json:"city"`      // 城市
    CreatedAt time.Time      `json:"createdAt"` // 创建时间
}

// CabinetQueryReq 电柜查询请求
type CabinetQueryReq struct {
    PaginationReq

    Serial *string        `json:"serial" query:"serial"`
    Name   *string        `json:"name" query:"name"`
    CityId *uint64        `json:"cityId" query:"cityId"`
    Brand  *string        `json:"brand" query:"brand"`
    Status *CabinetStatus `json:"status" query:"status"`
}
