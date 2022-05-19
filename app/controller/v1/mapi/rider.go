// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-20
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/labstack/echo/v4"
)

type rider struct{}

var Rider = new(rider)

func (*rider) List(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse()
}