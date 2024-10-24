// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-22
// Based on aurservd by liasica, magicrolan@qq.com.

package common

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/payment/alipay"
)

type order struct{}

var Order = new(order)

func (*order) Paytest(c echo.Context) (err error) {
	ctx := app.Context(c)
	result := new(model.OrderCreateRes)
	prepay, no, _ := alipay.NewApp().AppPayDemo()
	result.OutTradeNo = no
	result.Prepay = prepay

	return ctx.SendResponse(result)
}
