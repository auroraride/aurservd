// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-27
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type business struct{}

var Business = new(business)

// Active
// @ID           RiderBusinessActive
// @Router       /rider/v1/business/active [POST]
// @Summary      R7001 激活骑士卡
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.BusinessCabinetReq  true  "业务请求"
// @Success      200 {object}   model.BusinessCabinetStatus  "请求成功"
func (*business) Active(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.BusinessCabinetReq](c)
    return ctx.SendResponse(
        service.NewRiderBusiness(ctx.Rider).Active(req),
    )
}

// Unsubscribe
// @ID           RiderBusinessUnsubscribe
// @Router       /rider/v1/business/unsubscribe [POST]
// @Summary      R7002 退租
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.BusinessCabinetReq  true  "业务请求"
// @Success      200 {object}   model.BusinessCabinetStatus  "请求成功"
func (*business) Unsubscribe(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.BusinessCabinetReq](c)
    return ctx.SendResponse(
        service.NewRiderBusiness(ctx.Rider).Unsubscribe(req),
    )
}

// Pause
// @ID           RiderBusinessPause
// @Router       /rider/v1/business/pause [POST]
// @Summary      R7003 寄存
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.BusinessCabinetReq  true  "业务请求"
// @Success      200 {object}   model.BusinessCabinetStatus  "请求成功"
func (*business) Pause(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.BusinessCabinetReq](c)
    return ctx.SendResponse(
        service.NewRiderBusiness(ctx.Rider).Pause(req),
    )
}

// Continue
// @ID           RiderBusinessContinue
// @Router       /rider/v1/business/continue [POST]
// @Summary      R7004 取消寄存
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.BusinessCabinetReq  true  "业务请求"
// @Success      200 {object}   model.BusinessCabinetStatus  "请求成功"
func (*business) Continue(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.BusinessCabinetReq](c)
    return ctx.SendResponse(
        service.NewRiderBusiness(ctx.Rider).Continue(req),
    )
}

// Status
// @ID           RiderBusinessStatus
// @Router       /rider/v1/business/status [GET]
// @Summary      R7005 业务状态
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query   model.BusinessCabinetStatusReq  true  "业务请求"
// @Success      200 {object}   model.BusinessCabinetStatusRes  "请求成功"
func (*business) Status(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.BusinessCabinetStatusReq](c)
    return ctx.SendResponse(
        service.NewRiderBusiness(ctx.Rider).Status(req),
    )
}
