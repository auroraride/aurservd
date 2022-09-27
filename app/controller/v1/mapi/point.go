// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-27
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type point struct{}

var Point = new(point)

// Modify
// @ID           ManagerPointModify
// @Router       /manager/v1/point/modify [POST]
// @Summary      M7016 修改骑手积分
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.PointModifyReq  true  "请求参数"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*point) Modify(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PointModifyReq](c)
    service.NewPointWithModifier(ctx.Modifier).Modify(req)
    return ctx.SendResponse()
}

// Log
// @ID           ManagerPointLog
// @Router       /manager/v1/point/log [GET]
// @Summary      M7017 积分变动日志
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Success      200  {object}  model.PaginationRes{items=[]model.PointLogListRes}  "请求成功"
func (*point) Log(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.PointLogListReq](c)
    return ctx.SendResponse(service.NewPointWithModifier(ctx.Modifier).LogList(req))
}