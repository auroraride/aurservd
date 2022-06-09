// Copyright (C) liasica. 2021-present.
//
// Created at 2022/3/1
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type branch struct {
}

var Branch = new(branch)

// List
// @ID           BranchList
// @Router       /manager/v1/branch [GET]
// @Summary      M3001 网点列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query  model.BranchListReq  true  "desc"
// @Success      200  {object}  model.PaginationRes{items=[]model.BranchItem}  "请求成功"
func (*branch) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BranchListReq](c)
    return ctx.SendResponse(service.NewBranch().List(req))
}

// Selector
// @ID           BranchSelector
// @Router       /manager/v1/branch/selector [GET]
// @Summary      M3005 网点选择列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  model.ItemListRes{items=[]model.BranchSampleItem}  "请求成功"
func (*branch) Selector(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewBranch().Selector())
}

// Create
// @ID           BranchCreate
// @Router       /manager/v1/branch [POST]
// @Summary      M3002 新增网点
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.BranchCreateReq  true  "网点数据"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*branch) Create(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BranchCreateReq](c)
    service.NewBranchWithModifier(ctx.Modifier).Create(req)
    return ctx.SendResponse()
}

// Modify
// @ID           BranchModify
// @Router       /manager/v1/branch/{id} [PUT]
// @Summary      M3003 编辑网点
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.BranchModifyReq  true  "网点数据"
// @Param        id  path  int  true  "网点ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*branch) Modify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BranchModifyReq](c)
    service.NewBranchWithModifier(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}

// AddContract
// @ID           BranchAddContract
// @Router       /manager/v1/{id}/contract [POST]
// @Summary      M3004 新增合同
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.BranchContract  true  "合同数据"
// @Param        id  path  int  true  "网点ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*branch) AddContract(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BranchContract](c)
    service.NewBranchWithModifier(ctx.Modifier).AddContract(req.BranchID, req)
    return ctx.SendResponse()
}
