// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-12, by aurb

package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/snag"
)

func Warestore() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				am *ent.AssetManager
				ep *ent.Employee
			)
			// 查看是否是登录态
			token := strings.TrimSpace(c.Request().Header.Get(app.HeaderWarestoreToken))
			if token != "" {
				// 查找登录用户
				am = biz.NewWarestore().TokenVerify(token)
			}
			return next(app.NewWarestoreContext(c, am, ep))
		}
	}
}

func WarestoreAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.(*app.WarestoreContext)
			if ctx.AssetManager == nil && ctx.Employee == nil {
				snag.Panic(snag.StatusUnauthorized)
			}

			return next(ctx)
		}
	}
}
