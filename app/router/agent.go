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

	// A1 用户
	auth.GET("/profile", aapi.Agent.Profile)  // 个人资料
	auth.GET("/openid", aapi.Agent.GetOpenid) // 获取微信openid

	// A2 骑手
	auth.GET("/rider", aapi.Rider.List)               // 列表
	auth.GET("/rider/:id", aapi.Rider.Detail)         // 骑手详情
	auth.POST("/rider", aapi.Rider.Create)            // 添加骑手
	auth.GET("/rider/info", aapi.Rider.RiderInfo)     // 通过二维码获取骑手信息
	auth.POST("/rider/invite", aapi.Rider.Invite)     // 邀请骑手二维码
	auth.POST("/rider/alter", aapi.Rider.Alter)       // 增加/减少骑手时长
	auth.POST("/rider/reactive", aapi.Rider.Reactive) // 重新激活骑手

	// A3 账户
	auth.GET("/prepayment/overview", aapi.Prepayment.Overview)
	auth.GET("/prepayment", aapi.Prepayment.List)
	auth.POST("/prepayment/pay/wxmp", aapi.Prepayment.WechatMiniprogramPay, rawDump)

	// A4 账单
	auth.GET("/bill/usage", aapi.Bill.Usage)
	auth.GET("/bill/historical", aapi.Bill.Historical)

	// A5 电柜
	auth.GET("/cabinet", aapi.Cabinet.List)                  // 电柜列表
	auth.GET("/cabinet/detail/:serial", aapi.Cabinet.Detail) // 电柜详情
	auth.GET("/cabinet/section", aapi.Cabinet.Section)       // 电柜选择
	auth.POST("/cabinet/maintain", aapi.Cabinet.Maintain)    // 电柜维护
	auth.POST("/cabinet/binopen", aapi.Cabinet.BinOpen)      // 电柜开仓

	// A6 物资
	auth.GET("/stock", aapi.Stock.Detail)               // 出入库详情
	auth.GET("/stock/battery", aapi.Stock.BatteryStock) //  电池物资
	auth.GET("/stock/ebike", aapi.Stock.EBikeStock)     // 电车物资

	// A7 骑士卡 / 订阅
	auth.POST("/subscribe/active", aapi.Subscribe.Active)            // 激活骑手
	auth.GET("/subscribe/alter", aapi.Subscribe.AlterList)           // 申请加时列表
	auth.POST("/subscribe/alter/review", aapi.Subscribe.AlterReivew) // 审批加时
	auth.POST("/subscribe/halt", aapi.Subscribe.Halt)                // 强制退租

	// A8 业务
	auth.GET("/business/exchange", aapi.Business.Exchange) // 换电列表
	auth.GET("/business/price", aapi.Business.Price)       // 价格列表

	// A9 统计
	auth.GET("/statistics/overview", aapi.Statistics.Overview)

	// AA 电池
	auth.GET("/battery", aapi.Battery.List)                // 电池列表
	auth.GET("/battery/selection", aapi.Battery.Selection) // 电池搜索
	auth.GET("/battery/model", aapi.Battery.Model)         // 电池型号列表

	// AB 电车
	auth.GET("/bike", aapi.Bike.List)                    // 电车列表
	auth.GET("/bike/unallocated", aapi.Bike.Unallocated) // 搜索未分配车辆

	// AZ 杂项
	auth.POST("/misc/feedback", aapi.Misc.Feedback)            // 意见反馈
	auth.POST("/misc/feedback/image", aapi.Misc.FeedbackImage) // 意见反馈 - 上传图片
}
