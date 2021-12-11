// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package common

import (
    "bytes"
    "errors"
    "github.com/auroraride/aurservd/internal/app/response"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/service"
    "github.com/labstack/echo/v4"
    "net/http"
)

// CaptchaGenerate 生成图片验证码
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

    c.Response().Header().Set(ar.HeaderCaptchaID, id)

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
    id := c.Request().Header.Get(ar.HeaderCaptchaID)
    if !service.NewCaptcha().Verify(id, r.Code, false) {
        return errors.New("验证码校验失败")
    }
    return response.New(c).Success().Send()
}