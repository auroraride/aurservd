// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-18
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type selection struct {
}

var Selection = new(selection)

// Plan
// @ID           ManagerSelectionPlan
// @Router       /manager/v1/section/plan [GET]
// @Summary      MB001 筛选骑士卡
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.PlanSelectionReq  false  "骑士卡筛选项"
// @Success      200  {object}  []model.CascaderOptionLevel3  "请求成功"
func (*selection) Plan(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PlanSelectionReq](c)
    return ctx.SendResponse(service.NewSelection().Plan(req))
}

// Rider
// @ID           ManagerSelectionRider
// @Router       /manager/v1/section/rider [GET]
// @Summary      MB002 筛选骑手
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.RiderSelectionReq  true  "骑手筛选项"
// @Success      200  {object}  []model.SelectOption  "请求成功"
func (*selection) Rider(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderSelectionReq](c)
    return ctx.SendResponse(service.NewSelection().Rider(req))
}

// Store
// @ID           ManagerSelectionStore
// @Router       /manager/v1/section/store [GET]
// @Summary      MB003 筛选门店
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []model.CascaderOptionLevel2  "请求成功"
func (*selection) Store(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewSelection().Store())
}

// Employee
// @ID           ManagerSelectionEmployee
// @Router       /manager/v1/section/employee [GET]
// @Summary      MB004 筛选店员
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []model.CascaderOptionLevel2  "请求成功"
func (*selection) Employee(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewSelection().Employee())
}

// City
// @ID           ManagerSelectionCity
// @Router       /manager/v1/section/city [GET]
// @Summary      MB005 筛选启用的城市
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []model.CascaderOptionLevel2  "请求成功"
func (*selection) City(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewSelection().City())
}
