// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-07, by lisicen

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
)

type advert struct{}

var Advert = new(advert)

// List
// @ID		AdvertList
// @Router	/rider/v2/advert [GET]
// @Summary	广告列表
// @Tags	Advert - 广告
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Success	200				{object}	definition.AdvertDetail	"请求成功"
func (a *advert) List(c echo.Context) error {
	ctx := app.ContextX[app.RiderContext](c)
	res := biz.NewAdvert().All()
	return ctx.SendResponse(res)
}
