// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-01
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type rider struct{}

var Rider = new(rider)

// List
// @ID           AgentRiderList
// @Router       /agent/v1/rider [GET]
// @Summary      A2001 骑手列表
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        query  query   model.AgentRiderListReq  false  "查询条件"
// @Success      200  {object}  model.PaginationRes{items=[]model.AgentRider}  "请求成功"
func (*rider) List(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentRiderListReq](c)
	return ctx.SendResponse(service.NewRiderAgentWithAgent(ctx.Agent, ctx.Enterprise).List(ctx.Enterprise.ID, req))
}

// Create
// @ID           AgentRiderCreate
// @Router       /agent/v1/rider [POST]
// @Summary      A2002 添加骑手
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body     model.EnterpriseRiderCreateReq  true  "骑手信息"
// @Success      200  {object}  model.EnterpriseRider  "请求成功"
func (*rider) Create(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.EnterpriseRiderCreateReq](c)
	service.NewEnterpriseRider().CreateByAgent(req, ctx.Agent, ctx.Stations)
	return ctx.SendResponse(model.StatusResponse{Status: true})
}

// Alter
// @ID           AgentRiderAlter
// @Router       /agent/v1/rider/alter [POST]
// @Summary      A2003 延长天数
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body     model.SubscribeAlter  true  "请求详情"
// @Success      200  {object}  model.RiderItemSubscribe  "请求成功"
func (*rider) Alter(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.SubscribeAlter](c)
	return ctx.SendResponse(service.NewSubscribeWithAgent(ctx.Agent, ctx.Enterprise).AlterDays(req))
}

// Detail
// @ID           AgentRider
// @Router       /agent/v1/rider/{id} [GET]
// @Summary      A2004 骑手详情
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        id  path  uint64  true  "骑手ID"
// @Success      200  {object}  model.AgentRider  "请求成功"
func (*rider) Detail(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewRiderAgentWithAgent(ctx.Agent, ctx.Enterprise).Detail(req, ctx.Enterprise.ID))
}

// Active 激活骑手
// @ID           AgentRiderActive
// @Router       /agent/v1/rider/active [POST]
// @Summary      A2005 激活骑手
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body     model.RiderActiveBatteryReq  true  "请求详情"
// @Success      200  {object}  string  "请求成功"
func (*rider) Active(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.RiderActiveBatteryReq](c)
	service.NewEnterprise().Active(req, ctx.Agent)
	return ctx.SendResponse(model.StatusResponse{Status: true})
}

// SubscribeApplyList  申请加时列表
// @ID           AgentSubscribeApplyList
// @Router       /agent/v1/rider/subscribe/apply [GET]
// @Summary      A2006 申请加时列表
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        query  query   model.SubscribeAlterApplyReq  true  "查询条件"
// @Success      200  {object}  model.PaginationRes{items=[]model.SubscribeAlterApplyListRsp}  "请求成功"
func (*rider) SubscribeApplyList(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.SubscribeAlterApplyReq](c)
	return ctx.SendResponse(service.NewEnterpriseWithAgent(ctx.Agent, ctx.Enterprise).SubscribeApplyList(req, ctx.Enterprise))
}

// ReviewApply  审核加时
// @ID           AgentReviewApply
// @Router       /agent/v1/rider/subscribe/apply [POST]
// @Summary      A2007 审核加时
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body   model.SubscribeAlterReviewReq  true  "审核请求"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*rider) ReviewApply(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.SubscribeAlterReviewReq](c)
	service.NewEnterpriseWithAgent(ctx.Agent, ctx.Enterprise).SubscribeApplyReviewApply(req)
	return ctx.SendResponse(model.StatusResponse{Status: true})
}
