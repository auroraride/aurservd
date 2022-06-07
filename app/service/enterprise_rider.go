// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/pkg/snag"
)

type enterpriseRiderService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *model.Employee
}

func NewEnterpriseRider() *enterpriseRiderService {
    return &enterpriseRiderService{
        ctx: context.Background(),
    }
}

func NewEnterpriseRiderWithRider(r *ent.Rider) *enterpriseRiderService {
    s := NewEnterpriseRider()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewEnterpriseRiderWithModifier(m *model.Modifier) *enterpriseRiderService {
    s := NewEnterpriseRider()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewEnterpriseRiderWithEmployee(e *model.Employee) *enterpriseRiderService {
    s := NewEnterpriseRider()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func (s *enterpriseRiderService) Create(req *model.EnterpriseRiderCreateReq) {
    // 查询是否存在
    if ar.Ent.Rider.QueryNotDeleted().Where(rider.Phone(req.Phone)).ExistX(s.ctx) {
        snag.Panic("此手机号已存在")
    }

    // tx, _ := ar.Ent.Tx(s.ctx)
    // tx.Rider
}
