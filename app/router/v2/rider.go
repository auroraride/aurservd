// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-10
// Based on aurservd by liasica, magicrolan@qq.com.

package v2

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/controller/v2/rapi"
)

func LoadRiderV2Routes(root *echo.Group) {
	g := root.Group("rider/v2")

	// 实名认证
	g.POST("/certification", rapi.Rider.Certification)
}
