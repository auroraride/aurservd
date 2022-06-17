// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-25
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type order struct{}

var Order = new(order)

// Create
// @ID           RiderOrderCreate
// @Router       /rider/v1/order [POST]
// @Summary      R3006 支付请求
// @Tags         [R]骑手接口
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
// @Summary      R3007 申请退款
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body  model.RefundReq  true  "desc"
// @Success      200  {object}  model.RefundRes  "请求成功"
func (*order) Refund(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.RefundReq](c)
    return ctx.SendResponse(service.NewRefundWithRider(ctx.Rider).Refund(ctx.Rider.ID, req))
}

// List
// @ID           RiderOrderList
// @Router       /rider/v1/order [GET]
// @Summary      R3008 骑士卡购买历史
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query   model.PaginationReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
// @Success      200  {object}  model.PaginationRes{items=[]model.RiderOrder}  "请求成功"
func (*order) List(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.PaginationReq](c)
    return ctx.SendResponse(
        service.NewRiderOrder().List(ctx.Rider.ID, model.PaginationReqFromPointer(req)),
    )
}

// Detail
// @ID           RiderOrderDetail
// @Router       /rider/v1/order/{id} [GET]
// @Summary      R3009 订单详情
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        id  path  int  true  "订单ID"
// @Success      200  {object}  model.RiderOrder  "请求成功"
func (*order) Detail(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.IDParamReq](c)
    srv := service.NewRiderOrderWithRider(ctx.Rider)
    return ctx.SendResponse(srv.Detail(srv.Query(ctx.Rider.ID, req.ID)))
}
