// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "context"
    "errors"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
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
            signinMessage := "需要登录"
            token := c.Request().Header.Get(app.HeaderRiderToken)
            id, err := ar.Cache.Get(context.Background(), token).Uint64()
            if err != nil {
                c.Set("errCode", response.StatusUnauthorized)
                return errors.New(signinMessage)
            }
            s := service.NewRider()
            var u *ent.Rider
            u, err = s.GetRiderById(id)
            if err != nil || u == nil {
                c.Set("errCode", response.StatusUnauthorized)
                return errors.New(signinMessage)
            }

            if s.IsBlocked(u) {
                c.Set("errCode", response.StatusForbidden)
                s.Signout(u)
                return errors.New("你已被封禁")
            }

            // 重载context
            return next(&app.RiderContext{
                GlobalContext: c.(*app.GlobalContext),
                Rider:         u,
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
                ctx.Set("errCode", response.StatusRequireContact)
                return errors.New("需要补充紧急联系人")
            }
            if p == nil || model.PersonAuthStatus(p.Status).RequireAuth() {
                ctx.Set("errCode", response.StatusRequireAuth)
                return errors.New("需要实名验证")
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
                ctx.Set("errCode", response.StatusLocked)
                ctx.Set("errData", ar.Map{"url": s.GetFaceUrl(ctx)})
                return errors.New("需要人脸验证")
            }
            return next(ctx)
        }
    }
}
