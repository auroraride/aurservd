// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-23
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/service"
)

type city struct{}

var City = new(city)

// List
// @ID		CityList
// @Router	/rider/v1/city [GET]
// @Summary	R2003 获取已开通城市
// @Tags	City - 城市
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	[]model.CityWithLocation	"请求成功"
func (*city) List(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(service.NewCity().OpenedCities())
}
