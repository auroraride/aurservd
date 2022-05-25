// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "github.com/auroraride/aurservd/app/controller/v1/common"
)

func loadCommonRoutes() {
    g := root.Group("common")

    g.GET("/captcha", common.CaptchaGenerate)

    g.POST("/captcha", common.CaptchaVerify)

    g.POST("/sms", common.SendSmsCode)

    g.GET("/oss/token", common.Oss.Token)

    g.GET("/basic", common.Basic.Get)

    g.Static("/demo", "public")
}
