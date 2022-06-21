// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-21
// Based on aurservd by liasica, magicrolan@qq.com.

package common

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type battery struct{}

var Battery = new(battery)

// Model
// @ID           ManagerBatteryModel
// @Router       /common/battery/model [GET]
// @Summary      C4 获取生效中的电池型号
// @Tags         [C]公共接口
// @Accept       json
// @Produce      json
// @Success      200 {object} []string "型号列表"
func (*battery) Model(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewBattery().Models())
}