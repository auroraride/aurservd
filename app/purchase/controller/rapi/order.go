package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/purchase/internal/model"
	"github.com/auroraride/aurservd/app/purchase/internal/service"
)

type order struct{}

var Order = new(order)

func (*order) Create(c echo.Context) error {
	ctx, req := app.RiderContextAndBinding[model.OrderCreateReq](c)
	return ctx.SendResponse(service.NewOrder().Create(ctx.Request().Context(), ctx.Rider, req))
}
