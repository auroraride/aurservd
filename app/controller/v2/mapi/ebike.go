package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type ebike struct{}

var Ebike = new(ebike)

// BatchModify
// @ID		ManagerEbikeBatchModify
// @Router	/manager/v2/ebike/batch [PUT]
// @Summary	批量修改电车
// @Tags	电车
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.EbikeBatchModifyReq	true	"电车信息"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*ebike) BatchModify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.EbikeBatchModifyReq](c)
	return ctx.SendResponse(biz.NewEbikeBiz().BatchModify(req))
}

// DeleteBrand
// @ID		ManagerEbikeBrandDelete
// @Router	/manager/v2/ebike/brand/:id [DELETE]
// @Summary	删除品牌
// @Tags	电车
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		uint64					true	"品牌ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*ebike) DeleteBrand(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.EbikeBrandDeleteReq](c)
	return ctx.SendResponse(biz.NewEbikeBiz().DeleteBrand(req))
}
