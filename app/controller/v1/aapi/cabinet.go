// Copyright (C) liasica. 2023-present.
//
// Created at 2023-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/auroraride/adapter/defs/cabdef"
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type cabinet struct{}

var Cabinet = new(cabinet)

// List
// @ID		AgentCabinetList
// @Router	/agent/v1/cabinet [GET]
// @Summary	A5001 电柜列表
// @Tags	[A]代理接口
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string											true	"代理校验token"
// @Param	query			query		model.AgentCabinetListReq						false	"请求参数"
// @Success	200				{object}	model.PaginationRes{items=[]model.AgentCabinet}	"请求成功"
func (*cabinet) List(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentCabinetListReq](c)
	return ctx.SendResponse(service.NewAgentCabinet().List(ctx, req))
}

// Detail
// @ID		AgentCabinetDetail
// @Router	/agent/v1/cabinet/detail/{serial} [GET]
// @Summary	A5002 电柜详情
// @Tags	[A]代理接口
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string				true	"代理校验token"
// @Param	serial			path		string				true	"电柜编号"
// @Param	query			query		model.LngLat		false	"请求参数"
// @Success	200				{object}	model.AgentCabinet	"请求成功"
func (*cabinet) Detail(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentCabinetDetailReq](c)
	return ctx.SendResponse(service.NewAgentCabinet().Detail(ctx, req))
}

// Section
// @ID		AgentCabinetSection
// @Router	/agent/v1/cabinet/section [GET]
// @Summary	A5003 电柜选择
// @Tags	[A]代理接口
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string							true	"代理校验token"
// @Success	200				{object}	[]model.CascaderOptionLevel2	"请求成功"
func (*cabinet) Section(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewSelection().Cabinet(&model.CabinetSelectionReq{EnterpriseID: ctx.Agent.EnterpriseID}))
}

// Maintain
// @ID		AgentCabinetMaintain
// @Router	/agent/v1/cabinet/maintain [POST]
// @Summary	A5004 电柜维护
// @Tags	[A]代理接口
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string					true	"代理校验token"
// @Param	query			query		model.AgentMaintainReq	false	"请求参数"
// @Success	200				{object}	model.StatusResponse	"请求成功"
func (*cabinet) Maintain(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentMaintainReq](c)
	service.NewAgentCabinet().Maintain(ctx, req)
	return ctx.SendResponse()
}

// BinOpen
// @ID		AgentCabinetBinOpen
// @Router	/agent/v1/cabinet/binopen [POST]
// @Summary	A5005 电柜开仓
// @Tags	[A]代理接口
// @Accept	json
// @Produce	json
// @Param	X-Agent-Token	header		string						true	"代理校验token"
// @Param	query			query		model.AgentBinOperateReq	false	"请求参数"
// @Success	200				{object}	[]cabdef.BinOperateResult	"请求成功"
func (*cabinet) BinOpen(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentBinOperateReq](c)
	return ctx.SendResponse(service.NewAgentCabinet().BinOpen(ctx, req, cabdef.OperateDoorOpen))
}
