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

    // 设置
    g.GET("/setting", mapi.Setting.List)
    g.PUT("/setting/:key", mapi.Setting.Modify)
    g.POST("/inventory", mapi.Inventory.CreateOrModify)
    g.GET("/inventory", mapi.Inventory.List)
    g.DELETE("/inventory", mapi.Inventory.Delete)
    g.GET("/inventory/transferable", mapi.Inventory.Transferable)

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
    g.POST("/csc/irv", mapi.Csc.IvrShiguangju) // 逾期催费

    // 电池
    g.GET("/battery/model", mapi.Battery.ListModels)
    g.POST("/battery/model", mapi.Battery.CreateModel)

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

    // 骑士卡
    g.GET("/plan", mapi.Plan.List)
    g.POST("/plan", mapi.Plan.Create)
    g.PUT("/plan/:id", mapi.Plan.UpdateEnable)
    g.DELETE("/plan/:id", mapi.Plan.Delete)

    // 骑手
    g.GET("/rider", mapi.Rider.List)
    g.POST("/rider/ban", mapi.Rider.Ban)
    g.POST("/rider/block", mapi.Rider.Block)
    g.POST("/subscribe/alter", mapi.Subscribe.Alter)
    g.GET("/rider/log", mapi.Rider.Log)
    g.POST("/subscribe/pause", mapi.Subscribe.Pause)
    g.POST("/subscribe/continue", mapi.Subscribe.Continue)
    g.POST("/subscribe/halt", mapi.Subscribe.Halt)
    g.POST("/rider/deposit", mapi.Rider.Deposit)
    g.POST("/rider/modify", mapi.Rider.Modify)
    g.DELETE("/rider/:id", mapi.Rider.Delete)
    g.POST("/rider/followup", mapi.Rider.FollowUpCreate)
    g.GET("/rider/followup", mapi.Rider.FollowUpList)

    // 业务
    g.GET("/order", mapi.Order.List)
    g.POST("/order/refund", mapi.Order.RefundAudit)
    g.GET("/exchange", mapi.Exchange.List)

    // 企业
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
    g.POST("/enterprise/price", mapi.Enterprise.ModifyPrice)
    g.DELETE("/enterprise/price/:id", mapi.Enterprise.DeletePrice)
    g.POST("/enterprise/contract", mapi.Enterprise.ModifyContract)
    g.DELETE("/enterprise/contract/:id", mapi.Enterprise.DeleteContract)

    // 店员
    g.POST("/employee", mapi.Employee.Create)
    g.PUT("/employee/:id", mapi.Employee.Modify)
    g.GET("/employee", mapi.Employee.List)
    g.DELETE("/employee/:id", mapi.Employee.Delete)
    g.GET("/employee/activity", mapi.Employee.Activity)
    g.POST("/employee/enable", mapi.Employee.Enable)
    g.GET("/employee/attendance", mapi.Attendance.List)

    // 物资
    g.POST("/stock", mapi.Stock.Create)
    g.GET("/stock", mapi.Stock.List)
    g.GET("/stock/overview", mapi.Stock.Overview)

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
    g.GET("/selection/wxemployees", mapi.Selection.WxEmployees)

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
}
