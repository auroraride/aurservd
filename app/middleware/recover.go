// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
    mw "github.com/labstack/echo/v4/middleware"
    log "github.com/sirupsen/logrus"
    "runtime/debug"
)

func Recover() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            defer func() {
                if r := recover(); r != nil {
                    switch r.(type) {
                    case *snag.Error:
                        log.Error(r)
                        c.Error(r.(*snag.Error))
                    case *ent.ValidationError:
                        log.Error(string(debug.Stack()))
                        c.Error(r.(*ent.ValidationError).Unwrap())
                    case error:
                        log.Error(string(debug.Stack()))
                        c.Error(r.(error))
                    default:
                        log.Error(string(debug.Stack()))
                        _ = mw.Recover()(next)(c)
                    }
                }
            }()
            return next(c)
        }
    }
}
