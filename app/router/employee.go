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

    g.Use(middleware.EmployeeMiddleware())
    g.POST("/signin", eapi.Employee.Signin)
    g.GET("/qrcode", eapi.Employee.Qrcode)

    g.GET("/subscribe/active", eapi.Subscribe.Inactive, middleware.EmployeeDutyMiddleware())
    g.POST("/subscribe/active", eapi.Subscribe.Active, middleware.EmployeeDutyMiddleware())
    g.GET("/business/rider", eapi.Business.Rider, middleware.EmployeeDutyMiddleware())
    g.POST("/business/pause", eapi.Business.Pause, middleware.EmployeeDutyMiddleware())
    g.POST("/business/continue", eapi.Business.Continue, middleware.EmployeeDutyMiddleware())

    // 打卡考勤
    g.POST("/attendance/precheck", eapi.Attendance.Precheck)
    g.POST("/attendance", eapi.Attendance.Create)
}
