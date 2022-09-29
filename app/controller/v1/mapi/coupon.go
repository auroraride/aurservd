// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-29
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type coupon struct{}

var Coupon = new(coupon)

// TemplateList
// @ID           ManagerCouponTemplateList
// @Router       /manager/v1/coupon/template [GET]
// @Summary      MH001 模板列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.CouponTemplateListReq  false  "筛选选项"
// @Success      200  {object}  model.PaginationRes{items=[]model.CouponTemplateListRes}  "请求成功"
func (*coupon) TemplateList(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CouponTemplateListReq](c)
    return ctx.SendResponse(service.NewCouponTemplateWithModifier(ctx.Modifier).List(req))
}

// TemplateCreate
// @ID           ManagerCouponTemplateCreate
// @Router       /manager/v1/coupon/template [POST]
// @Summary      MH002 创建模板
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.CouponTemplateCreateReq  true  "模板内容"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*coupon) TemplateCreate(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CouponTemplateCreateReq](c)
    service.NewCouponTemplateWithModifier(ctx.Modifier).Create(req)
    return ctx.SendResponse()
}

// TemplateStatus
// @ID           ManagerCouponTemplateStatus
// @Router       /manager/v1/coupon/template/status [POST]
// @Summary      MH003 模板启用/禁用
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.CouponTemplateStatusReq  true  "模板信息"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*coupon) TemplateStatus(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CouponTemplateStatusReq](c)
    service.NewCouponTemplateWithModifier(ctx.Modifier).Status(req)
    return ctx.SendResponse()
}

// Generate
// @ID           ManagerCouponGenerate
// @Router       /manager/v1/coupon/generate [POST]
// @Summary      MH004 生成优惠券
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.CouponGenerateReq  true  "优惠券信息"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*coupon) Generate(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CouponGenerateReq](c)
    service.NewCouponWithModifier(ctx.Modifier).Generate(req)
    return ctx.SendResponse()
}

// Assembly
// @ID           ManagerCouponAssembly
// @Router       /manager/v1/coupon/assembly [GET]
// @Summary      MH005 发券记录
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.CouponAssemblyListReq  false  "筛选条件"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*coupon) Assembly(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CouponAssemblyListReq](c)
    return ctx.SendResponse(service.NewCouponAssemblyWithModifier(ctx.Modifier).List(req))
}

// List
// @ID           ManagerCouponList
// @Router       /manager/v1/coupon [GET]
// @Summary      MH006 优惠券列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*coupon) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CouponListReq](c)
    return ctx.SendResponse(service.NewCouponWithModifier(ctx.Modifier).List(req))
}

// Allocate
// @ID           ManagerCouponAllocate
// @Router       /manager/v1/coupon/allocate [POST]
// @Summary      MH007 分配优惠券
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.CouponAllocateReq  true  "优惠券和骑手信息"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*coupon) Allocate(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CouponAllocateReq](c)
    service.NewCouponWithModifier(ctx.Modifier).Allocate(req)
    return ctx.SendResponse()
}
