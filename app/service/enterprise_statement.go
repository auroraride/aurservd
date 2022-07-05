// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/enterprise"
    "github.com/auroraride/aurservd/internal/ent/enterprisestatement"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
    "math"
    "time"
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
        orm: ent.Database.EnterpriseStatement,
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

func (s *enterpriseStatementService) GetBill(req *model.StatementBillReq) *model.StatementBillRes {
    // 判定结算日是否早于当天
    end := tools.NewTime().ParseDateStringX(req.End)
    if !end.Before(carbon.Now().StartOfDay().Carbon2Time()) {
        snag.Panic("结算日必须早于当天")
    }

    // 判定企业
    e, _ := ent.Database.Enterprise.QueryNotDeleted().Where(enterprise.ID(req.ID)).WithCity().First(s.ctx)
    if e == nil {
        snag.Panic("未找到企业")
    }

    if e.Payment != model.EnterprisePaymentPostPay {
        snag.Panic("只有后付费企业才可请求")
    }

    // 查询未结账账单
    es, bills := NewEnterprise().CalculateStatement(e, end)
    if es.Start.Sub(end).Seconds() > 0 {
        msg := fmt.Sprintf("无账单信息, 账单开始日期: %s, 所选截止日期: %s", es.Start.Format(carbon.DateLayout), end.Format(carbon.DateLayout))
        snag.Panic(msg)
    }

    uid := uuid.New().String()

    res := &model.StatementBillRes{
        ID:   e.ID,
        UUID: uid,
        City: model.City{
            ID:   e.CityID,
            Name: e.Edges.City.Name,
        },
        ContactName:  e.ContactName,
        ContactPhone: e.ContactPhone,
        Start:        es.Start.Format(carbon.DateLayout),
        End:          req.End,
        Bills:        bills,
        StatementID:  es.ID,
    }

    srv := NewEnterprise()

    overview := make(map[string]*model.BillOverview)
    var cost float64
    td := tools.NewDecimal()
    for _, bill := range bills {
        cost = td.Sum(cost, bill.Cost)
        res.Days += bill.Days

        key := srv.PriceKey(bill.City.ID, bill.Model)
        ov, ok := overview[key]
        if !ok {
            ov = &model.BillOverview{
                Model:  bill.Model,
                Number: 0,
                Price:  bill.Price,
                Days:   0,
                Cost:   0,
                City:   bill.City,
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

    // 缓存账单
    cache.Set(s.ctx, uid, res, 600*time.Second)

    return res
}

// Bill 后付费企业结账
func (s *enterpriseStatementService) Bill(req *model.StatementClearBillReq) {
    // 查找账单
    br := new(model.StatementBillRes)
    err := cache.Get(s.ctx, req.UUID).Scan(br)
    if err != nil {
        log.Error(err)
        snag.Panic("未找到账单信息")
        return
    }

    start := tools.NewTime().ParseDateStringX(br.Start)
    end := tools.NewTime().ParseDateStringX(br.End)

    if start.After(end) {
        snag.Panic("账单信息错误")
    }

    es, _ := s.orm.QueryNotDeleted().
        Where(
            enterprisestatement.ID(br.StatementID),
            enterprisestatement.SettledAtIsNil(),
        ).
        First(s.ctx)
    if es == nil {
        snag.Panic("未找到账单信息")
    }

    // 下个账单开始日
    next := carbon.Time2Carbon(end).StartOfDay().AddDay().Carbon2Time()

    tx, _ := ent.Database.Tx(s.ctx)
    _, err = tx.EnterpriseStatement.
        UpdateOneID(br.StatementID).
        SetSettledAt(time.Now()).
        SetEnd(end).
        SetRiderNumber(len(br.Bills)).
        SetDays(br.Days).
        SetCost(br.Cost).
        SetNillableRemark(req.Remark).
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 更新所有订阅
    riders := 0
    for _, bill := range br.Bills {
        var sub *ent.Subscribe
        sub, err = tx.Subscribe.UpdateOneID(bill.SubscribeID).SetLastBillDate(end).Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)

        if sub.EndAt == nil || !sub.EndAt.Before(next) {
            riders += 1
        }

        // 保存账单
        _, err = tx.EnterpriseBill.Create().
            SetEnterpriseID(es.EnterpriseID).
            SetStatementID(es.ID).
            SetStart(start).
            SetEnd(end).
            SetDays(bill.Days).
            SetPrice(bill.Price).
            SetCost(bill.Cost).
            SetCityID(bill.City.ID).
            SetNillableStationID(bill.StationID).
            SetRiderID(bill.RiderID).
            SetSubscribeID(bill.SubscribeID).
            SetModel(bill.Model).
            Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)
    }

    // 创建新账单
    _, err = tx.EnterpriseStatement.Create().
        SetEnterpriseID(es.EnterpriseID).
        SetStart(next).
        SetCost(tools.NewDecimal().Sub(es.Cost, br.Cost)).
        SetDays(es.Days - br.Days).
        SetRiderNumber(riders).
        SetBalance(es.Balance).
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _ = tx.Commit()

    cache.Del(s.ctx, req.UUID)
}
