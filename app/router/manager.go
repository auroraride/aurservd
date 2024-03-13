// Copyright (C) liasica. 2021-present.
//
// Created at 2022/2/25
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
	"github.com/auroraride/aurservd/app/controller/v1/mapi"
	"github.com/auroraride/aurservd/app/middleware"
)

func loadManagerRoutes() {
	g := root.Group("manager/v1")

	g.POST("/user/signin", mapi.Manager.Signin) // 登录

	g.Use(middleware.ManagerMiddleware())
	g.POST("/user", mapi.Manager.Create)
	g.GET("/user", mapi.Manager.List)
	g.DELETE("/user/:id", mapi.Manager.Delete)
	g.PUT("/user/:id", mapi.Manager.Modify)
	g.GET("/user/profile", mapi.Manager.Profile)

	// 设置
	g.GET("/setting", mapi.Setting.List)
	g.PUT("/setting/:key", mapi.Setting.Modify)
	g.POST("/inventory", mapi.Inventory.CreateOrModify)
	g.GET("/inventory", mapi.Inventory.List)
	g.DELETE("/inventory", mapi.Inventory.Delete)
	g.GET("/inventory/transferable", mapi.Inventory.Transferable)
	g.GET("/setting/legal/:name", mapi.Setting.LegalRead)
	g.POST("/setting/legal", mapi.Setting.LegalSave)

	// 设置-新手引导
	g.GET("/guide", mapi.Setting.GuideList)          // 列表
	g.GET("/guide/:id", mapi.Setting.GuideGet)       // 详情
	g.POST("/guide", mapi.Setting.GuideSave)         // 新增
	g.PUT("/guide", mapi.Setting.GuideModify)        // 编辑
	g.DELETE("/guide/:id", mapi.Setting.GuideDelete) // 删除

	// 设置-广告活动
	g.GET("/activity", mapi.Setting.ActivityList)          // 列表
	g.GET("/activity/:id", mapi.Setting.ActivityDetail)    // 详情
	g.POST("/activity", mapi.Setting.ActivityCreate)       // 新增
	g.PUT("/activity", mapi.Setting.ActivityModify)        // 编辑
	g.DELETE("/activity/:id", mapi.Setting.ActivityDelete) // 删除

	// 设置-问题分类
	g.GET("/question/category", mapi.QuestionCategory.List) // 列表
	g.POST("/question/category", mapi.QuestionCategory.Create)
	g.PUT("/question/category", mapi.QuestionCategory.Modify)
	g.GET("/question/category/:id", mapi.QuestionCategory.Detail)
	g.DELETE("/question/category/:id", mapi.QuestionCategory.Delete)

	// 设置-常见问题
	g.GET("/question", mapi.Question.List)
	g.POST("/question", mapi.Question.Create)
	g.PUT("/question", mapi.Question.Modify)
	g.GET("/question/:id", mapi.Question.Detail)
	g.DELETE("/question/:id", mapi.Question.Delete)

	// 城市
	g.GET("/city", mapi.City.List)       // 城市列表
	g.PUT("/city/:id", mapi.City.Modify) // 启用或关闭城市

	// 网点
	g.GET("/branch", mapi.Branch.List)                      // 网点列表
	g.GET("/branch/selector", mapi.Branch.Selector)         // 网点简单列表
	g.POST("/branch", mapi.Branch.Create)                   // 新增网点
	g.PUT("/branch/:id", mapi.Branch.Modify)                // 编辑网点
	g.POST("/branch/:id/contract", mapi.Branch.AddContract) // 添加合同
	g.POST("/branch/contract/sheet", mapi.Branch.Sheet)
	g.GET("/branch/nearby", mapi.Branch.Nearby)

	// 门店
	g.GET("/store", mapi.Store.List)
	g.POST("/store", mapi.Store.Create)
	g.PUT("/store/:id", mapi.Store.Modify)
	g.DELETE("/store/:id", mapi.Store.Delete)

	// 客服工具
	g.POST("/csc/irv", mapi.Csc.BatchReminder) // 逾期催费

	// 电池
	g.GET("/battery/model", mapi.Battery.ListModels)
	g.POST("/battery/model", mapi.Battery.CreateModel)
	g.DELETE("/battery/model", mapi.Battery.DeleteModel)
	g.GET("/battery", mapi.Battery.List)
	g.POST("/battery", mapi.Battery.Create)
	g.PUT("/battery/:id", mapi.Battery.Modify)
	g.POST("/battery/batch", mapi.Battery.BatchCreate)
	g.POST("/battery/bind", mapi.Battery.Bind)
	g.POST("/battery/unbind", mapi.Battery.Unbind)
	g.GET("/battery/xc/:sn", mapi.Battery.Detail)
	g.GET("/battery/xc/statistics/:sn", mapi.Battery.Statistics)
	g.GET("/battery/xc/position/:sn", mapi.Battery.Position)
	g.GET("/battery/xc/fault", mapi.Battery.Fault)

	// 电柜
	g.POST("/cabinet", mapi.Cabinet.Create)
	g.GET("/cabinet", mapi.Cabinet.List)
	g.PUT("/cabinet/:id", mapi.Cabinet.Modify)
	g.DELETE("/cabinet/:id", mapi.Cabinet.Delete)
	g.GET("/cabinet/:id", mapi.Cabinet.Detail)
	g.POST("/cabinet/door-operate", mapi.Cabinet.DoorOperate)
	g.POST("/cabinet/reboot", mapi.Cabinet.Reboot)
	g.GET("/cabinet/fault", mapi.Cabinet.Fault)
	g.PUT("/cabinet/fault/:id", mapi.Cabinet.FaultDeal)
	g.GET("/cabinet/data", mapi.Cabinet.Data)
	g.POST("/cabinet/transfer", mapi.Cabinet.Transfer)
	g.POST("/cabinet/maintain", mapi.Cabinet.Maintain)
	g.POST("/cabinet/openbind", mapi.Cabinet.OpenBind)
	g.POST("/cabinet/bin/deactivate", mapi.Cabinet.Deactivate)
	g.POST("/cabinet/interrupt", mapi.Cabinet.Interrupt)

	// 骑士卡
	g.GET("/plan", mapi.Plan.List)
	g.POST("/plan", mapi.Plan.Create)
	g.PUT("/plan/:id", mapi.Plan.UpdateEnable)
	g.DELETE("/plan/:id", mapi.Plan.Delete)
	g.GET("/plan/introduce/notset", mapi.Plan.IntroduceNotset)
	g.GET("/plan/introduce", mapi.Plan.IntroduceList)
	g.POST("/plan/introduce", mapi.Plan.IntroduceCreate)
	g.PUT("/plan/introduce/:id", mapi.Plan.IntroduceModify)
	g.POST("/plan/time", mapi.Plan.Time)

	// 骑手
	g.GET("/rider", mapi.Rider.List)
	g.POST("/rider/ban", mapi.Rider.Ban)
	g.POST("/rider/block", mapi.Rider.Block)
	g.POST("/rider/deposit", mapi.Rider.Deposit)
	g.POST("/rider/modify", mapi.Rider.Modify)
	g.DELETE("/rider/:id", mapi.Rider.Delete)
	g.POST("/rider/followup", mapi.Rider.FollowUpCreate)
	g.GET("/rider/followup", mapi.Rider.FollowUpList)
	g.POST("/rider/exchange-limit", mapi.Rider.ExchangeLimit)
	g.POST("/rider/exchange-frequency", mapi.Rider.ExchangeFrequency)

	// 订阅
	g.POST("/subscribe/alter", mapi.Subscribe.Alter)
	g.GET("/rider/log", mapi.Rider.Log)
	g.POST("/subscribe/pause", mapi.Subscribe.Pause)
	g.POST("/subscribe/continue", mapi.Subscribe.Continue)
	g.POST("/subscribe/halt", mapi.Subscribe.Halt)
	g.POST("/subscribe/active", mapi.Subscribe.Active)
	g.POST("/subscribe/suspend", mapi.Subscribe.Suspend)
	g.POST("/subscribe/unsuspend", mapi.Subscribe.UnSuspend)
	g.POST("/subscribe/ebike/change", mapi.Subscribe.EbikeChange)
	g.GET("/subscribe/reminder", mapi.Subscribe.Reminder)
	g.POST("/subscribe/ebike/unbind", mapi.Subscribe.EbikeUnbind)

	// 业务
	g.GET("/order", mapi.Order.List)
	g.POST("/order/refund", mapi.Order.RefundAudit)
	g.GET("/exchange", mapi.Exchange.List)
	g.GET("/business", mapi.Business.List)
	g.GET("/business/pause", mapi.Business.Pause)
	g.GET("/business/reserve", mapi.Business.Reserve)
	g.GET("/business/suspend", mapi.Business.Suspend)

	// 团签
	g.POST("/enterprise", mapi.Enterprise.Create)
	g.PUT("/enterprise/:id", mapi.Enterprise.Modify)
	g.GET("/enterprise", mapi.Enterprise.List)
	g.GET("/enterprise/:id", mapi.Enterprise.Detail)
	g.POST("/enterprise/:id/prepayment", mapi.Enterprise.Prepayment)
	g.POST("/enterprise/station", mapi.Enterprise.CreateStation)
	g.PUT("/enterprise/station/:id", mapi.Enterprise.ModifyStation)
	g.GET("/enterprise/station", mapi.Enterprise.ListStation)
	g.POST("/enterprise/rider", mapi.Enterprise.CreateRider)
	g.GET("/enterprise/rider", mapi.Enterprise.ListRider)
	g.GET("/enterprise/bill", mapi.Statement.GetBill)
	g.POST("/enterprise/bill", mapi.Statement.Bill)
	g.GET("/enterprise/bill/historical", mapi.Statement.Historical)
	g.GET("/enterprise/bill/statement", mapi.Statement.Statement)
	g.GET("/enterprise/bill/usage", mapi.Statement.Usage)
	g.POST("/enterprise/price", mapi.Enterprise.Price)
	g.DELETE("/enterprise/price/:id", mapi.Enterprise.DeletePrice)
	g.POST("/enterprise/contract", mapi.Enterprise.ModifyContract)
	g.DELETE("/enterprise/contract/:id", mapi.Enterprise.DeleteContract)
	g.GET("/enterprise/agent", mapi.Enterprise.AgentList)
	g.POST("/enterprise/agent", mapi.Enterprise.AgentCreate)
	g.PUT("/enterprise/agent/:id", mapi.Enterprise.AgentModify)
	g.DELETE("/enterprise/agent/:id", mapi.Enterprise.AgentDelete)
	g.GET("/enterprise/prepayment", mapi.Enterprise.PrepaymentList)                        // 充值列表
	g.POST("/enterprise/bind/cabinet", mapi.Enterprise.BindCabinet)                        // 绑定电柜
	g.GET("/enterprise/unbind/cabinet/:id", mapi.Enterprise.UnbindCabinet)                 // 解绑电柜
	g.GET("/enterprise/subscribe/alter/:enterpriseId", mapi.Enterprise.SubscribeAlterList) // 加时申请列表
	g.POST("/enterprise/subscribe/alter", mapi.Enterprise.SubscribeApply)                  // 审批订阅申请

	// 店员
	g.POST("/employee", mapi.Employee.Create)
	g.PUT("/employee/:id", mapi.Employee.Modify)
	g.GET("/employee", mapi.Employee.List)
	g.DELETE("/employee/:id", mapi.Employee.Delete)
	g.GET("/employee/activity", mapi.Employee.Activity)
	g.POST("/employee/enable", mapi.Employee.Enable)
	g.GET("/employee/attendance", mapi.Attendance.List)
	g.POST("/employee/offwork", mapi.Employee.OffWork)

	// 物资
	g.POST("/stock", mapi.Stock.Create)
	g.GET("/stock/battery/overview", mapi.Stock.BatteryOverview)
	g.GET("/stock/store", mapi.Stock.StoreList)
	g.GET("/stock/cabinet", mapi.Stock.CabinetList)
	g.GET("/stock/detail", mapi.Stock.Detail)
	g.GET("/stock/enterprise/list", mapi.Stock.EnterpriseList) // 团签物资列表

	// 选择项目
	g.GET("/selection/plan", mapi.Selection.Plan)
	g.GET("/selection/rider", mapi.Selection.Rider)
	g.GET("/selection/store", mapi.Selection.Store)
	g.GET("/selection/employee", mapi.Selection.Employee)
	g.GET("/selection/city", mapi.Selection.City)
	g.GET("/selection/branch", mapi.Selection.Branch)
	g.GET("/selection/enterprise", mapi.Selection.Enterprise)
	g.GET("/selection/cabinet", mapi.Selection.Cabinet)
	g.GET("/selection/role", mapi.Selection.Role)
	g.GET("/selection/wxemployees", mapi.Selection.WxEmployee)
	g.GET("/selection/planmodel", mapi.Selection.PlanModel)
	g.GET("/selection/cabinetmodel", mapi.Selection.CabinetModel)
	g.GET("/selection/model", mapi.Selection.Model)
	g.GET("/selection/coupon/template", mapi.Selection.CouponTemplate)
	g.GET("/selection/ebike/brand", mapi.Selection.EbikeBrand)
	g.GET("/selection/battery/serial", mapi.Selection.BatterySerial)
	g.GET("/selection/commission/plan", mapi.Selection.CommissionPlan)

	// 救援
	g.GET("/assistance", mapi.Assistance.List)
	g.GET("/assistance/:id", mapi.Assistance.Detail)
	g.GET("/assistance/nearby", mapi.Assistance.Nearby)
	g.POST("/assistance/allocate", mapi.Assistance.Allocate)
	g.POST("/assistance/free", mapi.Assistance.Free)
	g.POST("/assistance/refuse", mapi.Assistance.Refuse)

	// 角色权限
	g.GET("/permission", mapi.Permission.List)
	g.GET("/permission/role", mapi.Permission.ListRole)
	g.POST("/permission/role", mapi.Permission.CreateRole)
	g.PUT("/permission/role/:id", mapi.Permission.ModifyRole)
	g.DELETE("/permission/role/:id", mapi.Permission.DeleteRole)

	// 导入数据
	g.POST("/import/rider", mapi.Import.Rider)

	// 导出数据
	export := g.Group("/export")
	export.GET("", mapi.Export.List)
	export.GET("/download/:sn", mapi.Export.Download)
	export.POST("/rider", mapi.Export.Rider)
	export.POST("/statement/detail", mapi.Export.StatementDetail)
	export.POST("/statement/usage", mapi.Export.StatementUsage)
	export.POST("/order", mapi.Export.Order)
	export.POST("/commission", mapi.Export.Commission)
	export.POST("/business", mapi.Export.Business)
	export.POST("/stock-detail", mapi.Export.StockDetail)
	export.POST("/exchange", mapi.Export.Exchange)

	// 积分
	g.POST("/point/modify", mapi.Point.Modify)
	g.GET("/point/log", mapi.Point.Log)
	g.POST("/point/batch", mapi.Point.Batch)

	// 优惠券
	g.GET("/coupon/template", mapi.Coupon.TemplateList)
	g.POST("/coupon/template", mapi.Coupon.TemplateCreate)
	g.POST("/coupon/template/status", mapi.Coupon.TemplateStatus)
	g.POST("/coupon/generate", mapi.Coupon.Generate)
	g.GET("/coupon/assembly", mapi.Coupon.Assembly)
	g.GET("/coupon", mapi.Coupon.List)
	g.POST("/coupon/allocate", mapi.Coupon.Allocate)

	// 电车
	g.GET("/ebike/brand", mapi.Ebike.BrandList)
	g.POST("/ebike/brand", mapi.Ebike.BrandCreate)
	g.PUT("/ebike/brand/:id", mapi.Ebike.BrandModify)
	g.GET("/ebike", mapi.Ebike.List)
	g.POST("/ebike", mapi.Ebike.Create)
	g.PUT("/ebike/:id", mapi.Ebike.Modify)
	g.POST("/ebike/batch", mapi.Ebike.BatchCreate)

	// 合同
	g.GET("/contract", mapi.Contract.List)

	// 反馈
	g.GET("/enterprise/feedback", mapi.Enterprise.FeedbackList)

	// 运维
	g.GET("/maintainer", mapi.Maintainer.List)
	g.POST("/maintainer", mapi.Maintainer.Create)
	g.PUT("/maintainer/:id", mapi.Maintainer.Modify)

	// 营销
	// PM1 推广会员
	g.GET("/promotion/member", mapi.Member.List)                         // 会员列表
	g.GET("/promotion/member/:id", mapi.Member.Detail)                   // 会员详情
	g.GET("/promotion/member/team/:id", mapi.Member.TeamList)            // 会员团队列表
	g.PUT("/promotion/member/:id", mapi.Member.Update)                   // 修改会员
	g.POST("/promotion/member/setcommission", mapi.Member.SetCommission) // 设置会员佣金方案

	// PM2 收益
	g.GET("/promotion/earnings/:id", mapi.PromotionEarnings.List)       // 收益明细列表
	g.POST("/promotion/earnings/cancel", mapi.PromotionEarnings.Cancel) // 取消收益

	// PM3 提现
	g.GET("/promotion/withdrawal", mapi.PromotionWithdrawal.List)                      // 提现记录列表
	g.POST("/promotion/withdrawal/alter/review", mapi.PromotionWithdrawal.AlterReview) // 提现审核
	g.GET("/promotion/withdrawal/export", mapi.PromotionWithdrawal.Export)             // 导出提现申请

	// PM4 等级
	g.GET("/promotion/level", mapi.PromotionLevel.List)                // 等级列表
	g.POST("/promotion/level", mapi.PromotionLevel.Create)             // 创建等级
	g.PUT("/promotion/level/:id", mapi.PromotionLevel.Update)          // 修改等级
	g.DELETE("/promotion/level/:id", mapi.PromotionLevel.Delete)       // 删除等级
	g.GET("/promotion/level/selection", mapi.PromotionLevel.Selection) // 等级选择

	// PM5 成长
	g.GET("/promotion/growth/:id", mapi.PromotionGrowth.List) // 成长记录列表

	// PM6 任务
	g.GET("/promotion/task", mapi.PromotionLevelTask.List)       // 任务列表
	g.PUT("/promotion/task/:id", mapi.PromotionLevelTask.Update) // 修改任务

	// PM7 成就
	g.GET("/promotion/achievement", mapi.PromotionAchievement.List)               // 成就列表
	g.POST("/promotion/achievement", mapi.PromotionAchievement.Create)            // 创建成就
	g.PUT("/promotion/achievement/:id", mapi.PromotionAchievement.Update)         // 修改成就
	g.DELETE("/promotion/achievement/:id", mapi.PromotionAchievement.Delete)      // 删除成就
	g.POST("/promotion/achievement/upload", mapi.PromotionAchievement.UploadIcon) // 上传icon

	// PM8 佣金方案
	g.GET("/promotion/commission", mapi.PromotionCommission.List)                             // 佣金方案列表
	g.GET("/promotion/commission/:id", mapi.PromotionCommission.Detail)                       // 佣金方案详情
	g.GET("/promotion/commission/history/:id", mapi.PromotionCommission.HistoryList)          // 历史佣金方案列表
	g.POST("/promotion/commission", mapi.PromotionCommission.Create)                          // 创建佣金方案
	g.PUT("/promotion/commission/:id", mapi.PromotionCommission.Update)                       // 修改佣金方案
	g.DELETE("/promotion/commission/:id", mapi.PromotionCommission.Delete)                    // 删除佣金方案
	g.POST("/promotion/commission/enable", mapi.PromotionCommission.Enable)                   // 更新佣金方案状态
	g.GET("/promotion/commission/selection", mapi.PromotionCommission.Selection)              // 方案选择
	g.GET("/promotion/commission/task/selection", mapi.PromotionCommission.TaskSelection)     // 方案任务选择
	g.GET("/promotion/commission/plan/list/:id", mapi.PromotionCommission.CommissionPlanList) //  佣金方案骑士卡列表

	// PM9 推广设置
	g.GET("/promotion/setting/:key", mapi.PromotionSetting.Setting) // 获取会员设置
	g.PUT("/promotion/setting/:key", mapi.PromotionSetting.Update)  // 修改会员设置

	// PMA 推广关系
	g.GET("/promotion/progress/list/:id", mapi.PromotionReferrals.ProgressList) // 推广关系列表

}
