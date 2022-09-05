// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/enterpriseprepayment"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
)

type prepaymentService struct {
    ctx        context.Context
    modifier   *model.Modifier
    rider      *ent.Rider
    enterprise *ent.Enterprise
    agent      *ent.Agent
    orm        *ent.EnterprisePrepaymentClient
}

func NewPrepayment() *prepaymentService {
    return &prepaymentService{
        ctx: context.Background(),
        orm: ent.Database.EnterprisePrepayment,
    }
}

func NewPrepaymentWithAgent(ag *ent.Agent, en *ent.Enterprise) *prepaymentService {
    s := NewPrepayment()
    s.agent = ag
    s.enterprise = en
    return s
}

func NewPrepaymentWithModifier(m *model.Modifier) *prepaymentService {
    s := NewPrepayment()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *prepaymentService) Overview(en *ent.Enterprise) (res model.PrepaymentOverview) {
    res.Balance = en.Balance
    var result []struct {
        EnterpriseID uint64  `json:"enterprise_id"`
        Amount       float64 `json:"amount"`
        Times        int     `json:"times"`
    }
    _ = ent.Database.EnterprisePrepayment.
        QueryNotDeleted().
        Where(enterpriseprepayment.EnterpriseID(en.ID)).
        GroupBy(enterpriseprepayment.FieldEnterpriseID).
        Aggregate(
            ent.As(ent.Sum(enterpriseprepayment.FieldAmount), "amount"),
            ent.As(ent.Count(), "times"),
        ).
        Scan(s.ctx, &result)
    if len(result) == 0 {
        return
    }
    res.Times = result[0].Times
    res.Amount = result[0].Amount
    res.Cost = tools.NewDecimal().Sub(res.Amount, res.Balance)
    return
}

func (s *prepaymentService) List(req *model.PrepaymentListReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().
        Where(enterpriseprepayment.EnterpriseID(req.EnterpriseID)).
        Order(ent.Desc(enterpriseprepayment.FieldCreatedAt))
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.EnterprisePrepayment) model.PrepaymentListRes {
        res := model.PrepaymentListRes{
            Amount: item.Amount,
            Time:   item.CreatedAt.Format(carbon.DateTimeLayout),
            Remark: item.Remark,
        }
        res.Name = "平台管理员"
        return res
    })
}
