// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package common

import (
    "github.com/auroraride/aurservd/app"
    "github.com/labstack/echo/v4"
)

type basic struct{}

var Basic = new(basic)

func (*basic) Get(c echo.Context) (err error) {
    ctx := app.Context(c)

    return ctx.SendResponse()
}