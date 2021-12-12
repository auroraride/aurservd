// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
    "errors"
    "github.com/auroraride/aurservd/app"
    "github.com/labstack/echo/v4"
)

func DeviceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(ctx echo.Context) error {
        c := ctx.(*app.Context)
        sn := c.Request().Header.Get(app.HeaderDeviceSerial)
        dt := c.Request().Header.Get(app.HeaderDeviceType)
        if sn == "" || dt == "" {
            return errors.New("设备校验失败")
        }
        var err error
        c.Device, err = app.NewDevice(sn, dt)
        if err != nil {
            return err
        }
        return next(ctx)
    }
}
