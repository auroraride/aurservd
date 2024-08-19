// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-19, by aurb

package assetapi

import (
	"sort"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	assetPerm "github.com/auroraride/aurservd/app/assetpermission"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type assetPermission struct{}

var AssetPermission = new(assetPermission)

// List
// @ID		AssetManagerPermissionList
// @Router	/manager/v2/asset/permission [GET]
// @Summary	权限列表
// @Tags	权限 - AssetPermission
// @Accept	json
// @Produce	json
// @Param	X-AssetManager-Token	header		string				true	"管理员校验token"
// @Success	200						{object}	[]permission.Group	"请求成功"
func (*assetPermission) List(c echo.Context) (err error) {
	ctx := app.Context(c)
	items := make([]*assetPerm.Group, 0)
	for _, g := range assetPerm.Groups {
		items = append(items, g)
	}

	sort.Slice(items, func(i, j int) bool {
		return strings.Compare(items[i].Name, items[j].Name) < 0
	})

	return ctx.SendResponse(items)
}

// ListRole
// @ID		AssetManagerPermissionListRole
// @Router	/manager/v2/asset/permission/role [GET]
// @Summary	角色列表
// @Tags	权限 - AssetPermission
// @Accept	json
// @Produce	json
// @Param	X-AssetManager-Token	header		string			true	"管理员校验token"
// @Success	200						{object}	[]model.Role	"请求成功"
func (*assetPermission) ListRole(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(biz.NewAssetRole().List())
}

// CreateRole
// @ID		AssetManagerPermissionCreateRole
// @Router	/manager/v2/asset/permission/role [POST]
// @Summary	创建角色
// @Tags	权限 - AssetPermission
// @Accept	json
// @Produce	json
// @Param	X-AssetManager-Token	header		string							true	"管理员校验token"
// @Param	body					body		definition.AssetRoleCreateReq	true	"角色字段"
// @Success	200						{object}	definition.AssetRole			"请求成功"
func (*assetPermission) CreateRole(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.AssetRoleCreateReq](c)
	return ctx.SendResponse(
		biz.NewAssetRoleWithModifier(ctx.Modifier).Create(req),
	)
}

// ModifyRole
// @ID		AssetManagerPermissionModifyRole
// @Router	/manager/v2/asset/permission/role/{id} [PUT]
// @Summary	修改角色
// @Tags	权限 - AssetPermission
// @Accept	json
// @Produce	json
// @Param	X-AssetManager-Token	header		string							true	"管理员校验token"
// @Param	id						path		uint64							true	"角色ID"
// @Param	body					body		definition.AssetRoleModifyReq	true	"角色详情"
// @Success	200						{object}	definition.AssetRole			"请求成功"
func (*assetPermission) ModifyRole(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.AssetRoleModifyReq](c)
	return ctx.SendResponse(biz.NewAssetRoleWithModifier(ctx.Modifier).Modify(req))
}

// DeleteRole
// @ID		AssetManagerPermissionDeleteRole
// @Router	/manager/v2/asset/permission/role/{id} [DELETE]
// @Summary	删除角色
// @Tags	权限 - AssetPermission
// @Accept	json
// @Produce	json
// @Param	X-AssetManager-Token	header		string					true	"管理员校验token"
// @Param	id						path		uint64					true	"角色ID"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*assetPermission) DeleteRole(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	biz.NewAssetRoleWithModifier(ctx.Modifier).Delete(req)
	return ctx.SendResponse()
}
