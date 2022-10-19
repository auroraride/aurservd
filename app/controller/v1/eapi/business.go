// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-14
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type business struct{}

var Business = new(business)

// Rider
// @ID           EmployeeBusinessRider
// @Router       /employee/v1/business/rider [GET]
// @Summary      E2001 骑手业务详情
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token   header  string  true  "店员校验token"
// @Param        qrcode query       string  true  "骑手二维码"
// @Success      200    {object}    model.SubscribeBusiness  "业务详情返回"
func (*business) Rider(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.QRQueryReq](c)
    id := service.NewRider().ParseQrcode(req.Qrcode)
    return ctx.SendResponse(service.NewBusinessWithEmployee(ctx.Employee).Detail(id))
}

// Inactive
// @ID           EmployeeSubscribeInactive
// @Router       /employee/v1/subscribe/active [GET]
// @Summary      E2002 未激活骑士卡详情
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        qrcode  query  string  true  "二维码详情, 可带`SUBSCRIBE:`, 也可不带"
// @Success      200  {object}  model.SubscribeActiveInfo  "请求成功"
func (*business) Inactive(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.QRQueryReq](c)
    detail, _ := service.NewBusinessEmployeeWithEmployee(ctx.Employee).Inactive(req.Qrcode)
    return ctx.SendResponse(detail)
}

// Active
// @ID           RiderOrderActive
// @Router       /employee/v1/subscribe/active [POST]
// @Summary      E2003 激活骑士卡
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.QRPostReq  true  "二维码可带`SUBSCRIBE:`, 也可不带"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*business) Active(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.QRPostReq](c)
    service.NewBusinessEmployeeWithEmployee(ctx.Employee).Active(req)
    return ctx.SendResponse()
}

// Pause
// @ID           EmployeeBusinessPause
// @Router       /employee/v1/business/pause [POST]
// @Summary      E2004 寄存电池
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        body  body     model.BusinessSubscribeID  true  "寄存请求"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*business) Pause(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.BusinessSubscribeID](c)
    service.NewBusinessRiderWithEmployee(ctx.Employee).Pause(req.SubscribeID)
    return ctx.SendResponse()
}

// Continue
// @ID           EmployeeBusinessContinue
// @Router       /employee/v1/business/continue [POST]
// @Summary      E2005 结束寄存电池
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        body  body     model.BusinessSubscribeID  true  "结束寄存请求"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*business) Continue(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.BusinessSubscribeID](c)
    service.NewBusinessRiderWithEmployee(ctx.Employee).Continue(req.SubscribeID)
    return ctx.SendResponse()
}

// UnSubscribe
// @ID           EmployeeBusinessUnSubscribe
// @Router       /employee/v1/business/unsubscribe [POST]
// @Summary      E2006 退租
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token   header  string  true  "店员校验token"
// @Param        body  body     model.UnsubscribeEmployeeReq  true  "退租请求"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*business) UnSubscribe(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.UnsubscribeEmployeeReq](c)
    service.NewBusinessRiderWithEmployee(ctx.Employee).UnSubscribe(req.SubscribeID)
    return ctx.SendResponse()
}

// List
// @ID           EmployeeBusinessList
// @Router       /employee/v1/business [GET]
// @Summary      E2007 业务记录
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        query  query   model.BusinessListReq  true  "列表请求筛选参数"
// @Success      200  {object}  model.PaginationRes{items=[]model.BusinessEmployeeListRes}  "请求成功"
func (*business) List(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.BusinessListReq](c)
    return ctx.SendResponse(service.NewBusinessWithEmployee(ctx.Employee).ListEmployee(req))
}
