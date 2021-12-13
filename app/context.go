// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/labstack/echo/v4"
)

// http headers
const (
    // HeaderCaptchaID 图片验证码ID
    HeaderCaptchaID = "X-Captcha-Id"
    // HeaderDeviceSerial 骑手设备序列号 (由此判定是否更换了设备)
    HeaderDeviceSerial = "X-Device-Serial"
    // HeaderDeviceType 骑手设备类型
    HeaderDeviceType = "X-Device-Type"
    // HeaderRiderToken 骑手token
    HeaderRiderToken = "X-Rider-Token"
)

type GlobalContext struct {
    echo.Context

    Device *Device
}

type RiderContext struct {
    *GlobalContext

    Rider *ent.Rider
}

// BindValidate 绑定并校验数据
func (c *GlobalContext) BindValidate(ptr interface{}) {
    err := c.Bind(ptr)
    if err != nil {
        panic(response.NewError(err))
    }
    err = c.Validate(ptr)
    if err != nil {
        panic(response.NewError(err))
    }
}
