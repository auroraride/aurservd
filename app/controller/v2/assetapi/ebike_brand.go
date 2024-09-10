// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-21, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type ebikeBrand struct{}

var EbikeBrand = new(ebikeBrand)

// List
// @ID		EbikeBrandList
// @Router	/manager/v2/asset/ebike/brand [GET]
// @Summary	品牌列表
// @Tags	电车型号 - EbikeBrand
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string											true	"管理员校验token"
// @Param	query					query		definition.EbikeBrandListReq					true	"品牌详情"
// @Success	200						{object}	model.PaginationRes{items=[]model.EbikeBrand}	"请求成功"
func (*ebikeBrand) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.EbikeBrandListReq](c)
	return ctx.SendResponse(biz.NewEbikeBrand().List(req))
}

// Create
// @ID		EbikeBrandCreate
// @Router	/manager/v2/asset/ebike/brand [POST]
// @Summary	创建品牌
// @Tags	电车型号 - EbikeBrand
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string						true	"管理员校验token"
// @Param	body					body		model.EbikeBrandCreateReq	true	"品牌详情"
// @Success	200						{object}	model.StatusResponse		"请求成功"
func (*ebikeBrand) Create(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.EbikeBrandCreateReq](c)
	return ctx.SendResponse(biz.NewEbikeBrandWithModifier(ctx.Modifier).Create(req))
}

// Modify
// @ID		EbikeBrandModify
// @Router	/manager/v2/asset/ebike/brand/{id} [PUT]
// @Summary	修改品牌
// @Tags	电车型号 - EbikeBrand
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string						true	"管理员校验token"
// @Param	id						path		uint64						true	"品牌ID"
// @Param	body					body		model.EbikeBrandModifyReq	true	"品牌详情"
// @Success	200						{object}	model.StatusResponse		"请求成功"
func (*ebikeBrand) Modify(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.EbikeBrandModifyReq](c)
	return ctx.SendResponse(biz.NewEbikeBrandWithModifier(ctx.Modifier).Modify(req))
}

// Delete
// @ID		EbikeBrandDelete
// @Router	/manager/v2/asset/ebike/brand/{id} [DELETE]
// @Summary	删除品牌
// @Tags	电车型号 - EbikeBrand
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	id						path		uint64					true	"品牌ID"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*ebikeBrand) Delete(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.EbikeBrandDeleteReq](c)
	return ctx.SendResponse(biz.NewEbikeBrandWithModifier(ctx.Modifier).Delete(req))
}
