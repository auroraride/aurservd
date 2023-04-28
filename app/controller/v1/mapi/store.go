// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/labstack/echo/v4"
)

type store struct{}

var Store = new(store)

// List
// @ID           StoreList
// @Router       /manager/v1/store [GET]
// @Summary      M3006 列举门店
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query  model.StoreListReq  true  "desc"
// @Success      200  {object}  model.PaginationRes{items=[]model.StoreItem}  "请求成功"
func (*store) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.StoreListReq](c)
	return ctx.SendResponse(service.NewStore().List(req))
}

// Create
// @ID           StoreCreate
// @Router       /manager/v1/store [POST]
// @Summary      M3007 创建门店
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.StoreCreateReq  true  "desc"
// @Success      200  {object}  model.StoreItem  "请求成功"
func (*store) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.StoreCreateReq](c)
	return ctx.SendResponse(service.NewStoreWithModifier(ctx.Modifier).Create(req))
}

// Modify
// @ID           StoreModify
// @Router       /manager/v1/store/{id} [PUT]
// @Summary      M3008 修改门店
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  int  true  "门店ID"
// @Param        body  body  model.StoreModifyReq  true  "desc"
// @Success      200  {object}  model.StoreItem  "请求成功"
func (*store) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.StoreModifyReq](c)
	return ctx.SendResponse(service.NewStoreWithModifier(ctx.Modifier).Modify(req))
}

// Delete
// @ID           StoreDelete
// @Router       /manager/v1/store/{id} [DELETE]
// @Summary      M3009 删除门店
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  int  true  "门店ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*store) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	service.NewStoreWithModifier(ctx.Modifier).Delete(req)
	return ctx.SendResponse()
}
