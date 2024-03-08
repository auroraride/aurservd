// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-03-08, by lisicen

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
)

type activity struct{}

var Activity = new(activity)

// List
// @ID		ActivityList
// @Router	/rider/v2/activity [GET]
// @Summary	活动列表
// @Tags	Activity - 活动
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	definition.ActivityDetail	"请求成功"
func (a *activity) List(c echo.Context) error {
	ctx := app.ContextX[app.RiderContext](c)
	res := biz.NewActivity().All()
	return ctx.SendResponse(res)
}
