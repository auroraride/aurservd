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

// Detail
// @ID           AgentCabinetDetail
// @Router       /agent/v1/cabinet/{serial} [GET]
// @Summary      A5002 电柜详情
// @Tags         [A]代理接口
// @Accept       json
// @Produce      json
// @Param        X-Agent-Token  header  string  true  "代理校验token"
// @Param        serial  path  string  true  "电柜编号"
// @Success      200  {object}  model.StatusResponse  "请求成功"
func (*cabinet) Detail(c echo.Context) (err error) {
	ctx, req := app.AgentContextAndBinding[model.AgentCabinetDetailReq](c)
	return ctx.SendResponse(service.NewAgentCabinet().Detail(req.Serial, ctx.Agent, ctx.Stations))
}