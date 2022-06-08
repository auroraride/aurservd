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
// @ID           ManagerEmployeeSignin
// @Router       /employee/v1/signin [POST]
// @Summary      E1001 登录
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.EmployeeSignReq  true  "desc"
// @Success      200  {object}  model.EmployeeProfile  "请求成功"
func (*employee) Signin(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.EmployeeSignReq](c)
    return ctx.SendResponse(service.NewEmployee().Signin(req))
}

func (*employee) Qrcode(c echo.Context) (err error) {
    ctx := app.ContextX[app.EmployeeContext](c)

    return ctx.SendResponse()
}