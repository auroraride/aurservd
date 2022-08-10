// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type export struct{}

var Export = new(export)

func (*export) List(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse()
}

func (*export) Download(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse()
}

func (*export) Rider(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.RiderListExport](c)
    return ctx.SendResponse(service.NewRiderWithModifier(ctx.Modifier).ListExport(req))
}

func (*export) StatementDetail(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.StatementBillDetailExport](c)
    return ctx.SendResponse(service.NewEnterpriseStatementWithModifier(ctx.Modifier).DetailExport(req))
}

func (*export) StatementUsage(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.StatementUsageExport](c)
    return ctx.SendResponse(service.NewEnterpriseStatementWithModifier(ctx.Modifier).UsageExport(req))
}
