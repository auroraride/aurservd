package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type rider struct{}

var Rider = new(rider)

// Direction
// @ID		RiderDirection
// @Router	/rider/v2/direction [GET]
// @Summary	路径规划
// @Tags	Rider - 骑手
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	query			query		definition.RiderDirectionReq	true	"请求参数"
// @Success	200				{object}	definition.RiderDirectionRes	"请求成功"
func (*rider) Direction(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.RiderDirectionReq](c)
	return ctx.SendResponse(biz.NewRiderBiz().Direction(req))
}