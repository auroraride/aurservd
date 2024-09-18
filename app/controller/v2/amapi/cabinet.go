// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-02, by aurb

package amapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type cabinet struct{}

var Cabinet = new(cabinet)

// Asset
// @ID		CabinetAsset
// @Router	/manager/v2/asset/cabinet_assets [GET]
// @Summary	电柜物资
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string														true	"管理员校验token"
// @Param	query					query		definition.CabinetAssetListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]definition.CabinetAssetDetail}	"请求成功"
func (*cabinet) Asset(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.CabinetAssetListReq](c)
	return ctx.SendResponse(biz.NewCabinetAsset().Assets(req))
}

// AssetDetail
// @ID		CabinetAssetDetail
// @Router	/manager/v2/asset/cabinet_assets/{id} [GET]
// @Summary	物资详情
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	id						path		string							true	"ID"
// @Success	200						{object}	definition.CabinetTotalDetail	"请求成功"
func (*cabinet) AssetDetail(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewCabinetAsset().AssetDetail(req.ID))
}
