// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type stock struct{}

var Stock = new(stock)

// Create
// @ID           ManagerStockCreate
// @Router       /manager/v1/stock [POST]
// @Summary      M1015 调拨物资
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.StockTransferReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*stock) Create(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.StockTransferReq](c)
    service.NewStockWithModifier(ctx.Modifier).Transfer(req)
    return ctx.SendResponse()
}

// List
// @ID           ManagerStockList
// @Router       /manager/v1/stock [GET]
// @Summary      M1017 门店物资详细
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.StockListReq  true  "desc"
// @Success      200  {object}  model.PaginationRes{items=[]model.StockListRes}  "请求成功"
func (*stock) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.StockListReq](c)
    return ctx.SendResponse(service.NewStockWithModifier(ctx.Modifier).List(req))
}

// Overview
// @ID           ManagerStockOverview
// @Router       /manager/v1/stock/overview [GET]
// @Summary      M1016 物资管理概览
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  model.StockOverview  "请求成功"
func (*stock) Overview(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewStock().Overview())
}
