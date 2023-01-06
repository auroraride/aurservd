// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-06
// Based on aurservd by liasica, magicrolan@qq.com.

package controller

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type maintain struct{}

var Maintain = new(maintain)

func (*maintain) Update(c echo.Context) (err error) {
    ctx := app.Context(c)

    // 标记为维护中
    service.NewMaintain().SetMaintain(true)

    // 是否有进行中的换电业务

    // 是否有进行中的其他业务

    return ctx.SendResponse()
}