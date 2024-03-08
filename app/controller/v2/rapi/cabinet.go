package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type cabinet struct{}

var Cabinet = new(cabinet)

// List
// @ID		CabinetList
// @Router	/rider/v2/cabinet [GET]
// @Summary	电柜列表
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	query	query		definition.CabinetByRiderReq	true	"根据骑手获取电柜请求参数"
// @Success	200		{object}	[]definition.CabinetByRiderRes	"请求成功"
func (*cabinet) List(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.CabinetByRiderReq](c)
	return ctx.SendResponse(biz.NewCabinet().ListByRider(ctx.Rider, req))
}

// Detail
// @ID		CabinetDetail
// @Router	/rider/v2/cabinet/{serial} [GET]
// @Summary	电柜详情
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	query	query		definition.CabinetDetailRes	true	"电柜详情请求参数"
// @Success	200		{object}	model.CabinetDetailRes		"请求成功"
func (*cabinet) Detail(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.CabinetDetailRes](c)
	return ctx.SendResponse(biz.NewCabinet().DetailBySerial(req.Serial))
}
