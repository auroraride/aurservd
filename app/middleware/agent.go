// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-01
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/snag"
)

func Agent() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				ag       *ent.Agent
				en       *ent.Enterprise
				stations ent.EnterpriseStations
			)

			// 查看是否是登录态
			token := strings.TrimSpace(c.Request().Header.Get(app.HeaderAgentToken))
			if token != "" {
				// 查找登录用户
				ag, en, stations = service.NewAgent().TokenVerify(token)
			}
			return next(app.NewAgentContext(c, ag, en, stations))
		}
	}
}

func AgentAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.(*app.AgentContext)
			if ctx.Agent == nil || ctx.Enterprise == nil {
				snag.Panic(snag.StatusUnauthorized)
			}
			return next(ctx)
		}
	}
}
