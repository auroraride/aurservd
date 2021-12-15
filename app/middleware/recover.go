// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "github.com/auroraride/aurservd/app"
    "github.com/labstack/echo/v4"
    mw "github.com/labstack/echo/v4/middleware"
)

func Recover() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            defer func() {
                if r := recover(); r != nil {
                    err, ok := r.(*app.Error)
                    if ok {
                        c.Error(err)
                    } else {
                        _ = mw.Recover()(next)(c)
                    }
                }
            }()
            return next(c)
        }
    }
}
