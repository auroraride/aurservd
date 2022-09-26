// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-26
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
)

type couponService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
}

func NewCoupon() *couponService {
    return &couponService{
        ctx: context.Background(),
    }
}

func NewCouponWithRider(r *ent.Rider) *couponService {
    s := NewCoupon()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewCouponWithModifier(m *model.Modifier) *couponService {
    s := NewCoupon()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewCouponWithEmployee(e *ent.Employee) *couponService {
    s := NewCoupon()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    s.employeeInfo = &model.Employee{
        ID:    s.employee.ID,
        Name:  s.employee.Name,
        Phone: s.employee.Phone,
    }
    return s
}
