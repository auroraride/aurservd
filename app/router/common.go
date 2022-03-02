// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "github.com/auroraride/aurservd/app/controller/v1/common"
)

func (r *router) commonRoutes() {
    g := r.Group("/common")

    g.GET("/captcha", common.CaptchaGenerate)

    g.POST("/captcha", common.CaptchaVerify)

    g.POST("/sms", common.SendSmsCode)

    g.GET("/oss/token", common.Oss.Token)

    g.Static("/demo", "public")
}
