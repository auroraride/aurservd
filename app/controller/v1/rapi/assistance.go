// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-23
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type assistance struct{}

var Assistance = new(assistance)

// Breakdown
// @ID           RiderAssistanceBreakdown
// @Router       /rider/v1/assistance/breakdown [GET]
// @Summary      R5001 获取救援原因
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Success      200 {object}  []string  "请求成功"
func (*assistance) Breakdown(c echo.Context) (err error) {
    ctx := app.Context(c)
    return ctx.SendResponse(service.NewAssistance().Breakdown())
}

// Create
// @ID           RiderAssistanceCreate
// @Router       /rider/v1/assistance [POST]
// @Summary      R5002 发起救援
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        body  body     model.AssistanceCreateReq  true  "救援参数"
// @Success      200 {object}   model.AssistanceCreateRes  "请求成功"
func (*assistance) Create(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.AssistanceCreateReq](c)
    return ctx.SendResponse(service.NewAssistanceWithRider(ctx.Rider).Create(req))
}