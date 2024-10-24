// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/11
// Based on aurservd by liasica, magicrolan@qq.com.

package v1

import (
	"strconv"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/app"
	"github.com/labstack/echo/v4"

	inapp "github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/controller/v1/rapi"
	"github.com/auroraride/aurservd/app/middleware"
	"github.com/auroraride/aurservd/internal/ent"
)

func LoadRiderV1Routes(root *echo.Group) {

	g := root.Group("rider/v1")

	// socket
	g.Any("/socket", rapi.Socket.Rider)

	rawDump := app.NewDumpLoggerMiddleware().WithConfig(&app.DumpConfig{
		RequestHeader:  true,
		ResponseHeader: true,
	})
	g.Any("/callback", rapi.Callback.RiderCallback, rawDump)                         // 骑手api回调中心
	g.Any("/callback/esign", rapi.Callback.ESignCallback, rawDump)                   // esign回调中心
	g.Any("/callback/alipay", rapi.Callback.AlipayCallback, rawDump)                 // 骑手支付宝回调中心
	g.Any("/callback/wechatpay", rapi.Callback.WechatPayCallback, rawDump)           // 骑手微信支付回调中心
	g.Any("/callback/wechatpay/refund", rapi.Callback.WechatRefundCallback, rawDump) // 骑手微信退款回调中心

	// 引入骑手api需要的中间件
	dumpSkipPaths := map[string]bool{
		"/rider/v1/city":                   true,
		"/rider/v1/socket":                 true,
		"/rider/callback":                  true,
		"/rider/callback/esign":            true,
		"/rider/callback/alipay":           true,
		"/rider/callback/wechatpay":        true,
		"/rider/callback/wechatpay/refund": true,
		"/rider/v1/branch":                 true,
	}
	dumpReqHeaders := map[string]struct{}{
		inapp.HeaderCaptchaID:    {},
		inapp.HeaderDeviceSerial: {},
		inapp.HeaderDeviceType:   {},
		inapp.HeaderPushId:       {},
	}
	dump := app.NewDumpLoggerMiddleware().WithConfig(&app.DumpConfig{
		ResponseBodySkipper: func(c echo.Context) bool {
			return dumpSkipPaths[c.Path()]
		},
		RequestHeader: true,
		RequestHeaderSkipper: func(s string) bool {
			_, ok := dumpReqHeaders[s]
			return !ok
		},
		Extra: func(c echo.Context) []byte {
			if r, ok := c.Get("rider").(*ent.Rider); ok && r != nil {
				buf := adapter.NewBuffer()
				defer adapter.ReleaseBuffer(buf)

				// {"id":123,"phone":"1111","name":"test"}
				buf.WriteString(`{"id":`)
				buf.WriteString(strconv.FormatUint(r.ID, 10))
				buf.WriteString(`,"phone":"`)
				buf.WriteString(r.Phone)
				buf.WriteString(`","name":"`)
				buf.WriteString(r.Name)
				buf.WriteString(`"}`)

				return buf.Bytes()
			}
			return nil
		},
	})

	g.Use(
		middleware.DeviceMiddleware(),
		middleware.RiderMiddleware(),
		dump,
	)

	g.POST("/signin", rapi.Rider.Signin)                  // 登录
	g.POST("/authenticator", rapi.Rider.Authenticator)    // 认证
	g.GET("/authenticator/:token", rapi.Rider.AuthResult) // 获取实名认证结果
	g.GET("/face/:token", rapi.Rider.FaceResult)          // 获取人脸验证结果
	g.POST("/contact", rapi.Rider.Contact)                // 编辑紧急联系人

	// 检测是否需要实名验证以及补充紧急联系人
	g.Use(middleware.RiderRequireAuthAndContact())

	// 检测是否需要人脸识别
	g.Use(middleware.RiderFaceMiddleware())

	g.GET("/demo", rapi.Rider.Demo)       // 测试空白页面
	g.GET("/profile", rapi.Rider.Profile) // 获取用户信息
	g.GET("/deposit", rapi.Rider.Deposit)
	g.GET("/deregister", rapi.Rider.Deregister) // 注销账户

	// 已开通城市
	g.GET("/city", rapi.City.List)

	// 合同
	contract := g.Group("/contract")
	contract.POST("/sign", rapi.Contract.Sign)
	contract.GET("/:sn", rapi.Contract.SignResult)

	// 获取网点
	g.GET("/branch", rapi.Branch.List)
	g.GET("/branch/riding", rapi.Branch.Riding)
	g.GET("/branch/facility/:fid", rapi.Branch.Facility)

	// 业务
	g.GET("/plan", rapi.Plan.List)
	g.GET("/plan/renewly", rapi.Plan.Renewly)
	g.POST("/order", rapi.Order.Create)
	g.POST("/order/refund", rapi.Order.Refund) // 申请退款
	g.GET("/order", rapi.Order.List)
	g.GET("/order/:id", rapi.Order.Detail)
	g.GET("/enterprise/battery", rapi.Enterprise.Battery)
	g.POST("/enterprise/subscribe", rapi.Enterprise.Subscribe)
	g.GET("/enterprise/subscribe", rapi.Enterprise.SubscribeStatus)
	g.GET("/order/status", rapi.Order.Status)
	g.POST("/business/active", rapi.Business.Active)
	g.POST("/business/unsubscribe", rapi.Business.Unsubscribe)
	g.POST("/business/pause", rapi.Business.Pause)
	g.POST("/business/continue", rapi.Business.Continue)
	g.GET("/business/status", rapi.Business.Status)
	g.GET("/business/pause/info", rapi.Business.PauseInfo)
	g.GET("/business/allocated/:id", rapi.Business.Allocated)
	g.GET("/business/subscribe/signed/:id", rapi.Business.SubscribeSigned)

	// 代理、团签
	g.POST("/enterprise/join", rapi.Enterprise.JoinEnterprise)                    // 加入团签
	g.GET("/enterprise/info", rapi.Enterprise.RiderEnterpriseInfo)                // 小程序加入团签信息
	g.POST("/enterprise/exit", rapi.Enterprise.ExitEnterprise)                    // 退出团签
	g.POST("/enterprise/subscribe/alter", rapi.Enterprise.SubscribeAlter)         // 申请加时
	g.GET("/enterprise/subscribe/alter/list", rapi.Enterprise.SubscribeAlterList) // 申请列表

	// 电柜
	cabinet := g.Group("/cabinet")
	cabinet.GET("/process/:serial", rapi.Cabinet.GetProcess)
	cabinet.POST("/process", rapi.Cabinet.Process)
	cabinet.GET("/process/status", rapi.Cabinet.ProcessStatus)
	cabinet.POST("/report", rapi.Cabinet.Report)
	cabinet.GET("/fault", rapi.Cabinet.Fault)

	g.POST("/exchange/store", rapi.Exchange.Store)
	g.GET("/exchange/overview", rapi.Exchange.Overview)
	g.GET("/exchange/log", rapi.Exchange.Log)

	// 救援
	g.GET("/assistance/breakdown", rapi.Assistance.Breakdown)
	g.POST("/assistance", rapi.Assistance.Create)
	g.POST("/assistance/cancel", rapi.Assistance.Cancel)
	g.GET("/assistance/current", rapi.Assistance.Current)
	g.GET("/assistance", rapi.Assistance.List)

	// 设定
	g.GET("/setting/app", rapi.Setting.App)
	g.GET("/setting/question", rapi.Setting.Question)

	// 预约
	g.GET("/reserve", rapi.Reserve.Unfinished)
	g.POST("/reserve", rapi.Reserve.Create)
	g.DELETE("/reserve/:id", rapi.Reserve.Cancel)

	// 钱包
	g.GET("/wallet/overview", rapi.Wallet.Overview)
	g.GET("/wallet/pointlog", rapi.Wallet.PointLog)
	g.GET("/wallet/points", rapi.Wallet.Points)
	g.GET("/wallet/coupons", rapi.Wallet.Coupons)

	// 电池
	g.GET("/battery", rapi.Battery.Detail)
}
