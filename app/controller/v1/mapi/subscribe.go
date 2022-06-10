// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-02
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type subscribe struct{}

var Subscribe = new(subscribe)

// Alter
// @ID           ManagerSubscribeAlter
// @Router       /manager/v1/subscribe/alter [POST]
// @Summary      M7004 修改订阅时间
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body  model.SubscribeAlter  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Alter(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.SubscribeAlter](c)
    return ctx.SendResponse(service.NewSubscribeWithModifier(ctx.Modifier).AlterDays(req))
}

// Pause
// @ID           ManagerRiderPause
// @Router       /manager/v1/subscribe/pause [POST]
// @Summary      M7006 暂停计费
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "订阅ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Pause(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDPostReq](c)
    service.NewRiderMgrWithModifier(ctx.Modifier).PauseSubscribe(req.ID)
    return ctx.SendResponse()
}

// Continue
// @ID           ManagerRiderContinue
// @Router       /manager/v1/subscribe/continue [POST]
// @Summary      M7007 继续计费
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "订阅ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Continue(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDPostReq](c)
    service.NewRiderMgrWithModifier(ctx.Modifier).ContinueSubscribe(req.ID)
    return ctx.SendResponse()
}

// Halt
// @ID           ManagerSubscribeHalt
// @Router       /manager/v1/subscribe/halt [POST]
// @Summary      M7008 强制退租
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        id  path  uint64  true  "订阅ID"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Halt(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.IDPostReq](c)
    service.NewRiderMgrWithModifier(ctx.Modifier).HaltSubscribe(req.ID)
    return ctx.SendResponse()
}