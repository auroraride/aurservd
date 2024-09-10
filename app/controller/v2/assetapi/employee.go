// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-23, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type employee struct{}

var Employee = new(employee)

// List
// @ID		EmployeeList
// @Router	/manager/v2/asset/employee [GET]
// @Summary	列表
// @Tags	门店店员 - Employee
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string													true	"管理员校验token"
// @Param	query					query		definition.EmployeeListReq								true	"desc"
// @Success	200						{object}	model.PaginationRes{items=[]definition.EmployeeListRes}	"请求成功"
func (*employee) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.EmployeeListReq](c)
	return ctx.SendResponse(biz.NewEmployee().List(req))
}

// Create
// @ID		EmployeeCreate
// @Router	/manager/v2/asset/employee [POST]
// @Summary	创建
// @Tags	门店店员 - Employee
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	body					body		definition.EmployeeCreateReq	true	"请求参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*employee) Create(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.EmployeeCreateReq](c)
	return ctx.SendResponse(biz.NewEmployeeWithModifier(ctx.Modifier).Create(req))
}

// Delete
// @ID		EmployeeDelete
// @Router	/manager/v2/asset/employee/{id} [DELETE]
// @Summary	删除
// @Tags	门店店员 - Employee
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	id						path		string					true	"仓库ID"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*employee) Delete(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewEmployeeWithModifier(ctx.Modifier).Delete(req.ID))
}

// Modify
// @ID		EmployeeModify
// @Router	/manager/v2/asset/employee/{id} [PUT]
// @Summary	修改
// @Tags	门店店员 - Employee
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	id						path		int								true	"ID"
// @Param	body					body		definition.EmployeeModifyReq	true	"请求参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*employee) Modify(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.EmployeeModifyReq](c)
	return ctx.SendResponse(biz.NewEmployeeWithModifier(ctx.Modifier).Modify(req))
}
