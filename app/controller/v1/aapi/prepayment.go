// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-04
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type prepayment struct{}

var Prepayment = new(prepayment)

// Overview
// @ID           AgentPrepaymentOverview
// @Router       /agent/v1/prepayment/overview [GET]
// @Summary      A3001 账户详情
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Success      200  {object}  model.PrepaymentOverview  "请求成功"
func (*prepayment) Overview(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewPrepayment(ctx.Agent, ctx.Enterprise).Overview(ctx.Enterprise))
}

// List
// @ID           AgentPrepaymentList
// @Router       /agent/v1/prepayment [GET]
// @Summary      A3002 充值记录
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        query  query  model.PrepaymentListReq  true  "请求参数"
// @Success      200  {object}  model.PaginationRes{items=[]model.PrepaymentListRes}  "请求成功"
func (*prepayment) List(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.PrepaymentListReq](c)
	return ctx.SendResponse(service.NewPrepayment(ctx.Agent, ctx.Enterprise).List(ctx.Agent.EnterpriseID, req))
}

// WechatMiniprogramPay
// @ID           AgentPrepaymentWechatMiniprogramPay
// @Router       /agent/v1/prepayment/pay/wxmp [POST]
// @Summary      A3003 小程序储值
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body  model.AgentPrepayReq  true  "请求参数"
// @Success      200  {object}  model.AgentPrepayRes  "请求成功"
func (*prepayment) WechatMiniprogramPay(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentPrepayReq](c)
	return ctx.SendResponse(service.NewPrepayment().WechatMiniprogramPay(ctx.Agent, req))
}
