// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-27, by aurb

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
)

type asset struct{}

var Asset = new(asset)

// Assets
// @ID		AgentAssets
// @Router	/agent/v1/assets [GET]
// @Summary	资产数据
// @Tags	Assets - 资产管理
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string														true	"仓管校验token"
// @Param	query			query		definition.WarestoreAssetsReq								true	"查询参数"
// @Success	200				{object}	model.PaginationRes{items=[]definition.WarestoreAssetRes}	"请求成功"
func (*asset) Assets(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[definition.WarestoreAssetsReq](c)
	return ctx.SendResponse(biz.NewWarestore().Assets(definition.AssetSignInfo{Agent: ctx.Agent}, req))
}

// AssetsCommon
// @ID		AgentAssetsCommon
// @Router	/agent/v1/assets/common [GET]
// @Summary	资产电车、电池等数据
// @Tags	Assets - 资产管理
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string											true	"仓管校验token"
// @Param	query			query		definition.WarestoreAssetsCommonReq				true	"查询参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.AssetListRes}	"请求成功"
func (*asset) AssetsCommon(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[definition.WarestoreAssetsCommonReq](c)
	return ctx.SendResponse(biz.NewWarestore().AssetsCommon(definition.AssetSignInfo{Agent: ctx.Agent}, req))
}
