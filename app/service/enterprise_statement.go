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
    "github.com/auroraride/aurservd/internal/ent/enterprisebill"
    "github.com/auroraride/aurservd/internal/ent/enterprisestatement"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    log "github.com/sirupsen/logrus"
    "math"
    "path/filepath"
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

func NewEnterpriseStatementWithModifier(m *model.Modifier) *enterpriseStatementService {
    s := NewEnterpriseStatement()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *enterpriseStatementService) SetContext(ctx context.Context) {
    s.ctx = ctx
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
            SetStart(model.DateFromTime(e.CreatedAt)).
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

    if e.Payment != model.EnterprisePaymentPostPay && !req.Force {
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

    start := model.DateFromStringX(br.Start)
    end := model.DateFromStringX(br.End)

    if start.After(end.Time) {
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
    next := end.Tomorrow()

    tx, _ := ent.Database.Tx(s.ctx)
    _, err = tx.EnterpriseStatement.
        UpdateOneID(br.StatementID).
        SetSettledAt(time.Now()).
        SetEnd(&end).
        SetRiderNumber(len(br.Bills)).
        SetDays(br.Days).
        SetCost(br.Cost).
        SetRemark(req.Remark).
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 更新所有订阅
    riders := 0
    for _, bill := range br.Bills {
        var sub *ent.Subscribe
        sub, err = tx.Subscribe.UpdateOneID(bill.SubscribeID).SetLastBillDate(&end).Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)

        if sub.EndAt == nil || !sub.EndAt.Before(next.Time) {
            riders += 1
        }

        // 保存账单
        _, err = tx.EnterpriseBill.Create().
            SetEnterpriseID(es.EnterpriseID).
            SetStatementID(es.ID).
            SetStart(model.DateFromStringX(bill.Start)).
            SetEnd(model.DateFromStringX(bill.End)).
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
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _ = tx.Commit()

    cache.Del(s.ctx, req.UUID)
}

// Historical 获取企业结账历史
func (s *enterpriseStatementService) Historical(req *model.StatementBillHistoricalListReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().
        Where(
            enterprisestatement.EnterpriseID(req.EnterpriseID),
            enterprisestatement.SettledAtNotNil(),
        ).Order(ent.Desc(enterprisestatement.FieldSettledAt))

    return model.ParsePaginationResponse(
        q,
        req.PaginationReq,
        func(item *ent.EnterpriseStatement) model.StatementBillHistoricalListRes {
            return model.StatementBillHistoricalListRes{
                ID:        item.ID,
                Cost:      item.Cost,
                Remark:    item.Remark,
                Creator:   item.Creator,
                Days:      item.Days,
                Start:     item.Start.Format(carbon.DateLayout),
                End:       item.End.Format(carbon.DateLayout),
                SettledAt: item.SettledAt.Format(carbon.DateTimeLayout),
            }
        },
    )
}

// Statement 账单详情
func (s *enterpriseStatementService) Statement(req *model.StatementBillDetailReq, w echo.Context) []model.StatementDetail {
    es, _ := ent.Database.EnterpriseStatement.QueryNotDeleted().
        Where(
            enterprisestatement.ID(req.ID),
            enterprisestatement.SettledAtNotNil(),
        ).
        WithEnterprise().
        First(s.ctx)
    if es == nil {
        snag.Panic("未找到账单")
    }

    e := es.Edges.Enterprise

    items, _ := ent.Database.EnterpriseBill.
        QueryNotDeleted().
        WithRider(func(query *ent.RiderQuery) {
            query.WithPerson()
        }).
        WithCity().
        WithStation().
        Where(enterprisebill.StatementID(req.ID)).
        All(s.ctx)

    res := make([]model.StatementDetail, len(items))
    for i, item := range items {
        r := item.Edges.Rider
        p := item.Edges.Rider.Edges.Person
        c := item.Edges.City
        t := item.Edges.Station
        res[i] = model.StatementDetail{
            Rider: model.RiderBasic{
                ID:    r.ID,
                Phone: r.Phone,
                Name:  p.Name,
            },
            City: model.City{
                ID:   c.ID,
                Name: c.Name,
            },
            Start: item.Start.Format(carbon.DateLayout),
            End:   item.End.Format(carbon.DateLayout),
            Days:  item.Days,
            Model: item.Model,
            Price: item.Price,
            Cost:  item.Cost,
        }
        if t != nil {
            res[i].Station = &model.EnterpriseStation{
                ID:   t.ID,
                Name: t.Name,
            }
        }
    }

    if req.Export {
        if len(items) == 0 {
            snag.Panic("无详细账单信息")
        }
        sheet := fmt.Sprintf("%s-%s", es.Start.Format(carbon.ShortDateLayout), es.End.Format(carbon.ShortDateLayout))
        fp := filepath.Join("runtime/export/statement", fmt.Sprintf("%s%s账单明细.xlsx", e.Name, sheet))
        ex := tools.NewExcelExistsExport(w, fp, sheet)
        if ex == nil {
            return nil
        }

        // 设置数据
        var rows [][]any
        rows = append(rows, []any{"姓名", "电话", "城市", "站点", "开始日期", "结束日期", "使用天数", "电池型号", "日单价", "费用"})
        for _, x := range res {
            so := ""
            if x.Station != nil {
                so = x.Station.Name
            }
            rows = append(rows, []any{
                x.Rider.Name,
                x.Rider.Phone,
                x.City.Name,
                so,
                x.Start,
                x.End,
                x.Days,
                x.Model,
                fmt.Sprintf("%.2f", x.Price),
                fmt.Sprintf("%.2f", x.Cost),
            })
        }

        ex.AddValues(rows).Done().Export()
        return nil
    }

    return res
}

func (s *enterpriseStatementService) usageItem(item *ent.Subscribe, end time.Time) model.StatementUsageRes {
    r := item.Edges.Rider
    c := item.Edges.City
    p := r.Edges.Person
    if item.EndAt != nil {
        end = *item.EndAt
    }
    start := *item.StartAt
    if item.LastBillDate != nil && item.LastBillDate.Before(start) {
        start = item.LastBillDate.Time
    }
    del := ""
    if r.DeletedAt != nil {
        del = r.DeletedAt.Format(carbon.DateTimeLayout)
    }
    res := model.StatementUsageRes{
        StatementDetail: model.StatementDetail{
            Start: item.StartAt.Format(carbon.DateLayout),
            End:   end.Format(carbon.DateLayout),
            Days:  tools.NewTime().UsedDays(end, *item.StartAt),
            Model: item.Model,
            City: model.City{
                ID:   c.ID,
                Name: c.Name,
            },
            Rider: model.RiderBasic{
                ID:    r.ID,
                Phone: r.Phone,
                Name:  p.Name,
            },
        },
        Status:    model.SubscribeStatusText(item.Status),
        DeletedAt: del,
    }
    a := item.Edges.Station
    if a != nil {
        res.Station = &model.EnterpriseStation{
            ID:   a.ID,
            Name: a.Name,
        }
    }
    return res
}

func (s *enterpriseStatementService) Usage(req *model.StatementUsageReq) *model.PaginationRes {
    q := ent.Database.Subscribe.QueryNotDeleted().
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).WithCity().WithStation().
        Where(subscribe.EnterpriseID(req.ID), subscribe.StartAtNotNil()).
        WithBills(func(bq *ent.EnterpriseBillQuery) {
            bq.Order(ent.Asc(enterprisebill.FieldEnd))
        })

    if req.Start != "" {
        q.Where(subscribe.StartAtGTE(tools.NewTime().ParseDateStringX(req.Start)))
    }
    end := carbon.Now().StartOfDay().Carbon2Time()
    if req.End != "" {
        end = tools.NewTime().ParseDateStringX(req.End)
        next := tools.NewTime().ParseNextDateStringX(req.End)
        q.Where(
            subscribe.StartAtLT(next),
            subscribe.Or(
                subscribe.EndAtIsNil(),
                subscribe.EndAtLT(next),
            ),
        )
    }

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Subscribe) model.StatementUsageRes {
        return s.usageItem(item, end)
    })
}
