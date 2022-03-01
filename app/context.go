// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/labstack/echo/v4"
)

type Context struct {
    echo.Context

    Device *Device
}

// BindValidate 绑定并校验数据
func (c *Context) BindValidate(ptr interface{}) {
    err := c.Bind(ptr)
    if err != nil {
        snag.Panic(err)
    }
    err = c.Validate(ptr)
    if err != nil {
        snag.Panic(err)
    }
}
