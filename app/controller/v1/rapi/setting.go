// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-18
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/labstack/echo/v4"
)

type setting struct{}

var Setting = new(setting)

// App
// @ID           RiderSettingApp
// @Router       /rider/v1/setting/app [GET]
// @Summary      R6001 获取APP设置
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200 {object}  model.SettingRiderApp  "请求成功"
func (*setting) App(c echo.Context) (err error) {
    ctx := app.Context(c)

    return ctx.SendResponse(model.SettingRiderApp{
        AssistanceFee: cache.Float64(model.SettingRescueFee),
    })
}