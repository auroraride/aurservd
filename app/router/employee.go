// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "github.com/auroraride/aurservd/app/controller/v1/eapi"
    "github.com/auroraride/aurservd/app/middleware"
)

func loadEmployeeRoutes() {
    g := root.Group("employee/v1")

    // socket
    g.Any("/socket", eapi.Socket.Employee)

    g.Use(middleware.EmployeeMiddleware())

    // 打卡考勤
    g.POST("/attendance/precheck", eapi.Attendance.Precheck)
    g.POST("/attendance", eapi.Attendance.Create)

    g.POST("/signin", eapi.Employee.Signin)
    g.GET("/qrcode", eapi.Employee.Qrcode)
    g.GET("/profile", eapi.Employee.Profile)
    g.POST("/store/status", eapi.Store.Status)

    g.GET("/subscribe/active", eapi.Subscribe.Inactive, middleware.EmployeeDutyMiddleware())
    g.POST("/subscribe/active", eapi.Subscribe.Active, middleware.EmployeeDutyMiddleware())
    g.GET("/business/rider", eapi.Business.Rider, middleware.EmployeeDutyMiddleware())
    g.POST("/business/pause", eapi.Business.Pause, middleware.EmployeeDutyMiddleware())
    g.POST("/business/continue", eapi.Business.Continue, middleware.EmployeeDutyMiddleware())
    g.POST("/business/unsubscribe", eapi.Business.UnSubscribe, middleware.EmployeeDutyMiddleware())
    g.GET("/business", eapi.Business.List)
    g.GET("/exchange", eapi.Exchange.List)
    g.GET("/stock/overview", eapi.Stock.Overview, middleware.EmployeeDutyMiddleware())
    g.GET("/stock", eapi.Stock.List, middleware.EmployeeDutyMiddleware())

    // 物资
    g.GET("/exception/setting", eapi.Exception.Setting, middleware.EmployeeDutyMiddleware())
    g.POST("/exception", eapi.Exception.Create, middleware.EmployeeDutyMiddleware())

    // 骑手
    g.GET("/rider", eapi.Rider.Detail, middleware.EmployeeDutyMiddleware())
    g.GET("/rider/exchange", eapi.Rider.Exchange, middleware.EmployeeDutyMiddleware())

    // 救援
    g.GET("/assistance/:id", eapi.Assistance.Detail)
    g.POST("/assistance/process", eapi.Assistance.Process)
    g.POST("/assistance/pay", eapi.Assistance.Pay)
    g.GET("/assistance/pay", eapi.Assistance.PayStatus)
    g.GET("/assistance", eapi.Assistance.List)
    g.GET("/assistance/overview", eapi.Assistance.Overview)

    // 电车
    g.GET("/ebike/unallocated", eapi.Ebike.Unallocated)
    g.POST("/ebike/allocate", eapi.Ebike.Allocate)
    g.GET("/ebike/allocate/info", eapi.Ebike.Info)
    g.GET("/ebike/allocate", eapi.Ebike.List)
}
