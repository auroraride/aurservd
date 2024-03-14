package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
)

type guide struct{}

var Guide = new(guide)

// List
// @ID		GuideList
// @Router	/rider/v2/guide [GET]
// @Summary	新手指引列表
// @Tags	Guide - 新手引导
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	[]definition.GuideDetail	"请求成功"
func (g *guide) List(c echo.Context) error {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(biz.NewGuide().All())
}
