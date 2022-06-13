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

// Detail
// @ID           EmployeeSubscribeDetail
// @Router       /employee/v1/subscribe/detail [GET]
// @Summary      E2001 待激活骑士卡详情
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        qrcode  query  string  true  "二维码详情, 可带`SUBSCRIBE:`, 也可不带"
// @Success      200  {object}  model.SubscribeActiveInfo  "请求成功"
func (*subscribe) Detail(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.QRQueryReq](c)
    return ctx.SendResponse(
        service.NewEmployeeSubscribeWithEmployee(ctx.Employee).Inactive(req.Qrcode),
    )
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
