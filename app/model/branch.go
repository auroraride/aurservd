// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// BranchListReq 后台网点列表请求
type BranchListReq struct {
    PaginationReq

    CityID *uint64 `json:"cityId" query:"cityId" trans:"城市"`
}

// Branch 网点请求体
type Branch struct {
    ID uint64 `json:"id,omitempty" param:"id" swaggerignore:"true"`

    CityID    *uint64           `json:"cityId" validate:"required" trans:"城市"`
    Name      *string           `json:"name" validate:"required" trans:"网点名称"`
    Lng       *float64          `json:"lng" validate:"required" trans:"经度"`
    Lat       *float64          `json:"lat" validate:"required" trans:"纬度"`
    Address   *string           `json:"address" validate:"required" trans:"详细地址"`
    Photos    []string          `json:"photos" validate:"required" trans:"网点照片"`
    Contracts []*BranchContract `json:"contracts,omitempty"`

    Creator      *Modifier `json:"creator,omitempty"`
    LastModifier *Modifier `json:"lastModifier,omitempty"`
}

// BranchContract 网点合同请求体
type BranchContract struct {
    ID       uint64 `json:"id,omitempty" swaggerignore:"true"`
    BranchID uint64 `json:"branchId,omitempty" param:"id" swaggerignore:"true"`

    LandlordName      string   `json:"landlordName" validate:"required" trans:"房东姓名"`
    IDCardNumber      string   `json:"idCardNumber" validate:"required" trans:"房东身份证"`
    Phone             string   `json:"phone" validate:"required,phone" trans:"房东手机号"`
    BankNumber        string   `json:"bankNumber" validate:"required" trans:"房东银行卡号"`
    Pledge            float64  `json:"pledge" validate:"required" trans:"押金"`
    Rent              float64  `json:"rent" validate:"required" trans:"租金"`
    Lease             uint     `json:"lease" validate:"required" trans:"租期月数"`
    ElectricityPledge float64  `json:"electricityPledge" validate:"required" trans:"电费押金"`
    Electricity       float64  `json:"electricity" validate:"required" trans:"电费单价"`
    Area              float64  `json:"area" validate:"required" trans:"网点面积"`
    StartTime         string   `json:"startTime" validate:"required" trans:"租期开始时间"`
    EndTime           string   `json:"endTime" validate:"required" trans:"租期结束时间"`
    File              string   `json:"file" validate:"required" trans:"合同文件"`
    Sheets            []string `json:"sheets" validate:"required" trans:"底单"`
}
