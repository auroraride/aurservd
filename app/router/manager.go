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
    g.POST("/user/add", mapi.Manager.Add) // 新增管理员

    // 城市
    g.GET("/city", mapi.City.List)       // 城市列表
    g.PUT("/city/:id", mapi.City.Modify) // 启用或关闭城市

    // 网点
    g.GET("/branch", mapi.Branch.List)                      // 网点列表
    g.GET("/branch/selector", mapi.Branch.Selector)         // 网点简单列表
    g.POST("/branch", mapi.Branch.Add)                      // 新增网点
    g.PUT("/branch/:id", mapi.Branch.Modify)                // 编辑网点
    g.POST("/branch/:id/contract", mapi.Branch.AddContract) // 添加合同

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
}
