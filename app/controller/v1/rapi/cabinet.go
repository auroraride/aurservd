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

// Report
// @ID           CabinetReport
// @Router       /rider/v1/path [GET]
// @Summary      R40001 电柜故障上报
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body  model.CabinetFaultReportReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*cabinet) Report(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.CabinetFaultReportReq](c)

    return ctx.SendResponse(
        model.StatusResponse{Status: service.NewCabinetFault().Report(ctx.Rider, req)},
    )
}
