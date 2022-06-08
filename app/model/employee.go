// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type Employee struct {
    ID    uint64 `json:"id"`
    Name  string `json:"name"`  // 店员名称
    Phone string `json:"phone"` // 店员电话
}

type EmployeeCreateReq struct {
    CityID uint64 `json:"cityId" validate:"required" trans:"城市ID"`
    Name   string `json:"name" validate:"required" trans:"姓名"`
    Phone  string `json:"phone" validate:"required,phone" trans:"手机号"`
}

type EmployeeModifyReq struct {
    ID     *uint64 `json:"id" validate:"required" param:"id" trans:"店员ID"`
    CityID *uint64 `json:"cityId"` // 城市ID
    Name   *string `json:"name"`   // 姓名
    Phone  *string `json:"phone"`  // 手机号
}

type EmployeeListReq struct {
    PaginationReq
    StoreID *uint64 `json:"storeId"` // 门店ID
    Keyword *string `json:"name"`    // 查询店员电话或姓名
    CityID  *uint64 `json:"cityId"`  // 城市ID
    Start   *string `json:"start"`   // 业绩统计开始时间
    End     *string `json:"end"`     // 业绩统计结束时间
}

type EmployeeListRes struct {
    ID              uint64  `json:"id"`
    Name            string  `json:"name"`            // 姓名
    Phone           string  `json:"phone"`           // 电话
    City            City    `json:"city"`            // 城市
    Store           *Store  `json:"store,omitempty"` // 当前上班门店, 字段为空的时候是休息状态
    ExchangeTimes   int     `json:"exchangeTimes"`   // 换电次数
    AssistanceTimes int     `json:"assistanceTimes"` // 救援次数
    AssistanceMiles float64 `json:"assistanceMiles"` // 救援里程(米)
}

type EmployeeDeleteReq struct {
    ID uint64 `json:"id" validate:"required" param:"id" trans:"店员ID"`
}

type EmployeeSignReq struct {
    Phone   string `json:"phone" validate:"required" trans:"电话"`
    SmsId   string `json:"smsId" validate:"required" trans:"短信ID"`
    SmsCode string `json:"smsCode" validate:"required" trans:"短信验证码"`
}

type EmployeeProfile struct {
    ID     uint64 `json:"id"`
    Token  string `json:"token"`  // 认证token
    Qrcode string `json:"qrcode"` // 二维码, 未上班或外出中二维码失效
    Phone  string `json:"phone"`  // 电话
    Name   string `json:"name"`   // 姓名
}
