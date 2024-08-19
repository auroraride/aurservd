// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-12, by aurb

package app

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/internal/ent"
)

type WarestoreContext struct {
	*BaseContext
	AssetManager *ent.AssetManager
	Employee     *ent.Employee
}

// NewWarestoreContext 新建代理上下文
func NewWarestoreContext(c echo.Context, am *ent.AssetManager, ep *ent.Employee) *WarestoreContext {
	ctx := &WarestoreContext{
		BaseContext:  Context(c),
		AssetManager: am,
		Employee:     ep,
	}
	return ctx
}

// WarestoreContextAndBinding 代理上下文绑定数据
func WarestoreContextAndBinding[T any](c echo.Context) (*WarestoreContext, *T) {
	return ContextBindingX[WarestoreContext, T](c)
}
