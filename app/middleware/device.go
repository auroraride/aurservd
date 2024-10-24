// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package middleware

import (
	"errors"
	"strings"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/labstack/echo/v4"
)

func DeviceMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			c := app.Context(ctx)
			sn := splitString(c.Request().Header.Get(app.HeaderDeviceSerial))
			dt := splitString(c.Request().Header.Get(app.HeaderDeviceType))
			if sn == "" || dt == "" {
				return errors.New("设备校验失败")
			}
			var err error
			c.Device, err = model.NewDevice(sn, dt)
			if err != nil {
				return err
			}
			return next(ctx)
		}
	}
}

func splitString(str string) (s string) {
	arr := strings.Split(str, ",")
	for _, s = range arr {
		if len(s) > 0 {
			return
		}
	}
	return
}
