// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-01, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type store struct{}

var Store = new(store)

// Asset
// @ID		StoreAsset
// @Router	/manager/v2/asset/store_assets [GET]
// @Summary	门店物资
// @Tags	Store - 门店
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string														true	"管理员校验token"
// @Param	query					query		definition.StoreAssetListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]definition.StoreAssetDetail}	"请求成功"
func (*store) Asset(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.StoreAssetListReq](c)
	return ctx.SendResponse(biz.NewStoreAsset().Assets(req))
}

// AssetDetail
// @ID		StoreAssetDetail
// @Router	/manager/v2/asset/store_assets/{id} [GET]
// @Summary	仓库物资详情
// @Tags	Store - 门店
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	id						path		string							true	"ID"
// @Success	200						{object}	[]definition.CommonAssetDetail	"请求成功"
func (*store) AssetDetail(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewStoreAsset().AssetDetail(req.ID))
}
