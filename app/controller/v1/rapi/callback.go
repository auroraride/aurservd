// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/13
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import "github.com/labstack/echo/v4"

type callback struct {}

var Callback = new(callback)

func (*callback) RiderCallback(c echo.Context) error {
    return nil
}
