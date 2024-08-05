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
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.AssetTransferCreateReq	true	"调拨参数"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*assetTransfer) Transfer(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetTransferCreateReq](c)
	platform := model.AssetTransferTypePlatform
	req.Type = &platform
	return ctx.SendResponse(service.NewAssetTransfer().Transfer(ctx.Request().Context(), req, ctx.Modifier))
}

// TransferList
// @ID		AssetTransferList
// @Router	/manager/v2/asset/transfer [GET]
// @Summary	资产调拨列表
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string													true	"管理员校验token"
// @Param	query			query		model.AssetTransferListReq									true	"查询参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.AssetTransferListRes}	"请求成功"
func (*assetTransfer) TransferList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetTransferListReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().TransferList(ctx.Request().Context(), req))
}

// TransferDetail
// @ID		AssetTransferDetail
// @Router	/manager/v2/asset/transfer/{id} [GET]
// @Summary	资产调拨详情
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	id				path		uint64							true	"调拨ID"
// @Success	200				{object}	model.AssetTransferDetailReq	"请求成功"
func (*assetTransfer) TransferDetail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetTransferDetailReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().TransferDetail(ctx.Request().Context(), req))
}

// TransferCancel
// @ID		AssetTransferCancel
// @Router	/manager/v2/asset/transfer/cancel/{id} [PUT]
// @Summary	取消资产调拨
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	id				path		uint64							true	"调拨ID"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*assetTransfer) TransferCancel(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssetTransferDetailReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().TransferCancel(ctx.Request().Context(), req, ctx.Modifier))
}

// GetTransferBySN
// @ID		GetTransferBySN
// @Router	/manager/v2/asset/transfer/sn/{sn} [GET]
// @Summary	根据调拨单号获取调拨详情
// @Tags	资产
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	sn				path		string							true	"调拨单号"
// @Success	200				{object}	model.AssetTransferDetailReq	"请求成功"
func (*assetTransfer) GetTransferBySN(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.GetTransferBySNReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().GetTransferBySN(ctx.Request().Context(), req))
}
