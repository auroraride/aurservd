// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-05
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type enterprise struct{}

var Enterprise = new(enterprise)

// Create
// @ID           ManagerEnterpriseCreate
// @Router       /manager/v1/enterprise [POST]
// @Summary      M90001 创建企业
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.EnterpriseDetail  true  "desc"
// @Success      200  {object}  int  "请求成功"
func (*enterprise) Create(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseDetail](c)
    return ctx.SendResponse(service.NewEnterpriseWithModifier(ctx.Modifier).Create(req))
}

// Modify
// @ID           ManagerEnterpriseModify
// @Router       /manager/v1/enterprise/{id} [PUT]
// @Summary      M90002 修改企业
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.EnterpriseDetailWithID  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*enterprise) Modify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseDetailWithID](c)
    service.NewEnterpriseWithModifier(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}

// List
// @ID           ManagerEnterpriseList
// @Router       /manager/v1/enterprise [GET]
// @Summary      M90003 列举企业
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query  model.EnterpriseListReq  true  "desc"
// @Success      200  {object}  model.PaginationRes{items=[]model.EnterpriseListRes} "请求成功"
func (*enterprise) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EnterpriseListReq](c)
    return ctx.SendResponse(
        service.NewEnterpriseWithModifier(ctx.Modifier).List(req),
    )
}