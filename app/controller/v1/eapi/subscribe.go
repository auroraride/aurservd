// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-30
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type subscribe struct{}

var Subscribe = new(subscribe)

func (*subscribe) Detail(c echo.Context) (err error) {
    ctx := app.Context(c)
    
    return ctx.SendResponse()
}

// Active
// @ID           RiderOrderActive
// @Router       /employee/v1/subscribe/active [POST]
// @Summary      E2002 激活骑士卡
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body  model.QRPostReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Active(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.QRPostReq](c)
    service.NewEmployeeSubscribeWithEmployee(ctx.Employee).Active(req)
    return ctx.SendResponse()
}
