package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type cabinet struct{}

var Cabinet = new(cabinet)

// ListEC
// @ID		CabinetListEC
// @Router	/manager/v2/cabinet/ec [GET]
// @Summary	查询电柜能耗
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	query			query		definition.CabinetECReq		true	"查询参数"
// @Success	200				{object}	[]definition.CabinetECRes	"成功"
func (*cabinet) ListEC(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.CabinetECReq](c)
	return ctx.SendResponse(biz.NewCabinet().ListECInfo(*req))
}

// ListECMonth
// @ID		CabinetListECMonth
// @Router	/manager/v2/cabinet/ec/month [GET]
// @Summary	查询电柜月能耗
// @Tags	电柜
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	query			query		definition.CabinetECMonthReq	true	"查询参数"
// @Success	200				{object}	[]definition.CabinetECRes		"成功"
func (*cabinet) ListECMonth(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.CabinetECMonthReq](c)
	return ctx.SendResponse(biz.NewCabinet().ListECMonth(req))
}
