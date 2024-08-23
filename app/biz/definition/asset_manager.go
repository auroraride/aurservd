// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-17, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

// AssetManagerSigninReq 管理员登录请求
type AssetManagerSigninReq struct {
	Phone    string `json:"phone" validate:"required,phone" trans:"手机号"`
	Password string `json:"password" validate:"required" trans:"密码"`
}

// AssetManagerSigninRes 管理员登录返回
type AssetManagerSigninRes struct {
	ID          uint64   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Token       string   `json:"token,omitempty"`
	Phone       string   `json:"phone,omitempty"`
	Permissions []string `json:"permissions"` // 权限列表
	Super       bool     `json:"super"`       // 是否超级管理员
}

// AssetManagerCreateReq 管理员新增
type AssetManagerCreateReq struct {
	AssetManagerSigninReq
	Name   string `json:"name" validate:"required" trans:"姓名"`
	RoleID uint64 `json:"roleId" validate:"required" trans:"角色ID"`
}

type AssetManagerListReq struct {
	model.PaginationReq
	Keyword     *string `json:"keyword" query:"keyword"`         // 搜索关键词 姓名/手机号
	Warestore   *bool   `json:"warestore" query:"warestore"`     // 是否仓管人员
	WarehouseID *uint64 `json:"warehouseId" query:"warehouseId"` // 仓库id
	Enable      *bool   `json:"enable" query:"enable"`           // 仓管是否启用
}

type AssetManagerListRes struct {
	ID         uint64                   `json:"id"`
	Name       string                   `json:"name"`                 // 姓名
	Phone      string                   `json:"phone"`                // 手机号
	Role       AssetRole                `json:"role"`                 // 角色
	MiniEnable bool                     `json:"miniEnable"`           // 仓管是否启用
	MiniLimit  uint                     `json:"miniLimit"`            // 仓管限制范围
	Warehouses []*AssetManagerWarehouse `json:"warehouses,omitempty"` // 仓管仓库信息
}

type AssetManagerWarehouse struct {
	ID       uint64 `json:"id"`       // 门店ID
	Name     string `json:"name"`     // 门店名称
	CityName string `json:"cityName"` // 城市名称
}

type AssetManagerModifyReq struct {
	model.IDParamReq
	Password     string   `json:"password"`     // 密码
	RoleID       uint64   `json:"roleId"`       // 角色ID
	Phone        string   `json:"phone"`        // 电话
	Name         string   `json:"name"`         // 姓名
	WarehouseIDs []uint64 `json:"warehouseIds"` // 仓库IDS
	MiniEnable   *bool    `json:"miniEnable"`   // 仓管是否启用
	MiniLimit    *uint    `json:"miniLimit"`    // 仓管限制范围
}
