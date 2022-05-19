// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-13
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type battery struct{}

var Battery = new(battery)

// ListModels
// @ID           BatteryModels
// @Router       /manager/v1/battery/model [GET]
// @Summary      M40001 获取电池型号
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  model.ItemListRes{items=[]model.BatteryModel}  "请求成功"
func (*battery) ListModels(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewBattery().ListModels())
}

// CreateModel
// @ID           BatteryCreateModel
// @Router       /manager/v1/battery/model [POST]
// @Summary      M40002 创建电池型号
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.BatteryModelCreateReq  true  "电池型号数据"
// @Success      200  {object}  model.ItemRes{item=model.BatteryModel}  "请求成功"
func (*battery) CreateModel(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.BatteryModelCreateReq](c)
    return ctx.SendResponse(model.ItemRes{Item: service.NewBattery().CreateModel(ctx.Modifier, req)})
}
