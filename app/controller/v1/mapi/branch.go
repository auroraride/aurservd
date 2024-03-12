// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type branch struct {
}

var Branch = new(branch)

// List
// @ID		BranchList
// @Router	/manager/v1/branch [GET]
// @Summary	网点列表
// @Tags	网点
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string											true	"管理员校验token"
// @Param	query			query		model.BranchListReq								true	"desc"
// @Success	200				{object}	model.PaginationRes{items=[]model.BranchItem}	"请求成功"
func (*branch) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BranchListReq](c)
	return ctx.SendResponse(service.NewBranch().List(req))
}

// Selector
// @ID		BranchSelector
// @Router	/manager/v1/branch/selector [GET]
// @Summary	网点选择列表
// @Tags	网点
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string												true	"管理员校验token"
// @Success	200				{object}	model.ItemListRes{items=[]model.BranchSampleItem}	"请求成功"
func (*branch) Selector(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewBranch().Selector())
}

// Create
// @ID		BranchCreate
// @Router	/manager/v1/branch [POST]
// @Summary	新增网点
// @Tags	网点
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.BranchCreateReq	true	"网点数据"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*branch) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BranchCreateReq](c)
	service.NewBranchWithModifier(ctx.Modifier).Create(req)
	return ctx.SendResponse()
}

// Modify
// @ID		BranchModify
// @Router	/manager/v1/branch/{id} [PUT]
// @Summary	编辑网点
// @Tags	网点
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.BranchModifyReq	true	"网点数据"
// @Param	id				path		int						true	"网点ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*branch) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BranchModifyReq](c)
	service.NewBranchWithModifier(ctx.Modifier).Modify(req)
	return ctx.SendResponse()
}

// AddContract
// @ID		BranchAddContract
// @Router	/manager/v1/{id}/contract [POST]
// @Summary	新增合同
// @Tags	网点
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.BranchContract	true	"合同数据"
// @Param	id				path		int						true	"网点ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*branch) AddContract(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BranchContract](c)
	service.NewBranchWithModifier(ctx.Modifier).AddContract(req.BranchID, req)
	return ctx.SendResponse()
}

// Sheet
// @ID		ManagerBranchSheet
// @Router	/manager/v1/branch/contract/sheet [POST]
// @Summary	修改合同底单
// @Tags	网点
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		model.BranchContractSheetReq	true	"合同底单修改请求"
// @Success	200				{object}	model.StatusResponse			"请求成功"
func (*branch) Sheet(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BranchContractSheetReq](c)
	service.NewBranchWithModifier(ctx.Modifier).Sheet(req)
	return ctx.SendResponse()
}

// Nearby
// @ID		ManagerBranchNearby
// @Router	/manager/v1/branch/nearby [GET]
// @Summary	查找附近的网点
// @Tags	网点
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	query			query		model.BranchDistanceListReq		false	"筛选选项"
// @Success	200				{object}	[]model.BranchDistanceListRes	"请求成功"
func (*branch) Nearby(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BranchDistanceListReq](c)
	return ctx.SendResponse(service.NewBranch().ListByDistanceManager(req))
}
