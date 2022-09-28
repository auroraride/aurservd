// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-28
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type couponAssemblyService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
}

func NewCouponAssembly() *couponAssemblyService {
    return &couponAssemblyService{
        ctx: context.Background(),
    }
}

func NewCouponAssemblyWithRider(r *ent.Rider) *couponAssemblyService {
    s := NewCouponAssembly()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewCouponAssemblyWithModifier(m *model.Modifier) *couponAssemblyService {
    s := NewCouponAssembly()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewCouponAssemblyWithEmployee(e *ent.Employee) *couponAssemblyService {
    s := NewCouponAssembly()
    if e != nil {
        s.employee = e
        s.employeeInfo = &model.Employee{
            ID:    e.ID,
            Name:  e.Name,
            Phone: e.Phone,
        }
        s.ctx = context.WithValue(s.ctx, "employee", s.employeeInfo)
    }
    return s
}
