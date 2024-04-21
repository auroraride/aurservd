package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
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
// @Success	200	{object}	[]string	"请求成功"
func (*selection) Model(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(service.NewSelection().Models())
}
