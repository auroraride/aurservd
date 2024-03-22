package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type fault struct{}

var Fault = new(fault)

// List
// @ID		FaultList
// @Router	/manager/v1/fault [GET]
// @Summary	故障列表
// @Tags	故障上报
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string											true	"骑手校验token"
// @Param	query			query		definition.FaultListReq							true	"故障列表请求"
// @Success	200				{object}	model.PaginationRes{items=[]definition.Fault}	"请求成功"
func (*fault) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.FaultListReq](c)
	return ctx.SendResponse(biz.NewFaultBiz().List(req))
}

// Modify
// @ID		FaultModify
// @Router	/manager/v1/fault/:id [PUT]
// @Summary	故障更改状态
// @Tags	故障上报
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	id				path		int								true	"故障ID"
// @Param	body			body		definition.FaultModifyStatusReq	true	"故障修改请求"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*fault) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.FaultModifyStatusReq](c)
	return ctx.SendResponse(biz.NewFaultBiz().ModifyStatus(req))
}
