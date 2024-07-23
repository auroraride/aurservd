// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-15, by Jorjan

package mapi

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
// @Router	/manager/v1/material [GET]
// @Summary	列表
// @Tags	物资 - Material
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		definition.MaterialListReq	true	"desc"
// @Success	200				{object}	[]definition.MaterialDetail	"请求成功"
func (*material) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.MaterialListReq](c)
	return ctx.SendResponse(biz.NewMaterial().List(req))
}

// Detail
// @ID		MaterialDetail
// @Router	/manager/v1/material/{id} [GET]
// @Summary	详情
// @Tags	物资 - Material
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	id				path		string						true	"仓库ID"
// @Success	200				{object}	definition.MaterialDetail	"请求成功"
func (*material) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewMaterial().Detail(req.ID))
}

// Create
// @ID		MaterialCreate
// @Router	/manager/v1/material [POST]
// @Summary	创建
// @Tags	物资 - Material
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.MaterialCreateReq	true	"desc"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*material) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.MaterialCreateReq](c)
	return ctx.SendResponse(biz.NewMaterialWithModifier(ctx.Modifier).Create(req))
}

// Delete
// @ID		MaterialDelete
// @Router	/manager/v1/material/{id} [DELETE]
// @Summary	删除
// @Tags	物资 - Material
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		string					true	"仓库ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*material) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewMaterialWithModifier(ctx.Modifier).Delete(req.ID))
}

// Modify
// @ID		MaterialModify
// @Router	/manager/v1/material/{id} [PUT]
// @Summary	修改
// @Tags	物资 - Material
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.MaterialModifyReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*material) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.MaterialModifyReq](c)
	return ctx.SendResponse(biz.NewMaterialWithModifier(ctx.Modifier).Modify(req))
}
