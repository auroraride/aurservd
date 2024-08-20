// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-19, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

const (
	SignTokenWarehouse = "WAREHOUSE" // 仓库平台
	SignTokenStore     = "STORE"     // 门店平台
)

type PlatType uint8

const (
	PlatTypeWarehouse PlatType = iota + 1 // 仓库平台登录
	PlatTypeStore                         // 门店平台登录
)

func (f PlatType) Value() uint8 {
	return uint8(f)
}

type WarestorePeopleSigninReq struct {
	Phone    string   `json:"phone" validate:"required" trans:"电话"`
	Password string   `json:"password" validate:"required" trans:"密码"`
	PlatType PlatType `json:"platType" validate:"required,oneof=1 2" trans:"登录平台类型"`
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
	ID           uint64 `json:"id"`
	Phone        string `json:"phone"`        // 手机号
	Name         string `json:"name"`         // 姓名
	RoleName     string `json:"roleName"`     // 角色名称
	Duty         bool   `json:"duty"`         // 上下班 `true`上班 `false`下班
	DutyLocation string `json:"dutyLocation"` // 上班位置
}

// TransferListReq 调拨记录筛选条件
type TransferListReq struct {
	model.PaginationReq
	model.AssetTransferFilter
}
