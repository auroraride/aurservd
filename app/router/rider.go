// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "github.com/auroraride/aurservd/app/controller/v1/rapi"
    "github.com/auroraride/aurservd/app/middleware"
)

// rideRoute 骑手路由
func (r *router) rideRoute() {
    g := r.Group("/rider")
    // 引入骑手api需要的中间件
    g.Use(middleware.DeviceMiddleware(), middleware.RiderMiddleware())

    g.POST("/signin", rapi.Rider.Signin)               // 登录
    g.POST("/contact", rapi.Rider.Contact)             // 添加紧急联系人
    g.POST("/authenticator", rapi.Rider.Authenticator) // 认证
    // 更换设备扫脸
}
