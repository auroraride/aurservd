// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package common

import (
    "errors"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/labstack/echo/v4"
    "github.com/lithammer/shortuuid/v4"
)

// SendSmsCode
// @ID           SendSmsCode
// @Router       /commom/sms [POST]
// @Summary      C2.发送短信验证码
// @Description  上传文件必须，单次获取有效时间为1个小时
// @Tags         [C]公共接口
// @Param        body  body  model.SmsReq  true  "请求参数"
// @Param        X-Captcha-Id  header  string  true  "Captcha验证码ID"
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.SmsRes  "请求成功"
func SendSmsCode(c echo.Context) error {
    ctx, req := app.ContextBinding[model.SmsReq](c)
    id := ctx.Request().Header.Get(app.HeaderCaptchaID)
    var smsId string
    var err error
    debugPhones := ar.Config.App.Debug.Phone

    if !debugPhones[req.Phone] && !service.NewCaptcha().Verify(id, req.CaptchaCode, false) {
        return errors.New("图形验证码校验失败")
    }
    if debugPhones[req.Phone] {
        smsId = shortuuid.New()
    } else {
        // 发送短信
        smsId, err = service.NewSms().SendCode(req.Phone)
        if err != nil {
            return err
        }
    }

    return ctx.SendResponse(model.SmsResponse{Id: smsId})
}
