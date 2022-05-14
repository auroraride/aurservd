// Copyright (C) liasica. 2021-present.
//
// Created at 2021-12-15
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/labstack/echo/v4"
)

type contract struct {
}

var Contract = new(contract)

// Sign 签署合同
func (*contract) Sign(c echo.Context) error {
    ctx := c.(*app.RiderContext)
    return ctx.SendResponse(ar.Map{"url": service.NewContract().Sign(ctx.Rider)})
}

// SignResult 获取合同签署结果
func (*contract) SignResult(c echo.Context) error {
    ctx, req := app.ContextBindingX[app.RiderContext, model.ContractSignResultReq](c)
    return ctx.SendResponse(service.NewContract().Result(ctx.Rider, req.Sn))
}
