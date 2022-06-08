// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-08
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type employee struct{}

var Employee = new(employee)

// Create
// @ID           ManagerEmployeeCreate
// @Router       /manager/v1/employee [POST]
// @Summary      MA010 新增店员
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.EmployeeCreateReq  true  "desc"
// @Success      200  {object}  uint64  "请求成功"
func (*employee) Create(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EmployeeCreateReq](c)
    return ctx.SendResponse(service.NewEmployeeWithModifier(ctx.Modifier).Create(req).ID)
}

// Modify
// @ID           ManagerEmployeeModify
// @Router       /manager/v1/employee/{id} [PUT]
// @Summary      MA011 修改店员
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "店员ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*employee) Modify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EmployeeModifyReq](c)
    service.NewEmployeeWithModifier(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}

// List
// @ID           ManagerEmployeeList
// @Router       /manager/v1/employee [GET]
// @Summary      MA012 列举店员
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query  model.EmployeeListReq  true  "desc"
// @Success      200  {object}  model.PaginationRes{items=[]model.EmployeeListRes}  "请求成功"
func (*employee) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EmployeeListReq](c)
    return ctx.SendResponse(service.NewEmployeeWithModifier(ctx.Modifier).List(req))
}

// Delete
// @ID           ManagerEmployeeDelete
// @Router       /manager/v1/employee/{id} [DELETE]
// @Summary      MA013 删除店员
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "店员ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*employee) Delete(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EmployeeDeleteReq](c)
    service.NewEmployeeWithModifier(ctx.Modifier).Delete(req)
    return ctx.SendResponse()
}
