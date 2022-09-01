// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-01
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
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
                    en = ag.Edges.Enterprise
                }
            }

            if needLogin && (ag == nil || en == nil) {
                snag.Panic(snag.StatusUnauthorized)
            }

            return next(app.NewAgentContext(c, ag, en))
        }
    }
}
