// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-14
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "errors"
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type allocate struct{}

var Allocate = new(allocate)

// UnallocatedEbike
// @ID           EmployeeUnallocatedEbike
// @Router       /employee/v1/allocate/ebike [GET]
// @Summary      E6001 获取未分配电车信息
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        keyword  query  string  true  "关键词"
// @Success      200  {object}  model.Ebike  "电车信息"
func (*allocate) UnallocatedEbike(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.KeywordQueryReq](c)
    return ctx.SendResponse(service.NewAllocate(ctx.Employee).UnallocatedEbikeInfo(req.Keyword))
}

// Allocate
// @ID           EmployeeAllocate
// @Router       /employee/v1/allocate [POST]
// @Summary      E6002 分配车辆 (废弃)
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
func (*allocate) Allocate(c echo.Context) (err error) {
    return errors.New("接口已废弃")
    // return ctx.SendResponse(service.NewAllocate(ctx.Employee).AllocateEbike(req))
}

// Info
// @ID           EmployeeAllocateInfo
// @Router       /employee/v1/allocate/info/{id} [GET]
// @Summary      E6003 分配详情
// @Description  骑手签约成功后通过socket推送门店消息
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        id   path      uint64  true  "分配ID"
// @Success      200  {object}  model.AllocateDetail  "请求成功"
func (*allocate) Info(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.IDParamReq](c)
    return ctx.SendResponse(service.NewAllocate(ctx.Employee).Info(req))
}

// List
// @ID           EmployeeAllocateList
// @Router       /employee/v1/allocate [GET]
// @Summary      E6004 分配记录
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        query  query   model.AllocateEmployeeListReq  false  "筛选选项"
// @Success      200  {object}  model.PaginationRes{items=[]model.AllocateDetail}  "请求成功"
func (*allocate) List(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.AllocateEmployeeListReq](c)
    return ctx.SendResponse(service.NewAllocate(ctx.Employee).EmployeeList(req))
}
