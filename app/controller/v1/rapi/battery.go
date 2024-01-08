// Copyright (C) liasica. 2023-present.
//
// Created at 2023-01-04
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/service"
)

type battery struct{}

var Battery = new(battery)

// Detail
// @ID           RiderBatteryDetail
// @Router       /rider/v1/battery [GET]
// @Summary      RA001 获取电池详情
// @Tags         Battery - 电池
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200 {object}  model.BatteryDetail  "请求成功"
func (*battery) Detail(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(service.NewBattery(ctx.Rider).RiderDetail(ctx.Rider.ID))
}
