// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-10
// Based on aurservd by liasica, magicrolan@qq.com.

package v2

import (
	"strconv"

	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/app"
	"github.com/labstack/echo/v4"

	inapp "github.com/auroraride/aurservd/app"
	v1 "github.com/auroraride/aurservd/app/controller/v1/rapi"
	"github.com/auroraride/aurservd/app/controller/v2/rapi"
	pr "github.com/auroraride/aurservd/app/purchase/controller/rapi"

	"github.com/auroraride/aurservd/app/middleware"
	"github.com/auroraride/aurservd/internal/ent"
)

func LoadRiderV2Routes(root *echo.Group) {
	g := root.Group("rider/v2")

	rawDump := app.NewDumpLoggerMiddleware().WithConfig(&app.DumpConfig{
		RequestHeader:  true,
		ResponseHeader: true,
	})
	g.Any("/callback/alipay/auth/freeze", rapi.Callback.AlipayFandAuthFreeze, rawDump)     //  骑手支付宝资金授权冻结回调
	g.Any("/callback/alipay/auth/unfreeze", rapi.Callback.AlipayFandAuthUnfreeze, rawDump) //  骑手支付宝资金授权解冻回调
	g.Any("/callback/alipay/trade/pay", rapi.Callback.AlipayTradePay, rawDump)             // 冻结转支付回调

	g.Any("/callback/alipay/mini/pay", rapi.Callback.AlipayMiniProgramPay, rawDump) // 支付宝小程序支付回调

	// 购买商品回调
	g.Any("/callback/purchase/alipay", pr.Callback.PurchaseAlipay, rawDump)    // 支付宝购买商品回调
	g.Any("/callback/purchase/wechatpay", pr.Callback.PurchaseWechat, rawDump) // 微信购买商品回调

	// 记录请求日志
	dumpSkipPaths := map[string]bool{}
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
		middleware.RiderMiddlewareV2(),
		dump,
	)

	// 中间件
	logged := middleware.RiderAuthMiddlewareV2          // 骑手登录认证中间件
	person := middleware.RiderCertificationMiddlewareV2 // 实人认证中间件，包含登录认证

	// 不需要登录的接口
	g.POST("/signin", rapi.Rider.Signin) // 登录
	g.GET("/city", v1.City.List)         // 已开通城市

	g.GET("/branch", rapi.Branch.List)                 // 网点列表
	g.GET("/branch/facility/:fid", v1.Branch.Facility) // 网点设施

	g.GET("/setting/version", rapi.Setting.LatestVersion) // App最新版本
	g.GET("/selection/model", rapi.Selection.Model)       // 电池型号选择
	g.GET("/selection/brand", rapi.Selection.Brand)       // 电车型号选择

	g.GET("/activity", rapi.Activity.List) // 活动

	g.GET("/question/category", rapi.QuestionCategory.All) // 问题分类
	g.GET("/question", rapi.Question.All)                  // 问题列表

	g.GET("/instructions/:key", rapi.Instructions.Detail) // 买前必读 积分 优惠券使用说明

	g.GET("/mini/openid", rapi.Rider.GetOpenid) // 获取openid

	// 电柜
	cabinet := g.Group("/cabinet")
	cabinet.GET("", rapi.Cabinet.List)           // 电柜列表
	cabinet.GET("/:serial", rapi.Cabinet.Detail) // 详情

	// 门店
	store := g.Group("/store")
	store.GET("", rapi.Store.List)                                     // 门店列表
	store.GET("/:id", rapi.Store.Detail)                               // 门店详情
	store.GET("/subscribe/:id", rapi.Store.StoreBySubscribe, logged()) // 根据订阅查询骑手激活门店信息

	// 骑手登录认证中间件
	// auth := g.Group("", middleware.RiderAuthMiddlewareV2())
	certification := g.Group("/certification", logged())

	certification.GET("/ocr/client", rapi.Person.CertificationOcrClient)   // 获取人身核验OCR参数
	certification.GET("/ocr/cloud", rapi.Person.CertificationOcrCloud)     // 获取阿里云OCR签名
	certification.POST("/face", rapi.Person.CertificationFace)             // 提交身份信息并获取人脸核身参数
	certification.GET("/face/result", rapi.Person.CertificationFaceResult) // 获取人脸核身结果
	certification.POST("/supplement", rapi.Person.CertificationSupplement) // 补充实名信息

	// 骑手
	g.GET("/profile", rapi.Rider.Profile, logged())                        // 获取用户信息
	g.GET("/deposit", v1.Rider.Deposit, logged())                          // 获取押金信息
	g.DELETE("/deregister", v1.Rider.Deregister, logged())                 // 注销账户
	g.POST("/change/phone", rapi.Rider.ChangePhone, logged())              // 修改手机号
	g.POST("/contact", v1.Rider.Contact, logged())                         // 编辑紧急联系人
	g.POST("/mobpush", rapi.Rider.SetMobPushId, logged())                  // 绑定骑手推送ID
	g.POST("/report/phone/device", rapi.Rider.ReportPhoneDevice, logged()) // 上报手机设备信息

	// 骑士卡
	plan := g.Group("/plan")
	plan.GET("", rapi.Plan.List)                       // 套餐列表
	plan.GET("/renewly", v1.Plan.Renewly, person())    // 续费列表
	plan.GET("/:id", rapi.Plan.Detail)                 // 套餐详情
	plan.GET("/store", rapi.Plan.ByStore)              // 门店套餐列表
	plan.GET("/store/detail", rapi.Plan.ByStoreDetail) // 门店套餐详情

	// 电池
	battery := g.Group("/battery", person())
	battery.GET("", v1.Battery.Detail) // 电池详情

	// 救援
	assistance := g.Group("/assistance", person())
	assistance.GET("/breakdown", v1.Assistance.Breakdown) // 获取救援原因
	assistance.POST("", v1.Assistance.Create)             // 创建救援
	assistance.POST("/cancel", v1.Assistance.Cancel)      // 取消救援
	assistance.GET("/current", v1.Assistance.Current)     // 当前救援
	assistance.GET("", v1.Assistance.List)                // 救援列表

	// 代理、团签
	enterprise := g.Group("/enterprise", person())
	enterprise.GET("/battery", v1.Enterprise.Battery)                         // 获取可用电池
	enterprise.POST("/subscribe", v1.Enterprise.Subscribe)                    // 选择电池
	enterprise.GET("/subscribe", v1.Enterprise.SubscribeStatus)               // 订阅激活状态
	enterprise.POST("/join", v1.Enterprise.JoinEnterprise)                    // 加入团签
	enterprise.GET("/info", v1.Enterprise.RiderEnterpriseInfo)                // 小程序加入团签信息
	enterprise.POST("/exit", v1.Enterprise.ExitEnterprise)                    // 退出团签
	enterprise.POST("/subscribe/alter", v1.Enterprise.SubscribeAlter)         // 申请加时
	enterprise.GET("/subscribe/alter/list", v1.Enterprise.SubscribeAlterList) // 申请列表

	// 订单
	order := g.Group("/order", person())
	order.GET("", v1.Order.List)                            // 订单列表
	order.GET("/:id", v1.Order.Detail)                      // 订单详情
	order.POST("", rapi.Order.Create)                       // 创建订单
	order.POST("/refund", rapi.Order.Refund)                // 申请退款
	order.POST("/deposit/credit", rapi.Order.DepositCredit) // 押金订单

	// 设置
	setting := g.Group("/setting")
	setting.GET("/app", v1.Setting.App)           // App设置
	setting.GET("/question", v1.Setting.Question) // 问题分类

	// 钱包
	wallet := g.Group("/wallet", logged())
	wallet.GET("/overview", v1.Wallet.Overview) // 钱包概览
	wallet.GET("/pointlog", v1.Wallet.PointLog) // 积分明细
	wallet.GET("/points", v1.Wallet.Points)     // 积分详情
	wallet.GET("/coupons", v1.Wallet.Coupons)   // 优惠券列表

	// 意见反馈
	feedback := g.Group("/feedback", logged())
	feedback.POST("", rapi.Feedback.Create) // 创建反馈
	feedback.GET("", rapi.Feedback.List)    // 反馈列表

	// 故障上报
	fault := g.Group("/fault", logged())
	fault.POST("", rapi.Fault.Create)          // 故障上报
	fault.GET("/cause", rapi.Fault.FaultCause) // 故障原因

	// 合同
	contract := g.Group("/contract", person())
	contract.POST("/sign", rapi.Contract.Sign)     // 签署合同
	contract.POST("/create", rapi.Contract.Create) // 创建合同
	contract.GET("/:docId", rapi.Contract.Detail)  // 查看合同
	// 业务
	business := g.Group("/business", person())
	business.POST("/active", rapi.Business.Active)         // 激活骑士卡
	business.POST("/unsubscribe", v1.Business.Unsubscribe) // 退租
	business.POST("/pause", v1.Business.Pause)             // 寄存
	business.POST("/continue", v1.Business.Continue)       // 取消寄存
	business.GET("/status", v1.Business.Status)            // 业务状态
	business.GET("/pause/info", v1.Business.PauseInfo)     // 寄存信息

	// 换电
	g.GET("/exchange/overview", v1.Exchange.Overview, logged())           // 换电概览
	g.POST("/exchange/store", v1.Exchange.Store, logged())                // 门店换电
	g.GET("/exchange", v1.Exchange.Log, logged())                         // 换电记录
	g.GET("/exchange/process/:serial", v1.Cabinet.GetProcess, person())   // 电柜换电 - 获取换电信息
	g.POST("/exchange/process", v1.Cabinet.Process, person())             // 电柜换电 - 开始流程
	g.GET("/exchange/process/status", v1.Cabinet.ProcessStatus, person()) // 电柜换电 - 获取流程状态

	// 预约
	reserve := g.Group("/reserve", person())
	reserve.GET("", v1.Reserve.Unfinished)    // 未完成预约
	reserve.POST("", v1.Reserve.Create)       // 创建预约
	reserve.DELETE("/:id", v1.Reserve.Cancel) // 取消预约

	// 订阅
	subscribe := g.Group("/subscribe", person())
	subscribe.PUT("/store", rapi.Subscribe.StoreModify)      // 车电套餐修改激活门店
	subscribe.GET("/status", rapi.Subscribe.SubscribeStatus) // 查询订阅是否激活

	// 商品
	goods := g.Group("/goods")
	goods.GET("", rapi.Goods.List)       // 商品列表
	goods.GET("/:id", rapi.Goods.Detail) // 商品详情

	// 协议
	agreement := g.Group("/agreement")
	agreement.GET("/enterprise/price/:id", rapi.Agreement.QueryAgreementByEnterprisePriceID) // 根据企业价格ID查询协议

	// 车电品牌
	g.GET("/ebike/brand/:id", rapi.Ebike.EbikeBrandDetail) // 车电品牌详情

	// 买车
	p := g.Group("/purchase", person())
	p.GET("/order", pr.Order.List)       // 订单列表
	p.GET("/order/:id", pr.Order.Detail) // 订单详情
	p.POST("/order", pr.Order.Create)    // 创建订单
	p.POST("/pay", pr.Payment.Pay)       // 支付

	// 买车合同
	c := p.Group("/contract")
	c.POST("/sign", pr.Contract.Sign)     // 签署合同
	c.POST("/create", pr.Contract.Create) // 创建合同
	c.GET("/:docId", pr.Contract.Detail)  // 查看合同
}
