package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	am "github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/purchase/internal/model"
	"github.com/auroraride/aurservd/app/purchase/internal/service"
)

type order struct{}

var Order = new(order)

// List
// @ID		OrderList
// @Router	/manager/v2/purchase/order [GET]
// @Summary	订单列表
// @Tags	Order - 购车订单
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string													true	"管理员校验token"
// @Param	query			query		model.PurchaseOrderListReq								true	"desc"
// @Success	200				{object}	model.PaginationRes{items=[]model.PurchaseOrderListRes}	"请求成功"
func (*order) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.PurchaseOrderListReq](c)
	return ctx.SendResponse(service.NewOrder().List(req))
}

// Detail
// @ID		OrderDetail
// @Router	/manager/v2/purchase/order/{id} [GET]
// @Summary	订单详情
// @Tags	Order - 购车订单
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	id				path		uint64						true	"订单ID"
// @Success	200				{object}	model.PurchaseOrderDetail	"请求成功"
func (*order) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[am.IDParamReq](c)
	return ctx.SendResponse(service.NewOrder().Detail(req.ID))
}

// Active
// @ID		OrderActive
// @Router	/manager/v2/purchase/order/active [POST]
// @Summary	订单激活
// @Tags	Order - 购车订单
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.PurchaseOrderActiveReq	true	"desc"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*order) Active(c echo.Context) error {
	ctx, req := app.ManagerContextAndBinding[model.PurchaseOrderActiveReq](c)
	return ctx.SendResponse(service.NewOrder().Active(ctx.Request().Context(), req, ctx.Modifier))
}

// Follow
// @ID		OrderFollow
// @Router	/manager/v2/purchase/order/follow [POST]
// @Summary	订单跟进
// @Tags	Order - 购车订单
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.PurchaseOrderFollowReq	true	"desc"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*order) Follow(c echo.Context) error {
	ctx, req := app.ManagerContextAndBinding[model.PurchaseOrderFollowReq](c)
	return ctx.SendResponse(service.NewOrder().Follow(ctx.Request().Context(), req, ctx.Modifier))
}

// Cancel
// @ID		OrderCancel
// @Router	/manager/v2/purchase/order/cancel/{id} [PUT]
// @Summary	取消订单
// @Tags	Order - 购车订单
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		uint64					true	"订单ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*order) Cancel(c echo.Context) error {
	ctx, req := app.ManagerContextAndBinding[am.IDParamReq](c)
	return ctx.SendResponse(service.NewOrder().Cancel(ctx.Request().Context(), req.ID, ctx.Modifier))
}
