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
	Keyword *string `json:"keyword" query:"keyword"` // 搜索关键词 姓名/手机号
}

type AssetManagerListRes struct {
	ID    uint64    `json:"id"`
	Name  string    `json:"name"`  // 姓名
	Phone string    `json:"phone"` // 手机号
	Role  AssetRole `json:"role"`  // 角色
}

type AssetManagerModifyReq struct {
	model.IDParamReq
	Password string `json:"password"` // 密码
	RoleID   uint64 `json:"roleId"`   // 角色ID
	Phone    string `json:"phone"`    // 电话
	Name     string `json:"name"`     // 姓名
}
