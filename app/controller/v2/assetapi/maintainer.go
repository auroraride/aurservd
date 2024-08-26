// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-02, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type maintainer struct{}

var Maintainer = new(maintainer)

// Asset
// @ID		MaintainerAsset
// @Router	/manager/v2/asset/maintainer_assets [GET]
// @Summary	运维物资
// @Tags	Maintainer - 运维
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string															true	"管理员校验token"
// @Param	query					query		definition.MaintainerAssetListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]definition.MaintainerAssetDetail}	"请求成功"
func (*maintainer) Asset(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.MaintainerAssetListReq](c)
	return ctx.SendResponse(biz.NewMaintainerAsset().Assets(req))
}

// AssetDetail
// @ID		MaintainerAssetDetail
// @Router	/manager/v2/asset/maintainer_assets/{id} [GET]
// @Summary	物资详情
// @Tags	Maintainer - 运维
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	id						path		string							true	"ID"
// @Success	200						{object}	definition.CommonAssetDetail	"请求成功"
func (*maintainer) AssetDetail(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewMaintainerAsset().AssetDetail(req.ID))
}
