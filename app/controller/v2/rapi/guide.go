package rapi

import (
	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/labstack/echo/v4"
)

type guide struct{}

var Guide = new(guide)

// List
// @ID		guideList
// @Router	/rider/v2/guide [Get]
// @Summary	新手引导
// @Tags	Guide - 新手引导
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string				true	"骑手校验token"
// @Success	200				{object}	model.GuideDetail	"请求成功"
func (g *guide) List(c echo.Context) error {
	ctx := app.ContextX[app.RiderContext](c)
	res := biz.NewGuide().All()
	return ctx.SendResponse(res)
}
