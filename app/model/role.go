// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-06
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type Role struct {
    ID          uint64   `json:"id"`
    Name        string   `json:"name"`                  // 角色名称
    Permissions []string `json:"permissions,omitempty"` // 允许的权限列表 (权限key), 未分配权限的无此字段
    Builtin     bool     `json:"builtin"`               // 是否内置角色, 内置角色无法编辑
    Super       bool     `json:"super"`                 // 是否超级管理员, 超级管理员权限无法编辑且拥有全部权限
}

type RoleCreateReq struct {
    Name string `json:"name" validate:"required" trans:"角色名称"`
}

type RoleModifyReq struct {
    IDParamReq
    Name        string   `json:"name,omitempty"`        // 角色名称
    Permissions []string `json:"permissions,omitempty"` // 权限列表
}
