package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type selection struct {
}

var Selection = new(selection)

// Model
// @ID		SelectionModel
// @Router	/rider/v2/selection/model [GET]
// @Summary	获取电池型号选择
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	query	query		model.SelectionCabinetModelByCityReq	true	"请求参数"
// @Success	200		{object}	[]string								"请求成功"
func (*selection) Model(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.SelectionCabinetModelByCityReq](c)
	return ctx.SendResponse(service.NewSelection().ModelByCity(req))
}

// Brand
// @ID		SelectionBrand
// @Router	/rider/v2/selection/brand [GET]
// @Summary	获取车型列表选择
// @Tags	Selection - 车型列表
// @Accept	json
// @Produce	json
// @Param	query	query		model.SelectionBrandByCityReq	true	"请求参数"
// @Success	200		{object}	[]model.SelectOption			"请求成功"
func (*selection) Brand(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.SelectionBrandByCityReq](c)
	return ctx.SendResponse(service.NewSelection().EbikeBrandByCity(req))
}
