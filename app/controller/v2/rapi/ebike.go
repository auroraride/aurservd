package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type ebike struct{}

var Ebike = new(ebike)

// EbikeBrandDetail
// @ID		EbikeBrandDetail
// @Router	/rider/v2/ebike/brand/{id} [GET]
// @Summary	电车品牌详情
// @Tags	Ebike - 电车
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	id				path		uint64						true	"电车品牌ID"
// @Param	query			query		definition.EbikeDetailReq	true	"请求参数"
// @Success	200				{object}	definition.EbikeDetailRes	"请求成功"
func (*ebike) EbikeBrandDetail(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[definition.EbikeDetailReq](c)
	return ctx.SendResponse(biz.NewEbikeBiz().EbikeBrandDetail(req))
}
