// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
)

// ManagerContext 管理员上下文
type ManagerContext struct {
	*BaseContext

	Manager  *ent.Manager
	Modifier *model.Modifier
	Operator *model.OperatorInfo
}

// NewManagerContext 新建管理员上下文
func NewManagerContext(c echo.Context, mgr *ent.Manager, m *model.Modifier) *ManagerContext {
	ctx := &ManagerContext{
		BaseContext: Context(c),
		Manager:     mgr,
		Modifier:    m,
	}
	if mgr != nil {
		ctx.Operator = &model.OperatorInfo{
			Type:  model.OperatorTypeManager,
			ID:    mgr.ID,
			Phone: mgr.Phone,
			Name:  mgr.Name,
		}
	}
	return ctx
}

// GetManagerContext 获取管理端上下文
func GetManagerContext(c echo.Context) *ManagerContext {
	return c.(*ManagerContext)
}

// ManagerContextAndBinding 管理端上下文绑定数据
func ManagerContextAndBinding[T any](c echo.Context) (*ManagerContext, *T) {
	return ContextBindingX[ManagerContext, T](c)
}
