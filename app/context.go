// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
    "github.com/labstack/echo/v4"
)

// http headers
const (
    // HeaderCaptchaID 图片验证码ID
    HeaderCaptchaID    = "X-Captcha-Id"
    HeaderDeviceSerial = "X-Device-Serial"
    HeaderDeviceType   = "X-Device-Type"
)

type Context struct {
    echo.Context

    Device *Device
}

func (c *Context) BindValidate(ptr interface{}) error {
    err := c.Bind(ptr)
    if err != nil {
        return err
    }
    return c.Validate(ptr)
}
