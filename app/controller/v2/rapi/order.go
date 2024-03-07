package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
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
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	body			body		model.OrderCreateReq	true	"订单创建请求"
// @Success	200				{object}	model.OrderCreateRes	"请求成功"
func (*order) Create(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.OrderCreateReq](c)
	return ctx.SendResponse(service.NewOrderWithRider(ctx.Rider).Create(req))
}

// DepositFree 免押金支付
// @ID		OrderDepositFree
// @Router	/rider/v2/order/depositfree [POST]
// @Summary	免押金支付
// @Tags	Order - 订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	body			body		definition.OrderDepositFreeReq	true	"免押金支付请求"
// @Success	200				{object}	definition.OrderDepositFreeRes	"请求成功"
func (*order) DepositFree(c echo.Context) (err error) {
	// ctx, req := app.RiderContextAndBinding[definition.OrderDepositFreeReq](c)
	// return ctx.SendResponse(biz.NewOrder().DepositFree(ctx.Rider, req))
	return nil
}
