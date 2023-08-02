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
	service.NewEnterpriseRider().Create(req)
	return ctx.SendResponse()
}

// Alter
// @ID           AgentRiderAlter
// @Router       /agent/v1/rider/alter [POST]
// @Summary      A2003 增加/减少骑手时长
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body     model.SubscribeAlter  true  "请求详情"
// @Success      200  {object}  model.RiderItemSubscribe  "请求成功"
func (*rider) Alter(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.SubscribeAlter](c)
	service.NewSubscribeWithAgent(ctx.Agent, ctx.Enterprise).AlterDays(&model.SubscribeAlterReq{
		SubscribeAlter: *req,
		EnterpriseID:   ctx.Agent.EnterpriseID,
		AgentID:        ctx.Agent.ID,
	})
	return ctx.SendResponse()
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

// RiderInfo
// @ID           AgentRiderInfo
// @Router       /agent/v1/rider/info [GET]
// @Summary      A2005 通过二维码换取骑手信息
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        qrcode  query   string  true  "二维码"
// @Success      200  {object}  model.RiderSampleInfo  "请求成功"
func (*rider) RiderInfo(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.QRQueryReq](c)
	id := service.NewRider().ParseQrcode(req.Qrcode)
	return ctx.SendResponse(service.NewRider().GetRiderNameById(id))
}

// Invite 邀请骑手
// @ID           AgentRiderInvite
// @Router       /agent/v1/rider/invite [POST]
// @Summary      A2006 邀请骑手
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body   model.EnterpriseRiderInviteReq  true "邀请骑手"
// @Success      200  {object}  string  "请求成功"
func (*rider) Invite(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.EnterpriseRiderInviteReq](c)
	return ctx.SendResponse(service.NewminiProgram().Invite(ctx.Enterprise, req))
}

// Reactive 重新激活骑手
// @ID           AgentRiderReactive
// @Router       /agent/v1/rider/reactive [POST]
// @Summary      A2007 重新激活骑手
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body   model.ReactiveSubscribeReq  true "重新激活骑手"
// @Success      200  {object}  string  "请求成功"
func (*rider) Reactive(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.ReactiveSubscribeReq](c)
	service.NewSubscribe().ReactiveSubscribe(ctx, req)
	return ctx.SendResponse()
}

// Delete
// @ID           AgentRiderDelete
// @Router       /agent/v1/rider/{id} [DELETE]
// @Summary      A2008 删除骑手
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        id  path  uint64  true  "骑手ID"
// @Success      200  {object}  string  "请求成功"
func (*rider) Delete(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.IDParamReq](c)
	service.NewRiderAgentWithAgent(ctx.Agent, ctx.Enterprise).Delete(req, ctx.Agent.EnterpriseID)
	return ctx.SendResponse()
}
