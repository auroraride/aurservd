// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-26
// Based on aurservd by liasica, magicrolan@qq.com.

package rapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type plan struct{}

var Plan = new(plan)

// List
// @ID           RiderPlanList
// @Router       /rider/v1/plan [GET]
// @Summary      R3002 获取骑士卡
// @Tags         [R]骑手接口
// @Accept       json
// @Produce      json
// @Param        X-Rider-Token  header  string  true  "骑手校验token"
// @Param        query  query  model.PlanListRiderReq  true  "desc"
// @Success      200  {object}  []model.RiderPlanItem  "请求成功"
func (*plan) List(c echo.Context) (err error) {
    ctx, req := app.RiderContextAndBinding[model.PlanListRiderReq](c)
    return ctx.SendResponse(
        service.NewPlanWithRider(ctx.Rider).RiderList(req),
    )
}
