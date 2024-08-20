// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-17, by aurb

package assetapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/biz"
	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/pkg/snag"
)

type assetManager struct{}

var AssetManager = new(assetManager)

// Signin
// @ID		AssetManagerSignin
// @Router	/manager/v2/asset/user/signin [POST]
// @Summary	管理员登录
// @Tags	管理 - AssetManager
// @Accept	json
// @Produce	json
// @Param	body	body		definition.AssetManagerSigninReq	true	"desc"
// @Success	200		{object}	model.ManagerSigninRes				"请求成功"
func (*assetManager) Signin(c echo.Context) (err error) {
	ctx, req := app.ContextBinding[definition.AssetManagerSigninReq](c)
	data, err := biz.NewAssetManager().Signin(req)
	if err != nil {
		snag.Panic(err)
	}
	return ctx.SendResponse(data)
}

// Create
// @ID		AssetManagerCreate
// @Router	/manager/v2/asset/user [POST]
// @Summary	新增管理员
// @Tags	管理 - AssetManager
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string								true	"管理员校验token"
// @Param	body					body		definition.AssetManagerCreateReq	true	"desc"
// @Success	200						{object}	model.StatusResponse				"请求成功"
func (*assetManager) Create(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.AssetManagerCreateReq](c)
	err = biz.NewAssetManager().Create(req)
	if err != nil {
		snag.Panic(err)
	}
	return ctx.SendResponse()
}

// List
// @ID		AssetManagerList
// @Router	/manager/v2/asset/user [GET]
// @Summary	列举管理员
// @Tags	管理 - AssetManager
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string														true	"管理员校验token"
// @Param	query					query		definition.AssetManagerListReq								true	"desc"
// @Success	200						{object}	model.PaginationRes{items=[]definition.AssetManagerListRes}	"请求成功"
func (*assetManager) List(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.AssetManagerListReq](c)
	return ctx.SendResponse(biz.NewAssetManagerWithModifier(ctx.Modifier).List(req))
}

// Delete
// @ID		AssetManagerDelete
// @Router	/manager/v2/asset/user/{id} [DELETE]
// @Summary	删除管理员
// @Tags	管理 - AssetManager
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string					true	"管理员校验token"
// @Param	id						path		uint64					true	"管理员ID"
// @Success	200						{object}	model.StatusResponse	"请求成功"
func (*assetManager) Delete(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[model.IDParamReq](c)
	biz.NewAssetManagerWithModifier(ctx.Modifier).Delete(req)
	return ctx.SendResponse()
}

// Modify
// @ID		AssetManagerModify
// @Router	/manager/v2/asset/user/{id} [PUT]
// @Summary	编辑管理员
// @Tags	管理 - AssetManager
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string								true	"管理员校验token"
// @Param	id						path		uint64								true	"管理员ID"
// @Param	body					body		definition.AssetManagerModifyReq	true	"编辑属性"
// @Success	200						{object}	model.StatusResponse				"请求成功"
func (*assetManager) Modify(c echo.Context) (err error) {
	ctx, req := app.AssetManagerContextAndBinding[definition.AssetManagerModifyReq](c)
	biz.NewAssetManagerWithModifier(ctx.Modifier).Modify(req)
	return ctx.SendResponse()
}

// Profile
// @ID		AssetManagerProfile
// @Router	/manager/v2/asset/user/profile [GET]
// @Summary	管理员信息
// @Tags	管理 - AssetManager
// @Accept	json
// @Produce	json
// @Param	X-Asset-Manager-Token	header		string	true	"管理员校验token"
// @Success	200						{object}	model.ManagerSigninRes
func (*assetManager) Profile(c echo.Context) (err error) {
	ctx := app.ContextX[app.AssetManagerContext](c)
	return ctx.SendResponse(biz.NewAssetManager().Profile(ctx.AssetManager))
}
