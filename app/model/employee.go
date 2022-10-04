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

type EmployeeActivityListFilter struct {
    StoreID uint64 `json:"storeId" query:"storeId"` // 门店ID
    Keyword string `json:"name" query:"keyword"`    // 查询店员电话或姓名
    CityID  uint64 `json:"cityId" query:"cityId"`   // 城市ID
    Start   string `json:"start" query:"start"`     // 业绩统计开始时间
    End     string `json:"end" query:"end"`         // 业绩统计结束时间
}

type EmployeeActivityListReq struct {
    PaginationReq
    EmployeeActivityListFilter
}

type EmployeeActivityExportReq struct {
    EmployeeActivityListFilter
    Remark string `json:"remark" validate:"required" trans:"备注"`
}

type EmployeeActivityListRes struct {
    ID               uint64  `json:"id"`
    Name             string  `json:"name"`             // 姓名
    Phone            string  `json:"phone"`            // 电话
    City             City    `json:"city"`             // 城市
    Store            *Store  `json:"store,omitempty"`  // 当前上班门店, 字段为空的时候是休息状态
    ExchangeTimes    int     `json:"exchangeTimes"`    // 换电次数
    Amount           float64 `json:"amount"`           // 业绩金额
    AssistanceTimes  int     `json:"assistanceTimes"`  // 救援次数
    AssistanceMeters float64 `json:"assistanceMeters"` // 救援里程(米)
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
    ID     uint64           `json:"id"`
    Token  string           `json:"token"`           // 认证token
    Qrcode string           `json:"qrcode"`          // 二维码, 未上班或外出中二维码失效
    Phone  string           `json:"phone"`           // 电话
    Name   string           `json:"name"`            // 姓名
    Onduty bool             `json:"onduty"`          // 是否上班
    Store  *StoreWithStatus `json:"store,omitempty"` // 上班门店, 未上班为空, 业务办理禁止进入
}

type EmployeeQrcodeRes struct {
    Qrcode string `json:"qrcode"`
}

type EmployeeListReq struct {
    PaginationReq

    Status  uint8   `json:"status" query:"status" enums:"0,1,2"` // 启用状态筛选 0:全部 1:启用 2:禁用
    Keyword *string `json:"keyword" query:"keyword"`             // 搜索关键词, 手机号或姓名
    CityID  *uint64 `json:"cityId" query:"cityId"`               // 城市ID
}

type EmployeeListRes struct {
    ID     uint64 `json:"id"`
    Enable bool   `json:"enable"`          // 是否启用
    Name   string `json:"name"`            // 姓名
    Phone  string `json:"phone"`           // 电话
    City   City   `json:"city"`            // 城市
    Store  *Store `json:"store,omitempty"` // 当前上班门店
}

type EmployeeEnableReq struct {
    ID     uint64 `json:"id" validate:"required" trans:"店员ID"`
    Enable bool   `json:"enable"` // 修改店员启用状态 true:启用 false:禁用
}
