// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-26
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
	"github.com/auroraride/aurservd/app/controller/kit"
)

func loadKitRoutes() {
	g := root.Group("kit")
	g.GET("/cabinet/name/:serial", kit.Cabinet.Name)
}
