// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-24, by aurb

package wapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/service"
)

type selection struct{}

var Selection = new(selection)

// Warehouse
// @ID		SelectionWarehouse
// @Router	/warestore/v2/selection/warehouse [GET]
// @Summary	城市-仓库
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string							true	"仓管校验token"
// @Success	200					{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Warehouse(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(biz.NewWarehouse().ListByCity())
}

// ManagerWarehouse
// @ID		SelectionManagerWarehouse
// @Router	/warestore/v2/selection/manager_warehouse [GET]
// @Summary	仓管所属仓库列表
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string							true	"仓管校验token"
// @Success	200					{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) ManagerWarehouse(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(biz.NewWarehouse().ListByManager(ctx.AssetManager))
}

// Store
// @ID		SelectionStore
// @Router	/warestore/v2/selection/store [GET]
// @Summary	城市-门店
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string							true	"仓管校验token"
// @Success	200					{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Store(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(service.NewSelection().Store())
}

// EmployeeStore
// @ID		SelectionEmployeeStores
// @Router	/warestore/v2/selection/employee_store [GET]
// @Summary	店员所属门店列表
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string					true	"仓管校验token"
// @Success	200					{object}	[]model.SelectOption	"请求成功"
func (*selection) EmployeeStore(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(biz.NewStore().ListByEmployee(ctx.Employee))
}

// Enterprise
// @ID		SelectionEnterprise
// @Router	/warestore/v2/selection/enterprise [GET]
// @Summary	城市-团签企业(代理筛选)
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string							true	"仓管校验token"
// @Success	200					{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Enterprise(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(service.NewSelection().Enterprise())
}

// Maintainer
// @ID		SelectionMaintainer
// @Router	/warestore/v2/selection/maintainer [GET]
// @Summary	运维人员
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string					true	"仓管校验token"
// @Success	200					{object}	[]model.SelectOption	"请求成功"
func (*selection) Maintainer(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(biz.NewSelection().MaintainerList())
}

// Station
// @ID		SelectionStation
// @Router	/warestore/v2/selection/station [GET]
// @Summary	企业-站点
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string							true	"仓管校验token"
// @Success	200					{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Station(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(biz.NewSelection().StationList())
}

// Model
// @ID		SelectionModel
// @Router	/warestore/v2/selection/model [GET]
// @Summary	电池型号
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string					true	"仓管校验token"
// @Success	200					{object}	[]model.SelectOption	"请求成功"
func (*selection) Model(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(service.NewBatteryModel().SelectionModels())
}

// EbikeBrand
// @ID		SelectionEbikeBrand
// @Router	/warestore/v2/selection/ebike/brand [GET]
// @Summary	车辆型号列表
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string					true	"仓管校验token"
// @Success	200					{object}	[]model.SelectOption	"请求成功"
func (*selection) EbikeBrand(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(service.NewSelection().EbikeBrand())
}

// CityStation
// @ID		SelectionCityStation
// @Router	/warestore/v2/selection/city_station [GET]
// @Summary	企业-站点
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Warestore-Token	header		string							true	"仓管校验token"
// @Success	200					{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) CityStation(c echo.Context) (err error) {
	ctx := app.ContextX[app.WarestoreContext](c)
	return ctx.SendResponse(biz.NewSelection().CityStationList())
}
