// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    EnterpriseStatusLack         uint8 = iota // 未合作
    EnterpriseStatusCollaborated              // 合作中
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

type EnterprisePriceWithCity struct {
    Voltage float64 `json:"voltage" enums:"60,72"` // 电压型号 暂时固定为 60 或 72
    Price   float64 `json:"price"`                 // 单价(元/天)
    City    City    `json:"city"`                  // 城市
}

// EnterpriseBasic 企业基础字段
type EnterpriseBasic struct {
    ID   uint64 `json:"id"`   // 企业ID
    Name string `json:"name"` // 企业名称
}

// EnterpriseDetail 企业详细字段
type EnterpriseDetail struct {
    Name         *string              `json:"name" validate:"required" trans:"企业名称"`
    Status       *uint8               `json:"status" enums:"0,1,2" validate:"required,min=0,max=2" trans:"合作状态"` // 0:未合作 1:合作中 2:已暂停
    ContactName  *string              `json:"contactName" validate:"required" trans:"联系人"`
    ContactPhone *string              `json:"contactPhone" validate:"required" trans:"联系电话"`
    IdcardNumber *string              `json:"idcardNumber" validate:"required" trans:"身份证号"`
    CityID       *uint64              `json:"cityId" validate:"required" trans:"所在城市"`
    Address      *string              `json:"address" validate:"required" trans:"企业地址"`
    Payment      *uint8               `json:"payment" validate:"required,min=1,max=2" enums:"1,2" trans:"付费方式"` // 1:预付费 2:后付费
    Deposit      *float64             `json:"deposit" validate:"required" trans:"押金"`
    Contracts    []EnterpriseContract `json:"contracts,omitempty" validate:"required,min=1" trans:"合同"`
    Prices       []EnterprisePrice    `json:"prices,omitempty" validate:"required,min=1" trans:"价格列表"`
}

type EnterpriseDetailWithID struct {
    *EnterpriseDetail
    ID uint64 `json:"id" param:"id" validate:"required" trans:"企业ID"`
}

type EnterpriseListReq struct {
    PaginationReq
    CityID         *uint64 `json:"cityId" query:"cityId"`                 // 城市ID
    ContactKeyword *string `json:"contactKeyword" query:"contactKeyword"` // 联系人 姓名/电话/身份证 关键词
    Name           *string `json:"name" query:"name"`                     // 公司名称
    Status         *uint8  `json:"status" query:"status"`                 // 合作状态
    Payment        *uint8  `json:"payment" query:"payment" enums:"1,2"`   // 支付方式 1预付费 2后付费
    Start          *string `json:"start" query:"start"`                   // 合同到期时间晚于
    End            *string `json:"end" query:"end"`                       // 合同到期时间早于
    // StatementStart *string `json:"statementStart" query:"statementStart"` // 计费时间早于
    // StatementEnd   *string `json:"statementEnd" query:"statementEnd"`     // 计费时间晚于
}

type EnterpriseRes struct {
    ID           uint64                    `json:"id"`                    // 企业ID
    Balance      float64                   `json:"balance"`               // 可用余额
    Name         string                    `json:"name"`                  // 企业名称
    Status       uint8                     `json:"status" enums:"0,1,2" ` // 合作状态 0:未合作 1:已合作 2:已暂停
    ContactName  string                    `json:"contactName"`           // 联系人
    ContactPhone string                    `json:"contactPhone"`          // 联系电话
    IdcardNumber string                    `json:"idcardNumber"`          // 身份证号
    Address      string                    `json:"address"`               // 企业地址
    Payment      uint8                     `json:"payment"`               // 付费方式 1:预付费 2:后付费
    Deposit      float64                   `json:"deposit"`               // 押金
    Riders       int                       `json:"riders"`                // 骑手数量
    Contracts    []EnterpriseContract      `json:"contracts,omitempty"`   // 合同
    Prices       []EnterprisePriceWithCity `json:"prices,omitempty"`      // 价格列表
    City         City                      `json:"city"`                  // 城市
}

type EnterprisePrepaymentReq struct {
    ID     uint64  `json:"id" validate:"required" param:"id" trans:"企业ID"`
    Remark string  `json:"remark" validate:"required" trans:"备注"`
    Amount float64 `json:"amount" validate:"required" trans:"金额"`
}

type EnterpriseStationCreateReq struct {
    EnterpriseID uint64 `json:"enterpriseId" validate:"required"  trans:"企业ID"`
    Name         string `json:"name" validate:"required" trans:"站点名称"`
}

type EnterpriseStationModifyReq struct {
    Name string `json:"name" validate:"required" trans:"站点名称"`
    ID   uint64 `json:"id" validate:"required" param:"id" trans:"站点ID"`
}

type EnterpriseStation struct {
    ID   uint64 `json:"id"`   // 站点ID
    Name string `json:"name"` // 站点名称
}

type EnterpriseStationListReq struct {
    EnterpriseID uint64 `json:"enterpriseId" query:"enterpriseId" trans:"企业ID"`
}

type EnterprisePriceVoltageListReq struct {
    CityID uint64 `json:"cityId" validate:"required" query:"cityId" trans:"城市ID"`
}

type EnterprisePriceVoltageListRes struct {
    ID      uint64  `json:"id"`                    // 型号ID
    Voltage float64 `json:"voltage" enums:"60,72"` // 电压型号, 暂时固定为 60 或 72
}

type EnterpriseRiderSubscribeChooseReq struct {
    ID uint64 `json:"id" validate:"required" trans:"型号ID"`
}

type EnterpriseRiderSubscribeChooseRes struct {
    Qrcode string `json:"qrcode"` // 二维码, 格式为SUBSCRIBE:订阅ID, 后续使用订阅ID请求状态
}

type EnterpriseRiderSubscribeStatusReq struct {
    ID uint64 `json:"id" validate:"required" query:"id" trans:"订阅ID"`
}
