// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package common

import (
    "errors"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/response"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type smsReq struct {
    Phone       string `json:"phone" validate:"required"`       // 手机号
    CaptchaCode string `json:"captchaCode" validate:"required"` // captcha 验证码
}

// SendSmsCode 发送短信验证码
func SendSmsCode(ctx echo.Context) error {
    c := ctx.(*app.Context)
    r := new(smsReq)
    err := c.BindValidate(r)
    if err != nil {
        return err
    }
    id := c.Request().Header.Get(app.HeaderCaptchaID)
    if !service.NewCaptcha().Verify(id, r.CaptchaCode, false) {
        return errors.New("图形验证码校验失败")
    }
    // 发送短信
    var smsId string
    smsId, err = service.NewSms().SendCode(r.Phone)
    return response.New(c).SetData(map[string]string{"id": smsId}).Success().Send()
}
