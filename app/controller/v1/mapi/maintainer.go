// Copyright (C) liasica. 2023-present.
//
// Created at 2023-08-10
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type maintainer struct{}

var Maintainer = new(maintainer)

// List
// @ID           ManagerMaintainerList
// @Router       /manager/v1/maintainer [GET]
// @Summary      MJ001 运维人员列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query  model.MaintainerListReq  false  "筛选参数"
// @Success      200  {object}  model.Pagination{items=[]model.Maintainer}  "请求成功"
func (*maintainer) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.MaintainerListReq](c)
	return ctx.SendResponse(service.NewMaintainer().List(req))
}

// Create
// @ID           ManagerMaintainerCreate
// @Router       /manager/v1/maintainer [POST]
// @Summary      MJ002 创建运维人员
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.MaintainerCreateReq  true  "请求参数"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*maintainer) Create(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.MaintainerCreateReq](c)
	service.NewMaintainer().Create(req)
	return ctx.SendResponse()
}

// Modify
// @ID           ManagerMaintainerModify
// @Router       /manager/v1/maintainer/{id} [POST]
// @Summary      MJ003 修改运维人员信息
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.MaintainerCreateReq  true  "请求参数"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*maintainer) Modify(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.MaintainerModifyReq](c)
	service.NewMaintainer().Modify(req)
	return ctx.SendResponse()
}
