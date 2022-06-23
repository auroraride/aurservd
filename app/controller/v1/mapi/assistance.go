// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-23
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type assistance struct{}

var Assistance = new(assistance)

func (*assistance) Nearby(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDQueryReq](c)

    return ctx.SendResponse(service.NewAssistanceWithModifier(ctx.Modifier).Nearby(req))
}