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

type enterprise struct{}

var Enterprise = new(enterprise)

// Asset
// @ID		EnterpriseAsset
// @Router	/manager/v2/asset/enterprise_assets [GET]
// @Summary	团签物资
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string															true	"管理员校验token"
// @Param	query					query		definition.EnterpriseAssetListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]definition.EnterpriseAssetDetail}	"请求成功"
func (*enterprise) Asset(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.EnterpriseAssetListReq](c)
	return ctx.SendResponse(biz.NewEnterpriseAsset().Assets(req))
}

// AssetDetail
// @ID		EnterpriseAssetDetail
// @Router	/manager/v2/asset/enterprise_assets/{id} [GET]
// @Summary	团签物资详情
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	id						path		string							true	"ID"
// @Success	200						{object}	definition.CommonAssetDetail	"请求成功"
func (*enterprise) AssetDetail(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewEnterpriseAsset().AssetDetail(req.ID))
}

// AssetsExport
// @ID		EnterpriseAssetsExport
// @Router	/manager/v2/asset/enterprise_assets/export [POST]
// @Summary	团签物资导出
// @Tags	Enterprise - 团签
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string								true	"管理员校验token"
// @Param	body					body		definition.EnterpriseAssetListReq	true	"查询参数"
// @Success	200						{object}	model.ExportRes						"成功"
func (*enterprise) AssetsExport(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.EnterpriseAssetListReq](c)
	return ctx.SendResponse(biz.NewEnterpriseAssetWithModifier(ctx.Modifier).AssetsExport(req))
}
