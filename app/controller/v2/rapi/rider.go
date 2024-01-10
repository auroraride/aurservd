// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-10
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
)

type rider struct{}

var Rider = new(rider)

func (*rider) Certification(c echo.Context) (err error) {
	ctx := app.Context(c)

	return ctx.SendResponse()
}
