// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
    "github.com/labstack/echo/v4"
)

type Context struct {
    echo.Context
}

func (c *Context) BindValidate(ptr interface{}) error {
    err := c.Bind(ptr)
    if err != nil {
        return err
    }
    return c.Validate(ptr)
}
