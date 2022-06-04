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

// GetProcess
// @ID           RiderCabinetGetProcess
// @Router       /rider/v1/cabinet/process/{serial} [GET]
// @Summary      R40001 获取换电信息
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        serial  path  string  true  "电柜二维码"
// @Success      200  {object}  model.RiderCabinetInfo  "请求成功"
func (*cabinet) GetProcess(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.RiderCabinetOperateInfoReq](c)
    return ctx.SendResponse(service.NewRiderCabinetWithRider(ctx.Rider).GetProcess(req))
}

// Process
// @ID           RiderCabinetProcess
// @Router       /rider/v1/cabinet/process [POST]
// @Summary      R40002 操作换电
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body  model.RiderCabinetOperateReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*cabinet) Process(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.RiderCabinetOperateReq](c)
    service.NewRiderCabinetWithRider(ctx.Rider).Process(req)
    return ctx.SendResponse()
}

// ProcessStatus
// @ID           RiderCabinetProcessStatus
// @Router       /rider/v1/cabinet/process/status [GET]
// @Summary      R40003 换电状态
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query  model.RiderCabinetOperateStatusReq  true  "desc"
// @Success      200  {object}  model.RiderCabinetOperateRes  "请求成功"
func (*cabinet) ProcessStatus(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.RiderCabinetOperateStatusReq](c)

    return ctx.SendResponse(
        service.NewRiderCabinetWithRider(ctx.Rider).ProcessStatus(req),
    )
}

// Report
// @ID           CabinetReport
// @Router       /rider/v1/cabinet/report [POST]
// @Summary      R40004 电柜故障上报
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
