// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type export struct{}

var Export = new(export)

// List
// @ID		ManagerExportList
// @Router	/manager/v1/export [GET]
// @Summary	导出列表
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string												true	"管理员校验token"
// @Param	query			query		model.ExportListReq									false	"分页信息"
// @Success	200				{object}	model.PaginationRes{items=[]model.ExportListRes}	"请求成功"
func (*export) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ExportListReq](c)
	return ctx.SendResponse(service.NewExportWithModifier(ctx.Modifier).List(ctx.Manager, req))
}

// Download
// @ID		ManagerExportDownload
// @Router	/manager/v1/export/download/{sn} [GET]
// @Summary	下载文件
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	sn				path		string					true	"编号"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*export) Download(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ExportDownloadReq](c)
	return ctx.Attachment(service.NewExportWithModifier(ctx.Modifier).Download(req))
}

// Rider
// @ID		ManagerExportRider
// @Router	/manager/v1/export/rider [POST]
// @Summary	导出骑手
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.RiderListExport	false	"筛选条件"
// @Success	200				{object}	model.ExportRes			"请求成功"
func (*export) Rider(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.RiderListExport](c)
	return ctx.SendResponse(service.NewRiderWithModifier(ctx.Modifier).ListExport(req))
}

// StatementDetail
// @ID		ManagerExportStatementDetail
// @Router	/manager/v1/export/statement/detail [POST]
// @Summary	导出企业账单
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.StatementBillDetailExport	true	"筛选条件"
// @Success	200				{object}	model.ExportRes					"请求成功"
func (*export) StatementDetail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.StatementBillDetailExport](c)
	return ctx.SendResponse(service.NewEnterpriseStatementWithModifier(ctx.Modifier).DetailExport(req))
}

// StatementUsage
// @ID		ManagerExportStatementUsage
// @Router	/manager/v1/export/statement/usage [POST]
// @Summary	导出企业使用明细
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.StatementUsageExport	true	"筛选条件"
// @Success	200				{object}	model.ExportRes				"请求成功"
func (*export) StatementUsage(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.StatementUsageExport](c)
	return ctx.SendResponse(service.NewEnterpriseStatementWithModifier(ctx.Modifier).UsageExport(req))
}

// Order
// @ID		ManagerExportOrder
// @Router	/manager/v1/export/order [POST]
// @Summary	导出订单
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.OrderListExport	true	"筛选条件"
// @Success	200				{object}	model.ExportRes			"请求成功"
func (*export) Order(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.OrderListExport](c)
	return ctx.SendResponse(service.NewOrderWithModifier(ctx.Modifier).Export(req))
}

// Commission
// @ID		ManagerExportCommission
// @Router	/manager/v1/export/commission [POST]
// @Summary	导出业绩
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.EmployeeActivityExportReq	true	"筛选条件"
// @Success	200				{object}	model.ExportRes					"请求成功"
func (*export) Commission(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.EmployeeActivityExportReq](c)
	return ctx.SendResponse(service.NewEmployeeWithModifier(ctx.Modifier).ActivityExport(req))
}

// Business
// @ID		ManagerExportBusiness
// @Router	/manager/v1/export/business [POST]
// @Summary	导出业务记录
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.BusinessExportReq	true	"筛选条件"
// @Success	200				{object}	model.ExportRes			"请求成功"
func (*export) Business(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BusinessExportReq](c)
	return ctx.SendResponse(service.NewBusinessWithModifier(ctx.Modifier).Export(req))
}

// StockDetail
// @ID		ManagerExportStockDetail
// @Router	/manager/v1/export/stock-detail [POST]
// @Summary	导出出入库明细
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.StockDetailExportReq	true	"筛选条件"
// @Success	200				{object}	model.ExportRes				"请求成功"
// func (*export) StockDetail(c echo.Context) (err error) {
// 	ctx, req := app.ManagerContextAndBinding[model.StockDetailExportReq](c)
// 	return ctx.SendResponse(service.NewStockWithModifier(ctx.Modifier).Export(req))
// }

// Exchange
// @ID		ManagerExportExchange
// @Router	/manager/v1/export/exchange [POST]
// @Summary	导出换电明细
// @Tags	导出
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.ExchangeListExport	true	"筛选条件"
// @Success	200				{object}	model.ExportRes				"请求成功"
func (*export) Exchange(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ExchangeListExport](c)
	return ctx.SendResponse(service.NewExchangeWithModifier(ctx.Modifier).Export(req))
}
