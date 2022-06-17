// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-17
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type exceptionService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
}

func NewException() *exceptionService {
    return &exceptionService{
        ctx: context.Background(),
    }
}

func NewExceptionWithRider(r *ent.Rider) *exceptionService {
    s := NewException()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewExceptionWithModifier(m *model.Modifier) *exceptionService {
    s := NewException()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewExceptionWithEmployee(e *ent.Employee) *exceptionService {
    s := NewException()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func (s *exceptionService) Setting() model.ExceptionEmployeeSetting {
    set := NewSetting().GetSetting(model.SettingException)
    return model.ExceptionEmployeeSetting{
        Items: NewInventory().ListInventory(model.InventoryListReq{
            Count:    true,
            Transfer: true,
            Purchase: true,
        }),
        Reasons: set.([]interface{}),
    }
}
