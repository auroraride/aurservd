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

type assetCheck struct{}

var AssetCheck = new(assetCheck)

// GetAssetBySN
// @ID		AssetCheckGetAssetBySN
// @Router	/maintainer/v1/check/sn/{sn} [GET]
// @Summary	通过SN查询资产
// @Tags	AssetCheck - 资产盘点
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string								true	"仓管校验token"
// @Param	sn					path		string								true	"sn"
// @Param	query				query		definition.AssetCheckByAssetSnReq	true	"查询参数"
// @Success	200					{object}	model.AssetCheckByAssetSnRes		"请求成功"
func (*assetCheck) GetAssetBySN(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[definition.AssetCheckByAssetSnReq](c)
	return ctx.SendResponse(biz.NewAssetCheck().GetAssetBySN(definition.AssetSignInfo{Maintainer: ctx.Maintainer}, req))
}
