// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-08
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "context"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
)

var (
    employeeAuthSkipper = map[string]bool{
        "/employee/v1/signin": true,
        "/employee/v1/speech": true,
    }
)

func EmployeeMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            url := c.Request().URL.Path
            var emr *ent.Employee
            // 获取employee
            if !employeeAuthSkipper[url] {
                // 获取骑手, 判定是否需要登录
                token := c.Request().Header.Get(app.HeaderEmployeeToken)
                id, err := cache.Get(context.Background(), token).Uint64()
                if err != nil {
                    snag.Panic(snag.StatusUnauthorized)
                }

                s := service.NewEmployee()
                emr, err = s.GetEmployeeByID(id)
                if err != nil || emr == nil {
                    snag.Panic(snag.StatusUnauthorized)
                }
                if !emr.Enable {
                    s.Signout(emr)
                    snag.Panic(snag.StatusForbidden, ar.BannedMessage)
                }
            }

            // 重载context
            return next(app.NewEmployeeContext(c, emr))
        }
    }
}

func EmployeeDutyMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            ctx := app.ContextX[app.EmployeeContext](c)
            if ctx.Employee.Edges.Store == nil {
                snag.Panic("未上班, 无法处理业务")
            }
            return next(ctx)
        }
    }
}
