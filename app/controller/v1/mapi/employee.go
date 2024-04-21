// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-08
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type employee struct{}

var Employee = new(employee)

// Create
// @ID		ManagerEmployeeCreate
// @Router	/manager/v1/employee [POST]
// @Summary	新增店员
// @Tags	店员
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.EmployeeCreateReq	true	"desc"
// @Success	200				{object}	uint64					"请求成功"
func (*employee) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EmployeeCreateReq](c)
	return ctx.SendResponse(service.NewEmployeeWithModifier(ctx.Modifier).Create(req).ID)
}

// Modify
// @ID		ManagerEmployeeModify
// @Router	/manager/v1/employee/{id} [PUT]
// @Summary	修改店员
// @Tags	店员
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		uint64					true	"店员ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*employee) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EmployeeModifyReq](c)
	service.NewEmployeeWithModifier(ctx.Modifier).Modify(req)
	return ctx.SendResponse()
}

// List
// @ID		ManagerEmployeeList
// @Router	/manager/v1/employee [GET]
// @Summary	列举店员
// @Tags	店员
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string												true	"管理员校验token"
// @Param	query			query		model.EmployeeListReq								true	"筛选选项"
// @Success	200				{object}	model.PaginationRes{items=[]model.EmployeeListRes}	"请求成功"
func (*employee) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EmployeeListReq](c)
	return ctx.SendResponse(service.NewEmployeeWithModifier(ctx.Modifier).List(req))
}

// Delete
// @ID		ManagerEmployeeDelete
// @Router	/manager/v1/employee/{id} [DELETE]
// @Summary	删除店员
// @Tags	店员
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		uint64					true	"店员ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*employee) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EmployeeDeleteReq](c)
	service.NewEmployeeWithModifier(ctx.Modifier).Delete(req)
	return ctx.SendResponse()
}

// Activity
// @ID		ManagerEmployeeActivity
// @Router	/manager/v1/employee/activity [GET]
// @Summary	店员业绩
// @Tags	店员
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	query			query		model.EmployeeActivityListReq	true	"店员业绩列表筛选"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*employee) Activity(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EmployeeActivityListReq](c)
	return ctx.SendResponse(service.NewEmployeeWithModifier(ctx.Modifier).Activity(req))
}

// Enable
// @ID		ManagerEmployeeEnable
// @Router	/manager/v1/emoloyee/enable [POST]
// @Summary	启用/禁用店员
// @Tags	店员
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.EmployeeEnableReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*employee) Enable(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EmployeeEnableReq](c)
	service.NewEmployeeWithModifier(ctx.Modifier).Enable(req)
	return ctx.SendResponse()
}

// OffWork
// @ID		ManagerEmployeeOffWork
// @Router	/manager/v1/employee/offwork [POST]
// @Summary	强制下班
// @Tags	店员
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				body		uint64					true	"店员ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*employee) OffWork(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDPostReq](c)
	service.NewEmployeeWithModifier(ctx.Modifier).OffWork(req)
	return ctx.SendResponse()
}
