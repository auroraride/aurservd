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
// @Param	X-Maintainer-Token	header		string											true	"仓管校验token"
// @Param	query				query		definition.WarestoreAssetsCommonReq				true	"查询参数"
// @Success	200					{object}	model.PaginationRes{items=[]model.AssetListRes}	"请求成功"
func (*asset) AssetsCommon(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[definition.WarestoreAssetsCommonReq](c)
	return ctx.SendResponse(biz.NewWarestore().AssetsCommon(definition.AssetSignInfo{Maintainer: ctx.Maintainer}, req))
}
