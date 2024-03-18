package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type order struct{}

var Order = new(order)

// Create
// @ID		OrderCreate
// @Router	/rider/v2/order [POST]
// @Summary	支付请求
// @Tags	Order - 订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		definition.OrderCreateReq	true	"订单创建请求"
// @Success	200				{object}	model.OrderCreateRes		"请求成功"
func (*order) Create(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.OrderCreateReq](c)
	return ctx.SendResponse(biz.NewOrderBiz().Create(ctx.Rider, req))
}

// DepositCredit 信用免押
// @ID		OrderDepositCredit
// @Router	/rider/v2/order/deposit/credit [POST]
// @Summary	免押金支付
// @Tags	Order - 订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string								true	"骑手校验token"
// @Param	body			body		definition.OrderDepositCreditReq	true	"免押金支付请求"
// @Success	200				{object}	definition.OrderDepositCreditRes	"请求成功"
func (*order) DepositCredit(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.OrderDepositCreditReq](c)
	return ctx.SendResponse(biz.NewOrderBiz().DepositCredit(ctx.Rider, req))
}
