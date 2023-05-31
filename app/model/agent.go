// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-31
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type AgentCreateReq struct {
	Phone        string `json:"phone" validate:"required" trans:"电话"`
	Name         string `json:"name" validate:"required" trans:"姓名"`
	EnterpriseID uint64 `json:"enterpriseId" validate:"required" trans:"团签ID"`
}

type AgentModifyReq struct {
	ID    uint64 `json:"id" param:"id" validate:"required" trans:"代理账号ID"`
	Phone string `json:"phone"` // 电话
	Name  string `json:"name"`  // 姓名
}

type AgentListReq struct {
	PaginationReq
	EnterpriseID uint64 `json:"enterpriseId" query:"enterpriseId" validate:"required" trans:"团签ID"`
}

type AgentMeta struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type AgentSigninReq struct {
	Phone   string `json:"phone" validate:"required" trans:"电话"`
	SmsId   string `json:"smsId" validate:"required" trans:"短信ID"`
	SmsCode string `json:"smsCode" validate:"required" trans:"短信验证码"`
}

type AgentProfile struct {
	ID         uint64              `json:"id"`
	Phone      string              `json:"phone"`              // 手机号
	Name       string              `json:"name"`               // 姓名
	Contract   string              `json:"contract,omitempty"` // 合同URL, 可能为空
	Enterprise Enterprise          `json:"enterprise"`         // 企业
	Balance    float64             `json:"balance"`            // 可用余额
	Riders     int                 `json:"riders"`             // 骑手数量
	Using      int                 `json:"using"`              // 计费中骑手数
	Yesterday  float64             `json:"yesterday"`          // 昨日使用
	Stations   []EnterpriseStation `json:"stations"`           // 站点列表
	Prices     []EnterprisePrice   `json:"prices"`             // 价格列表
	Days       []int               `json:"days"`               // 日期选项
}

type AgentSigninRes struct {
	Profile AgentProfile `json:"profile"`
	Token   string       `json:"token"`
}
