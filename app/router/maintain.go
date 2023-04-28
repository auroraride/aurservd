// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-06
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import "github.com/auroraride/aurservd/app/controller"

func loadMaintainRoutes() {
	g := root.Group("maintain")
	g.GET("/update", controller.Maintain.Update)
}
