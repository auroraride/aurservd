// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-25
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type assistance struct{}

var Assistance = new(assistance)

// Detail
// @ID           EmployeeAssistanceDetail
// @Router       /employee/v1/assistance/{id} [GET]
// @Summary      E5001 获取救援详情
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        id  path  uint64  true  "救援ID"
// @Success      200  {object}  model.AssistanceEmployeeDetailRes  "请求成功"
func (*assistance) Detail(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.IDParamReq](c)
    return ctx.SendResponse(
        service.NewAssistanceWithEmployee(ctx.Employee).EmployeeDetail(req.ID),
    )
}

// Process
// @ID           EmployeeAssistanceProcess
// @Router       /employee/v1/assistance/process [POST]
// @Summary      E5002 处理救援
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        body  body     model.AssistanceProcessReq  true  "救援处理详情"
// @Success      200  {object}  model.AssistanceProcessRes  "请求成功"
func (*assistance) Process(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.AssistanceProcessReq](c)
    return ctx.SendResponse(service.NewAssistanceWithEmployee(ctx.Employee).Process(req))
}

// Pay
// @ID           EmployeeAssistancePay
// @Router       /employee/v1/assistance/pay [POST]
// @Summary      E5003 救援支付
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        body  body     model.AssistancePayReq  true  "支付信息"
// @Success      200  {object}  model.AssistancePayRes  "请求成功"
func (*assistance) Pay(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.AssistancePayReq](c)
    return ctx.SendResponse(service.NewAssistanceWithEmployee(ctx.Employee).Pay(req))
}

// PayStatus
// @ID           EmployeeAssistancePayStatus
// @Router       /employee/v1/assistance/pay [GET]
// @Summary      E5004 救援支付状态
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        outTradeNo     query  string  true  "订单编号"
// @Success      200 {object}   model.OrderStatusRes  "请求成功"
func (*assistance) PayStatus(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.OrderStatusReq](c)
    return ctx.SendResponse(service.NewOrder().QueryStatus(req))
}

// List
// @ID           EmployeeAssistanceList
// @Router       /rider/v1/assistance [GET]
// @Summary      E5005 救援列表
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        query  query   model.PaginationReq  true  "分页参数"
// @Success      200  {object}  model.Pagination{items=[]model.AssistanceSimpleListRes}  "请求成功"
func (*assistance) List(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.PaginationReq](c)
    return ctx.SendResponse(service.NewAssistanceWithEmployee(ctx.Employee).SimpleList(*req))
}
