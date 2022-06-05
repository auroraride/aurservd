// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    EnterpriseStatusLack         uint8 = iota // 未合作
    EnterpriseStatusCollaborated              // 已合作
    EnterpriseStatusSuspensed                 // 已暂停
)

const (
    EnterprisePaymentPrepay  uint8 = iota + 1 // 预付费
    EnterprisePaymentPostPay                  // 后付费
)

type EnterpriseContract struct {
    Start string `json:"start" validate:"required"`
    End   string `json:"end" validate:"required"`
    File  string `json:"file" validate:"required"`
}

type EnterprisePrice struct {
    CityID  uint64  `json:"cityId" validate:"required" trans:"城市"`
    Voltage float64 `json:"voltage" validate:"required" enums:"60,72" trans:"电压型号"` // 暂时固定为 60 或 72
    Price   float64 `json:"price" validate:"required" trans:"单价(元/天)"`
}

type Enterprise struct {
    ID   uint64 `json:"id"`   // 企业ID
    Name string `json:"name"` // 企业名称
}

type EnterprisePostReq struct {
    Name         *string              `json:"name" validate:"required" trans:"企业名称"`
    Status       *uint8               `json:"status" enums:"0,1,2" validate:"required,min=0,max=2" trans:"合作状态"` // 0:未合作 1:已合作 2:已暂停
    ContactName  *string              `json:"contactName" validate:"required" trans:"联系人"`
    ContactPhone *string              `json:"contactPhone" validate:"required" trans:"联系电话"`
    IdcardNumber *string              `json:"idcardNumber" validate:"required" trans:"身份证号"`
    CityID       *uint64              `json:"cityId" validate:"required" trans:"所在城市"`
    Address      *string              `json:"address" validate:"required" trans:"企业地址"`
    Payment      *uint8               `json:"payment" validate:"required,min=1,max=2" enums:"1,2" trans:"付费方式"` // 1:预付费 2:后付费
    Deposit      *float64             `json:"deposit" validate:"required" trans:"押金"`
    Contracts    []EnterpriseContract `json:"contracts" validate:"required,min=1" trans:"合同"`
    Prices       []EnterprisePrice    `json:"prices" validate:"required,min=1" trans:"价格列表"`
}

type EnterpriseModifyReq struct {
    *EnterprisePostReq
    ID uint64 `json:"id" param:"id" validate:"required" trans:"企业ID"`
}
