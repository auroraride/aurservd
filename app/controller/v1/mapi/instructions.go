package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type instructions struct{}

var Instructions = new(instructions)

// Modify
// @ID		ManagerInstructionsModify
// @Router	/manager/v1/instructions [POST]
// @Summary	修改说明
// @Tags	说明
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string								true	"管理员校验token"
// @Param	body			body		definition.InstructionsCreateReq	true	"desc"
// @Success	200				{object}	model.StatusResponse				"请求成功"
func (*instructions) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.InstructionsCreateReq](c)
	return ctx.SendResponse(biz.NewInstructionsWithModifierBiz(ctx.Modifier).Modify(req))
}

// Detail
// @ID		ManagerInstructionsDetail
// @Router	/manager/v1/instructions/:key [GET]
// @Summary	说明详情
// @Tags	说明
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	key				path		string						true	"说明key"
// @Success	200				{object}	definition.InstructionsRes	"请求成功"
func (*instructions) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.InstructionsDetailReq](c)
	return ctx.SendResponse(biz.NewInstructions().Detail(req.Key))
}
