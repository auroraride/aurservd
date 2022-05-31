// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import "github.com/auroraride/aurservd/app/controller/v1/eapi"

func loadEmployeeRoutes() {
    g := root.Group("employee/v1")

    g.POST("/subscribe/active", eapi.Subscribe.Active)
}
