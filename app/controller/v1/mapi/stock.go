// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type stock struct{}

var Stock = new(stock)

// Create
// @ID		ManagerStockCreate
// @Router	/manager/v1/stock [POST]
// @Summary	ME001 调拨物资
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.StockTransferReq	true	"desc"
// @Success	200				{object}	[]string				"请求成功"
func (*stock) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.StockTransferReq](c)
	return ctx.SendResponse(service.NewStockWithModifier(ctx.Modifier).Transfer(req))
}

// BatteryOverview
// @ID		ManagerStockBatteryOverview
// @Router	/manager/v1/stock/battery/overview [GET]
// @Summary	ME002 电池概览
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	query			query		model.StockOverviewReq			true	"查询目标"
// @Success	200				{object}	model.StockBatteryOverviewRes	"使用中数量 = 激活 + 结束寄存 - 寄存 - 退租"
func (*stock) BatteryOverview(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.StockOverviewReq](c)
	return ctx.SendResponse(service.NewStock().BatteryOverview(req))
}

// StoreList
// @ID		ManagerStockStoreList
// @Router	/manager/v1/stock/store [GET]
// @Summary	ME003 门店物资列表
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string											true	"管理员校验token"
// @Param	query			query		model.StockListReq								true	"desc"
// @Success	200				{object}	model.PaginationRes{items=[]model.StockListRes}	"请求成功"
func (*stock) StoreList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.StockListReq](c)
	return ctx.SendResponse(service.NewStockWithModifier(ctx.Modifier).StoreList(req))
}

// CabinetList
// @ID		ManagerStockCabinetList
// @Router	/manager/v1/stock/cabinet [GET]
// @Summary	ME004 电柜物资列表
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string													true	"管理员校验token"
// @Param	query			query		model.StockCabinetListReq								true	"desc"
// @Success	200				{object}	model.PaginationRes{items=[]model.StockCabinetListRes}	"请求成功"
func (*stock) CabinetList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.StockCabinetListReq](c)
	return ctx.SendResponse(service.NewStockWithModifier(ctx.Modifier).CabinetList(req))
}

// Detail
// @ID		ManagerStockDetail
// @Router	/manager/v1/stock/detail [GET]
// @Summary	ME005 出入库明细
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string												true	"管理员校验token"
// @Param	query			query		model.StockDetailReq								false	"筛选条件"
// @Success	200				{object}	model.PaginationRes{items=[]model.StockDetailRes}	"请求成功"
func (*stock) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.StockDetailReq](c)
	return ctx.SendResponse(service.NewStockWithModifier(ctx.Modifier).Detail(req))
}

// EnterpriseList 团签物资列表
// @ID		ManagerStockEnterpriseList
// @Router	/manager/v1/stock/enterprise/list [GET]
// @Summary	ME006 团签物资列表
// @Tags	[M]管理接口
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string											true	"管理员校验token"
// @Param	query			query		model.StockListReq								true	"desc"
// @Success	200				{object}	model.PaginationRes{items=[]model.StockListRes}	"请求成功"
func (*stock) EnterpriseList(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.StockListReq](c)
	return ctx.SendResponse(service.NewStockWithModifier(ctx.Modifier).EnterpriseList(req))
}
