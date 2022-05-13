// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-09
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    mw "github.com/auroraride/aurservd/app/middleware"
    "github.com/labstack/echo/v4"
)

func loadKaixinRoutes() {
    e.POST("/kaixin/battery", func(c echo.Context) error {
        return c.JSON(200, map[string]any{"state": "ok", "msg": "ok"})
    }, mw.BodyDumpWithConfig(mw.BodyDumpConfig{
        WithRequestHeaders:  true,
        WithResponseHeaders: true,
    }))
}
