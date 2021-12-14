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
    "github.com/auroraride/aurservd/internal/baidu"
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

            // 重载context
            ctx := &app.RiderContext{
                GlobalContext: c.(*app.GlobalContext),
                Rider:         u,
            }

            // TODO 判定骑手是否新设备

            // TODO 判断是否需要认证和补充联系人
            return next(ctx)
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
            r := ctx.Rider
            p := r.Edges.Person
            if ctx.Device.Serial != r.LastDevice {
                ctx.Set("errCode", response.StatusLocked)
                ctx.Set("errData", ar.Map{"url": baidu.New().GetFaceUrl(p.Name, p.IcNumber)})
                return errors.New("需要人脸验证")
            }
            return next(ctx)
        }
    }
}
