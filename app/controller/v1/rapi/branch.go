// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type branch struct{}

var Branch = new(branch)

func (*branch) List(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.BranchWithDistanceReq](c)
    return ctx.SendResponse(service.NewBranch().ListByDistance(req))
}
