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
	"github.com/auroraride/aurservd/app/middleware"
	"github.com/auroraride/aurservd/internal/ent"
)

func LoadRiderV2Routes(root *echo.Group) {
	g := root.Group("rider/v2")

	rawDump := app.NewDumpLoggerMiddleware().WithConfig(&app.DumpConfig{
		RequestHeader:  true,
		ResponseHeader: true,
	})
	g.Any("/callback/fand/auth/freeze", rapi.Callback.AlipayFandAuthFreeze, rawDump) //  骑手支付宝资金授权冻结回调中心

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

	// 不需要登录的接口
	g.POST("/signin", v1.Rider.Signin) // 登录
	g.GET("/city", v1.City.List)       // 已开通城市

	g.GET("/branch", v1.Branch.List)                   // 网点列表
	g.GET("/branch/facility/:fid", v1.Branch.Facility) // 网点设施

	g.GET("/setting/version", rapi.Setting.LatestVersion) // App最新版本
	g.GET("/selection/model", rapi.Selection.Model)       // 电池型号选择

	g.GET("/guide", rapi.Guide.List)       // 新手引导
	g.GET("/activity", rapi.Activity.List) // 活动

	g.GET("/question/category", rapi.QuestionCategory.All) // 问题分类
	g.GET("/question", rapi.Question.All)                  // 问题列表

	g.GET("/instructions/:key", rapi.Instructions.Detail) // 买前必读 积分 优惠券使用说明

	// 骑手登录认证中间件
	auth := g.Group("", middleware.RiderAuthMiddlewareV2())

	auth.GET("/certification/ocr/client", rapi.Person.CertificationOcrClient)   // 获取人身核验OCR参数
	auth.GET("/certification/ocr/cloud", rapi.Person.CertificationOcrCloud)     // 获取阿里云OCR签名
	auth.POST("/certification/face", rapi.Person.CertificationFace)             // 提交身份信息并获取人脸核身参数
	auth.GET("/certification/face/result", rapi.Person.CertificationFaceResult) // 获取人脸核身结果

	auth.GET("/profile", v1.Rider.Profile)       // 获取用户信息
	auth.GET("/deposit", v1.Rider.Deposit)       // 获取押金信息
	auth.GET("/deregister", v1.Rider.Deregister) // 注销账户

	// 骑士卡
	auth.GET("/plan", v1.Plan.List)            // 套餐列表
	auth.GET("/plan/renewly", v1.Plan.Renewly) // 续费列表

	auth.GET("/cabinet", rapi.Cabinet.List)           // 电柜列表
	auth.GET("/cabinet/:serial", rapi.Cabinet.Detail) // 详情

	auth.GET("/battery", v1.Battery.Detail) // 电池详情

	// 救援
	auth.GET("/assistance/breakdown", v1.Assistance.Breakdown) // 获取救援原因
	auth.POST("/assistance", v1.Assistance.Create)             // 创建救援
	auth.POST("/assistance/cancel", v1.Assistance.Cancel)      // 取消救援
	auth.GET("/assistance/current", v1.Assistance.Current)     // 当前救援
	auth.GET("/assistance", v1.Assistance.List)                // 救援列表

	// 代理、团签
	auth.GET("/enterprise/battery", v1.Enterprise.Battery)
	auth.POST("/enterprise/subscribe", v1.Enterprise.Subscribe)
	auth.GET("/enterprise/subscribe", v1.Enterprise.SubscribeStatus)
	auth.POST("/enterprise/join", v1.Enterprise.JoinEnterprise)                    // 加入团签
	auth.GET("/enterprise/info", v1.Enterprise.RiderEnterpriseInfo)                // 小程序加入团签信息
	auth.POST("/enterprise/exit", v1.Enterprise.ExitEnterprise)                    // 退出团签
	auth.POST("/enterprise/subscribe/alter", v1.Enterprise.SubscribeAlter)         // 申请加时
	auth.GET("/enterprise/subscribe/alter/list", v1.Enterprise.SubscribeAlterList) // 申请列表

	auth.GET("/order", v1.Order.List)          // 订单列表
	auth.GET("/order/:id", v1.Order.Detail)    // 订单详情
	auth.GET("/order/status", v1.Order.Status) // 订单状态

	auth.GET("/exchange/overview", v1.Exchange.Overview) // 换电概览
	auth.GET("/exchange/log", v1.Exchange.Log)           // 换电记录

	// 设置
	auth.GET("/setting/app", v1.Setting.App)           // App设置
	auth.GET("/setting/question", v1.Setting.Question) // 问题分类

	// 钱包
	auth.GET("/wallet/overview", v1.Wallet.Overview) // 钱包概览
	auth.GET("/wallet/pointlog", v1.Wallet.PointLog) // 积分明细
	auth.GET("/wallet/points", v1.Wallet.Points)     // 积分详情
	auth.GET("/wallet/coupons", v1.Wallet.Coupons)   // 优惠券列表

	// 意见反馈
	auth.POST("/feedback", rapi.Feedback.Create) // 创建反馈
	auth.GET("/feedback", rapi.Feedback.List)    // 反馈列表

	// 地图
	auth.GET("/direction", rapi.Rider.Direction) // 获取地图路径规划

	// 故障上报
	auth.POST("/fault", rapi.Fault.Create)          // 故障上报
	auth.GET("/fault/cause", rapi.Fault.FaultCause) // 故障原因

	// 实人认证中间件（包含骑手登录认证）和联系方式认证
	cert := g.Group("", middleware.RiderCertificationMiddlewareV2(), middleware.RiderRequireAuthAndContactV2())
	// 合同
	contract := cert.Group("/contract")
	contract.POST("/sign", v1.Contract.Sign)     // 签署合同
	contract.GET("/:sn", v1.Contract.SignResult) // 合同签署结果

	// 业务
	cert.POST("/business/active", v1.Business.Active)                       // 激活骑士卡
	cert.POST("/business/unsubscribe", v1.Business.Unsubscribe)             // 退租
	cert.POST("/business/pause", v1.Business.Pause)                         // 寄存
	cert.POST("/business/continue", v1.Business.Continue)                   // 取消寄存
	cert.GET("/business/status", v1.Business.Status)                        // 业务状态
	cert.GET("/business/pause/info", v1.Business.PauseInfo)                 // 寄存信息
	cert.GET("/business/allocated/:id", v1.Business.Allocated)              // 长连接轮询是否已分配
	cert.GET("/business/subscribe/signed/:id", v1.Business.SubscribeSigned) // 连接轮询是否已签约

	// 订单
	cert.POST("/order", rapi.Order.Create)                 // 创建订单
	cert.POST("/order/refund", v1.Order.Refund)            // 申请退款
	cert.POST("/deposit/credit", rapi.Order.DepositCredit) // 押金订单

	// 电柜
	cabinet := cert.Group("/cabinet")
	cabinet.GET("/process/:serial", v1.Cabinet.GetProcess)   // 获取换电信息
	cabinet.POST("/process", v1.Cabinet.Process)             // 换电
	cabinet.GET("/process/status", v1.Cabinet.ProcessStatus) // 换电状态
	cabinet.POST("/exchange/store", v1.Exchange.Store)       // 门店换电

	// 预约
	cert.GET("/reserve", v1.Reserve.Unfinished)    // 未完成预约
	cert.POST("/reserve", v1.Reserve.Create)       // 创建预约
	cert.DELETE("/reserve/:id", v1.Reserve.Cancel) // 取消预约
}
