package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type assetTransfer struct{}

var AssetTransfer = new(assetTransfer)

// Transfer
// @ID		AssetTransfer
// @Router	/manager/v2/asset/transfer [POST]
// @Summary	资产调拨
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	body					body		model.AssetTransferCreateReq	true	"调拨参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*assetTransfer) Transfer(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetTransferCreateReq](c)
	req.OperatorType = model.OperatorTypeAssetManager
	req.OperatorID = ctx.Modifier.ID
	return ctx.SendResponse(service.NewAssetTransfer().Transfer(ctx.Request().Context(), req, ctx.Modifier))
}

// TransferList
// @ID		AssetTransferList
// @Router	/manager/v2/asset/transfer [GET]
// @Summary	资产调拨列表
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string													true	"管理员校验token"
// @Param	query					query		model.AssetTransferListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]model.AssetTransferListRes}	"请求成功"
func (*assetTransfer) TransferList(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetTransferListReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().TransferList(ctx.Request().Context(), req))
}

// TransferDetail
// @ID		AssetTransferDetail
// @Router	/manager/v2/asset/transfer/{id} [GET]
// @Summary	资产调拨详情
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string						true	"管理员校验token"
// @Param	id						path		uint64						true	"调拨ID"
// @Success	200						{object}	[]model.AssetTransferDetail	"请求成功"
func (*assetTransfer) TransferDetail(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetTransferDetailReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().TransferDetail(ctx.Request().Context(), req))
}

// TransferCancel
// @ID		AssetTransferCancel
// @Router	/manager/v2/asset/transfer/cancel/{id} [PUT]
// @Summary	取消资产调拨
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	id						path		uint64					true	"调拨ID"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*assetTransfer) TransferCancel(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetTransferDetailReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().TransferCancel(ctx.Request().Context(), req, ctx.Modifier))
}

// GetTransferBySN
// @ID		GetTransferBySN
// @Router	/manager/v2/asset/transfer/sn/{sn} [GET]
// @Summary	根据调拨单号获取调拨详情
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	sn						path		string							true	"调拨单号"
// @Success	200						{object}	model.AssetTransferDetailReq	"请求成功"
func (*assetTransfer) GetTransferBySN(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.GetTransferBySNReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().GetTransferBySN(ctx.Request().Context(), req))
}

// TransferReceive
// @ID		AssetTransferReceive
// @Router	/manager/v2/asset/transfer/receive [POST]
// @Summary	接收资产调拨
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string								true	"管理员校验token"
// @Param	body					body		model.AssetTransferReceiveBatchReq	true	"接收参数"
// @Success	200						{object}	model.StatusResponse				"请求成功"
func (*assetTransfer) TransferReceive(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetTransferReceiveBatchReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().TransferReceive(ctx.Request().Context(), req, ctx.Modifier))
}

// TransferFlow
// @ID		AssetTransferFlow
// @Router	/manager/v2/asset/transfer/flow [GET]
// @Summary	资产流转明细
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string						true	"管理员校验token"
// @Param	query					query		model.AssetTransferFlowReq	true	"查询参数"
// @Success	200						{object}	[]model.AssetTransferFlow	"请求成功"
func (*assetTransfer) TransferFlow(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetTransferFlowReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().Flow(ctx.Request().Context(), req))
}

// TransferDetailsList
// @ID		AssetTransferDetailsList
// @Router	/manager/v2/asset/transfer/details [GET]
// @Summary	调拨详情列表(出入库明细)
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string															true	"管理员校验token"
// @Param	query					query		model.AssetTransferDetailListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]model.AssetTransferDetailListRes}	"请求成功"
func (*assetTransfer) TransferDetailsList(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetTransferDetailListReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().TransferDetailsList(ctx.Request().Context(), req))
}

// Modify
// @ID		AssetTransferModify
// @Router	/manager/v2/asset/transfer/{id} [PUT]
// @Summary	修改资产调拨
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	id						path		uint64							true	"调拨ID"
// @Param	body					body		model.AssetTransferModifyReq	true	"修改参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*assetTransfer) Modify(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.AssetTransferModifyReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().Modify(ctx.Request().Context(), req, ctx.Modifier))
}
