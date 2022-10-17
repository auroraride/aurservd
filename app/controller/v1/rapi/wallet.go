// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-17
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type wallet struct{}

var Wallet = new(wallet)

// Overview
// @ID           RiderWalletOverview
// @Router       /rider/v1/wallet/overview [GET]
// @Summary      R9001 钱包概览
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200 {object}  model.WalletOverview  "请求成功"
func (*wallet) Overview(c echo.Context) (err error) {
    ctx := app.ContextX[app.RiderContext](c)
    return ctx.SendResponse(service.NewWallet(ctx.Rider).Overview())
}