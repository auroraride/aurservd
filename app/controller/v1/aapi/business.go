// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type business struct{}

var Business = new(business)

// Exchange
// @ID           AgentExchangeList
// @Router       /agent/v1/business/exchange [GET]
// @Summary      A8001 换电列表
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        query  query   model.AgentExchangeListReq  true  "查询条件"
// @Success      200  {object}  model.PaginationRes{items=[]model.ExchangeManagerListRes}  "请求成功"
func (*business) Exchange(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentExchangeListReq](c)
	// TODO 筛选站点
	return ctx.SendResponse(service.NewExchange().List(&model.ExchangeManagerListReq{
		PaginationReq: req.PaginationReq,
		ExchangeListFilter: model.ExchangeListFilter{
			ExchangeListBasicFilter: model.ExchangeListBasicFilter{
				Aimed:   model.BusinessAimedEnterprise,
				Start:   req.Start,
				End:     req.End,
				Keyword: req.Keyword,
			},
			CabinetID:    req.CabinetID,
			EnterpriseID: ctx.Agent.EnterpriseID,
		},
	}))
}

// Price
// @ID           AgentBusinessPrice
// @Router       /agent/v1/business/price [GET]
// @Summary      A8002 价格列表
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Success      200  {object}  []model.EnterprisePriceWithCity
func (*business) Price(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewEnterprise().PriceList(ctx.Enterprise.ID))
}
