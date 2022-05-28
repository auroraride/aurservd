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
// @Summary      R30006 支付请求
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body  model.OrderCreateReq  true  "desc"
// @Success      200  {object}  model.OrderCreateRes  "请求成功"
func (*order) Create(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.OrderCreateReq](c)

    return ctx.SendResponse(
        service.NewOrderWithRider(ctx.Rider).Create(req),
    )
}

// NotActived
// @ID           RiderOrderNotActived
// @Router       /rider/v1/order/not-active [GET]
// @Summary      R30007 未激活骑士卡信息
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200  {object}  model.OrderNotActived  "请求成功"
func (*order) NotActived(c echo.Context) (err error) {
    ctx := app.ContextX[app.RiderContext](c)

    return ctx.SendResponse(
        service.NewRiderOrder().NotActivedDetail(ctx.Rider.ID),
    )
}
