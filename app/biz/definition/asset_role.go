// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-17, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

type AssetRole struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`                  // 角色名称
	Permissions []string `json:"permissions,omitempty"` // 允许的权限列表 (权限key), 未分配权限的无此字段
	Builtin     bool     `json:"builtin"`               // 是否内置角色, 内置角色无法编辑
	Super       bool     `json:"super"`                 // 是否超级管理员, 超级管理员权限无法编辑且拥有全部权限
}

type AssetRoleCreateReq struct {
	Name        string   `json:"name" validate:"required" trans:"角色名称"`
	Permissions []string `json:"permissions,omitempty"` // 权限列表, 可以创建之后编辑再添加权限
}

type AssetRoleModifyReq struct {
	model.IDParamReq
	Name        string   `json:"name,omitempty"`        // 角色名称
	Permissions []string `json:"permissions,omitempty"` // 权限列表
}
