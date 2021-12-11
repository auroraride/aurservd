// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import "github.com/auroraride/aurservd/internal/controller/v1/common"

func (r *router) commonRoute() {
    g := r.Group("/common")
    g.GET("/captcha", common.CaptchaGenerate)
    g.POST("/captcha", common.CaptchaVerify)
}
