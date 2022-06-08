// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type battery struct{}

var Battery = new(battery)

// ListVoltage
// @ID           RiderBatteryListVoltage
// @Router       /rider/v1/battery/voltage [GET]
// @Summary      R3001 电压型号
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200  {object}  []int  "请求成功"
func (*battery) ListVoltage(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewBattery().ListVoltages())
}