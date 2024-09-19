// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-08
// Based on aurservd by liasica, magicrolan@qq.com.

package app

import (
	"github.com/labstack/echo/v4"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
)

// EmployeeContext 店员上下文
type EmployeeContext struct {
	*BaseContext

	Employee *ent.Employee
	Operator *model.OperatorInfo
}

// NewEmployeeContext 新建店员上下文
func NewEmployeeContext(c echo.Context, emr *ent.Employee) *EmployeeContext {
	ctx := &EmployeeContext{
		BaseContext: Context(c),
		Employee:    emr,
	}
	if emr != nil {
		ctx.Operator = &model.OperatorInfo{
			Type:  model.OperatorTypeEmployee,
			ID:    emr.ID,
			Phone: emr.Phone,
			Name:  emr.Name,
		}
	}
	return ctx
}

// EmployeeContextAndBinding 店员上下文绑定数据
func EmployeeContextAndBinding[T any](c echo.Context) (*EmployeeContext, *T) {
	return ContextBindingX[EmployeeContext, T](c)
}
