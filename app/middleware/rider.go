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
    riderLoginSkipper = map[string]bool{
        "/rider/v1/signin":                    true,
        "/rider/v1/socket":                    true,
        "/rider/v1/callback":                  true,
        "/rider/v1/callback/esign":            true,
        "/rider/v1/callback/alipay":           true,
        "/rider/v1/callback/wechatpay":        true,
        "/rider/v1/callback/wechatpay/refund": true,
        "/rider/v1/city":                      true,
        "/rider/v1/branch":                    true,
        "/rider/v1/branch/riding":             true,
        "/rider/v1/branch/facility/:fid":      true,
        "/rider/v1/setting/question":          true,
    }
    riderAuthSkipper = map[string]bool{
        "/rider/v1/profile": true,
        "/rider/v1/reserve": true,
    }
    riderFaceSkipper = map[string]bool{
        "/rider/v1/profile": true,
        "/rider/v1/reserve": true,
    }
)

func init() {
    for k, v := range riderLoginSkipper {
        riderAuthSkipper[k] = v
        riderFaceSkipper[k] = v
    }
}

// riderLogin 获取骑手
func riderLogin(token, pushId string, needLogin bool) (u *ent.Rider) {
    var err error
    s := service.NewRider()
    id, _ := cache.Get(context.Background(), token).Uint64()
    u, err = s.GetRiderById(id)
    // 判定是否需要登录
    if needLogin && (err != nil || u == nil) {
        snag.Panic(snag.StatusUnauthorized)
    }

    if u != nil {
        // 延长token有效期
        s.ExtendTokenTime(u.ID, token)

        // 获取与判定是否需要更新骑手推送ID
        if u.PushID != pushId {
            _ = ent.Database.Rider.UpdateOneID(u.ID).SetPushID(pushId).Exec(context.Background())
        }

        // 用户被封禁
        if s.IsBanned(u) || s.IsBlocked(u) {
            s.Signout(u)
            snag.Panic(snag.StatusForbidden, ar.BannedMessage)
        }
    }

    return u
}

// RiderMiddleware 骑手中间件
func RiderMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            url := c.Path()
            token := c.Request().Header.Get(app.HeaderRiderToken)
            needLogin := !riderLoginSkipper[url]
            pushId := c.Request().Header.Get(app.HeaderPushId)
            u := riderLogin(token, pushId, needLogin)
            // 重载context
            return next(app.NewRiderContext(c, u, token))
        }
    }
}

// RiderRequireAuthAndContact 实名验证以及紧急联系人中间件
func RiderRequireAuthAndContact() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            url := c.Path()

            if riderAuthSkipper[url] {
                return next(c)
            }

            ctx := c.(*app.RiderContext)

            p, _ := ctx.Rider.QueryPerson().First(context.Background())
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
            url := c.Path()

            if riderFaceSkipper[url] {
                return next(c)
            }

            ctx := c.(*app.RiderContext)
            u := ctx.Rider
            s := service.NewRider()
            if s.IsNewDevice(u, ctx.Device) {
                // TODO 暂时跳过人脸校验
                // snag.Panic(snag.StatusLocked, ar.Map{"url": s.GetFaceUrl(ctx)})
            }
            return next(ctx)
        }
    }
}
