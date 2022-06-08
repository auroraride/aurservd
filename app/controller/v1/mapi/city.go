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
// @Summary      M2001 城市列表
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        status           query  model.CityListReq  false  "启用状态"
// @Success      200  {object}  []model.CityItem  "请求成功"
func (*city) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CityListReq](c)

    return ctx.SendResponse(service.NewCity().List(req))
}

// Modify
// @ID           CityModify
// @Router       /manager/v1/city/{id} [PUT]
// @Summary      M2002 修改城市
// @Description  desc
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id    path  int  true  "城市ID"
// @Param        body  body  model.CityModifyReq  true  "城市数据"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*city) Modify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.CityModifyReq](c)

    return ctx.SendResponse(
        model.CityModifyRes{Open: service.NewCityWithModifier(ctx.Modifier).Modify(req)},
    )
}
