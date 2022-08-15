// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-24
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type business struct{}

var Business = new(business)

// List
// @ID           ManagerBusinessList
// @Router       /manager/v1/business [GET]
// @Summary      MG001 骑手业务记录
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.BusinessListReq  true  "列表请求筛选参数"
// @Success      200  {object}  model.PaginationRes{items=[]model.BusinessListRes}  "请求成功"
func (*business) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BusinessListReq](c)
    return ctx.SendResponse(service.NewBusiness().ListManager(req))
}

// Pause
// @ID           ManagerBusinessPause
// @Router       /manager/v1/business/pause [GET]
// @Summary      MG002 寄存记录
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.BusinessPauseList  true  "列表请求筛选参数"
// @Success      200  {object}  model.PaginationRes{items=[]model.BusinessPauseListRes}  "请求成功"
func (*business) Pause(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BusinessPauseList](c)
    return ctx.SendResponse(service.NewBusinessWithModifier(ctx.Modifier).ListPause(req))
}

// Reserve
// @ID           ManagerBusinessReserve
// @Router       /manager/v1/business/reserve [GET]
// @Summary      MG004 预约记录
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.ReserveListReq  true  "列表请求筛选参数"
// @Success      200  {object}  model.PaginationRes{items=[]model.ReserveListRes}  "请求成功"
func (*business) Reserve(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.ReserveListReq](c)
    return ctx.SendResponse(service.NewReserveWithModifier(ctx.Modifier).List(req))
}

// Suspend
// @ID           ManagerBusinessSuspend
// @Router       /manager/v1/business/suspend [GET]
// @Summary      MG005 暂停记录
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.SuspendListReq  true  "列表请求筛选参数"
// @Success      200  {object}  model.PaginationRes{items=[]model.SuspendListRes}  "请求成功"
func (*business) Suspend(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.SuspendListReq](c)
    return ctx.SendResponse(service.NewSuspendWithModifier(ctx.Modifier).List(req))
}
