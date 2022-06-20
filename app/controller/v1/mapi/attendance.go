// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-20
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type attendance struct{}

var Attendance = new(attendance)

// List
// @ID           ManagerAttendanceList
// @Router       /manager/v1/employee/attendance [GET]
// @Summary      MA016 工作流
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        query  query   model.AttendanceListReq  true  "筛选请求"
// @Success      200  {object}  model.PaginationRes{items=[]model.AttendanceListRes}  "请求成功"
func (*attendance) List(c echo.Context) (err error) {
    ctx, req := app.ManagerContextAndBinding[model.AttendanceListReq](c)
    return ctx.SendResponse(service.NewAttendanceWithModifier(ctx.Modifier).List(req))
}