// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-17, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type selection struct{}

var Selection = new(selection)

// WarehouseByCity
// @ID		SelectionWarehouseByCity
// @Router	/manager/v2/asset/selection/warehouse_city [GET]
// @Summary	城市仓库列表
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Success	200						{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) WarehouseByCity(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(biz.NewWarehouse().ListByCity())
}

// City
// @ID		SelectionCity
// @Router	/manager/v2/asset/selection/city [GET]
// @Summary	筛选启用的城市
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Success	200						{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) City(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(service.NewSelection().City())
}

// EbikeBrand
// @ID		SelectionEbikeBrand
// @Router	/manager/v2/asset/selection/ebike/brand [GET]
// @Summary	车辆型号列表
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Success	200						{object}	[]model.SelectOption	"请求成功"
func (*selection) EbikeBrand(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(service.NewSelection().EbikeBrand())
}

// Store
// @ID		SelectionStore
// @Router	/manager/v2/asset/selection/store [GET]
// @Summary	筛选门店
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Success	200						{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Store(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(service.NewSelection().Store())
}

// Enterprise
// @ID		SelectionEnterprise
// @Router	/manager/v2/asset/selection/enterprise [GET]
// @Summary	筛选企业
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Success	200						{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Enterprise(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(service.NewSelection().Enterprise())
}

// AssetRole
// @ID		SelectionAssetRole
// @Router	/manager/v2/asset/selection/role [GET]
// @Summary	筛选角色
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Success	200						{object}	[]model.SelectOption	"请求成功"
func (*selection) AssetRole(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(biz.NewAssetRole().RoleSelection())
}

// Model
// @ID		SelectionModel
// @Router	/manager/v2/asset/selection/model [GET]
// @Summary	筛选电池型号
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string						true	"管理员校验token"
// @Param	query					query		model.SelectModelsReq	true	"查询参数"
// @Success	200						{object}	[]model.SelectOption		"请求成功"
func (*selection) Model(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.SelectModelsReq](c)
	return ctx.SendResponse(service.NewBatteryModel().SelectionModels(req))
}

// Maintainer
// @ID		SelectionMaintainer
// @Router	/manager/v2/asset/selection/maintainer [GET]
// @Summary	筛选运维人员
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Success	200						{object}	[]model.SelectOption	"请求成功"
func (*selection) Maintainer(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(biz.NewSelection().MaintainerList())
}

// Station
// @ID		SelectionStation
// @Router	/manager/v2/asset/selection/station [GET]
// @Summary	筛选站点
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Success	200						{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Station(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(biz.NewSelection().StationList())
}

// Material
// @ID		SelectionMaterialSelect
// @Router	/manager/v2/asset/selection/material [GET]
// @Summary	物资类型筛选
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	query					query		model.SelectMaterialReq	true	"查询参数"
// @Success	200						{object}	[]model.SelectOption	"请求成功"
func (*selection) Material(c echo.Context) (err error) {
	ctx, req := app.ContextBinding[model.SelectMaterialReq](c)
	return ctx.SendResponse(biz.NewSelection().MaterialSelect(req))
}
