// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package common

import (
    "errors"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

var (
    // debug
    debugPhones = map[string]bool{
        "18501358308": true,
        "19381630638": true,
    }
)

type smsReq struct {
    Phone       string `json:"phone" validate:"required"`       // 手机号
    CaptchaCode string `json:"captchaCode" validate:"required"` // captcha 验证码
}

// SendSmsCode 发送短信验证码
func SendSmsCode(c echo.Context) error {
    ctx := c.(*app.Context)
    req := new(smsReq)
    ctx.BindValidate(req)
    id := ctx.Request().Header.Get(app.HeaderCaptchaID)
    if !debugPhones[req.Phone] && !service.NewCaptcha().Verify(id, req.CaptchaCode, false) {
        return errors.New("图形验证码校验失败")
    }
    // 发送短信
    smsId, err := service.NewSms().SendCode(req.Phone)
    if err != nil {
        return err
    }
    return app.NewResponse(c).SetData(map[string]string{"id": smsId}).Success().Send()
}
