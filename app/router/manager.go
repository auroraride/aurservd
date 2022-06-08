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
    g.POST("/user", mapi.Manager.Create) // 新增管理员

    // 设置
    g.GET("/setting", mapi.Setting.List)
    g.PUT("/setting/:key", mapi.Setting.Modify)

    // 城市
    g.GET("/city", mapi.City.List)       // 城市列表
    g.PUT("/city/:id", mapi.City.Modify) // 启用或关闭城市

    // 网点
    g.GET("/branch", mapi.Branch.List)                      // 网点列表
    g.GET("/branch/selector", mapi.Branch.Selector)         // 网点简单列表
    g.POST("/branch", mapi.Branch.Create)                   // 新增网点
    g.PUT("/branch/:id", mapi.Branch.Modify)                // 编辑网点
    g.POST("/branch/:id/contract", mapi.Branch.AddContract) // 添加合同

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
    g.GET("/cabinet", mapi.Cabinet.Query)
    g.PUT("/cabinet/:id", mapi.Cabinet.Modify)
    g.DELETE("/cabinet/:id", mapi.Cabinet.Delete)
    g.GET("/cabinet/:id", mapi.Cabinet.Detail)
    g.POST("/cabinet/door-operate", mapi.Cabinet.DoorOperate)
    g.POST("/cabinet/reboot", mapi.Cabinet.Reboot)
    g.GET("/cabinet/fault", mapi.Cabinet.Fault)
    g.PUT("/cabinet/fault/:id", mapi.Cabinet.FaultDeal)

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

    // 业务
    g.GET("/order", mapi.Order.List)

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

    // 店员
    g.POST("/employee", mapi.Employee.Create)
    g.PUT("/employee/:id", mapi.Employee.Modify)
    g.GET("/employee", mapi.Employee.List)
    g.DELETE("/employee/:id", mapi.Employee.Delete)
}
