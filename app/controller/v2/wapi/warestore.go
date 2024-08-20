// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-12, by aurb

package wapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type warestore struct{}

var Warestore = new(warestore)

// Signin
// @ID		WarestoreSignin
// @Router	/warestore/v2/signin [POST]
// @Summary	登录
// @Tags	仓管接口 - Warestore
// @Accept	json
// @Produce	json
// @Param	body	body		definition.WarestorePeopleSigninReq	true	"登录请求"
// @Success	200		{object}	definition.WarestorePeopleSigninRes	"请求成功"
func (*warestore) Signin(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[definition.WarestorePeopleSigninReq](c)
	return ctx.SendResponse(biz.NewWarestore().Signin(req))
}

// AssetCount
// @ID		WarestoreAssetCount
// @Router	/warestore/v2/asset/count [GET]
// @Summary	资产统计
// @Tags	仓管接口 - Warestore
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string						true	"仓管校验token"
// @Success	200					{object}	definition.AssetCountRes	"请求成功"
func (*warestore) AssetCount(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(biz.NewWarestore().AssetCount(ctx.AssetManager, ctx.Employee))
}

// TransferList
// @ID		WarestoreTransferList
// @Router	/warestore/v2/transfer [GET]
// @Summary	调拨记录列表
// @Tags	仓管接口 - Warestore
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string													true	"仓管校验token"
// @Param	query				query		definition.TransferListReq								true	"接收参数"
// @Success	200					{object}	model.PaginationRes{items=[]model.AssetTransferListRes}	"请求成功"
func (*warestore) TransferList(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[definition.TransferListReq](c)
	return ctx.SendResponse(biz.NewWarestore().TransferList(ctx.AssetManager, ctx.Employee, req))
}

// TransferDetail
// @ID		WarestoreTransferDetail
// @Router	/warestore/v2/transfer/{id} [GET]
// @Summary	调拨记录详情
// @Tags	仓管接口 - Warestore
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string							true	"仓管校验token"
// @Param	id					path		uint64							true	"调拨ID"
// @Success	200					{object}	definition.TransferDetailRes	"请求成功"
func (*warestore) TransferDetail(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[model.AssetTransferDetailReq](c)
	return ctx.SendResponse(biz.NewWarestore().TransferDetail(ctx.Request().Context(), req))
}

// TransferReceive
// @ID		WarestoreTransferReceive
// @Router	/warestore/v2/transfer/receive [POST]
// @Summary	接收资产调拨
// @Tags	仓管接口 - Warestore
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string									true	"仓管校验token"
// @Param	body				body		definition.AssetTransferReceiveBatchReq	true	"接收参数"
// @Success	200					{object}	model.StatusResponse					"请求成功"
func (*warestore) TransferReceive(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[definition.AssetTransferReceiveBatchReq](c)
	return ctx.SendResponse(biz.NewWarestore().TransferReceive(ctx.AssetManager, ctx.Employee, req))
}

// TransferFlow
// @ID		WarestoreTransferFlow
// @Router	/warestore/v2/transfer/flow [GET]
// @Summary	资产流转明细
// @Tags	仓管接口 - Warestore
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string						true	"仓管校验token"
// @Param	query				query		model.AssetTransferFlowReq	true	"查询参数"
// @Success	200					{object}	[]model.AssetTransferFlow	"请求成功"
func (*warestore) TransferFlow(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[model.AssetTransferFlowReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().Flow(req))
}

// Assets
// @ID		WarestoreAssets
// @Router	/warestore/v2/assets [GET]
// @Summary	资产数据
// @Tags	仓管接口 - Warestore
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string							true	"仓管校验token"
// @Param	query				query		definition.WarestoreAssetsReq	true	"查询参数"
// @Success	200					{object}	[]definition.WarestoreAssetRes	"请求成功"
func (*warestore) Assets(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[definition.WarestoreAssetsReq](c)
	return ctx.SendResponse(biz.NewWarestore().Assets(ctx.AssetManager, ctx.Employee, req))
}
