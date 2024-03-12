// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-17
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type exception struct{}

var Exception = new(exception)

// Setting
// @ID		EmployeeExceptionSetting
// @Router	/employee/v1/exception/setting [GET]
// @Summary	E3001 物资异常配置
// @Tags	[E]店员接口
// @Accept	json
// @Produce	json
// @Param	X-Employee-Token	header		string							true	"店员校验token"
// @Success	200					{object}	model.ExceptionEmployeeSetting	"请求成功"
func (*exception) Setting(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(
		service.NewException().Setting(),
	)
}

// Create
// @ID		EmployeeExceptionCreate
// @Router	/employee/v1/exception [POST]
// @Summary	E3002 异常上报
// @Tags	[E]店员接口
// @Accept	json
// @Produce	json
// @Param	X-Employee-Token	header		string						true	"店员校验token"
// @Param	body				body		model.ExceptionEmployeeReq	true	"异常上报请求"
// @Success	200					{object}	model.StatusResponse		"请求成功"
func (*exception) Create(c echo.Context) (err error) {
	ctx, req := app.EmployeeContextAndBinding[model.ExceptionEmployeeReq](c)
	service.NewExceptionWithEmployee(ctx.Employee).Create(req)
	return ctx.SendResponse()
}
