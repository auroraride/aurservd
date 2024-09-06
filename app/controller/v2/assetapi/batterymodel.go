// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-13, by Jorjan

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type batterymodel struct{}

var BatteryModel = new(batterymodel)

// List
// @ID		BatteryModelList
// @Router	/manager/v2/asset/batterymodel [GET]
// @Summary	列表
// @Tags	电池型号 - BatteryModel
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string														true	"管理员校验token"
// @Param	query					query		model.BatteryModelListReq									true	"desc"
// @Success	200						{object}	model.PaginationRes{items=[]definition.BatteryModelDetail}	"请求成功"
func (*batterymodel) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.BatteryModelListReq](c)
	return ctx.SendResponse(service.NewBatteryModel().List(req))
}

// Create
// @ID		BatteryModelCreate
// @Router	/manager/v2/asset/batterymodel [POST]
// @Summary	创建
// @Tags	电池型号 - BatteryModel
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string						true	"管理员校验token"
// @Param	body					body		model.BatteryModelCreateReq	true	"desc"
// @Success	200						{object}	model.StatusResponse		"请求成功"
func (*batterymodel) Create(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.BatteryModelCreateReq](c)
	return ctx.SendResponse(service.NewBatteryModelWithModifier(ctx.Modifier).Create(req))
}

// Delete
// @ID		BatteryModelDelete
// @Router	/manager/v2/asset/batterymodel/{id} [DELETE]
// @Summary	删除
// @Tags	电池型号 - BatteryModel
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	id						path		string					true	"仓库ID"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*batterymodel) Delete(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewBatteryModelWithModifier(ctx.Modifier).Delete(req.ID))
}
