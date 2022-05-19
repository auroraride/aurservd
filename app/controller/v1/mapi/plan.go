// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type plan struct{}

var Plan = new(plan)

// Create
// @ID           PlanCreate
// @Router       /manager/v1/plan [POST]
// @Summary      M60001 创建骑士卡
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*plan) Create(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PlanCreateReq](c)
    service.NewPlan().CreatePlan(ctx.Modifier, req)
    return ctx.SendResponse()
}
