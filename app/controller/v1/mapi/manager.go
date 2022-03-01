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

type manager struct {
}

var (
    Manager = new(manager)
)

// Signin 管理员登录
func (*manager) Signin(c echo.Context) (err error) {
    req := new(model.ManagerSigninReq)
    ctx := c.(*app.Context)
    ctx.BindValidate(req)

    data, err := service.NewManager().Signin(req)
    if err != nil {
        return
    }
    return app.NewResponse(c).SetData(data).Send()
}

func (*manager) Add(c echo.Context) (err error) {
    req := new(model.ManagerAddReq)
    ctx := c.(*app.ManagerContext)
    ctx.BindValidate(req)

    err = service.NewManager().Add(req)
    if err != nil {
        return err
    }
    return app.NewResponse(c).Send()
}
