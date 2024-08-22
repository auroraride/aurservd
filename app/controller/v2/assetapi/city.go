// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-22, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type city struct{}

var City = new(city)

// List
// @ID		CityList
// @Router	/manager/v2/asset/city [GET]
// @Summary	城市列表
// @Tags	City - 城市
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string				true	"管理员校验token"
// @Param	status					query		model.CityListReq	false	"启用状态"
// @Success	200						{object}	[]model.CityItem	"请求成功"
func (*city) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.CityListReq](c)
	return ctx.SendResponse(service.NewCity().List(req))
}

// Modify
// @ID		CityModify
// @Router	/manager/v2/asset/city/{id} [PUT]
// @Summary	修改城市
// @Tags	City - 城市
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	id						path		int						true	"城市ID"
// @Param	body					body		model.CityModifyReq		true	"城市数据"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*city) Modify(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.CityModifyReq](c)
	return ctx.SendResponse(
		model.CityModifyRes{Open: service.NewCityWithModifier(ctx.Modifier).Modify(req)},
	)
}
