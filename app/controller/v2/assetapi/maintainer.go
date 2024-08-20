// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-02, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type maintainer struct{}

var Maintainer = new(maintainer)

// MaintainerAsset
// @ID		MaintainerAsset
// @Router	/manager/v2/asset/maintainer_assets [GET]
// @Summary	运维物资
// @Tags	运维 - Maintainer
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string															true	"管理员校验token"
// @Param	query					query		definition.MaintainerAssetListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]definition.MaintainerAssetDetail}	"请求成功"
func (*maintainer) MaintainerAsset(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.MaintainerAssetListReq](c)
	return ctx.SendResponse(biz.NewMaintainerAsset().Assets(req))
}
