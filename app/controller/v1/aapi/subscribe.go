// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-10
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type subscribe struct{}

var Subscribe = new(subscribe)

// Active
// @ID           AgentSubscribeActive
// @Router       /agent/v1/subscribe/active [POST]
// @Summary      A7001 激活订阅
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body     model.AgentSubscribeActiveReq  true  "请求详情"
// @Success      200  {object}  string  "请求成功"
func (*subscribe) Active(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentSubscribeActiveReq](c)
	service.NewEnterprise().Active(req, ctx.Enterprise.ID)
	return ctx.SendResponse()
}

// AlterList
// @ID           AgentSubscribeAlterList
// @Router       /agent/v1/subscribe/alter [GET]
// @Summary      A7002 申请加时列表
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        query  query   model.SubscribeAlterApplyReq  true  "查询条件"
// @Success      200  {object}  model.PaginationRes{items=[]model.SubscribeAlterApplyListRsp}  "请求成功"
func (*subscribe) AlterList(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.SubscribeAlterApplyReq](c)
	return ctx.SendResponse(service.NewAgentSubscribe(ctx.Agent, ctx.Enterprise).AlterList(ctx.Enterprise.ID, req))
}

// AlterReivew
// @ID           AgentSubscribeAlterReivew
// @Router       /agent/v1/subscribe/review [POST]
// @Summary      A7003 审核加时
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body   model.SubscribeAlterReviewReq  true  "审核请求"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*subscribe) AlterReivew(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.SubscribeAlterReviewReq](c)
	service.NewAgentSubscribe(ctx.Agent, ctx.Enterprise).AlterReview(&model.SubscribeAlterReviewReq{
		Ids:    req.Ids,
		Status: req.Status,
	})
	return ctx.SendResponse()
}
