// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-28, by aurb

package oapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type asset struct{}

var Asset = new(asset)

// AssetsCommon
// @ID		MaintainerAssetsCommon
// @Router	/maintainer/v1/assets/common [GET]
// @Summary	资产电池数据
// @Tags	Assets - 资产管理
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string											true	"运维校验token"
// @Param	query				query		definition.WarestoreAssetsCommonReq				true	"查询参数"
// @Success	200					{object}	model.PaginationRes{items=[]model.AssetListRes}	"请求成功"
func (*asset) AssetsCommon(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[definition.WarestoreAssetsCommonReq](c)
	return ctx.SendResponse(biz.NewWarestore().AssetsCommon(definition.AssetSignInfo{Maintainer: ctx.Maintainer}, req))
}

// AssetCount
// @ID		MaintainerAssetCount
// @Router	/maintainer/v1/assets/count [GET]
// @Summary	资产统计
// @Tags	Assets - 资产管理
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string						true	"运维校验token"
// @Success	200					{object}	definition.AssetCountRes	"请求成功"
func (*asset) AssetCount(c echo.Context) (err error) {
	ctx := app.ContextX[app.MaintainerContext](c)
	return ctx.SendResponse(biz.NewWarestore().AssetCount(definition.AssetSignInfo{Maintainer: ctx.Maintainer}))
}

// Assets
// @ID		AgentAssets
// @Router	/maintainer/v1/assets [GET]
// @Summary	资产数据
// @Tags	Assets - 资产管理
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string							true	"运维校验token"
// @Success	200					{object}	definition.WarestoreAssetDetail	"请求成功"
func (*asset) Assets(c echo.Context) (err error) {
	ctx := app.ContextX[app.MaintainerContext](c)
	return ctx.SendResponse(biz.NewWarestore().AssetsForMaintainer(ctx.Maintainer.ID))
}

// AssetBySn
// @ID		AssetBySn
// @Router	/maintainer/v1/assets/{sn} [GET]
// @Summary	通过sn查询资产数据（任意位置）
// @Tags	Assets - 资产管理
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string							true	"仓管校验token"
// @Param	sn					path		string							true	"sn"
// @Param	query				query		definition.AssetSnReq			true	"查询参数"
// @Success	200					{object}	model.AssetCheckByAssetSnRes	"请求成功"
func (*asset) AssetBySn(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[definition.AssetSnReq](c)
	return ctx.SendResponse(biz.NewWarestore().GetAssetBySN(req))
}
