// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package aapi

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/service"
)

type statistics struct{}

var Statistics = new(statistics)

// Overview
//	@ID			AgentStatisticsOverview
//	@Router		/agent/v1/statistics/overview [GET]
//	@Summary	A9001 统计概览
//	@Tags		[A]代理接口
//	@Accept		json
//	@Produce	json
//	@Param		X-Agent-Token	header		string								true	"代理校验token"
//	@Success	200				{object}	model.AgentStatisticsOverviewRes	"请求成功"
func (*statistics) Overview(c echo.Context) (err error) {
	ctx := app.ContextX[app.AgentContext](c)
	return ctx.SendResponse(service.NewAgentStatistics().Overview(ctx))
}
