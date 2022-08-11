// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-19
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type statement struct{}

var Statement = new(statement)

// GetBill
// @ID           ManagerStatementGetBill
// @Router       /manager/v1/enterprise/bill [GET]
// @Summary      M9011 获取账单
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.StatementBillReq  true  "账单请求参数"
// @Success      200  {object}  model.StatementBillRes  "请求成功"
func (*statement) GetBill(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.StatementBillReq](c)
    return ctx.SendResponse(service.NewEnterpriseStatementWithModifier(ctx.Modifier).GetBill(req))
}

// Bill
// @ID           ManagerStatementBill
// @Router       /manager/v1/enterprise/bill [POST]
// @Summary      M9012 结账
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.StatementClearBillReq  true  "结账请求"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*statement) Bill(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.StatementClearBillReq](c)
    service.NewEnterpriseStatementWithModifier(ctx.Modifier).Bill(req)
    return ctx.SendResponse()
}

// Historical
// @ID           ManagerStatementHistorical
// @Router       /manager/v1/enterprise/bill/historical [GET]
// @Summary      M9013 历史账单
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.StatementBillHistoricalListReq  true  "请求数据"
// @Success      200  {object}  model.PaginationRes{items=[]model.StatementBillHistoricalListRes}  "请求成功"
func (*statement) Historical(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.StatementBillHistoricalListReq](c)
    return ctx.SendResponse(service.NewEnterpriseStatementWithModifier(ctx.Modifier).Historical(req))
}

// Statement
// @ID           ManagerStatementStatement
// @Router       /manager/v1/enterprise/bill/statement [GET]
// @Summary      M9014 结账明细
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.StatementBillDetailReq  true  "请求详情"
// @Success      200  {object}  model.StatementDetail  "请求成功"
func (*statement) Statement(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.StatementBillDetailReq](c)
    return ctx.SendResponse(service.NewEnterpriseStatementWithModifier(ctx.Modifier).Statement(req.ID))
}

// Usage
// @ID           ManagerStatementUsage
// @Router       /manager/v1/enterprise/bill/usage [GET]
// @Summary      M9015 使用明细
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.StatementUsageReq  true  "筛选项"
// @Success      200  {object}  model.Pagination{items=[]model.StatementUsageRes}  "请求成功"
func (*statement) Usage(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.StatementUsageReq](c)
    return ctx.SendResponse(service.NewEnterpriseStatementWithModifier(ctx.Modifier).Usage(req))
}
