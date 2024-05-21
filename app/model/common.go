// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "github.com/jinzhu/copier"

// StatusResponse 状态返回
type StatusResponse struct {
	Status bool `json:"status"`
}

// SmsReq 短信请求
type SmsReq struct {
	Phone       string `json:"phone" validate:"required"`       // 手机号
	CaptchaCode string `json:"captchaCode" validate:"required"` // captcha 验证码
}

// SmsRes 短信发送返回
type SmsRes struct {
	Id string `json:"id"` // 任务ID
}

// AliyunOssStsRes 阿里云oss临时凭证
type AliyunOssStsRes struct {
	AccessKeySecret string `json:"accessKeySecret,omitempty"`
	Expiration      string `json:"expiration,omitempty"`
	AccessKeyId     string `json:"accessKeyId,omitempty"`
	StsToken        string `json:"stsToken,omitempty"`
	Bucket          string `json:"bucket,omitempty"`
	Region          string `json:"region,omitempty"`
}

// CaptchaReq 验证码请求
type CaptchaReq struct {
	Code string `json:"code"`
}

// SmsResponse 短信验证码返回
type SmsResponse struct {
	Id string `json:"id"`
}

// ItemListRes 列表返回
type ItemListRes struct {
	Items []interface{} `json:"items" kind:"slice"`
}

func SetItemListResItems[T any](res *ItemListRes, items []T) {
	res.Items = make([]any, len(items))
	_ = copier.Copy(&res.Items, items)
}

// ItemRes 单项返回
type ItemRes struct {
	Item interface{} `json:"item"`
}

// IDRes ID返回
type IDRes struct {
	ID uint64 `json:"id"`
}

// IDPostReq ID post请求
type IDPostReq struct {
	ID uint64 `json:"id" validate:"required"`
}

// IDParamReq id param 请求
type IDParamReq struct {
	ID uint64 `json:"id" param:"id" validate:"required"`
}

// IDQueryReq id param 请求
type IDQueryReq struct {
	ID uint64 `json:"id" query:"id" validate:"required"`
}

// QRPostReq 二维码POST请求
type QRPostReq struct {
	Qrcode string `json:"qrcode" validate:"required" trans:"二维码"`
}

type QRQueryReq struct {
	Qrcode string `json:"qrcode" query:"qrcode" validate:"required" trans:"二维码"`
}

type SelectOption struct {
	Value uint64 `json:"value"`          // 选择项值 (ID)
	Label string `json:"label"`          // 选择项名称
	Desc  string `json:"desc,omitempty"` // 描述
}

type SelectOptionGroup struct {
	Label   string         `json:"label"`   // 分组名称
	Options []SelectOption `json:"options"` // 分组选项
}

type CascaderOptionLevel2 struct {
	SelectOption
	Children []SelectOption `json:"children"` // 级联选择子项目
}

type CascaderOptionLevel3 struct {
	SelectOption
	Children *[]CascaderOptionLevel2 `json:"children"` // 级联选择子项目
}

type LngLat struct {
	Lng float64 `json:"lng" query:"lng"` // 经度
	Lat float64 `json:"lat" query:"lat"` // 纬度
}

type CascaderOption struct {
	Value    any                `json:"value"`              // 值
	Label    string             `json:"label"`              // 名
	Children *[]*CascaderOption `json:"children,omitempty"` // 子
}

type KeywordQueryReq struct {
	Keyword *string `json:"keyword" validate:"required" query:"keyword" trans:"关键词"`
}

const (
	SigninTypeSms  uint64 = iota + 1 // 短信登录
	SigninTypeAuth                   // 授权登录
)

const (
	AuthTypeWechat uint8 = iota + 1 // 微信
	AuthTypeAlipay                  // 支付宝
)
