// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-10
// Based on aurservd by liasica, magicrolan@qq.com.

package v2

import (
	"github.com/labstack/echo/v4"

	v2 "github.com/auroraride/aurservd/app/controller/v2/mapi"
	"github.com/auroraride/aurservd/app/middleware"
)

func LoadManagerV2Routes(root *echo.Group) {
	g := root.Group("manager/v2")

	g.Use(middleware.ManagerMiddleware())

	g.GET("/cabinet/ec", v2.Cabinet.ListEC)            // 电柜能耗详情
	g.GET("/cabinet/ec/month", v2.Cabinet.ListECMonth) // 电柜能耗

	// 导出数据
	export := g.Group("/export")
	export.POST("/cabinet/ec/month", v2.Export.ExportCabinetECMonth) // 导出电柜能耗每月
	export.POST("/cabinet/ec", v2.Export.ExportCabinetEc)            // 导出电柜能耗详情

	ebike := g.Group("/ebike")
	ebike.PUT("/batch", v2.Ebike.BatchModify)
	ebike.DELETE("/brand/:id", v2.Ebike.DeleteBrand)

}
