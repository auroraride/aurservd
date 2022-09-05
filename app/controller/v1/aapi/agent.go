// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-31
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
    "github.com/auroraride/aurservd/app"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/service"
    "github.com/labstack/echo/v4"
)

type agent struct{}

var Agent = new(agent)

// Signin
// @ID           AgentSignin
// @Router       /agent/v1/signin [POST]
// @Summary      A1001 登录
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理账号校验token"
// @Success      200  {object}  model.AgentSigninRes  "请求成功"
func (*agent) Signin(c echo.Context) (err error) {
    ctx, req := app.ContextBinding[model.AgentSigninReq](c)
    return ctx.SendResponse(service.NewAgent().Signin(req))
}

// Profile
// @ID           AgentAgentMeta
// @Router       /agent/v1/profile [GET]
// @Summary      A1002 代理资料
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Success      200  {object}  model.AgentProfile  "请求成功"
func (*agent) Profile(c echo.Context) (err error) {
    ctx := app.ContextX[app.AgentContext](c)
    return ctx.SendResponse(service.NewAgentWithAgent(ctx.Agent, ctx.Enterprise).Profile(ctx.Agent, ctx.Enterprise))
}
