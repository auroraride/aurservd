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
	return ctx.SendResponse(service.NewEnterpriseRiderWithAgent(ctx.Agent, ctx.Enterprise).CreateByAgent(req))
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
// @Success      200  {object}  model.AgentRiderDetail  "请求成功"
func (*rider) Detail(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.IDParamReq](c)
	return ctx.SendResponse(service.NewRiderAgentWithAgent(ctx.Agent, ctx.Enterprise).Detail(req, ctx.Enterprise.ID))
}
