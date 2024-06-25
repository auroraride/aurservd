package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
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

// Detail
// @ID		PlanDetail
// @Router	/rider/v2/plan/{id} [GET]
// @Summary	骑士卡详情
// @Tags	Plan - 骑士卡
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	id				path		uint64						true	"骑士卡ID"
// @Param	query			query		definition.PlanDetailReq	true	"请求参数"
// @Success	200				{object}	definition.PlanDetailRes	"请求成功"
func (*plan) Detail(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.PlanDetailReq](c)
	return ctx.SendResponse(biz.NewPlanBiz().Detail(req))
}

// ByStore
// @ID		PlanByStore
// @Router	/rider/v2/plan/store [GET]
// @Summary	门店租车套餐列表
// @Tags	Plan - 门店租车套餐
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	query			query		definition.StorePlanReq		true	"门店租车列表请求参数"
// @Success	200				{object}	[]definition.StoreEbikePlan	"请求成功"
func (*plan) ByStore(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.StorePlanReq](c)
	return ctx.SendResponse(biz.NewPlanBiz().ListByStore(req))
}

// ByStoreDetail
// @ID		PlanByStoreDetail
// @Router	/rider/v2/plan/store/detail [GET]
// @Summary	门店租车套餐详情
// @Tags	Plan - 门店租车套餐详情
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	query			query		definition.StorePlanDetailReq	true	"门店套餐详情请求参数"
// @Success	200				{object}	definition.StorePlanDetail		"请求成功"
func (*plan) ByStoreDetail(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.StorePlanDetailReq](c)
	return ctx.SendResponse(biz.NewPlanBiz().StorePlanDetail(ctx.Rider, req))
}

// EbikeBrand
// @ID		PlanEbikeBrand
// @Router	/rider/v2/plan/brand [GET]
// @Summary	车型列表
// @Tags	Plan - 车型列表
// @Accept	json
// @Produce	json
// @Success	200	{object}	[]model.SelectOption	"请求成功"
func (*plan) EbikeBrand(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse(service.NewSelection().EbikeBrand())
}
