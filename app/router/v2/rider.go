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

	// 骑手登录认证中间件
	auth := middleware.RiderAuthMiddlewareV2

	// 实人认证中间件（包含骑手登录认证）
	// cert := middleware.RiderCertificationMiddlewareV2

	// 获取人身核验OCR参数
	g.GET("/certification/ocr/client", rapi.Person.CertificationOcrClient, auth())

	// 获取阿里云OCR签名
	g.GET("/certification/ocr/cloud", rapi.Person.CertificationOcrCloud, auth())

	// 提交身份信息并获取人脸核身参数
	g.POST("/certification/face", rapi.Person.CertificationFace, auth())

	// 获取人脸核身结果
	g.GET("/certification/face/result", rapi.Person.CertificationFaceResult, auth())

	g.POST("/signin", v1.Rider.Signin)          // 登录
	g.GET("/profile", v1.Rider.Profile, auth()) // 获取用户信息
	g.GET("/deposit", v1.Rider.Deposit, auth())
	g.GET("/deregister", v1.Rider.Deregister, auth()) // 注销账户

	// 已开通城市
	g.GET("/city", v1.City.List, auth())

	// 合同
	contract := g.Group("/contract")
	contract.POST("/sign", v1.Contract.Sign, auth())
	contract.GET("/:sn", v1.Contract.SignResult, auth())

	// 获取网点
	g.GET("/branch", v1.Branch.List, auth())
	g.GET("/branch/riding", v1.Branch.Riding, auth())
	g.GET("/branch/facility/:fid", v1.Branch.Facility, auth())

	// 业务
	g.GET("/plan", v1.Plan.List, auth())
	g.GET("/plan/renewly", v1.Plan.Renewly, auth())
	g.POST("/order", rapi.Order.Create, auth())
	g.POST("/order/refund", v1.Order.Refund, auth()) // 申请退款
	g.GET("/order", v1.Order.List, auth())
	g.GET("/order/:id", v1.Order.Detail, auth())
	g.GET("/enterprise/battery", v1.Enterprise.Battery, auth())
	g.POST("/enterprise/subscribe", v1.Enterprise.Subscribe, auth())
	g.GET("/enterprise/subscribe", v1.Enterprise.SubscribeStatus, auth())
	g.GET("/order/status", v1.Order.Status, auth())
	g.POST("/business/active", v1.Business.Active, auth())
	g.POST("/business/unsubscribe", v1.Business.Unsubscribe, auth())
	g.POST("/business/pause", v1.Business.Pause, auth())
	g.POST("/business/continue", v1.Business.Continue, auth())
	g.GET("/business/status", v1.Business.Status, auth())
	g.GET("/business/pause/info", v1.Business.PauseInfo, auth())
	g.GET("/business/allocated/:id", v1.Business.Allocated, auth())
	g.GET("/business/subscribe/signed/:id", v1.Business.SubscribeSigned, auth())

	// 代理、团签
	g.POST("/enterprise/join", v1.Enterprise.JoinEnterprise, auth())                    // 加入团签
	g.GET("/enterprise/info", v1.Enterprise.RiderEnterpriseInfo, auth())                // 小程序加入团签信息
	g.POST("/enterprise/exit", v1.Enterprise.ExitEnterprise, auth())                    // 退出团签
	g.POST("/enterprise/subscribe/alter", v1.Enterprise.SubscribeAlter, auth())         // 申请加时
	g.GET("/enterprise/subscribe/alter/list", v1.Enterprise.SubscribeAlterList, auth()) // 申请列表

	// 电柜
	cabinet := g.Group("/cabinet")
	cabinet.GET("/process/:serial", v1.Cabinet.GetProcess, auth())
	cabinet.POST("/process", v1.Cabinet.Process, auth())
	cabinet.GET("/process/status", v1.Cabinet.ProcessStatus, auth())
	cabinet.POST("/report", v1.Cabinet.Report, auth())
	cabinet.GET("/fault", v1.Cabinet.Fault, auth())

	g.POST("/exchange/store", v1.Exchange.Store, auth())
	g.GET("/exchange/overview", v1.Exchange.Overview, auth())
	g.GET("/exchange/log", v1.Exchange.Log, auth())

	// 救援
	g.GET("/assistance/breakdown", v1.Assistance.Breakdown, auth())
	g.POST("/assistance", v1.Assistance.Create, auth())
	g.POST("/assistance/cancel", v1.Assistance.Cancel, auth())
	g.GET("/assistance/current", v1.Assistance.Current, auth())
	g.GET("/assistance", v1.Assistance.List, auth())

	// 设定
	g.GET("/setting/app", v1.Setting.App, auth())
	g.GET("/setting/question", v1.Setting.Question, auth())

	// 预约
	g.GET("/reserve", v1.Reserve.Unfinished, auth())
	g.POST("/reserve", v1.Reserve.Create, auth())
	g.DELETE("/reserve/:id", v1.Reserve.Cancel, auth())

	// 钱包
	g.GET("/wallet/overview", v1.Wallet.Overview, auth())
	g.GET("/wallet/pointlog", v1.Wallet.PointLog, auth())
	g.GET("/wallet/points", v1.Wallet.Points, auth())
	g.GET("/wallet/coupons", v1.Wallet.Coupons, auth())

	// 电池
	g.GET("/battery", v1.Battery.Detail, auth())
	g.GET("/selection/model", rapi.Selection.Model) // 电池型号选择

	// 电柜列表
	g.GET("/cabinet", rapi.Cabinet.List, auth())
	g.GET("/cabinet/:serial", rapi.Cabinet.Detail, auth())

	// 押金
	g.POST("/deposit/free", rapi.Order.DepositFree, auth())

	// 新手引导
	g.GET("/guide", rapi.Guide.List, auth())

	// 活动
	g.GET("/activity", rapi.Activity.List, auth())

	// 骑手端意见反馈
	g.POST("/feedback", rapi.Feedback.Create, auth())

	// 版本
	g.GET("/version", rapi.Version.Latest)
}
