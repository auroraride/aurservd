// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-03
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/labstack/echo/v4"
)

type order struct{}

var Order = new(order)

// List
// @ID		ManagerOrderList
// @Router	/manager/v1/order [GET]
// @Summary	M8001 订单列表
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string										true	"管理员校验token"
// @Param	query			query		model.OrderListReq							true	"desc"
// @Success	200				{object}	model.PaginationRes{items=[]model.Order}	"请求成功"
func (*order) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.OrderListReq](c)
	return ctx.SendResponse(
		service.NewOrderWithModifier(ctx.Modifier).List(req),
	)
}

// RefundAudit
// @ID		ManagerOrderRefundAudit
// @Router	/manager/v1/order/refund [POST]
// @Summary	M8002 退款审核
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.RefundAuditReq	true	"desc"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*order) RefundAudit(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.RefundAuditReq](c)
	service.NewRefundWithModifier(ctx.Modifier).RefundAudit(req)
	return ctx.SendResponse()
}

func (*order) Coupon(c echo.Context) (err error) {
	ctx := app.Context(c)
	return ctx.SendResponse()
}
