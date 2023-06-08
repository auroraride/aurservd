// Copyright (C) liasica. 2023-present.
//
// Created at 2023-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
)

type cabinet struct{}

var Cabinet = new(cabinet)

// List 电柜列表
// @ID           AgentCabinetList
// @Router       /agent/v1/cabinet [GET]
// @Summary      A5003 电柜列表
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        query  query   model.CabinetQueryReq  false  "筛选选项"
// @Success      200  {object}  model.PaginationRes{items=[]model.AgentCabinetDetailRes}  "请求成功"
func (*cabinet) List(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.CabinetQueryReq](c)
	return ctx.SendResponse(service.NewCabinet().List(&model.CabinetQueryReq{
		PaginationReq: model.PaginationReq{},
		Serial:        req.Serial,
		Name:          req.Name,
		CityID:        req.CityID,
		Brand:         req.Brand,
		Status:        req.Status,
		Model:         req.Model,
		Online:        req.Online,
		Intelligent:   req.Intelligent,
		EnterpriseID:  &ctx.Enterprise.ID,
	}))
}

// Detail
// @ID           AgentCabinetDetail
// @Router       /agent/v1/cabinet/{serial} [GET]
// @Summary      A5002 电柜详情
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        serial  path  string  true  "电柜编号"
// @Success      200  {object}  model.AgentCabinetDetailRes  "请求成功"
func (*cabinet) Detail(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentCabinetDetailReq](c)
	return ctx.SendResponse(service.NewAgentCabinet().Detail(req.Serial, ctx.Agent, ctx.Stations))
}

// CabinetFilter
// @ID           ManagerCabinetFilter
// @Router       /agent/v1/cabinet/filter [GET]
// @Summary      A5004 筛选电柜
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Success      200  {object}  []model.CascaderOptionLevel2  "请求成功"
func (*cabinet) CabinetFilter(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewSelection().EnterpriseCabinet(ctx.Enterprise.ID))
}
