// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-30
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/internal/ent"
)

type PromotionContext struct {
	*BaseContext
	Member *ent.PromotionMember
}

// NewPromotionContext 新建代理上下文
func NewPromotionContext(c echo.Context, mem *ent.PromotionMember) *PromotionContext {
	ctx := &PromotionContext{
		BaseContext: Context(c),
		Member:      mem,
	}
	return ctx
}

// PromotionContextAndBinding 代理上下文绑定数据
func PromotionContextAndBinding[T any](c echo.Context) (*PromotionContext, *T) {
	return ContextBindingX[PromotionContext, T](c)
}
