// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-01
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "github.com/auroraride/aurservd/app/controller/v1/aapi"
    "github.com/auroraride/aurservd/app/middleware"
)

func loadAgentRoutes() {
    g := root.Group("agent/v1")
    g.Use(middleware.AgentMiddleware())

    g.POST("/signin", aapi.Agent.Signin)
    g.GET("/profile", aapi.Agent.Profile)
}
