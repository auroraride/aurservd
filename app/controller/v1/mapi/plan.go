// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type plan struct{}

var Plan = new(plan)

// Create
// @ID           PlanCreate
// @Router       /manager/v1/plan [POST]
// @Summary      M60001 创建骑士卡
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*plan) Create(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PlanCreateReq](c)
    service.NewPlan().CreatePlan(ctx.Modifier, req)
    return ctx.SendResponse()
}

// UpdateEnable
// @ID           PlanUpdateEnable
// @Router       /manager/v1/plan/{id} [PUT]
// @Summary      M60002 上下架骑士卡
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  int  true  "骑士卡ID"
// @Param        body  body  model.PlanEnableModifyReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*plan) UpdateEnable(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PlanEnableModifyReq](c)
    service.NewPlan().UpdateEnable(ctx.Modifier, req)
    return ctx.SendResponse()
}

// Delete
// @ID           PlanDelete
// @Router       /manager/v1/plan/{id} [DELETE]
// @Summary      M60003 删除骑士卡
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  int  true  "骑士卡ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*plan) Delete(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
    service.NewPlan().Delete(ctx.Modifier, req)
    return ctx.SendResponse()
}

// List
// @ID           PlanList
// @Router       /manager/v1/plan [GET]
// @Summary      M60004 列举骑士卡
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query  model.PlanListReq  true  "desc"
// @Success      200  {object}  model.PaginationRes{items=[]model.PlanItemRes}  "请求成功"
func (*plan) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PlanListReq](c)
    return ctx.SendResponse(service.NewPlan().List(req))
}