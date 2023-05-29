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

	// 无须校验
	guide := g.Group("", middleware.Agent())
	guide.POST("/signin", aapi.Agent.Signin)

	// 需校验
	auth := g.Group("", middleware.Agent(), middleware.AgentAuth())
	auth.GET("/cabinet/:serial", aapi.Cabinet.Detail)
	auth.GET("/profile", aapi.Agent.Profile)

	// auth.GET("/rider", aapi.Rider.List)
	// auth.POST("/rider", aapi.Rider.Create)
	// auth.POST("/rider/alter", aapi.Rider.Alter)
	// auth.GET("/rider/:id", aapi.Rider.Detail)
	//
	// auth.GET("/prepayment/overview", aapi.Prepayment.Overview)
	// auth.GET("/prepayment", aapi.Prepayment.List)

	auth.GET("/bill/usage", aapi.Bill.Usage)
	auth.GET("/bill/historical", aapi.Bill.Historical)
}
