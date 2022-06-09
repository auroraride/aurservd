// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type setting struct{}

var Setting = new(setting)

// List
// @ID           ManagerSettingList
// @Router       /manager/v1/setting [GET]
// @Summary      M1010 列举设置
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  []model.SettingReq  "请求成功"
func (*setting) List(c echo.Context) (err error) {
    ctx := app.Context(c)

    return ctx.SendResponse(service.NewSetting().List())
}

// Modify
// @ID           ManagerSettingModify
// @Router       /manager/v1/setting/{key} [PUT]
// @Summary      M1011 调整设置
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        key  path  string  true  "设置项"
// @Param        body  body  model.SettingReq  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*setting) Modify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.SettingReq](c)
    service.NewSettingWithModifier(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}
