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

// SiteList 站点列表
// @ID           AgentAgentname
// @Router       /agent/v1/site/list [GET]
// @Summary      A1004 站点列表
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Success      200  {object}  []model.EnterpriseStation  "请求成功"
func (*agent) SiteList(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewEnterpriseStation().List(&model.EnterpriseStationListReq{
		EnterpriseID: ctx.Enterprise.ID,
	}))
}

// CityList 城市列表
// @ID           AgentCityList
// @Router       /agent/v1/city/list [GET]
// @Summary      A1005 城市列表
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Success      200  {object}  []model.CityListReq  "请求成功"
func (*agent) CityList(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewCity().List(&model.CityListReq{
		Status: model.CityStatusOpen,
	}))
}

// BatteryList 电池列表
// @ID           AgentBatteryList
// @Router       /agent/v1/battery/list [GET]
// @Summary      A1006 电池列表
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
