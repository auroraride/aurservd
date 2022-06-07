// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-06
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/enterprisestatement"
)

type statementService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *model.Employee
    orm      *ent.EnterpriseStatementClient
}

func NewStatement() *statementService {
    return &statementService{
        ctx: context.Background(),
        orm: ar.Ent.EnterpriseStatement,
    }
}

func NewStatementWithRider(r *ent.Rider) *statementService {
    s := NewStatement()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewStatementWithModifier(m *model.Modifier) *statementService {
    s := NewStatement()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewStatementWithEmployee(e *model.Employee) *statementService {
    s := NewStatement()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// Current 获取企业当前账单, 若无则新增
func (s *statementService) Current(enterpriseID uint64) *ent.EnterpriseStatement {
    res, _ := s.orm.QueryNotDeleted().Where(
        enterprisestatement.EnterpriseID(enterpriseID),
        enterprisestatement.SettledAtIsNil(),
    ).First(s.ctx)
    if res == nil {
        res, _ = s.orm.Create().SetEnterpriseID(enterpriseID).Save(s.ctx)
    }
    return res
}

func (s *statementService) GetBill(req *model.StatementBillReq) {

}
