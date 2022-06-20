// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/enterprise"
    "github.com/auroraride/aurservd/internal/ent/enterprisestatement"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "math"
)

type enterpriseStatementService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
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

func NewEnterpriseStatementWithEmployee(e *ent.Employee) *enterpriseStatementService {
    s := NewEnterpriseStatement()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// Current 获取企业当前账单, 若无则新增
func (s *enterpriseStatementService) Current(e *ent.Enterprise) *ent.EnterpriseStatement {
    res, _ := s.orm.QueryNotDeleted().Where(
        enterprisestatement.EnterpriseID(e.ID),
        enterprisestatement.SettledAtIsNil(),
    ).First(s.ctx)

    if res == nil {
        log.Infof("%d 未找到账单, 创建新账单", e.ID)
        res, _ = s.orm.Create().
            SetEnterpriseID(e.ID).
            SetBalance(e.Balance).
            SetStart(carbon.Time2Carbon(e.CreatedAt).StartOfDay().Carbon2Time()).
            Save(s.ctx)
    }
    return res
}

func (s *enterpriseStatementService) Bill(req *model.StatementBillReq) model.StatementBillRes {
    // 判定结算日是否早于当天
    end := tools.NewTime().ParseDateStringX(req.End)
    if !end.Before(carbon.Now().StartOfDay().Carbon2Time()) {
        snag.Panic("结算日必须早于当天")
    }

    // 判定企业
    e, _ := ar.Ent.Enterprise.QueryNotDeleted().Where(enterprise.ID(req.ID)).WithCity().First(s.ctx)
    if e == nil {
        snag.Panic("未找到企业")
    }

    if e.Payment != model.EnterprisePaymentPostPay {
        snag.Panic("只有后付费企业才可请求")
    }

    // 查询未结账账单
    sta, bills := NewEnterprise().CalculateStatement(e, end)

    res := model.StatementBillRes{
        ID: e.ID,
        City: model.City{
            ID:   e.CityID,
            Name: e.Edges.City.Name,
        },
        ContactName:  e.ContactName,
        ContactPhone: e.ContactPhone,
        Start:        sta.Start.Format(carbon.DateLayout),
        End:          req.End,
        Bills:        bills,
    }

    overview := make(map[string]*model.BillOverview)
    var cost float64
    td := tools.NewDecimal()
    for _, bill := range bills {
        cost = td.Sum(cost, bill.Cost)

        key := fmt.Sprintf("%d-%.2f", bill.City.ID, bill.Voltage)
        ov, ok := overview[key]
        if !ok {
            ov = &model.BillOverview{
                Voltage: bill.Voltage,
                Number:  0,
                Price:   bill.Price,
                Days:    0,
                Cost:    0,
                City:    bill.City,
            }
            overview[key] = ov
        }
        ov.Number += 1
        ov.Days += bill.Days
        ov.Cost = td.Sum(ov.Cost, bill.Cost)
    }

    for _, ov := range overview {
        res.Overview = append(res.Overview, ov)
    }

    res.Cost = math.Round(cost*100) / 100

    return res
}
