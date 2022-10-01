// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-01
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type ebike struct{}

var Ebike = new(ebike)

// BrandList
// @ID           ManagerEbikeBrandList
// @Router       /manager/v1/ebike/brand [GET]
// @Summary      MI001 品牌列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []model.EbikeBrand  "请求成功"
func (*ebike) BrandList(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewEbikeBrand().All())
}

// BrandCreate
// @ID           ManagerEbikeBrandCreate
// @Router       /manager/v1/ebike/brand [POST]
// @Summary      MI002 创建品牌
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.EbikeBrandCreateReq  true  "品牌详情"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*ebike) BrandCreate(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EbikeBrandCreateReq](c)
    service.NewEbikeBrand(ctx.Modifier).Create(req)
    return ctx.SendResponse()
}

// BrandModify
// @ID           ManagerEbikeBrandModify
// @Router       /manager/v1/ebike/brand/:id [PUT]
// @Summary      MI003 修改品牌
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id    path     uint64  true  "品牌ID"
// @Param        body  body     model.EbikeBrandModifyReq  true  "品牌详情"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*ebike) BrandModify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EbikeBrandModifyReq](c)
    service.NewEbikeBrand(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}
