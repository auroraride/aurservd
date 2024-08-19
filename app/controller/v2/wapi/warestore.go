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
	"github.com/auroraride/aurservd/app/service"
)

type warestore struct{}

var Warestore = new(warestore)

// Signin
// @ID		WarehouseSignin
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

// GetOpenid
// @ID		WarehouseGetOpenid
// @Router	/warestore/v2/openid [GET]
// @Summary	获取openid
// @Tags	仓管接口 - Warestore
// @Accept	json
// @Produce	json
// @Param	X-Warehouse-Token	header		string			true	"仓管校验token"
// @Param	code				query		string			true	"微信code"
// @Success	200					{object}	model.OpenidRes	"请求成功"
func (*warestore) GetOpenid(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[definition.OpenidReq](c)
	return ctx.SendResponse(service.NewminiProgram().GetAuth(req.Code))
}

// TransferList
// @ID		WarehouseTransferList
// @Router	/warestore/v2/transfer [GET]
// @Summary	调拨记录列表
// @Tags	仓管接口 - Warestore
// @Accept	json
// @Produce	json
// @Param	X-Warehouse-Token	header		string													true	"仓管校验token"
// @Param	body				body		definition.TransferListReq								true	"接收参数"
// @Success	200					{object}	model.PaginationRes{items=[]model.AssetTransferListRes}	"请求成功"
func (*warestore) TransferList(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[definition.TransferListReq](c)
	return ctx.SendResponse(biz.NewWarestore().TransferList(ctx.AssetManager, ctx.Employee, req))
}

// TransferDetail
// @ID		WarehouseTransferDetail
// @Router	/warestore/v2/transfer/{id} [GET]
// @Summary	调拨记录详情
// @Tags	仓管接口 - Warestore
// @Accept	json
// @Produce	json
// @Param	X-Warehouse-Token	header		string													true	"仓管校验token"
// @Param	body				body		model.AssetTransferDetailReq							true	"接收参数"
// @Success	200					{object}	model.PaginationRes{items=[]model.AssetTransferListRes}	"请求成功"
func (*warestore) TransferDetail(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[model.AssetTransferDetailReq](c)
	return ctx.SendResponse(service.NewAssetTransfer().TransferDetail(ctx.Request().Context(), req))
}

// TransferReceive
// @ID		WarehouseTransferReceive
// @Router	/warestore/v2/transfer/receive [POST]
// @Summary	接收资产调拨
// @Tags	仓管接口 - Warestore
// @Accept	json
// @Produce	json
// @Param	X-Warehouse-Token	header		string								true	"仓管校验token"
// @Param	body				body		model.AssetTransferReceiveBatchReq	true	"接收参数"
// @Success	200					{object}	model.StatusResponse				"请求成功"
func (*warestore) TransferReceive(c echo.Context) (err error) {
	ctx, req := app.WarestoreContextAndBinding[model.AssetTransferReceiveBatchReq](c)
	return ctx.SendResponse(biz.NewWarestore().TransferReceive(ctx.AssetManager, ctx.Employee, req))
}
