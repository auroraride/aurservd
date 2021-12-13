// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "context"
    "errors"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/labstack/echo/v4"
)

func RiderMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(ctx echo.Context) error {
            if ctx.Request().RequestURI == "/rider/signin" {
                return next(ctx)
            }

            // 获取骑手, 判定是否需要登录
            signinMessage := "需要登录"
            token := ctx.Request().Header.Get(app.HeaderRiderToken)
            id, err := ar.Cache.Get(context.Background(), token).Uint64()
            if err != nil {
                ctx.Set("errCode", response.StatusUnauthorized)
                return errors.New(signinMessage)
            }
            s := service.NewRider()
            var u *ent.Rider
            u, err = s.GetRiderById(id)
            if err != nil || u == nil {
                ctx.Set("errCode", response.StatusUnauthorized)
                return errors.New(signinMessage)
            }

            // 重载context
            c := &app.RiderContext{
                GlobalContext: ctx.(*app.GlobalContext),
                Rider:         u,
            }

            // TODO 判定骑手是否新设备

            // TODO 判断是否需要认证和补充联系人
            return next(c)
        }
    }
}
