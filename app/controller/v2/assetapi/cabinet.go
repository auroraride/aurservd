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

type cabinet struct{}

var Cabinet = new(cabinet)

// CabinetAsset
// @ID		CabinetAsset
// @Router	/manager/v2/asset/cabinet_assets [GET]
// @Summary	电柜物资
// @Tags	电柜 - Cabinet
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string														true	"管理员校验token"
// @Param	query					query		definition.CabinetAssetListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]definition.CabinetAssetDetail}	"请求成功"
func (*cabinet) CabinetAsset(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.CabinetAssetListReq](c)
	return ctx.SendResponse(biz.NewCabinetAsset().Assets(req))
}
