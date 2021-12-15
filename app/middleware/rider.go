// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "context"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
)

// RiderMiddleware 骑手中间件
func RiderMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            url := c.Request().RequestURI
            if url == "/rider/signin" {
                return next(c)
            }

            // 获取骑手, 判定是否需要登录
            token := c.Request().Header.Get(app.HeaderRiderToken)
            id, err := ar.Cache.Get(context.Background(), token).Uint64()
            if err != nil {
                snag.Panic(app.Response{Code: app.StatusUnauthorized, Message: ar.RiderRequireSignin})
            }
            s := service.NewRider()
            var u *ent.Rider
            u, err = s.GetRiderById(id)
            if err != nil || u == nil {
                snag.Panic(app.Response{Code: app.StatusUnauthorized, Message: ar.RiderRequireSignin})
            }

            // 用户被封禁
            if s.IsBlocked(u) {
                s.Signout(u)
                snag.Panic(app.Response{Code: app.StatusForbidden, Message: ar.RiderBlockedMessage})
            }

            // 延长token有效期
            s.ExtendTokenTime(u.ID, token)

            // 重载context
            return next(&app.RiderContext{
                Context: c.(*app.Context),
                Rider:   u,
            })
        }
    }
}

// RiderRequireAuthAndContact 实名验证以及紧急联系人中间件
func RiderRequireAuthAndContact() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            ctx := c.(*app.RiderContext)
            uri := c.Request().RequestURI
            p := ctx.Rider.Edges.Person
            if ctx.Rider.Contact == nil && uri != "/rider/contact" {
                snag.Panic(app.Response{Code: app.StatusRequireContact, Message: ar.RiderRequireContact})
            }
            if p == nil || model.PersonAuthStatus(p.Status).RequireAuth() {
                snag.Panic(app.Response{Code: app.StatusRequireAuth, Message: ar.RiderRequireAuth})
            }
            return next(ctx)
        }
    }
}

// RiderFaceMiddleware 检测是否需要人脸验证
func RiderFaceMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            ctx := c.(*app.RiderContext)
            u := ctx.Rider
            s := service.NewRider()
            if s.IsNewDevice(u, ctx.Device) {
                snag.Panic(app.Response{Code: app.StatusLocked, Message: ar.RiderRequireFace, Data: ar.Map{"url": s.GetFaceUrl(ctx)}})
            }
            return next(ctx)
        }
    }
}
