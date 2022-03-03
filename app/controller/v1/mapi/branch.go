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
// @Summary      M4.网点列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  model.PaginationRes{items=[]model.Branch}  "请求成功"
func (*branch) List(c echo.Context) (err error) {
    req := new(model.BranchListReq)
    app.GetManagerContext(c).BindValidate(req)

    return app.NewResponse(c).SetData(service.NewBranch().List(req)).Send()
}

// Add
// @ID           BranchAdd
// @Router       /manager/v1/branch [POST]
// @Summary      M5.新增网点
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.Branch  true  "网点数据"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*branch) Add(c echo.Context) (err error) {
    req := new(model.Branch)
    ctx := app.GetManagerContext(c)
    ctx.BindValidate(req)
    service.NewBranch().Add(req, ctx.Modifier)
    return app.NewResponse(c).Send()
}

// Modify
// @ID           BranchModify
// @Router       /manager/v1/branch/{id} [PUT]
// @Summary      M6.编辑网点
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.Branch  true  "网点数据"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*branch) Modify(c echo.Context) (err error) {
    req := new(model.Branch)
    ctx := app.GetManagerContext(c)
    ctx.BindValidate(req)
    service.NewBranch().Modify(req, ctx.Modifier)
    return app.NewResponse(c).Send()
}

// AddContract
// @ID           BranchAddContract
// @Router       /manager/v1/{id}/contract [POST]
// @Summary      M7.新增合同
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.BranchContract  true  "合同数据"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*branch) AddContract(c echo.Context) (err error) {
    req := new(model.BranchContract)
    ctx := app.GetManagerContext(c)
    ctx.BindValidate(req)
    service.NewBranch().AddContract(req.BranchID, req, ctx.Modifier)
    return app.NewResponse(c).Send()
}