// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-16
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/labstack/echo/v4"
)

type order struct{}

var Order = new(order)

// List
// @ID           EmployeeOrderList
// @Router       /employee/v1/order [GET]
// @Summary      E2009 订单记录
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*order) List(c echo.Context) (err error) {
    ctx := app.Context(c)

    return ctx.SendResponse()
}
