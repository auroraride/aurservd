// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-20, by liasica

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
)

type business struct{}

var Business = new(business)

func (*business) Exchange(c echo.Context) (err error) {
	ctx := app.Context(c)

	return ctx.SendResponse()
}
