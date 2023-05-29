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

	// g.GET("/rider", aapi.Rider.List)
	// g.POST("/rider", aapi.Rider.Create)
	// g.POST("/rider/alter", aapi.Rider.Alter)
	// g.GET("/rider/:id", aapi.Rider.Detail)
	//
	// g.GET("/prepayment/overview", aapi.Prepayment.Overview)
	// g.GET("/prepayment", aapi.Prepayment.List)
	//
	// g.GET("/bill/usage", aapi.Bill.Usage)
}
