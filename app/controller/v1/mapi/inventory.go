// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type inventory struct{}

var Inventory = new(inventory)

// CreateOrModify
// @ID		ManagerInventoryCreateOrModify
// @Router	/manager/v1/inventory [POST]
// @Summary	物资设定创建或更新
// @Tags	物资
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.Inventory			true	"desc"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*inventory) CreateOrModify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.Inventory](c)
	service.NewInventoryWithModifier(ctx.Modifier).CreateOrModify(req)
	return ctx.SendResponse()
}

// List
// @ID		ManagerInventoryList
// @Router	/manager/v1/inventory [GET]
// @Summary	列举物资设定
// @Tags	物资
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string			true	"管理员校验token"
// @Success	200				{object}	model.Inventory	"请求成功"
func (*inventory) List(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewInventory().List())
}

// Delete
// @ID		ManagerInventoryDelete
// @Router	/manager/v1/inventory [DELETE]
// @Summary	删除物资设定
// @Tags	物资
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.InventoryDelete	true	"desc"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*inventory) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.InventoryDelete](c)
	service.NewInventoryWithModifier(ctx.Modifier).Delete(req)
	return ctx.SendResponse()
}

// Transferable
// @ID		ManagerInventoryTransferable
// @Router	/manager/v1/inventory/transferable [GET]
// @Summary	可调拨物资清单
// @Tags	物资
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Success	200				{object}	[]model.InventoryItem	"请求成功"
func (*inventory) Transferable(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewInventory().ListInventory(model.InventoryListReq{Transfer: true}))
}
