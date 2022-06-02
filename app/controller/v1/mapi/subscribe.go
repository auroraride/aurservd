// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-02
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type subscribe struct{}

var Subscribe = new(subscribe)

// Alter
// @ID           ManagerSubscribeAlter
// @Router       /manager/v1/subscribe/alter [POST]
// @Summary      M70004 修改订阅时间
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.SubscribeAlter  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Alter(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.SubscribeAlter](c)
    return ctx.SendResponse(service.NewSubscribeWithModifier(ctx.Modifier).AlterDays(req))
}
