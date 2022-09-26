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

type couponTemplateService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
    orm          *ent.CouponTemplateClient
}

func NewCouponTemplate() *couponTemplateService {
    return &couponTemplateService{
        ctx: context.Background(),
        orm: ent.Database.CouponTemplate,
    }
}

func NewCouponTemplateWithRider(r *ent.Rider) *couponTemplateService {
    s := NewCouponTemplate()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewCouponTemplateWithModifier(m *model.Modifier) *couponTemplateService {
    s := NewCouponTemplate()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewCouponTemplateWithEmployee(e *ent.Employee) *couponTemplateService {
    s := NewCouponTemplate()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    s.employeeInfo = &model.Employee{
        ID:    s.employee.ID,
        Name:  s.employee.Name,
        Phone: s.employee.Phone,
    }
    return s
}

func (s *couponTemplateService) Create(req *model.CouponTemplate) {

}
