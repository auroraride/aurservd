// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-09
// Based on aurservd by liasica, magicrolan@qq.com.

package router

import (
    mw "github.com/auroraride/aurservd/app/middleware"
    "github.com/labstack/echo/v4"
)

func loadCabinetRoutes() {
    g := e.Group("cabinet", mw.BodyDumpWithConfig(mw.BodyDumpConfig{
        WithRequestHeaders:  true,
        WithResponseHeaders: true,
    }))

    // 凯信
    g.Any("/kaixin", func(c echo.Context) error {
        return c.JSON(200, map[string]any{"state": "ok", "msg": "ok"})
    })

    // 云动
    const yundongToken = "eDeRkenZymbliCEIsTeriAlMaticiPTi"

    g.Any("/yd/Cloud/CabinetStatusNotify", func(c echo.Context) error {
        return c.JSON(200, map[string]interface{}{
            "Status":  "2000", // 状态码2000表示成功
            "Data":    "",     // 响应数据
            "Message": "",     // 响应结果描述
        })
    })

    g.Any("/yd/Cloud/CabinetServiceNotify", func(c echo.Context) error {
        return c.JSON(200, map[string]interface{}{
            "Status":  "2000", // 状态码2000表示成功
            "Data":    "",     // 响应数据
            "Message": "",     // 响应结果描述
        })
    })
}
