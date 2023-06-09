// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-01
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
	"github.com/auroraride/adapter/app"

	"github.com/auroraride/aurservd/app/controller/v1/aapi"
	"github.com/auroraride/aurservd/app/middleware"
)

func loadAgentRoutes() {
	rawDump := app.NewDumpLoggerMiddleware().WithConfig(&app.DumpConfig{})

	g := root.Group("agent/v1")

	// 无须校验
	guide := g.Group("", middleware.Agent())
	guide.POST("/signin", aapi.Agent.Signin)

	// 需校验
	auth := g.Group("", middleware.Agent(), middleware.AgentAuth())
	auth.GET("/cabinet/:serial", aapi.Cabinet.Detail)
	auth.GET("/profile", aapi.Agent.Profile)
	auth.GET("/openid", aapi.Agent.GetOpenid)
	// 电池搜索列表
	auth.GET("/battery/list", aapi.Agent.BatteryList)
	// 电池型号列表
	auth.GET("/battery/model", aapi.Agent.BatteryModels)
	// 骑手列表
	auth.GET("/rider", aapi.Rider.List)
	// 骑手详情
	auth.GET("/rider/:id", aapi.Rider.Detail)
	// 激活骑手
	auth.POST("/rider/active", aapi.Rider.Active)
	// 添加骑手
	auth.POST("/rider", aapi.Rider.Create)
	// 通过二维码获取骑手信息
	auth.GET("/rider/info", aapi.Rider.RiderInfo)
	// 邀请骑手二维码
	auth.POST("/rider/invite", aapi.Rider.Invite)
	// 申请加时列表
	auth.GET("/subscribe/apply", aapi.Rider.SubscribeApplyList)
	// 审批加时
	auth.POST("/subscribe/apply", aapi.Rider.ReviewApply)
	// 换电记录列表
	auth.GET("/rider/exchange", aapi.Rider.ExchangeList)
	// 电柜列表
	auth.GET("/cabinet", aapi.Cabinet.List)
	// 筛选电柜
	auth.GET("/cabinet/filter", aapi.Cabinet.CabinetFilter)
	// 价格列表
	auth.GET("/price/list", aapi.Agent.PriceList)

	// 意见反馈
	auth.POST("/feedback", aapi.Agent.Feedback)
	// 上传图片
	auth.POST("/upload/image", aapi.Agent.UploadImage)

	auth.GET("/prepayment/overview", aapi.Prepayment.Overview)
	auth.GET("/prepayment", aapi.Prepayment.List)
	auth.POST("/prepayment/pay/wxmp", aapi.Prepayment.WechatMiniprogramPay, rawDump)

	auth.GET("/bill/usage", aapi.Bill.Usage)
	auth.GET("/bill/historical", aapi.Bill.Historical)

	auth.GET("/stock", aapi.Stock.Detail)

	// 首页汇总数据
	auth.GET("/index", aapi.Agent.Index)
}
