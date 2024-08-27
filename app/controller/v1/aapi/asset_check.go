// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-27, by aurb

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type assetCheck struct{}

var AssetCheck = new(assetCheck)

// Create
// @ID		AssetCheckCreate
// @Router	/agent/v1/check [POST]
// @Summary	创建盘点
// @Tags	AssetCheck - 资产盘点
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string							true	"仓管校验token"
// @Param	body			body		definition.AssetCheckCreateReq	true	"请求参数"
// @Success	200				{object}	definition.AssetCheckCreateRes	"请求成功"
func (*assetCheck) Create(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[definition.AssetCheckCreateReq](c)
	return ctx.SendResponse(biz.NewAssetCheck().Create(definition.AssetSignInfo{Agent: ctx.Agent}, req))
}

// GetAssetBySN
// @ID		AssetCheckGetAssetBySN
// @Router	/agent/v1/check/sn/{sn} [GET]
// @Summary	通过SN查询资产
// @Tags	AssetCheck - 资产盘点
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string								true	"仓管校验token"
// @Param	sn				path		string								true	"sn"
// @Param	query			query		definition.AssetCheckByAssetSnReq	true	"查询参数"
// @Success	200				{object}	model.AssetCheckByAssetSnRes		"请求成功"
func (*assetCheck) GetAssetBySN(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[definition.AssetCheckByAssetSnReq](c)
	return ctx.SendResponse(biz.NewAssetCheck().GetAssetBySN(definition.AssetSignInfo{Agent: ctx.Agent}, req))
}

// Detail
// @ID		AssetCheckDetail
// @Router	/agent/v1/check/{id} [GET]
// @Summary	盘点详情
// @Tags	AssetCheck - 资产盘点
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string					true	"仓管校验token"
// @Param	id				path		uint64					true	"盘点ID"
// @Success	200				{object}	model.AssetCheckListRes	"请求成功"
func (*assetCheck) Detail(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewAssetCheck().Detail(ctx.Request().Context(), req.ID))
}

// List
// @ID		AssetCheckList
// @Router	/agent/v1/check [GET]
// @Summary	盘点记录
// @Tags	AssetCheck - 资产盘点
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string													true	"管理员校验token"
// @Param	query			query		definition.AssetCheckListReq							true	"查询参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.AssetCheckListRes}	"请求成功"
func (*assetCheck) List(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[definition.AssetCheckListReq](c)
	return ctx.SendResponse(biz.NewAssetCheck().List(definition.AssetSignInfo{Agent: ctx.Agent}, req))
}

// AssetDetailList
// @ID		AssetCheckDetailList
// @Router	/agent/v1/check/asset/{id} [GET]
// @Summary	盘点资产明细
// @Tags	AssetCheck - 资产盘点
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string												true	"管理员校验token"
// @Param	id				path		uint64												true	"盘点ID"
// @Param	query			query		model.AssetCheckDetailListReq						true	"查询参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.AssetCheckDetail}	"请求成功"
func (*assetCheck) AssetDetailList(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AssetCheckDetailListReq](c)
	return ctx.SendResponse(service.NewAssetCheck().AssetDetailList(ctx.Request().Context(), req))
}
