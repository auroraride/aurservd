// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-20
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
)

func AutoToastMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			visible := c.Request().Header.Get(app.HeaderToastVisible)
			if visible != "NO" {
				visible = "YES"
			}
			c.Response().Header().Set(app.HeaderToastVisible, visible)
			return next(c)
		}
	}
}
