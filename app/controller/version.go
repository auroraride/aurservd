// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-30
// Based on aurservd by liasica, magicrolan@qq.com.

package controller

import (
    "fmt"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/labstack/echo/v4"
    "strings"
)

type version struct{}

var Version = new(version)

func (*version) Get(c echo.Context) (err error) {
    ctx := app.Context(c)
    plaform := ctx.QueryParam("plaform")
    a := ctx.QueryParam("app")
    m := ar.Map{
        "rider-android":    ar.Config.RiderApp.Android,
        "rider-ios":        ar.Config.RiderApp.IOS,
        "employee-android": ar.Config.EmployeeApp.Android,
        "employee-ios":     ar.Config.EmployeeApp.IOS,
    }
    if a == "" {
        a = "rider"
    }
    // TODO 读取单独文件
    key := fmt.Sprintf("%s-%s", strings.ToLower(a), strings.ToLower(plaform))
    return ctx.SendResponse(m[key])
}
