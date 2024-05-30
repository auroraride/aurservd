package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type export struct{}

var Export = new(export)

// ExportCabinetECMonth
// @ID		ExportCabinetECMonth
// @Router	/manager/v2/export/cabinet/ec/month [POST]
// @Summary	导出电柜能耗
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string								true	"管理员校验token"
// @Param	body			body		definition.CabinetECMonthExportReq	true	"查询参数"
// @Success	200				{object}	model.ExportRes						"成功"
func (*export) ExportCabinetECMonth(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.CabinetECMonthExportReq](c)
	return ctx.SendResponse(biz.NewCabinet().ECMonthExport(ctx.Modifier, req))
}

// ExportCabinetEc
// @ID		ExportCabinetEc
// @Router	/manager/v2/export/cabinet/ec [POST]
// @Summary	导出电柜能耗详情
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		definition.CabinetECReq	true	"查询参数"
// @Success	200				{object}	model.ExportRes			"成功"
func (*export) ExportCabinetEc(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.CabinetECReq](c)
	return ctx.SendResponse(biz.NewCabinet().ECExport(ctx.Modifier, req))
}
