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
// @Summary      E2003 骑手业务详情
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
    service.NewRiderMgrWithEmployee(ctx.Employee).PauseSubscribe(req.SubscribeID)
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
    service.NewRiderMgrWithEmployee(ctx.Employee).ContinueSubscribe(req.SubscribeID)
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
// @Param        body  body     model.BusinessSubscribeID  true  "退租请求"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*business) UnSubscribe(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.BusinessSubscribeID](c)
    service.NewRiderMgrWithEmployee(ctx.Employee).UnSubscribe(req.SubscribeID)
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
