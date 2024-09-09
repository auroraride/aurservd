// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-28, by aurb

package aapi

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
// @Router	/agent/v1/selection/warehouse [GET]
// @Summary	城市-仓库
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string							true	"仓管校验token"
// @Success	200				{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Warehouse(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(biz.NewWarehouse().ListByCity())
}

// Store
// @ID		SelectionStore
// @Router	/agent/v1/selection/store [GET]
// @Summary	城市-门店
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string							true	"仓管校验token"
// @Success	200				{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Store(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewSelection().Store())
}

// Maintainer
// @ID		SelectionMaintainer
// @Router	/agent/v1/selection/maintainer [GET]
// @Summary	运维人员
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string					true	"仓管校验token"
// @Success	200				{object}	[]model.SelectOption	"请求成功"
func (*selection) Maintainer(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(biz.NewSelection().MaintainerList())
}

// Station
// @ID		SelectionStation
// @Router	/agent/v1/selection/station [GET]
// @Summary	企业-站点
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string							true	"仓管校验token"
// @Success	200				{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) Station(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(biz.NewSelection().StationList())
}

// Model
// @ID		SelectionModel
// @Router	/agent/v1/selection/model [GET]
// @Summary	电池型号
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string					true	"仓管校验token"
// @Success	200				{object}	[]model.SelectOption	"请求成功"
func (*selection) Model(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewBatteryModel().SelectionModels())
}

// EbikeBrand
// @ID		SelectionEbikeBrand
// @Router	/agent/v1/selection/ebike/brand [GET]
// @Summary	车辆型号列表
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string					true	"仓管校验token"
// @Success	200				{object}	[]model.SelectOption	"请求成功"
func (*selection) EbikeBrand(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewSelection().EbikeBrand())
}

// CityStation
// @ID		SelectionCityStation
// @Router	/agent/v1/selection/city_station [GET]
// @Summary	城市-站点
// @Tags	Selection - 筛选
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string							true	"仓管校验token"
// @Success	200				{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*selection) CityStation(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(biz.NewSelection().CityStationList())
}
