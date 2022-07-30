// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "fmt"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
    "runtime/debug"
)

func Recover() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            defer func() {
                if r := recover(); r != nil {
                    stack := string(debug.Stack())
                    switch r.(type) {
                    case *snag.Error:
                        c.Error(r.(*snag.Error))
                        break
                    case *ent.ValidationError:
                        log.Error(stack)
                        c.Error(r.(*ent.ValidationError).Unwrap())
                        break
                    case error:
                        log.Error(stack)
                        c.Error(r.(error))
                        break
                    default:
                        log.Error(stack)
                        // _ = mw.Recover()(next)(c)
                        c.Error(fmt.Errorf("%v", r))
                        break
                    }
                }
            }()
            return next(c)
        }
    }
}
