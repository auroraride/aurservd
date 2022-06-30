// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-30
// Based on aurservd by liasica, magicrolan@qq.com.

package controller

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/labstack/echo/v4"
)

type version struct{}

var Version = new(version)

func (*version) Get(c echo.Context) (err error) {
    ctx := app.Context(c)
    plaform := ctx.QueryParam("plaform")
    var res ar.Version
    if plaform == "android" {
        res = ar.Config.Android
    } else {
        res = ar.Config.IOS
    }

    return ctx.SendResponse(res)
}
