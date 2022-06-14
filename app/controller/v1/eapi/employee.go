// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-08
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type employee struct{}

var Employee = new(employee)

// Signin
// @ID           EmployeeEmployeeSignin
// @Router       /employee/v1/signin [POST]
// @Summary      E1001 登录
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        body  body     model.EmployeeSignReq  true  "店员登录请求"
// @Success      200  {object}  model.EmployeeProfile  "店员登录信息"
func (*employee) Signin(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.EmployeeSignReq](c)
    return ctx.SendResponse(service.NewEmployee().Signin(req))
}

// Qrcode
// @ID           EmployeeEmployeeQrcode
// @Router       /employee/v1/qrcode [GET]
// @Summary      E1002 更换二维码
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Success      200  {object}  model.EmployeeQrcodeRes  "请求成功"
func (*employee) Qrcode(c echo.Context) (err error) {
    ctx := app.ContextX[app.EmployeeContext](c)
    return ctx.SendResponse(service.NewEmployeeWithEmployee(ctx.Employee).RefreshQrcode())
}

// Profile
// @ID           EmployeeEmployeeProfile
// @Router       /employee/v1/profile [GET]
// @Summary      E1005 店员资料
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token   header  string  true  "店员校验token"
// @Success      200  {object}      model.EmployeeProfile  "请求成功"
func (*employee) Profile(c echo.Context) (err error) {
    ctx := app.ContextX[app.EmployeeContext](c)
    return ctx.SendResponse(service.NewEmployee().Profile(ctx.Employee))
}
