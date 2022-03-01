// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type city struct{}

var City = new(city)

// List 城市列表
func (*city) List(c echo.Context) (err error) {
    req := new(model.CityListReq)
    app.GetManagerContext(c).BindValidate(req)

    return app.NewResponse(c).SetData(service.NewCity().List(req)).Send()
}

// Modify 修改城市
func (*city) Modify(c echo.Context) (err error) {
    req := new(model.CityModifyReq)
    app.GetManagerContext(c).BindValidate(req)
    return app.NewResponse(c).SetData(model.CityModifyRes{Open: service.NewCity().Modify(req)}).Send()
}