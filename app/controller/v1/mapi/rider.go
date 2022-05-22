// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type rider struct{}

var Rider = new(rider)

// List
// @ID           RiderList
// @Router       /manager/v1/rider [GET]
// @Summary      M70001 列举骑手
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query  model.RiderListReq  true  "请求体"
// @Success      200  {object}  model.PaginationRes{items=[]model.RiderItem}  "请求成功"
func (*rider) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderListReq](c)
    return ctx.SendResponse(service.NewRider().List(req))
}

// Ban
// @ID           RiderBan
// @Router       /manager/v1/rider/ban [POST]
// @Summary      M70002 封禁/解除封禁身份
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.PersonBanReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*rider) Ban(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PersonBanReq](c)
    service.NewPersonWithModifier(ctx.Modifier).Ban(req)
    return ctx.SendResponse()
}

// Block
// @ID           RiderBlock
// @Router       /manager/v1/rider/block [POST]
// @Summary      M70003 封禁/解除封禁骑手账户
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.RiderBlockReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*rider) Block(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderBlockReq](c)
    service.NewRiderWithModifier(ctx.Modifier).Block(req)
    return ctx.SendResponse()
}
