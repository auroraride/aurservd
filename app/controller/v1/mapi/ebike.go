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

// List
// @ID           ManagerEbikeList
// @Router       /manager/v1/ebike [GET]
// @Summary      MI004 电车列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.EbikeListReq  false  "筛选条件"
// @Success      200  {object}  model.PaginationRes{item=[]model.EbikeListRes}  "请求成功"
func (*ebike) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EbikeListReq](c)
    return ctx.SendResponse(service.NewEbike(ctx.Modifier).List(req))
}

// Create
// @ID           ManagerEbikeCreate
// @Router       /manager/v1/ebike [POST]
// @Summary      MI005 添加电车
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.EbikeCreateReq  true  "电车信息"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*ebike) Create(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EbikeCreateReq](c)
    service.NewEbike(ctx.Modifier).Create(req)
    return ctx.SendResponse()
}

// Modify
// @ID           ManagerEbikeModify
// @Router       /manager/v1/ebike/:id [PUT]
// @Summary      MI006 修改电车
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.EbikeModifyReq  true  "电车信息"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*ebike) Modify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.EbikeModifyReq](c)
    service.NewEbike(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}

// BatchCreate
// @ID           ManagerEbikeBatchCreate
// @Router       /manager/v1/ebike/batch [POST]
// @Summary      MI007 批量导入电车
// @Tags         [M]管理接口
// @Accept       mpfd
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        file  formData  file  true  "电车信息"
// @Success      200  {object}   []string  "请求成功"
func (*ebike) BatchCreate(c echo.Context) (err error) {
    ctx := app.ContextX[app.ManagerContext](c)
    return ctx.SendResponse(service.NewEbike(ctx.Modifier).BatchCreate(ctx.Context))
}
