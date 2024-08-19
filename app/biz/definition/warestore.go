// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-19, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

type PlatType uint8

const (
	PlatTypeWarehouse PlatType = iota + 1 // 仓库平台登录
	PlatTypeStore                         // 门店平台登录
)

func (f PlatType) Value() uint8 {
	return uint8(f)
}

type WarestorePeopleSigninReq struct {
	Phone      string   `json:"phone,omitempty" validate:"required_if=SigninType 1" trans:"电话"`
	SmsId      string   `json:"smsId,omitempty" validate:"required_if=SigninType 1" trans:"短信ID"`
	Code       string   `json:"code,omitempty" validate:"required_if=SigninType 1,required_if=SigninType 2" trans:"验证码"`
	SigninType uint64   `json:"signinType" validate:"required,oneof=1 2"`
	PlatType   PlatType `json:"platType" validate:"required,oneof=1 2" trans:"登录平台类型"`
}

type WarestorePeopleSigninRes struct {
	Profile WarestorePeopleProfile `json:"profile"`
	Token   string                 `json:"token"`
}

type OpenidReq struct {
	Code string `json:"code" query:"code"`
}

type OpenidRes struct {
	Openid string `json:"openid"`
}

type WarestorePeopleProfile struct {
	ID    uint64 `json:"id"`
	Phone string `json:"phone"` // 手机号
	Name  string `json:"name"`  // 姓名
}

// TransferListReq 调拨记录筛选条件
type TransferListReq struct {
	model.PaginationReq
	model.AssetTransferFilter
}
