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
)

type assetTransfer struct{}

var AssetTransfer = new(assetTransfer)

// Transfer
// @ID		AgentTransfer
// @Router	/agent/v1/transfer [POST]
// @Summary	创建调拨
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string								true	"仓管校验token"
// @Param	body			body		definition.AssetTransferCreateReq	true	"调拨参数"
// @Success	200				{object}	model.StatusResponse				"请求成功"
func (*assetTransfer) Transfer(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[definition.AssetTransferCreateReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().Transfer(definition.AssetSignInfo{Agent: ctx.Agent}, req))
}

// TransferList
// @ID		AgentTransferList
// @Router	/agent/v1/transfer [GET]
// @Summary	调拨记录列表
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string													true	"仓管校验token"
// @Param	query			query		definition.TransferListReq								true	"查询参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.AssetTransferListRes}	"请求成功"
func (*assetTransfer) TransferList(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[definition.TransferListReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().TransferList(definition.AssetSignInfo{Agent: ctx.Agent}, req))
}

// TransferDetail
// @ID		AgentTransferDetail
// @Router	/agent/v1/transfer/{id} [GET]
// @Summary	调拨记录详情
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string							true	"仓管校验token"
// @Param	id				path		uint64							true	"调拨ID"
// @Success	200				{object}	definition.TransferDetailRes	"请求成功"
func (*assetTransfer) TransferDetail(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AssetTransferDetailReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().TransferDetail(ctx.Request().Context(), req))
}

// TransferReceive
// @ID		AgentTransferReceive
// @Router	/agent/v1/transfer/receive [POST]
// @Summary	接收资产调拨/确认入库
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string									true	"仓管校验token"
// @Param	body			body		definition.AssetTransferReceiveBatchReq	true	"接收参数"
// @Success	200				{object}	model.StatusResponse					"请求成功"
func (*assetTransfer) TransferReceive(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[definition.AssetTransferReceiveBatchReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().TransferReceive(definition.AssetSignInfo{Agent: ctx.Agent}, req))
}

// TransferBySn
// @ID		AgentTransferBySn
// @Router	/agent/v1/transfer/sn/{sn} [GET]
// @Summary	根据调拨单号获取调拨详情
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string						true	"仓管校验token"
// @Param	sn				path		string						true	"资产SN"
// @Success	200				{object}	model.AssetTransferListRes	"请求成功"
func (*assetTransfer) TransferBySn(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.GetTransferBySNReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().GetTransferBySn(req))
}
