// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-13
// Based on aurservd by liasica, magicrolan@qq.com.

package eapi

type attendance struct{}

var Attendance = new(attendance)

// Precheck
// @ID		EmployeeAttendancePrecheck
// @Router	/employee/v1/attendance/precheck [POST]
// @Summary	E1003 打卡预检
// @Tags	[E]店员接口
// @Accept	json
// @Produce	json
// @Param	X-Employee-Token	header		string						true	"店员校验token"
// @Param	body				body		model.AttendancePrecheck	true	"预检请求"
// @Success	200					{object}	[]model.InventoryItem		"需盘点物资清单"
// func (*attendance) Precheck(c echo.Context) (err error) {
// 	ctx, req := app.EmployeeContextAndBinding[model.AttendancePrecheck](c)
// 	return ctx.SendResponse(service.NewAttendanceWithEmployee(ctx.Employee).Precheck(req))
// }

// Create
// @ID		EmployeeAttendanceCreate
// @Router	/employee/v1/attendance [POST]
// @Summary	E1004 考勤打卡
// @Tags	[E]店员接口
// @Accept	json
// @Produce	json
// @Param	X-Employee-Token	header		string						true	"店员校验token"
// @Param	body				body		model.AttendanceCreateReq	true	"打卡信息"
// @Success	200					{object}	model.StatusResponse		"请求成功"
// func (*attendance) Create(c echo.Context) (err error) {
// 	ctx, req := app.EmployeeContextAndBinding[model.AttendanceCreateReq](c)
// 	service.NewAttendanceWithEmployee(ctx.Employee).Create(req)
// 	return ctx.SendResponse()
// }
