// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-01
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
)

var (
	agentAuthSkipper = map[string]bool{
		"[POST]/agent/v1/signin": true,
	}
)

func AgentMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			p := fmt.Sprintf("[%s]%s", c.Request().Method, c.Path())
			needLogin := !agentAuthSkipper[p]

			var en *ent.Enterprise
			var ag *ent.Agent

			token := c.Request().Header.Get(app.HeaderAgentToken)
			if token != "" {
				if id, err := cache.Get(context.Background(), token).Uint64(); err == nil && id > 0 {
					// 获取代理和图签
					ag, _ = service.NewAgent().Query(id)
					if ag != nil {
						en = ag.Edges.Enterprise
					}
				}
			}

			if needLogin && (ag == nil || en == nil) {
				snag.Panic(snag.StatusUnauthorized)
			}

			return next(app.NewAgentContext(c, ag, en))
		}
	}
}

func Agent() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var en *ent.Enterprise
			var ag *ent.Agent

			// 查看是否是登录态
			str := c.Request().Header.Get("Authorization")
			if str != "" && strings.HasPrefix(str, app.AgentBearer) {
				token := strings.TrimLeft(str, app.AgentBearer)
				// 查找登录用户
				ag, en = service.NewAgent().TokenVerify(token)
				c.Set("agent", ag)
				c.Set("enterprise", en)
			}
			return next(app.NewAgentContext(c, ag, en))
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
