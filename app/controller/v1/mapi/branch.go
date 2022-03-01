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

// List 网点列表
func (*branch) List(c echo.Context) (err error) {
    req := new(model.BranchListReq)
    app.GetManagerContext(c).BindValidate(req)

    return app.NewResponse(c).SetData(service.NewBranch().List(req)).Send()
}

// Add 新增网点
func (*branch) Add(c echo.Context) (err error) {
    req := new(model.Branch)
    ctx := app.GetManagerContext(c)
    ctx.BindValidate(req)
    service.NewBranch().Add(req, ctx.Modifier)
    return app.NewResponse(c).Send()
}

// Modify 编辑网点
func (*branch) Modify(c echo.Context) (err error) {
    req := new(model.Branch)
    ctx := app.GetManagerContext(c)
    ctx.BindValidate(req)
    service.NewBranch().Modify(req, ctx.Modifier)
    return app.NewResponse(c).Send()
}
