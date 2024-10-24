package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type instructions struct{}

var Instructions = new(instructions)

// Detail
// @ID		RiderInstructionsDetail
// @Router	/rider/v2/instructions/{key} [GET]
// @Summary	买前必读 积分 优惠券使用说明
// @Tags	Instructions - 说明
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	key				path		string						true	"说明key"
// @Success	200				{object}	definition.InstructionsRes	"请求成功"
func (*instructions) Detail(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.InstructionsDetailReq](c)
	return ctx.SendResponse(biz.NewInstructions().Detail(req.Key))
}
