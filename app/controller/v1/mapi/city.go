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

// List
// ID            CityList
// @Router       /manager/v1/city [GET]
// @Summary      M2.城市列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        status           query  model.CityListReq  false  "启用状态"
// @Success      200  {object}  []model.CityItem  "请求成功"
func (*city) List(c echo.Context) (err error) {
    req := new(model.CityListReq)
    app.GetManagerContext(c).BindValidate(req)

    return app.NewResponse(c).SetData(service.NewCity().List(req)).Send()
}

// Modify
// @ID           CityModify
// @Router       /manager/v1/city [PUT]
// @Summary      M3.修改城市
// @Description  desc
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id    path  int  true  "城市ID"
// @Param        body  body  model.CityModifyReq  true  "城市数据"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*city) Modify(c echo.Context) (err error) {
    req := new(model.CityModifyReq)
    ctx := app.GetManagerContext(c)
    ctx.BindValidate(req)
    return app.NewResponse(c).SetData(
        model.CityModifyRes{Open: service.NewCity().Modify(req, ctx.Modifier)},
    ).Send()
}
