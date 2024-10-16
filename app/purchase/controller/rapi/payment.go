package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/purchase/internal/model"
	"github.com/auroraride/aurservd/app/purchase/internal/service"
)

type payment struct{}

var Payment = new(payment)

func (*payment) Pay(c echo.Context) error {
	ctx, req := app.RiderContextAndBinding[model.PaymentReq](c)
	return ctx.SendResponse(service.NewPayment().Pay(ctx.Request().Context(), req))
}
