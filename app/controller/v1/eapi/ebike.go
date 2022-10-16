// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-14
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type ebike struct{}

var Ebike = new(ebike)

// Unallocated
// @ID           EmployeeEbikeUnallocated
// @Router       /employee/v1/ebike/unallocated [GET]
// @Summary      E6001 获取未分配电车信息
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        keyword  query  string  true  "关键词"
// @Success      200  {object}  model.EbikeInfo  "电车信息"
func (*ebike) Unallocated(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.KeywordQueryReq](c)
    return ctx.SendResponse(service.NewEbikeAllocate(ctx.Employee).UnallocatedInfo(req.Keyword))
}

// Allocate
// @ID           EmployeeEbikeAllocate
// @Router       /employee/v1/ebike/allocate [POST]
// @Summary      E6002 分配车辆
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        body  body     model.EbikeAllocateReq  true  "分配请求"
// @Success      200  {object}  model.EbikeAllocateRes  "请求成功"
func (*ebike) Allocate(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.EbikeAllocateReq](c)
    return ctx.SendResponse(service.NewEbikeAllocate(ctx.Employee).Allocate(req))
}

// Info
// @ID           EmployeeEbikeInfo
// @Router       /employee/v1/ebike/allocate/info [GET]
// @Summary      E6003 车辆分配信息
// @Description  骑手签约成功后通过socket推送门店消息
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        allocateId     query  string  true  "分配ID"
// @Success      200  {object}  model.EbikeAllocateInfo  "请求成功"
func (*ebike) Info(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.EbikeAllocateIDQueryReq](c)
    return ctx.SendResponse(service.NewEbikeAllocate(ctx.Employee).Info(req))
}

// List
// @ID           EmployeeEbikeList
// @Router       /employee/v1/ebike/allocate [GET]
// @Summary      E6004 车辆分配记录
// @Tags         [E]店员接口
// @Accept       json
// @Produce      json
// @Param        X-Employee-Token  header  string  true  "店员校验token"
// @Param        query  query   model.EbikeAllocateEmployeeListReq  false  "筛选选项"
// @Success      200  {object}  model.PaginationRes{items=[]model.EbikeAllocateInfo}  "请求成功"
func (*ebike) List(c echo.Context) (err error) {
    ctx, req := app.EmployeeContextAndBinding[model.EbikeAllocateEmployeeListReq](c)
    return ctx.SendResponse(service.NewEbikeAllocate(ctx.Employee).EmployeeList(req))
}
