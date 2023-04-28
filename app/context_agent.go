// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-30
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/labstack/echo/v4"
)

type AgentContext struct {
	*BaseContext

	Enterprise *ent.Enterprise
	Agent      *ent.Agent
}

// NewAgentContext 新建代理上下文
func NewAgentContext(c echo.Context, ag *ent.Agent, en *ent.Enterprise) *AgentContext {
	ctx := &AgentContext{
		BaseContext: Context(c),
		Agent:       ag,
		Enterprise:  en,
	}
	return ctx
}

// AgentContextAndBinding 代理上下文绑定数据
func AgentContextAndBinding[T any](c echo.Context) (*AgentContext, *T) {
	return ContextBindingX[AgentContext, T](c)
}
