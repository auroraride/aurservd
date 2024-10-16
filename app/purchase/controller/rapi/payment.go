package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/purchase/internal/model"
	"github.com/auroraride/aurservd/app/purchase/internal/service"
)

type payment struct{}

var Payment = new(payment)

// Pay
// @ID		PaymentPay
// @Router	/rider/v2/purchase/pay [POST]
// @Summary	订单支付
// @Tags	Order - 购车订单
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	body			body		model.PaymentReq		true	"请求参数"
// @Success	200				{object}	model.PurchasePayRes	"请求成功"
func (*payment) Pay(c echo.Context) error {
	ctx, req := app.RiderContextAndBinding[model.PaymentReq](c)
	return ctx.SendResponse(service.NewPayment().Pay(ctx.Request().Context(), req))
}
