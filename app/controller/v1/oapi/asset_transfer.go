// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-28, by aurb

package oapi

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
// @ID		AssetTransfer
// @Router	/maintainer/v1/transfer [POST]
// @Summary	创建调拨
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string								true	"仓管校验token"
// @Param	body				body		definition.AssetTransferCreateReq	true	"调拨参数"
// @Success	200					{object}	model.StatusResponse				"请求成功"
func (*assetTransfer) Transfer(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[definition.AssetTransferCreateReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().Transfer(definition.AssetSignInfo{Maintainer: ctx.Maintainer}, req))
}

// TransferList
// @ID		AssetTransferList
// @Router	/maintainer/v1/transfer [GET]
// @Summary	调拨记录列表
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string													true	"仓管校验token"
// @Param	query				query		definition.TransferListReq								true	"查询参数"
// @Success	200					{object}	model.PaginationRes{items=[]model.AssetTransferListRes}	"请求成功"
func (*assetTransfer) TransferList(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[definition.TransferListReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().TransferList(definition.AssetSignInfo{Maintainer: ctx.Maintainer}, req))
}

// TransferDetail
// @ID		AssetTransferDetail
// @Router	/maintainer/v1/transfer/{id} [GET]
// @Summary	调拨记录详情
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string							true	"仓管校验token"
// @Param	id					path		uint64							true	"调拨ID"
// @Success	200					{object}	definition.TransferDetailRes	"请求成功"
func (*assetTransfer) TransferDetail(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[model.AssetTransferDetailReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().TransferDetail(definition.AssetSignInfo{Maintainer: ctx.Maintainer}, req))
}

// TransferReceive
// @ID		AssetTransferReceive
// @Router	/maintainer/v1/transfer/receive [POST]
// @Summary	接收资产调拨/确认入库
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string									true	"仓管校验token"
// @Param	body				body		definition.AssetTransferReceiveBatchReq	true	"接收参数"
// @Success	200					{object}	model.StatusResponse					"请求成功"
func (*assetTransfer) TransferReceive(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[definition.AssetTransferReceiveBatchReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().TransferReceive(definition.AssetSignInfo{Maintainer: ctx.Maintainer}, req))
}

// TransferBySn
// @ID		AssetTransferBySn
// @Router	/maintainer/v1/transfer/sn/{sn} [GET]
// @Summary	根据调拨单号获取调拨详情
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string						true	"仓管校验token"
// @Param	sn					path		string						true	"资产SN"
// @Success	200					{object}	model.AssetTransferListRes	"请求成功"
func (*assetTransfer) TransferBySn(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[model.GetTransferBySNReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().GetTransferBySn(req))
}

// TransferFlow
// @ID		AssetTransferFlow
// @Router	/maintainer/v1/transfer/flow [GET]
// @Summary	资产流转明细
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string						true	"仓管校验token"
// @Param	query				query		model.AssetTransferFlowReq	true	"查询参数"
// @Success	200					{object}	[]model.AssetTransferFlow	"请求成功"
func (*assetTransfer) TransferFlow(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[model.AssetTransferFlowReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().Flow(req))
}

// TransferDetailsList
// @ID		AssetTransferDetailsList
// @Router	/maintainer/v1/transfer/details [GET]
// @Summary	出入库明细
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string															true	"管理员校验token"
// @Param	query				query		definition.AssetTransferDetailListReq							true	"查询参数"
// @Success	200					{object}	model.PaginationRes{items=[]model.AssetTransferDetailListRes}	"请求成功"
func (*assetTransfer) TransferDetailsList(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[definition.AssetTransferDetailListReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().TransferDetailsList(definition.AssetSignInfo{Maintainer: ctx.Maintainer}, req))
}

// Modify
// @ID		AssetTransferModify
// @Router	/maintainer/v1/transfer/{id} [PUT]
// @Summary	修改资产调拨
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string								true	"管理员校验token"
// @Param	id					path		uint64								true	"调拨ID"
// @Param	body				body		definition.AssetTransferModifyReq	true	"修改参数"
// @Success	200					{object}	model.StatusResponse				"请求成功"
func (*assetTransfer) Modify(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[definition.AssetTransferModifyReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().Modify(definition.AssetSignInfo{Maintainer: ctx.Maintainer}, req))
}

// TransferCancel
// @ID		AssetTransferCancel
// @Router	/maintainer/v1/transfer/cancel/{id} [PUT]
// @Summary	取消资产调拨
// @Tags	AssetTransfer - 资产调拨
// @Accept	json
// @Produce	json
// @Param	X-Maintainer-Token	header		string					true	"管理员校验token"
// @Param	id					path		uint64					true	"调拨ID"
// @Success	200					{object}	model.StatusResponse	"请求成功"
func (*assetTransfer) TransferCancel(c echo.Context) (err error) {
	ctx, req := app.MaintainerContextAndBinding[model.AssetTransferDetailReq](c)
	return ctx.SendResponse(biz.NewAssetTransfer().TransferCancel(definition.AssetSignInfo{Maintainer: ctx.Maintainer}, req))
}
