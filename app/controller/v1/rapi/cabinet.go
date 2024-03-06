// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-18
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type cabinet struct{}

var Cabinet = new(cabinet)

// GetProcess
// @ID		CabinetGetProcess
// @Router	/rider/v1/cabinet/process/{serial} [GET]
// @Summary	R4001 获取换电信息
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	serial			path		string					true	"电柜二维码"
// @Success	200				{object}	model.RiderExchangeInfo	"请求成功"
func (*cabinet) GetProcess(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.RiderCabinetOperateInfoReq](c)
	return ctx.SendResponse(service.NewRiderExchange(ctx.Rider).GetProcess(req))
}

// Process
// @ID		CabinetProcess
// @Router	/rider/v1/cabinet/process [POST]
// @Summary	R4002 操作换电
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	body			body		model.RiderExchangeProcessReq	true	"desc"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*cabinet) Process(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.RiderExchangeProcessReq](c)
	service.NewRiderExchange(ctx.Rider).Start(req)
	return ctx.SendResponse()
}

// ProcessStatus
// @ID		CabinetProcessStatus
// @Router	/rider/v1/cabinet/process/status [GET]
// @Summary	R4003 换电状态
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string								true	"骑手校验token"
// @Param	query			query		model.RiderExchangeProcessStatusReq	true	"desc"
// @Success	200				{object}	model.RiderExchangeProcessRes		"请求成功"
func (*cabinet) ProcessStatus(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.RiderExchangeProcessStatusReq](c)

	return ctx.SendResponse(
		service.NewRiderExchange(ctx.Rider).GetProcessStatus(req),
	)
}

// Report
// @ID		CabinetReport
// @Router	/rider/v1/cabinet/report [POST]
// @Summary	R4004 电柜故障上报
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.CabinetFaultReportReq	true	"desc"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*cabinet) Report(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.CabinetFaultReportReq](c)
	return ctx.SendResponse(
		model.StatusResponse{Status: service.NewCabinetFault().Report(ctx.Rider, req)},
	)
}

// Fault
// @ID		CabinetFault
// @Router	/rider/v1/cabinet/fault [GET]
// @Summary	R4008 电柜故障列表
// @Tags	Cabinet - 电柜
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string		true	"骑手校验token"
// @Success	200				{object}	[]string	"请求成功"
func (*cabinet) Fault(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(service.NewSetting().GetSetting(model.SettingCabinetFaultKey))
}
