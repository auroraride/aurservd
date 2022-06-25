// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/labstack/echo/v4"
)

type socketapi struct{}

var Socket = new(socketapi)

func (*socketapi) Socket(c echo.Context) (err error) {
    ctx := app.Context(c)

    return ctx.SendResponse()
}