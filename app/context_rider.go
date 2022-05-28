// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/labstack/echo/v4"
)

// RiderContext 骑手上下文
type RiderContext struct {
    *BaseContext

    Rider *ent.Rider
    Token string
}

// NewRiderContext 创建骑手上下文
func NewRiderContext(c echo.Context, rider *ent.Rider, token string) *RiderContext {
    return &RiderContext{
        BaseContext: Context(c),
        Rider:       rider,
        Token:       token,
    }
}

// RiderContextAndBinding 骑手端上下文绑定数据
func RiderContextAndBinding[T any](c echo.Context) (*RiderContext, *T) {
    return ContextBindingX[RiderContext, T](c)
}
