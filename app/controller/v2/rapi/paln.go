package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/model"
)

type plan struct{}

var Plan = new(plan)

// List
// @ID		PlanList
// @Router	/rider/v2/plan [GET]
// @Summary	新购骑士卡
// @Tags	Plan - 骑士卡
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	query			query		model.PlanListRiderReq	true	"骑士卡列表请求参数"
// @Success	200				{object}	definition.PlanNewlyRes	"请求成功"
func (*plan) List(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.PlanListRiderReq](c)
	return ctx.SendResponse(biz.NewPlanBiz().RiderListNewly(ctx.Rider, req))
}

// Renewly
// @ID		PlanRenewly
// @Router	/rider/v2/plan/renewly [GET]
// @Summary	续费骑士卡
// @Tags	Plan - 骑士卡
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	model.RiderPlanRenewalRes	"请求成功"
func (*plan) Renewly() {}
