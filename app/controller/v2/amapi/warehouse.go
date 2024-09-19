// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-11, by Jorjan

package amapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type warehouse struct{}

var Warehouse = new(warehouse)

// List
// @ID		WarehouseList
// @Router	/manager/v2/asset/warehouse [GET]
// @Summary	仓库列表
// @Tags	Warehouse - 仓库
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string													true	"管理员校验token"
// @Param	query					query		definition.WareHouseListReq								true	"查询参数"
// @Success	200						{object}	model.PaginationRes{items=[]definition.WarehouseDetail}	"请求成功"
func (*warehouse) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.WareHouseListReq](c)
	return ctx.SendResponse(biz.NewWarehouse().List(req))
}

// Detail
// @ID		WarehouseDetail
// @Router	/manager/v2/asset/warehouse/{id} [GET]
// @Summary	仓库详情
// @Tags	Warehouse - 仓库
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string						true	"管理员校验token"
// @Param	id						path		string						true	"仓库ID"
// @Success	200						{object}	definition.WarehouseDetail	"请求成功"
func (*warehouse) Detail(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewWarehouse().Detail(req.ID))
}

// Create
// @ID		WarehouseCreate
// @Router	/manager/v2/asset/warehouse [POST]
// @Summary	创建仓库
// @Tags	Warehouse - 仓库
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	body					body		definition.WarehouseCreateReq	true	"请求参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*warehouse) Create(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.WarehouseCreateReq](c)
	return ctx.SendResponse(biz.NewWarehouseWithModifier(ctx.Modifier).Create(req))
}

// Delete
// @ID		WarehouseDelete
// @Router	/manager/v2/asset/warehouse/{id} [DELETE]
// @Summary	删除仓库
// @Tags	Warehouse - 仓库
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	id						path		string					true	"仓库ID"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*warehouse) Delete(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewWarehouseWithModifier(ctx.Modifier).Delete(req.ID))
}

// Modify
// @ID		WarehouseModify
// @Router	/manager/v2/asset/warehouse/{id} [PUT]
// @Summary	修改仓库
// @Tags	Warehouse - 仓库
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	body					body		definition.WarehouseModifyReq	true	"请求参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*warehouse) Modify(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.WarehouseModifyReq](c)
	return ctx.SendResponse(biz.NewWarehouseWithModifier(ctx.Modifier).Modify(req))
}

// Assets
// @ID		WarehouseAssets
// @Router	/manager/v2/asset/warehouse_assets [GET]
// @Summary	仓库物资
// @Tags	Warehouse - 仓库
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string															true	"管理员校验token"
// @Param	query					query		definition.WareHouseAssetListReq								true	"desc"
// @Success	200						{object}	model.PaginationRes{items=[]definition.WareHouseAssetDetail}	"请求成功"
func (*warehouse) Assets(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.WareHouseAssetListReq](c)
	return ctx.SendResponse(biz.NewWarehouse().Assets(req))
}

// AssetDetail
// @ID		WarehouseAssetDetail
// @Router	/manager/v2/asset/warehouse_assets/{id} [GET]
// @Summary	仓库物资详情
// @Tags	Warehouse - 仓库
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	id						path		string							true	"ID"
// @Success	200						{object}	definition.CommonAssetDetail	"请求成功"
func (*warehouse) AssetDetail(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewWarehouse().AssetsDetail(req.ID))
}
