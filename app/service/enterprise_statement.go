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
    "github.com/auroraride/aurservd/internal/ent/enterprisestatement"
)

type enterpriseStatementService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *model.Employee
    orm      *ent.EnterpriseStatementClient
}

func NewEnterpriseStatement() *enterpriseStatementService {
    return &enterpriseStatementService{
        ctx: context.Background(),
        orm: ar.Ent.EnterpriseStatement,
    }
}

func NewEnterpriseStatementWithRider(r *ent.Rider) *enterpriseStatementService {
    s := NewEnterpriseStatement()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewEnterpriseStatementWithModifier(m *model.Modifier) *enterpriseStatementService {
    s := NewEnterpriseStatement()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewEnterpriseStatementWithEmployee(e *model.Employee) *enterpriseStatementService {
    s := NewEnterpriseStatement()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// Current 获取企业当前账单, 若无则新增
func (s *enterpriseStatementService) Current(enterpriseID uint64) *ent.EnterpriseStatement {
    res, _ := s.orm.QueryNotDeleted().Where(
        enterprisestatement.EnterpriseID(enterpriseID),
        enterprisestatement.SettledAtIsNil(),
    ).First(s.ctx)
    if res == nil {
        res, _ = s.orm.Create().SetEnterpriseID(enterpriseID).Save(s.ctx)
    }
    return res
}
