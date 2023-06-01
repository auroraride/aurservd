// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-31
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
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
// @Param        body  body  model.AgentSigninReq  true  "登录请求"
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
	return ctx.SendResponse(service.NewAgent(ctx.Agent, ctx.Enterprise).Profile(ctx.Agent, ctx.Enterprise))
}

// GetOpenid 获取opienid
// @ID           AgentGetOpenid
// @Router       /agent/v1/getopenid [GET]
// @Summary      A1003 获取openid
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param		code  query  string  true  "微信code"
// @Success      200  code  string  "请求成功"
func (*agent) GetOpenid(c echo.Context) (err error) {
	ctx, req := app.ContextBinding[model.OpenidReq](c)
	return ctx.SendResponse(service.NewminiProgram().GetAuth(req.Code))
}

// BatteryList 电池列表
// @ID           AgentBatteryList
// @Router       /agent/v1/battery/list [GET]
// @Summary      A5001 电池列表
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        query  query   model.BatterySearchReq  true  "筛选项"
// @Success      200  {object}  []model.Battery
func (*agent) BatteryList(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.BatterySearchReq](c)
	return ctx.SendResponse(service.NewSelection().BatterySerialSearch(req))
}

// Invite 邀请骑手
// @ID           AgentRiderInvite
// @Router       /agent/v1/rider/invite [GET]
// @Summary      A2006 邀请骑手
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body   model.EnterpriseRiderInviteReq  true  "邀请骑手"
// @Success      200  {object}  []byte  "请求成功"
func (*rider) Invite(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.EnterpriseRiderInviteReq](c)
	return ctx.SendResponse(service.NewminiProgram().Invite(ctx.Enterprise, req))
}

// Index 首页数据
// @ID           AgentIndex
// @Router       /agent/v1/index [GET]
// @Summary      A2007 首页数据
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Success      200  {object}  []byte  "请求成功"
func (*agent) Index(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewAgent(ctx.Agent, ctx.Enterprise).Index(ctx.Agent, ctx.Enterprise))
}

// Feedback 意见反馈
// @ID           AgentFeedback
// @Router       /agent/v1/feedback [POST]
// @Summary      A1005 意见反馈
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        body  body   model.FeedbackReq  true  "反馈内容"
// @Success      200  {object}  bool  "请求成功"
func (*agent) Feedback(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.FeedbackReq](c)
	return ctx.SendResponse(service.NewFeedback().Create(req, ctx.Enterprise))
}

// UploadImage 意见反馈上传图片
// @ID           AgentUploadImage
// @Router       /agent/v1/upload/image [POST]
// @Summary      A1006 意见反馈上传图片
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        images  formData   file  true  "图片文件"
// @Success      200  {object} []string  "请求成功"
func (*agent) UploadImage(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewFeedback().UploadImage(ctx))
}
