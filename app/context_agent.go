// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-30
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
)

type AgentContext struct {
	*BaseContext

	Enterprise *ent.Enterprise
	Agent      *ent.Agent
	Stations   ent.EnterpriseStations // 站点列表
	Operator   *model.OperatorInfo
}

// NewAgentContext 新建代理上下文
func NewAgentContext(c echo.Context, ag *ent.Agent, en *ent.Enterprise, stations ent.EnterpriseStations) *AgentContext {
	ctx := &AgentContext{
		BaseContext: Context(c),
		Agent:       ag,
		Enterprise:  en,
		Stations:    stations,
	}
	if ag != nil {
		ctx.Operator = &model.OperatorInfo{
			Type:  model.OperatorTypeAgent,
			ID:    ag.ID,
			Phone: ag.Phone,
			Name:  ag.Name,
		}
	}
	return ctx
}

// AgentContextAndBinding 代理上下文绑定数据
func AgentContextAndBinding[T any](c echo.Context) (*AgentContext, *T) {
	return ContextBindingX[AgentContext, T](c)
}

// StationIDs 该代理站点ID列表
func (c *AgentContext) StationIDs() (ids []uint64) {
	for _, station := range c.Stations {
		ids = append(ids, station.ID)
	}
	return
}
