package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type stock struct{}

var Stock = new(stock)

// StockList
// @ID           AgentStockList
// @Router       /agent/v1/stock [GET]
// @Summary      AE001 团签物资列表
// @Tags         [A]代理商接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理商校验token"
// @Param        query  query   model.StockDetailReq  true  "desc"
// @Success      200  {object}  model.PaginationRes
func (a *stock) StockList(c echo.Context) error {
	ctx, req := app.AgentContextAndBinding[model.StockDetailReq](c)
	req.EnterpriseID = ctx.Agent.EnterpriseID
	return ctx.SendResponse(service.NewStock().Detail(req))
}

// StockDetail 通过ID查询明细
// @ID           AgentStockDetail
// @Router       /agent/v1/stock [GET]
// @Summary      AE002 团签物资明细
// @Tags         [A]代理商接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理商校验token"
// @Param        query  query   model.StockDetailByIdReq  true  "desc"
// @Success      200  {object}   model.StockDetailRes  "请求成功"
func (a *stock) StockDetail(c echo.Context) error {
	ctx, req := app.AgentContextAndBinding[model.StockDetailByIdReq](c)
	req.EnterpriseID = ctx.Agent.EnterpriseID
	return ctx.SendResponse(service.NewStock().DetailById(req))
}
