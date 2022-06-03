// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-04
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type exchange struct{}

var Exchange = new(exchange)

// Store
// @ID           RiderExchangeStore
// @Router       /rider/v1/store/exchange [POST]
// @Summary      R40005 门店换电
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body  model.ExchangeStoreReq  true  "desc"
// @Success      200  {object}  model.ExchangeStoreRes  "请求成功"
func (*exchange) Store(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.ExchangeStoreReq](c)
    return ctx.SendResponse(service.NewExchangeWithRider(ctx.Rider).Store(req))
}
