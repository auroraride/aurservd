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

    g.GET("/callback", rapi.Callback.RiderCallback) // 骑手api回调中心

    // 引入骑手api需要的中间件
    g.Use(middleware.DeviceMiddleware(), middleware.RiderMiddleware())

    g.POST("/signin", rapi.Rider.Signin)                  // 登录
    g.POST("/authenticator", rapi.Rider.Authenticator)    // 认证
    g.GET("/authenticator/:token", rapi.Rider.AuthResult) // 获取实名认证结果
    g.GET("/face/:token", rapi.Rider.FaceResult)          // 获取人脸验证结果

    // 检测是否需要实名验证以及补充紧急联系人
    g.Use(middleware.RiderRequireAuthAndContact())

    // 检测是否需要人脸识别
    g.Use(middleware.RiderFaceMiddleware())
    g.POST("/contact", rapi.Rider.Contact) // 编辑紧急联系人
}
