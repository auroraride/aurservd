// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package common

import (
    "bytes"
    "errors"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
    "net/http"
)

// CaptchaGenerate
// @ID           CaptchaGenerate
// @Router       /common/captcha [GET]
// @Summary      C1.生成图片验证码
// @Description  生成的图片验证码有效时间为10分钟
// @Tags         [C]公共接口
// @Accept       png
// @Produce      png
// @Success      200  {string}  string  "ok"
// @Header       200  {string}  X-Captcha-Id  true  "Captcha验证码ID"
func CaptchaGenerate(c echo.Context) error {
    id, item, err := service.NewCaptcha().DrawCaptcha()
    if err != nil {
        return err
    }

    b := new(bytes.Buffer)
    _, err = item.WriteTo(b)
    if err != nil {
        return err
    }

    c.Response().Header().Set(app.HeaderCaptchaID, id)

    return c.Stream(http.StatusOK, "image/png", b)
}

type captchaReq struct {
    Code string `json:"code"`
}

// CaptchaVerify 验证
func CaptchaVerify(c echo.Context) error {
    r := new(captchaReq)
    err := c.Bind(r)
    if err != nil {
        return err
    }
    id := c.Request().Header.Get(app.HeaderCaptchaID)
    if !service.NewCaptcha().Verify(id, r.Code, true) {
        return errors.New("验证码校验失败")
    }
    return app.NewResponse(c).Success().Send()
}