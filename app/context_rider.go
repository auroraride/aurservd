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

// RiderContext 骑手上下文
type RiderContext struct {
	*BaseContext

	Rider    *ent.Rider
	Token    string
	Operator *model.OperatorInfo
}

// NewRiderContext 创建骑手上下文
func NewRiderContext(c echo.Context, rider *ent.Rider, token string) *RiderContext {
	c.Set("rider", rider)
	ctx := &RiderContext{
		BaseContext: Context(c),
		Rider:       rider,
		Token:       token,
	}
	if rider != nil {
		ctx.Operator = &model.OperatorInfo{
			Type:  model.OperatorTypeRider,
			ID:    rider.ID,
			Phone: rider.Phone,
			Name:  rider.Name,
		}
	}
	return ctx
}

// RiderContextAndBinding 骑手端上下文绑定数据
func RiderContextAndBinding[T any](c echo.Context) (*RiderContext, *T) {
	return ContextBindingX[RiderContext, T](c)
}
