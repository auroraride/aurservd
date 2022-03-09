// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-03
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    "github.com/auroraride/aurservd/assets"
    "github.com/auroraride/aurservd/assets/docs"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/labstack/echo/v4"
)

// @title                极光出行API
// @version              1.0
// @description.markdown 极光出行所有API接口文档
// @BasePath             /
// doc https://github.com/swaggo/swag/issues/386 https://github.com/swaggo/swag/issues/548 https://github.com/go-openapi/runtime/blob/master/middleware/redoc.go
func loadRedocRoute() {
    g := e.Group("/docs")

    docs.SwaggerInfo.Host = ar.Config.App.Host

    g.GET("", func(c echo.Context) error {
        return c.HTML(200, assets.SwaggerRedocUI)
    })

    g.GET("/swagger.json", func(c echo.Context) error {
        return c.Blob(200, "application/json", []byte(docs.SwaggerInfo.ReadDoc()))
    })
}
