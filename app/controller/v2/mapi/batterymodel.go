// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-13, by Jorjan

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type batterymodel struct{}

var BatteryModel = new(batterymodel)

// List
// @ID		BatteryModelList
// @Router	/manager/v1/batterymodel [GET]
// @Summary	列表
// @Tags	电池型号 - BatteryModel
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	body			body		definition.BatteryModelListReq	true	"desc"
// @Success	200				{object}	[]definition.BatteryModelDetail	"请求成功"
func (*batterymodel) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.BatteryModelListReq](c)
	return ctx.SendResponse(biz.NewBatteryModel().List(req))
}

// Detail
// @ID		BatteryModelDetail
// @Router	/manager/v1/batterymodel/{id} [GET]
// @Summary	详情
// @Tags	电池型号 - BatteryModel
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string							true	"管理员校验token"
// @Param	id				path		string							true	"仓库ID"
// @Success	200				{object}	definition.BatteryModelDetail	"请求成功"
func (*batterymodel) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewBatteryModel().Detail(req.ID))
}

// Create
// @ID		BatteryModelCreate
// @Router	/manager/v1/batterymodel [POST]
// @Summary	创建
// @Tags	电池型号 - BatteryModel
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string								true	"管理员校验token"
// @Param	body			body		definition.BatteryModelCreateReq	true	"desc"
// @Success	200				{object}	model.StatusResponse				"请求成功"
func (*batterymodel) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.BatteryModelCreateReq](c)
	return ctx.SendResponse(biz.NewBatteryModelWithModifier(ctx.Modifier).Create(req))
}

// Delete
// @ID		BatteryModelDelete
// @Router	/manager/v1/batterymodel/{id} [DELETE]
// @Summary	删除
// @Tags	电池型号 - BatteryModel
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		string					true	"仓库ID"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*batterymodel) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewBatteryModelWithModifier(ctx.Modifier).Delete(req.ID))
}

// Modify
// @ID		BatteryModelModify
// @Router	/manager/v1/batterymodel/{id} [PUT]
// @Summary	修改
// @Tags	电池型号 - BatteryModel
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string								true	"管理员校验token"
// @Param	body			body		definition.BatteryModelModifyReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse				"请求成功"
func (*batterymodel) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[definition.BatteryModelModifyReq](c)
	return ctx.SendResponse(biz.NewBatteryModelWithModifier(ctx.Modifier).Modify(req))
}
