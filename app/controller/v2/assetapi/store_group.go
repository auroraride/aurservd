// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-22, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
)

type storeGroup struct{}

var StoreGroup = new(storeGroup)

// List
// @ID		StoreGroupList
// @Router	/manager/v2/asset/store_group [GET]
// @Summary	门店集合列表
// @Tags	StoreGroup - 仓库
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Success	200						{object}	[]definition.StoreGroupListRes	"请求成功"
func (*storeGroup) List(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(biz.NewStoreGroup().List())
}

// Create
// @ID		StoreGroupCreate
// @Router	/manager/v2/asset/store_group [POST]
// @Summary	门店集合创建
// @Tags	StoreGroup - 仓库
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string							true	"管理员校验token"
// @Param	body					body		definition.StoreGroupCreateRep	true	"创建参数"
// @Success	200						{object}	model.StatusResponse			"请求成功"
func (*storeGroup) Create(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.StoreGroupCreateRep](c)
	return ctx.SendResponse(biz.NewStoreGroupWithModifier(ctx.Modifier).Create(req))
}

// Delete
// @ID		StoreGroupDelete
// @Router	/manager/v2/asset/store_group/{id} [DELETE]
// @Summary	门店集合删除
// @Tags	StoreGroup - 仓库
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	id						path		uint64					true	"资产ID"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*storeGroup) Delete(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(biz.NewStoreGroupWithModifier(ctx.Modifier).Delete(req.ID))
}
