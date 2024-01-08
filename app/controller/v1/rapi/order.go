// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type order struct{}

var Order = new(order)

// Create
// @ID           RiderOrderCreate
// @Router       /rider/v1/order [POST]
// @Summary      R3005 支付请求
// @Tags         Order - 订单
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.OrderCreateReq  true  "订单创建请求"
// @Success      200  {object}  model.OrderCreateRes  "请求成功"
func (*order) Create(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.OrderCreateReq](c)

	return ctx.SendResponse(
		service.NewOrderWithRider(ctx.Rider).Create(req),
	)
}

// Refund
// @ID           RiderOrderRefund
// @Router       /rider/v1/order/refund [POST]
// @Summary      R3006 申请退款
// @Tags         Order - 订单
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.RefundReq  true  "desc"
// @Success      200  {object}  model.RefundRes  "请求成功"
func (*order) Refund(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.RefundReq](c)
	return ctx.SendResponse(service.NewRefundWithRider(ctx.Rider).Refund(ctx.Rider.ID, req))
}

// List
// @ID           RiderOrderList
// @Router       /rider/v1/order [GET]
// @Summary      R3007 骑士卡购买历史
// @Tags         Order - 订单
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query   model.PaginationReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
// @Success      200  {object}  model.PaginationRes{items=[]model.Order}  "请求成功"
func (*order) List(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.PaginationReq](c)
	return ctx.SendResponse(
		service.NewRiderOrder().List(ctx.Rider.ID, model.PaginationReqFromPointer(req)),
	)
}

// Detail
// @ID           RiderOrderDetail
// @Router       /rider/v1/order/{id} [GET]
// @Summary      R3008 订单详情
// @Tags         Order - 订单
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        id  path  int  true  "订单ID"
// @Success      200  {object}  model.Order  "请求成功"
func (*order) Detail(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewOrderWithRider(ctx.Rider).Detail(service.NewRiderOrderWithRider(ctx.Rider).Query(ctx.Rider.ID, req.ID)))
}

// Status
// @ID           RiderOrderStatus
// @Router       /rider/v1/order/status [GET]
// @Summary      R3009 订单支付状态
// @Tags         Order - 订单
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        outTradeNo     query  string  true  "订单编号"
// @Success      200 {object}   model.OrderStatusRes  "请求成功"
func (*order) Status(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.OrderStatusReq](c)
	return ctx.SendResponse(service.NewOrder().QueryStatus(req))
}
