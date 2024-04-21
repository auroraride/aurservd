package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type fault struct{}

var Fault = new(fault)

// Create
// @ID		FaultCreate
// @Router	/rider/v2/fault [POST]
// @Summary	故障上报
// @Tags	Fault - 故障
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		definition.FaultCreateReq	true	" 故障上报请求"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*fault) Create(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.FaultCreateReq](c)
	return ctx.SendResponse(biz.NewFaultBiz().Create(ctx.Rider, req))
}

// FaultCause
// @ID		FaultCause
// @Router	/rider/v2/fault/cause [GET]
// @Summary	故障原因
// @Tags	Fault - 故障
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	[]definition.FaultCauseRes	"请求成功"
func (*fault) FaultCause(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(biz.NewFaultBiz().FaultCause())
}
