// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-18
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type cabinet struct{}

var Cabinet = new(cabinet)

func (*cabinet) Report(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.CabinetFaultReportReq](c)

    return ctx.SendResponse(
        model.StatusResponse{Status: service.NewCabinetFault().Report(ctx.Rider, req)},
    )
}
