// Copyright (C) liasica. 2023-present.
//
// Created at 2023-05-30
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type stock struct{}

var Stock = new(stock)

// Detail
// @ID           AgentStockDetail
// @Router       /agent/v1/stock [GET]
// @Summary      A6001 出入库明细
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        query  query   model.StockDetailReq  false  "筛选条件"
// @Success      200  {object}  model.PaginationRes{items=[]model.StockDetailRes}  "请求成功"
func (*stock) Detail(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentStockDetailReq](c)
	return ctx.SendResponse(service.NewAgentStock().Detail(req))
}

// BatteryStock 电池物资
// @ID           AgentBatteryStock
// @Router       /agent/v1/stock/battery [GET]
// @Summary      A6002 电池物资
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Success      200  {object}  []model.BatteryStockSummaryRsp  "请求成功"
func (*stock) BatteryStock(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewStockSummary().BatteryStockSummary(ctx))
}

// EBikeStock 电车物资
// @ID           AgentBikeStock
// @Router       /agent/v1/stock/ebike [GET]
// @Summary      A6003 电车物资
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Success      200  {object}  []model.EbikeStockSummaryRsp  "请求成功"
func (*stock) EBikeStock(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewStockSummary().EbikeStockSummary(ctx))
}
