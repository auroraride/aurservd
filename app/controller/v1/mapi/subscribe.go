// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-02
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
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
// @Param        body  body     model.SubscribeAlter  true  "desc"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Alter(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.SubscribeAlter](c)
	return ctx.SendResponse(service.NewSubscribeWithModifier(ctx.Modifier).AlterDays(&model.SubscribeAlterReq{
		SubscribeAlter: *req,
	}))
}

// Pause
// @ID           ManagerRiderPause
// @Router       /manager/v1/subscribe/pause [POST]
// @Summary      M7006 暂停计费
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.BusinessSubscribeReq  true  "订阅信息"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Pause(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BusinessSubscribeReq](c)
	service.NewBusinessRider(nil).SetModifier(ctx.Modifier).SetCabinetID(req.CabinetID).SetStoreID(req.StoreID).Pause(req.ID)
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
// @Param        body  body     model.BusinessSubscribeReq  true  "订阅信息"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Continue(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BusinessSubscribeReq](c)
	service.NewBusinessRider(nil).SetModifier(ctx.Modifier).SetCabinetID(req.CabinetID).SetStoreID(req.StoreID).Continue(req.ID)
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
// @Param        body  body     model.BusinessSubscribeReq  true  "订阅信息"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Halt(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.BusinessSubscribeReq](c)
	service.NewBusinessRider(nil).
		SetModifier(ctx.Modifier).
		SetCabinetID(req.CabinetID).
		SetStoreID(req.StoreID).
		UnSubscribe(
			&model.BusinessSubscribeReq{
				ID: req.ID,
			},
		)
	return ctx.SendResponse()
}

// Active
// @ID           ManagerSubscribeActive
// @Router       /manager/v1/subscribe/active [POST]
// @Summary      M7009 激活订阅
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.ManagerSubscribeActive  true  "订阅信息"
// @Success      200  {object}  model.AllocateCreateRes  "请求成功"
func (*subscribe) Active(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ManagerSubscribeActive](c)
	return ctx.SendResponse(service.NewManagerSubscribe(ctx.Modifier).Active(req))
}

// Suspend
// @ID           ManagerSubscribeSuspend
// @Router       /manager/v1/subscribe/suspend [POST]
// @Summary      M7014 暂停扣费
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.SuspendReq  true  "请求字段"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) Suspend(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.SuspendReq](c)
	service.NewSuspendWithModifier(ctx.Modifier).Suspend(req)
	return ctx.SendResponse()
}

// UnSuspend
// @ID           ManagerSubscribeUnSuspend
// @Router       /manager/v1/subscribe/unsuspend [POST]
// @Summary      M7015 继续扣费
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body  body     model.SuspendReq  true  "请求字段"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) UnSuspend(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.SuspendReq](c)
	service.NewSuspendWithModifier(ctx.Modifier).UnSuspend(req)
	return ctx.SendResponse()
}

// EbikeChange
// @ID           ManagerSubscribeEbikeChange
// @Router       /manager/v1/subscribe/ebike/change [POST]
// @Summary      M7020 修改订阅车辆
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token header string true "管理员校验token"
// @Param        body body    model.ManagerSubscribeChangeEbike true "换车参数"
// @Success      200 {object} model.StatusResponse "请求成功"
func (*subscribe) EbikeChange(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ManagerSubscribeChangeEbike](c)
	service.NewManagerSubscribe(ctx.Modifier).ChangeEbike(req)
	return ctx.SendResponse()
}

// Reminder
// @ID           ManagerSubscribeReminder
// @Router       /manager/v1/subscribe/reminder [GET]
// @Summary      M7021 催费记录
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token header string true "管理员校验token"
// @Param        query query  model.ReminderListReq false "筛选选项"
// @Success      200 {object} model.PaginationRes{items=[]model.ReminderListRes} "请求成功"
func (*subscribe) Reminder(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ReminderListReq](c)
	return ctx.SendResponse(service.NewReminder(ctx.Modifier).List(req))
}

// EbikeUnbind
// @ID           ManagerSubscribeEbikeUnbind
// @Router       /manager/v1/subscribe/ebike/unbind [POST]
// @Summary      M7025 解绑骑手电车
// @Tags         [M]管理接口
// @Accept       json
// @Produce      json
// @Param        X-Manager-Token  header  string  true  "管理员校验token"
// @Param        body body model.ManagerSubscribeUnbindEbike true "请求参数"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) EbikeUnbind(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.ManagerSubscribeUnbindEbike](c)
	service.NewManagerSubscribe(ctx.Modifier).UnbindEbike(req)
	return ctx.SendResponse()
}
