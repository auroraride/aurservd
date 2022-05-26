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
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
)

var (
    riderAuthSkipper = map[string]bool{
        "/rider/v1/signin": true,
        // "/rider/v1/city":   true,
    }
)

// RiderMiddleware 骑手中间件
func RiderMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            url := c.Request().RequestURI
            if riderAuthSkipper[url] {
                return next(c)
            }

            // 获取骑手, 判定是否需要登录
            token := c.Request().Header.Get(app.HeaderRiderToken)
            id, err := cache.Get(context.Background(), token).Uint64()
            if err != nil {
                snag.Panic(snag.StatusUnauthorized)
            }
            s := service.NewRider()
            var u *ent.Rider
            u, err = s.GetRiderById(id)
            if err != nil || u == nil {
                snag.Panic(snag.StatusUnauthorized)
            }

            // 用户被封禁
            if s.IsBanned(u) {
                s.Signout(u)
                snag.Panic(snag.StatusForbidden, ar.BannedMessage)
            }

            // 延长token有效期
            s.ExtendTokenTime(u.ID, token)

            // 获取与判定是否需要更新骑手推送ID
            pushId := c.Request().Header.Get(app.HeaderPushId)
            if u.PushID != pushId {
                _ = ar.Ent.Rider.UpdateOneID(u.ID).SetPushID(pushId).Exec(context.Background())
            }

            // 重载context
            return next(app.NewRiderContext(c, u))
        }
    }
}

// RiderRequireAuthAndContact 实名验证以及紧急联系人中间件
func RiderRequireAuthAndContact() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            url := c.Request().RequestURI
            // if riderAuthSkipper[url] {
            //     return next(c)
            // }

            ctx := c.(*app.RiderContext)

            p := ctx.Rider.Edges.Person
            if ctx.Rider.Contact == nil && url != "/rider/contact" {
                snag.Panic(snag.StatusRequireContact)
            }
            if p == nil || model.PersonAuthStatus(p.Status).RequireAuth() {
                snag.Panic(snag.StatusRequireAuth, ar.Map{"url": service.NewRider().GetFaceAuthUrl(ctx)})
            }
            return next(ctx)
        }
    }
}

// RiderFaceMiddleware 检测是否需要人脸验证
func RiderFaceMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            // url := c.Request().RequestURI
            // if riderAuthSkipper[url] {
            //     return next(c)
            // }

            ctx := c.(*app.RiderContext)
            u := ctx.Rider
            s := service.NewRider()
            if s.IsNewDevice(u, ctx.Device) {
                snag.Panic(snag.StatusLocked, ar.Map{"url": s.GetFaceUrl(ctx)})
            }
            return next(ctx)
        }
    }
}
