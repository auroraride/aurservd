// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-23
// Based on aurservd by liasica, magicrolan@qq.com.

package mapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type assistance struct{}

var Assistance = new(assistance)

// List
// @ID		ManagerAssistanceList
// @Router	/manager/v1/assistance [GET]
// @Summary	救援列表
// @Tags	救援
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string													true	"管理员校验token"
// @Param	query			query		model.AssistanceListReq									false	"筛选项"
// @Success	200				{object}	model.PaginationRes{items=[]model.AssistanceListRes}	"请求成功"
func (*assistance) List(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssistanceListReq](c)
	return ctx.SendResponse(service.NewAssistance().List(req))
}

// Detail
// @ID		ManagerAssistanceDetail
// @Router	/manager/v1/assistance/{id} [GET]
// @Summary	救援详情
// @Tags	救援
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	id				path		uint64					true	"救援ID"
// @Success	200				{object}	model.AssistanceDetail	"救援详情"
func (*assistance) Detail(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewAssistance().Detail(req.ID))
}

// Nearby
// @ID		ManagerAssistanceNearby
// @Router	/manager/v1/assistance/nearby [GET]
// @Summary	附近门店
// @Tags	救援
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	id				query		uint64						true	"救援订单ID"
// @Success	200				{object}	[]model.AssistanceNearbyRes	"请求成功"
func (*assistance) Nearby(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.IDQueryReq](c)
	return ctx.SendResponse(service.NewAssistanceWithModifier(ctx.Modifier).Nearby(req))
}

// Allocate
// @ID		ManagerAssistanceAllocate
// @Router	/manager/v1/assistance/allocate [POST]
// @Summary	分配救援任务
// @Tags	救援
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.AssistanceAllocateReq	true	"分配参数"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*assistance) Allocate(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssistanceAllocateReq](c)
	service.NewAssistanceWithModifier(ctx.Modifier).Allocate(req)
	return ctx.SendResponse()
}

// Free
// @ID		ManagerAssistanceFree
// @Router	/manager/v1/assistance/free [POST]
// @Summary	救援免费
// @Tags	救援
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string					true	"管理员校验token"
// @Param	body			body		model.AssistanceFreeReq	true	"请求参数"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*assistance) Free(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssistanceFreeReq](c)
	service.NewAssistanceWithModifier(ctx.Modifier).Free(req)
	return ctx.SendResponse()
}

// Refuse
// @ID		ManagerAssistanceRefuse
// @Router	/manager/v1/assistance/refuse [POST]
// @Summary	拒绝救援
// @Tags	救援
// @Accept	json
// @Produce	json
// @Param	X-Manager-Token	header		string						true	"管理员校验token"
// @Param	body			body		model.AssistanceRefuseReq	true	"拒绝请求"
// @Success	200				{object}	model.StatusResponse		"请求成功"
func (*assistance) Refuse(c echo.Context) (err error) {
	ctx, req := app.ManagerContextAndBinding[model.AssistanceRefuseReq](c)
	service.NewAssistanceWithModifier(ctx.Modifier).Refuse(req)
	return ctx.SendResponse()
}
