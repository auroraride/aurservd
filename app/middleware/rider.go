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

func RiderMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(ctx echo.Context) error {
        signinMessage := "需要登录"
        // 获取骑手
        token := ctx.Request().Header.Get(app.HeaderRiderToken)
        id, err := ar.Cache.Get(context.Background(), token).Uint64()
        if err != nil {
            ctx.Set("errCode", response.StatusUnauthorized)
            return errors.New(signinMessage)
        }
        s := service.NewRider()
        var u *ent.Rider
        u, err = s.GetRiderById(id)
        if err != nil {
            ctx.Set("errCode", response.StatusUnauthorized)
            return errors.New(signinMessage)
        }

        // 重载context
        c := &app.RiderContext{
            Context: ctx.(*app.Context),
        }

        if c.Request().RequestURI != "/rider/authentication" {
            // 判断token权限
            perm := s.GetTokenPermission(u, c.Device)
            status, message := s.GetTokenPermissionResponse(perm)
            if status > 0 {
                ctx.Set("errCode", status)
                return errors.New(message)
            }
        }
        return next(c)
    }
}
