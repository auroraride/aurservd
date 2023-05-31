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

	// 骑手列表
	auth.GET("/rider", aapi.Rider.List)
	// 骑手详情
	auth.GET("/rider/:id", aapi.Rider.Detail)
	// 站点列表
	auth.GET("/site/list", aapi.Agent.SiteList)
	// 城市列表
	auth.GET("/city/list", aapi.Agent.CityList)
	// 电池搜索列表
	auth.GET("/battery/list", aapi.Agent.BatteryList)
	// 激活骑手
	auth.POST("/rider/activate", aapi.Rider.Active)
	// 添加骑手
	auth.POST("/rider", aapi.Rider.Create)
	// 邀请骑手
	auth.POST("/rider/invite", aapi.Rider.Invite)
	// 申请加时列表
	auth.GET("/subscribe/apply", aapi.Rider.SubscribeApplyList)
	// 审批加时
	auth.POST("/subscribe/apply", aapi.Rider.ReviewApply)
	// 换电记录列表
	auth.GET("/rider/exchange", aapi.Rider.ExchangeList)
	// 电柜列表
	auth.GET("/cabinet", aapi.Cabinet.List)
	// 出入库记录
	auth.GET("/stock", aapi.Stock.StockList)
	// 出入库详情
	auth.GET("/stock/:id", aapi.Stock.StockDetail)
	// 意见反馈
	auth.POST("/feedback", aapi.Agent.Feedback)
	// 上传图片
	auth.POST("/uploadImage", aapi.Agent.UploadImage)

	// g.POST("/rider/alter", aapi.Rider.Alter)
	// g.GET("/rider/:id", aapi.Rider.Detail)
	//
	// g.GET("/prepayment/overview", aapi.Prepayment.Overview)
	// g.GET("/prepayment", aapi.Prepayment.List)
	//
	// g.GET("/bill/usage", aapi.Bill.Usage)
}
