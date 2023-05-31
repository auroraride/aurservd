// Copyright (C) liasica. 2021-present.
//
// Created at 2021/12/12
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/pkg/snag"
)

type BaseContext struct {
	echo.Context

	Device *model.Device
}

// Context 获取上下文
func Context(c echo.Context) *BaseContext {
	switch v := c.(type) {
	case *ManagerContext:
		return v.BaseContext
	case *RiderContext:
		return v.BaseContext
	case *BaseContext:
		return v
	case *EmployeeContext:
		return v.BaseContext
	case *AgentContext:
		return v.BaseContext
	default:
		return nil
	}
}

// NewContext 创建上下文
func NewContext(c echo.Context) *BaseContext {
	return &BaseContext{
		Context: c,
	}
}

// BindValidate 绑定并校验数据
func (c *BaseContext) BindValidate(ptr any) {
	err := c.Bind(ptr)
	if err != nil {
		snag.Panic(err)
	}
	err = c.Validate(ptr)
	if err != nil {
		snag.Panic(err)
	}
}

// ContextBinding 绑定上下文并校验数据之后返回
func ContextBinding[T any](c echo.Context) (*BaseContext, *T) {
	ctx := Context(c)
	req := new(T)
	ctx.BindValidate(req)
	return ctx, req
}

type ContextWrapper interface {
	ManagerContext | RiderContext | EmployeeContext | AgentContext
}

func ContextX[T ContextWrapper](c any) *T {
	ctx, _ := c.(*T)
	return ctx
}

func ContextBindingX[K ContextWrapper, T any](c echo.Context) (*K, *T) {
	ctx := ContextX[K](c)
	req := new(T)
	Context(c).BindValidate(req)
	return ctx, req
}
