// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-15, by Jorjan

package amapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type material struct{}

var Material = new(material)

// List
// @ID		MaterialList
// @Router	/manager/v2/asset/material [GET]
// @Summary	列表
// @Tags	Material - 其他物资
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string													true	"管理员校验token"
// @Param	query					query		definition.MaterialListReq								true	"desc"
// @Success	200						{object}	model.PaginationRes{items=[]definition.MaterialDetail}	"请求成功"
func (*material) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.MaterialListReq](c)
	return ctx.SendResponse(biz.NewMaterial().List(req))
}

// Create
// @ID		MaterialCreate
// @Router	/manager/v2/asset/material [POST]
// @Summary	创建
// @Tags	Material - 其他物资
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	body					body		definition.MaterialCreateReq	true	"请求参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*material) Create(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.MaterialCreateReq](c)
	return ctx.SendResponse(biz.NewMaterialWithModifier(ctx.Modifier).Create(req))
}

// Delete
// @ID		MaterialDelete
// @Router	/manager/v2/asset/material/{id} [DELETE]
// @Summary	删除
// @Tags	Material - 其他物资
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	id						path		string					true	"仓库ID"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*material) Delete(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewMaterialWithModifier(ctx.Modifier).Delete(req.ID))
}

// Modify
// @ID		MaterialModify
// @Router	/manager/v2/asset/material/{id} [PUT]
// @Summary	修改
// @Tags	Material - 其他物资
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	id						path		int								true	"ID"
// @Param	body					body		definition.MaterialModifyReq	true	"请求参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*material) Modify(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.MaterialModifyReq](c)
	return ctx.SendResponse(biz.NewMaterialWithModifier(ctx.Modifier).Modify(req))
}
