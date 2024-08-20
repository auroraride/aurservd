// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-12, by aurb

package router

import (
	"github.com/auroraride/aurservd/app/controller/v2/wapi"
	"github.com/auroraride/aurservd/app/middleware"
)

func loadWarestoreRoutes() {
	g := root.Group("warestore/v2")

	// 无须校验
	guide := g.Group("", middleware.Warestore())
	guide.POST("/signin", wapi.Warestore.Signin)

	// 需校验
	auth := g.Group("", middleware.Warestore(), middleware.WarestoreAuth())

	// 资产调拨
	auth.GET("/transfer", wapi.Warestore.TransferList)             // 调拨列表
	auth.GET("/transfer/:id", wapi.Warestore.TransferDetail)       // 调拨详情
	auth.POST("/transfer/receive", wapi.Warestore.TransferReceive) // 调拨批量入库
}
