// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-05
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    perm "github.com/auroraride/aurservd/app/permission"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
    "sort"
    "strings"
)

type permission struct{}

var Permission = new(permission)

// List
// @ID           ManagerPermissionList
// @Router       /manager/v1/permission [GET]
// @Summary      MD001 权限列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []perm.Group  "请求成功"
func (*permission) List(c echo.Context) (err error) {
    ctx := app.ContextX[app.ManagerContext](c)
    items := make([]*perm.Group, 0)
    for _, g := range perm.Groups {
        items = append(items, g)
    }

    sort.Slice(items, func(i, j int) bool {
        return strings.Compare(items[i].Name, items[j].Name) < 0
    })

    return ctx.SendResponse(items)
}

// ListRole
// @ID           ManagerPermissionListRole
// @Router       /manager/v1/permission/role [GET]
// @Summary      MD002 角色列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []model.Role  "请求成功"
func (*permission) ListRole(c echo.Context) (err error) {
    ctx := app.ContextX[app.ManagerContext](c)
    return ctx.SendResponse(service.NewRole().List())
}

// CreateRole
// @ID           ManagerPermissionCreateRole
// @Router       /manager/v1/permission/role [POST]
// @Summary      MD003 创建角色
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.RoleCreateReq  true  "角色字段"
// @Success      200  {object}  model.Role  "请求成功"
func (*permission) CreateRole(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RoleCreateReq](c)
    return ctx.SendResponse(
        service.NewRoleWithModifier(ctx.Modifier).Create(req),
    )
}

// ModifyRole
// @ID           ManagerPermissionModifyRole
// @Router       /manager/v1/permission/role/{id} [PUT]
// @Summary      MD004 修改角色
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "角色ID"
// @Param        body  body     model.RoleModifyReq  true  "角色详情"
// @Success      200  {object}  model.Role  "请求成功"
func (*permission) ModifyRole(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RoleModifyReq](c)
    return ctx.SendResponse(service.NewRoleWithModifier(ctx.Modifier).Modify(req))
}

// DeleteRole
// @ID           ManagerPermissionDeleteRole
// @Router       /manager/v1/permission/role/{id} [DELETE]
// @Summary      MD005 删除角色
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "角色ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*permission) DeleteRole(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
    service.NewRoleWithModifier(ctx.Modifier).Delete(req)
    return ctx.SendResponse()
}
