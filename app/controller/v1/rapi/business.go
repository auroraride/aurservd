// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-27
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/labstack/echo/v4"
)

type business struct{}

var Business = new(business)

func (*business) Active(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse()
}

func (*business) Unsubscribe(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse()
}

func (*business) Pause(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse()
}

func (*business) Continue(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse()
}
