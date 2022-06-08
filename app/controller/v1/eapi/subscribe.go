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

// Active
// @ID           RiderOrderActive
// @Router       /employee/v1/subscribe/active [POST]
// @Summary      E2001 激活骑士卡
// @Tags         [E]门店接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body  model.QRPostReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Active(c echo.Context) (err error) {
    ctx, req := app.ContextBinding[model.QRPostReq](c)
    // TODO 真实employee
    service.NewEmployeeSubscribeWithEmployee(&model.Employee{
        ID:    1,
        Name:  "超级店员",
        Phone: "18888888888",
    }).Active(req)
    return ctx.SendResponse()
}
