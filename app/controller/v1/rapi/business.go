// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-27
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type business struct{}

var Business = new(business)

// Active
// @ID		BusinessActive
// @Router	/rider/v1/business/active [POST]
// @Summary	R7001 激活骑士卡
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.BusinessCabinetReq	true	"业务请求"
// @Success	200				{object}	model.BusinessCabinetStatus	"请求成功"
func (*business) Active(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.BusinessCabinetReq](c)
	return ctx.SendResponse(service.NewRiderBusiness(ctx.Rider).Active(req, model.RouteVersionV1))
}

// Unsubscribe
// @ID		BusinessUnsubscribe
// @Router	/rider/v1/business/unsubscribe [POST]
// @Summary	R7002 退租
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.BusinessCabinetReq	true	"业务请求"
// @Success	200				{object}	model.BusinessCabinetStatus	"请求成功"
func (*business) Unsubscribe(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.BusinessCabinetReq](c)
	return ctx.SendResponse(
		service.NewRiderBusiness(ctx.Rider).Unsubscribe(req),
	)
}

// Pause
// @ID		BusinessPause
// @Router	/rider/v1/business/pause [POST]
// @Summary	R7003 寄存
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.BusinessCabinetReq	true	"业务请求"
// @Success	200				{object}	model.BusinessCabinetStatus	"请求成功"
func (*business) Pause(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.BusinessCabinetReq](c)
	return ctx.SendResponse(
		service.NewRiderBusiness(ctx.Rider, ctx.Operator).Pause(req),
	)
}

// Continue
// @ID		BusinessContinue
// @Router	/rider/v1/business/continue [POST]
// @Summary	R7004 取消寄存
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Param	body			body		model.BusinessCabinetReq	true	"业务请求"
// @Success	200				{object}	model.BusinessCabinetStatus	"请求成功"
func (*business) Continue(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.BusinessCabinetReq](c)
	return ctx.SendResponse(
		service.NewRiderBusiness(ctx.Rider).Continue(req),
	)
}

// Status
// @ID		BusinessStatus
// @Router	/rider/v1/business/status [GET]
// @Summary	R7005 业务状态
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string							true	"骑手校验token"
// @Param	query			query		model.BusinessCabinetStatusReq	true	"业务请求"
// @Success	200				{object}	model.BusinessCabinetStatusRes	"请求成功"
func (*business) Status(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.BusinessCabinetStatusReq](c)
	return ctx.SendResponse(
		service.NewRiderBusiness(ctx.Rider).Status(req),
	)
}

// PauseInfo
// @ID		BusinessPauseInfo
// @Router	/rider/v1/business/pause/info [GET]
// @Summary	R7006 寄存信息
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string						true	"骑手校验token"
// @Success	200				{object}	model.BusinessPauseInfoRes	"请求成功"
func (*business) PauseInfo(c echo.Context) (err error) {
	ctx := app.ContextX[app.RiderContext](c)
	return ctx.SendResponse(service.NewRiderBusiness(ctx.Rider).PauseInfo())
}

// Allocated
// @ID			BusinessAllocated
// @Router		/rider/v1/business/allocated/{id} [GET]
// @Summary		R7009 长连接轮询是否已分配
// @Description	用以判定待激活骑士卡是否需要签约 (allocated = true)
// @Tags		Business - 业务
// @Accept		json
// @Produce		json
// @Param		X-Rider-Token	header		string					true	"骑手校验token"
// @Param		id				path		uint64					true	"订阅ID"
// @Success		200				{object}	model.AllocateRiderRes	"请求成功"
func (*business) Allocated(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewAllocate().LoopStatus(ctx.Rider.ID, req.ID))
}

// SubscribeSigned
// @ID		BusinessSubscribeSigned
// @Router	/rider/v1/business/subscribe/signed/{id} [GET]
// @Summary	R7010 长连接轮询是否已签约
// @Tags	Business - 业务
// @Accept	json
// @Produce	json
// @Param	X-Rider-Token	header		string					true	"骑手校验token"
// @Param	id				path		uint64					true	"订阅ID"
// @Success	200				{object}	model.SubscribeSigned	"请求成功"
func (*business) SubscribeSigned(c echo.Context) (err error) {
	ctx, req := app.RiderContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewSubscribe().Signed(ctx.Rider.ID, req.ID))
}
