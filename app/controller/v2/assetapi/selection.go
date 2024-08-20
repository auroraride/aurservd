// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-17, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
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
// @Success	200						{object}	[]definition.WarehouseByCityRes	"请求成功"
func (*selection) WarehouseByCity(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(biz.NewWarehouse().ListByCity())
}

// City
// ID       SelectionCity
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
