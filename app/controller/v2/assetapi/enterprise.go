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

type enterprise struct{}

var Enterprise = new(enterprise)

// EnterpriseAsset
// @ID		EnterpriseAsset
// @Router	/manager/v2/asset/enterprise_assets [GET]
// @Summary	团签物资
// @Tags	团签 - Enterprise
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string															true	"管理员校验token"
// @Param	query					query		definition.EnterpriseAssetListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]definition.EnterpriseAssetDetail}	"请求成功"
func (*enterprise) EnterpriseAsset(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.EnterpriseAssetListReq](c)
	return ctx.SendResponse(biz.NewEnterpriseAsset().Assets(req))
}
