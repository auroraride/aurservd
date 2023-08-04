package router

import (
	"github.com/auroraride/aurservd/app/controller/v1/papi"
	"github.com/auroraride/aurservd/app/middleware"
)

func loadPromotionRoutes() {
	g := root.Group("promotion/v1")

	// 无须校验
	guide := g.Group("", middleware.Promotion())

	guide.POST("/member/signin", papi.Member.Signin)                       // 登录
	guide.POST("/member/signup", papi.Member.Signup)                       // 注册
	guide.GET("/commission/rule", papi.PromotionCommission.CommissionRule) // 获取会员分佣规则

	auth := g.Group("", middleware.Promotion(), middleware.PromotionAuth())

	// 会员
	auth.GET("/member/profile", papi.Member.Profile)          // 会员个人资料
	auth.GET("/member/share/qrcode", papi.Member.ShareQrcode) // 获取推广二维码
	auth.POST("/member/avatar", papi.Member.UpdateAvatar)     // 修改头像
	auth.GET("/member/team", papi.Member.Team)                // 我的团队列表

	auth.POST("/bank/card", papi.PromotionBankCard.Create)       // 创建银行卡
	auth.PUT("/bank/card/:id", papi.PromotionBankCard.Update)    // 修改银行卡默认状态
	auth.GET("/bank/card", papi.PromotionBankCard.List)          // 获取银行卡列表
	auth.DELETE("/bank/card/:id", papi.PromotionBankCard.Delete) // 删除银行卡

	auth.POST("/auth/realname", papi.PromotionPerson.RealName) // 实名认证

	auth.GET("/withdrawal", papi.PromotionWithdrawal.List)                        // 会员提现列表
	auth.POST("/withdrawal/alter", papi.PromotionWithdrawal.Alter)                // 申请提现
	auth.POST("/withdrawal/fee", papi.PromotionWithdrawal.CalculateWithdrawalFee) // 计算提现手续费

	auth.GET("/earnings", papi.PromotionEarnings.List) // 会员收益列表

	// 设置
	auth.GET("/setting/:key", papi.PromotionSetting.Setting) // 获取设置列表

	// 统计
	auth.GET("/statistics/team/overview", papi.Statistics.TemaOverview)         // 首页我的团队统计
	auth.GET("/statistics/earnings/overview", papi.Statistics.EarningsOverview) // 首页我的收益统计
	auth.GET("/statistics/record/overview", papi.Statistics.RecordOverview)     // 首页邀请战绩统计
	auth.GET("/statistics/wallet/overview", papi.Statistics.WalletOverview)     // 我的钱包-统计
	auth.GET("/statistics/earnings/detail", papi.Statistics.EarningsDetail)     // 我的钱包-收益明细统计
	auth.GET("/statistics/team", papi.Statistics.Team)                          // 我的团队-统计
	auth.GET("/statistics/team/growth", papi.Statistics.TeamGrowth)             // 我的团队-增长趋势
}
