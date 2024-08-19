// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-01, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type store struct{}

var Store = new(store)

// StoreAsset
// @ID		StoreAsset
// @Router	/manager/v2/asset/store_assets [GET]
// @Summary	门店物资
// @Tags	门店 - Store
// @Accept	json
// @Produce	json
// @Param	X-AssetManager-Token	header		string														true	"管理员校验token"
// @Param	query					query		definition.StoreAssetListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]definition.StoreAssetDetail}	"请求成功"
func (*store) StoreAsset(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.StoreAssetListReq](c)
	return ctx.SendResponse(biz.NewStoreAsset().Assets(req))
}
