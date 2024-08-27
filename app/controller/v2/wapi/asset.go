// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-23, by aurb

package wapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type assets struct{}

var Assets = new(assets)

// AssetCount
// @ID		WarestoreAssetCount
// @Router	/warestore/v2/asset/count [GET]
// @Summary	资产统计
// @Tags	Assets - 资产管理
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string						true	"仓管校验token"
// @Success	200					{object}	definition.AssetCountRes	"请求成功"
func (*assets) AssetCount(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(biz.NewWarestore().AssetCount(ctx.AssetManager, ctx.Employee))
}

// Assets
// @ID		WarestoreAssets
// @Router	/warestore/v2/assets [GET]
// @Summary	资产数据
// @Tags	Assets - 资产管理
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string							true	"仓管校验token"
// @Param	query				query		definition.WarestoreAssetsReq	true	"查询参数"
// @Success	200					{object}	[]definition.WarestoreAssetRes	"请求成功"
func (*assets) Assets(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[definition.WarestoreAssetsReq](c)
	return ctx.SendResponse(biz.NewWarestore().Assets(definition.AssetSignInfo{
		AssetManager: ctx.AssetManager,
		Employee:     ctx.Employee,
	}, req))
}

// AssetsCommon
// @ID		WarestoreAssetsCommon
// @Router	/warestore/v2/assets/common [GET]
// @Summary	资产电车、电池等数据
// @Tags	Assets - 资产管理
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string											true	"仓管校验token"
// @Param	query				query		definition.WarestoreAssetsCommonReq				true	"查询参数"
// @Success	200					{object}	model.PaginationRes{items=[]model.AssetListRes}	"请求成功"
func (*assets) AssetsCommon(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[definition.WarestoreAssetsCommonReq](c)
	return ctx.SendResponse(biz.NewWarestore().AssetsCommon(definition.AssetSignInfo{
		AssetManager: ctx.AssetManager,
		Employee:     ctx.Employee,
	}, req))
}
