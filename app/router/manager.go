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
    g.GET("/branch", mapi.Branch.List)                      // 新增网点
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
}
