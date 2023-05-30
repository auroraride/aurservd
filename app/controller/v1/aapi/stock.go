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
	ctx, req := app.AgentContextAndBinding[model.StockDetailReq](c)
	return ctx.SendResponse(service.NewStock().Detail(req))
}