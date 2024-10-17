package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	am "github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/purchase/internal/model"
	"github.com/auroraride/aurservd/app/purchase/internal/service"
)

type order struct{}

var Order = new(order)

// Create
// @ID		OrderCreate
// @Router	/rider/v2/purchase/order [POST]
// @Summary	创建订单
// @Tags	Order - 购车订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	body			body		model.PurchaseOrderCreateReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*order) Create(c echo.Context) error {
	ctx, req := app.RiderContextAndBinding[model.PurchaseOrderCreateReq](c)
	return ctx.SendResponse(service.NewOrder().Create(ctx.Request().Context(), ctx.Rider, req))
}

// List
// @ID		OrderList
// @Router	/rider/v2/purchase/order [GET]
// @Summary	订单列表
// @Tags	Order - 购车订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string													true	"骑手校验token"
// @Param	query			query		model.PurchaseOrderListReq								true	"desc"
// @Success	200				{object}	model.PaginationRes{items=[]model.PurchaseOrderListRes}	"请求成功"
func (*order) List(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.PurchaseOrderListReq](c)
	return ctx.SendResponse(service.NewOrder().List(req))
}

// Detail
// @ID		OrderDetail
// @Router	/rider/v2/purchase/order/{id} [GET]
// @Summary	订单详情
// @Tags	Order - 购车订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"管理员校验token"
// @Param	id				path		uint64						true	"订单ID"
// @Success	200				{object}	model.PurchaseOrderDetail	"请求成功"
func (*order) Detail(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[am.IDParamReq](c)
	return ctx.SendResponse(service.NewOrder().Detail(req.ID))
}
