// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type order struct{}

var Order = new(order)

// Active
// @ID           RiderOrderActive
// @Router       /employee/v1/order/active [POST]
// @Summary      E20001 激活骑士卡
// @Tags         [E]门店接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body  model.QRPostReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*order) Active(c echo.Context) (err error) {
    ctx, req := app.ContextBinding[model.QRPostReq](c)
    // TODO 真实employee
    service.NewEmployeeOrderWithEmployee(&model.Employee{
        ID:    1,
        Name:  "超级店员",
        Phone: "18888888888",
    }).Active(req)
    return ctx.SendResponse()
}