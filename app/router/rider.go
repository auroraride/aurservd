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
    g.Use(middleware.DeviceMiddleware)
    // 登录
    g.POST("/signin", rapi.Rider.Signin)
    // 引入骑手认证中间件
    g.Use(middleware.RiderMiddleware)
    // 认证
    g.POST("/authentication", rapi.Rider.Authentication)
    // 更换设备扫脸
}
