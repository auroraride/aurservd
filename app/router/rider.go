// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "github.com/auroraride/aurservd/app/controller/v1/rapi"
    "github.com/auroraride/aurservd/app/middleware"
)

func (r *router) rideRoute() {
    g := r.Group("/rider")
    g.Use(middleware.DeviceMiddleware)
    // 骑手登录
    g.POST("/signin", rapi.Rider.Signin)
}
