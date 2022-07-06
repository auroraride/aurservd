// Copyright (C) liasica. 2021-present.
//
// Created at 2022/2/25
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "context"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
)

// ManagerMiddleware 后台中间件
func ManagerMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            url := c.Request().URL.Path
            if url == "/manager/v1/user/signin" {
                return next(c)
            }

            // 判定登录
            token := c.Request().Header.Get(app.HeaderManagerToken)
            id, err := cache.Get(context.Background(), token).Uint64()
            if err != nil {
                snag.Panic(snag.StatusUnauthorized)
            }
            s := service.NewManager()
            var m *ent.Manager
            m, err = s.GetManagerById(id)
            if err != nil || m == nil {
                snag.Panic(snag.StatusUnauthorized)
            }

            // 延长token有效期
            s.ExtendTokenTime(m.ID, token)

            // 重载context
            return next(app.NewManagerContext(c, m, &model.Modifier{
                ID:    m.ID,
                Name:  m.Name,
                Phone: m.Phone,
            }))
        }
    }
}
