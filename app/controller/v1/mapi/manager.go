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
	"github.com/auroraride/aurservd/pkg/snag"
)

type manager struct {
}

var (
	Manager = new(manager)
)

// Signin
// @ID           ManagerSignin
// @Router       /manager/v1/user/signin [POST]
// @Summary      M1001 用户登录
// @Description  管理员登录
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.ManagerSigninRes  "请求成功"
func (*manager) Signin(c echo.Context) (err error) {
	ctx, req := app.ContextBinding[model.ManagerSigninReq](c)
	data, err := service.NewManager().Signin(req)
	if err != nil {
		snag.Panic(err)
	}
	return ctx.SendResponse(data)
}

// Create
// @ID           ManagerCreate
// @Router       /manager/v1/user [POST]
// @Summary      M1002 新增管理员
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.ManagerCreateReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*manager) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ManagerCreateReq](c)
	err = service.NewManager().Create(req)
	if err != nil {
		snag.Panic(err)
	}
	return ctx.SendResponse()
}

// List
// @ID           ManagerManagerList
// @Router       /manager/v1/user [GET]
// @Summary      M1003 列举管理员
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query  model.ManagerListReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*manager) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ManagerListReq](c)
	return ctx.SendResponse(service.NewManagerWithModifier(ctx.Modifier).List(req))
}

// Delete
// @ID           ManagerManagerDelete
// @Router       /manager/v1/user/{id} [DELETE]
// @Summary      M1004 删除管理员
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "管理员ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*manager) Delete(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	service.NewManagerWithModifier(ctx.Modifier).Delete(req)
	return ctx.SendResponse()
}

// Modify
// @ID           ManagerModify
// @Router       /manager/v1/user/{id} [PUT]
// @Summary      M1005 编辑管理员
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path       uint64  true  "管理员ID"
// @Param        body  body     model.ManagerModifyReq  true  "编辑属性"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*manager) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ManagerModifyReq](c)
	service.NewManagerWithModifier(ctx.Modifier).Modify(req)
	return ctx.SendResponse()
}

// Profile
// @ID           ManagerProfile
// @Router       /manager/v1/user/profile [GET]
// @Summary      M1006 管理员信息
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  model.ManagerSigninRes
func (*manager) Profile(c echo.Context) (err error) {
	ctx := app.ContextX[app.ManagerContext](c)
	return ctx.SendResponse(service.NewManager().Profile(ctx.Manager))
}
