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

// Unfreeze
// @ID		OrderUnfreeze
// @Router	/rider/v2/order/unfreeze [POST]
// @Summary	解冻资金
// @Tags	Order - 订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string								true	"骑手校验token"
// @Param	body			body		definition.OrderDepositUnfreezeReq	true	"解冻押金请求"
// @Success	200				{object}	model.StatusResponse				"请求成功"
// func (*order) Unfreeze(c echo.Context) (err error) {
// 	ctx, req := app.RiderContextAndBinding[definition.OrderDepositUnfreezeReq](c)
// 	return ctx.SendResponse(biz.NewOrderBiz().FandAuthUnfreeze(req))
// }

// PaymentFreezeToPay
// @ID		OrderPaymentFreezeToPay
// @Router	/rider/v2/order/payment/freezetopay [POST]
// @Summary	冻结转支付
// @Tags	Order - 订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	body			body		definition.FreezeToPay	true	"冻结转支付请求"
// @Success	200				{object}	model.StatusResponse	"请求成功"
// func (*order) PaymentFreezeToPay(c echo.Context) (err error) {
// 	ctx, req := app.RiderContextAndBinding[definition.FreezeToPay](c)
// 	return ctx.SendResponse(biz.NewOrderBiz().PaymentFreezeToPay(req))
// }

// Refund
// @ID		OrderRefund
// @Router	/rider/v2/order/refund [POST]
// @Summary	申请退款
// @Tags	Order - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	body			body		definition.RefundReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*order) Refund(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.RefundReq](c)
	return ctx.SendResponse(biz.NewRefundBiz().Refund(ctx.Rider.ID, req))
}
