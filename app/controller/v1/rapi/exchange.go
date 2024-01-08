// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-04
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type exchange struct{}

var Exchange = new(exchange)

// Store
// @ID           ExchangeStore
// @Router       /rider/v1/exchange/store [POST]
// @Summary      R4005 门店换电
// @Tags         Exchange - 换电
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body  model.ExchangeStoreReq  true  "desc"
// @Success      200  {object}  model.ExchangeStoreRes  "请求成功"
func (*exchange) Store(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.ExchangeStoreReq](c)
	return ctx.SendResponse(service.NewExchangeWithRider(ctx.Rider).Store(req))
}

// Overview
// @ID           ExchangeOverview
// @Router       /rider/v1/exchange/overview [GET]
// @Summary      R4006 换电概览
// @Tags         Exchange - 换电
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200  {object}  model.ExchangeOverview  "请求成功"
func (*exchange) Overview(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(service.NewExchange().Overview(ctx.Rider.ID))
}

// Log
// @ID           ExchangeLog
// @Router       /rider/v1/exchange/log [GET]
// @Summary      R4007 换电记录
// @Tags         Exchange - 换电
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query   model.PaginationReq  true  "分页请求参数"
// @Success      200  {object}  model.PaginationRes{items=[]model.ExchangeRiderListRes}  "请求成功"
func (*exchange) Log(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.PaginationReq](c)
	return ctx.SendResponse(service.NewExchange().RiderList(ctx.Rider.ID, model.PaginationReqFromPointer(req)))
}
