// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type plan struct{}

var Plan = new(plan)

// List
// @ID		PlanList
// @Router	/rider/v1/plan [GET]
// @Summary	R3001 新购骑士卡
// @Tags	Plan - 骑士卡
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	query			query		model.PlanListRiderReq	true	"骑士卡列表请求参数"
// @Success	200				{object}	model.PlanNewlyRes		"请求成功"
func (*plan) List(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.PlanListRiderReq](c)
	return ctx.SendResponse(
		service.NewPlanWithRider(ctx.Rider).RiderListNewly(req),
	)
}

// Renewly
// @ID		PlanRenewly
// @Router	/rider/v1/plan/renewly [GET]
// @Summary	R3002 续费骑士卡
// @Tags	Plan - 骑士卡
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	model.RiderPlanRenewalRes	"请求成功"
func (*plan) Renewly(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(
		service.NewPlanWithRider(ctx.Rider).RiderListRenewal(),
	)
}
