// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-05-29, by Jorjan

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type goods struct{}

var Goods = new(goods)

// List
// @ID		GoodsList
// @Router	/manager/v1/goods [GET]
// @Summary	商品列表
// @Tags	商品
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		definition.GoodsListReq		true	"desc"
// @Success	200				{object}	[]definition.GoodsDetail	"请求成功"
func (*goods) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.GoodsListReq](c)
	return ctx.SendResponse(biz.NewGoods().List(req))
}

// Detail
// @ID		GoodsDetail
// @Router	/manager/v1/goods/{id} [GET]
// @Summary	商品详情
// @Tags	商品
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		string					true	"商品ID"
// @Success	200				{object}	definition.GoodsDetail	"请求成功"
func (*goods) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewGoods().Detail(req.ID))
}

// Create
// @ID		GoodsCreate
// @Router	/manager/v1/goods [POST]
// @Summary	创建商品
// @Tags	商品
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		definition.GoodsCreateReq	true	"desc"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*goods) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.GoodsCreateReq](c)
	return ctx.SendResponse(biz.NewGoodsWithModifierBiz(ctx.Modifier).Create(req))
}

// Delete
// @ID		GoodsDelete
// @Router	/manager/v1/goods/{id} [DELETE]
// @Summary	删除商品
// @Tags	商品
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		string					true	"商品ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*goods) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewGoodsWithModifierBiz(ctx.Modifier).Delete(req.ID))
}

// Modify
// @ID		GoodsModify
// @Router	/manager/v1/goods/{id} [PUT]
// @Summary	修改商品
// @Tags	商品
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		definition.GoodsModifyReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*goods) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.GoodsModifyReq](c)
	return ctx.SendResponse(biz.NewGoodsWithModifierBiz(ctx.Modifier).Modify(req))
}

// UpdateStatus
// @ID		GoodsUpdateStatus
// @Router	/manager/v1/goods/status/{id} [PUT]
// @Summary	上下架商品
// @Tags	商品
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.GoodsUpdateStatusReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*goods) UpdateStatus(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.GoodsUpdateStatusReq](c)
	return ctx.SendResponse(biz.NewGoodsWithModifierBiz(ctx.Modifier).UpdateStatus(req))
}
