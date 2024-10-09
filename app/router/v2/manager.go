// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-10
// Based on aurservd by liasica, magicrolan@qq.com.

package v2

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	v2 "github.com/auroraride/aurservd/app/controller/v2/mapi"
	"github.com/auroraride/aurservd/app/middleware"
	"github.com/auroraride/aurservd/internal/ar"
)

func LoadManagerV2Routes(root *echo.Group) {
	g := root.Group("manager/v2")

	// 重试token - 用户测试
	g.GET("/retry/token/SHEMBleCticKIdestAilknOGANtIoNAV", func(c echo.Context) error {
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.MapClaims{
				"username": "retryer",
				"exp":      time.Now().Add(time.Minute * 10).Unix(),
				"iat":      time.Now().Unix(),
			},
		)
		tokenString, _ := token.SignedString([]byte(ar.Config.App.RetryTokenSecret))

		return c.String(200, tokenString)
	})

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
