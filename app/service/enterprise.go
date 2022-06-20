// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-05
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "errors"
    "fmt"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/enterprise"
    "github.com/auroraride/aurservd/internal/ent/enterprisecontract"
    "github.com/auroraride/aurservd/internal/ent/enterpriseprice"
    "github.com/auroraride/aurservd/internal/ent/enterprisestatement"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "github.com/jinzhu/copier"
    "github.com/shopspring/decimal"
    log "github.com/sirupsen/logrus"
    "time"
)

type enterpriseService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    orm      *ent.EnterpriseClient
}

func NewEnterprise() *enterpriseService {
    return &enterpriseService{
        ctx: context.Background(),
        orm: ar.Ent.Enterprise,
    }
}

func NewEnterpriseWithRider(r *ent.Rider) *enterpriseService {
    s := NewEnterprise()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewEnterpriseWithModifier(m *model.Modifier) *enterpriseService {
    s := NewEnterprise()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *enterpriseService) Query(id uint64) (*ent.Enterprise, error) {
    return s.orm.QueryNotDeleted().Where(enterprise.ID(id)).Only(s.ctx)
}

func (s *enterpriseService) QueryX(id uint64) *ent.Enterprise {
    e, _ := s.Query(id)
    if e == nil {
        snag.Panic("未找到有效企业")
    }
    return e
}

// Create 创建企业
func (s *enterpriseService) Create(req *model.EnterpriseDetail) uint64 {

    tx, err := ar.Ent.Tx(s.ctx)
    if err != nil {
        snag.Panic(err)
    }

    e := &ent.Enterprise{}
    e, err = ent.EntitySetAttributes[ent.EnterpriseCreate, ent.Enterprise](tx.Enterprise.Create(), e, req).Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    s.SaveEnterprise(tx, e, req)
    _ = tx.Commit()

    return e.ID
}

// Modify 修改企业
func (s *enterpriseService) Modify(req *model.EnterpriseDetailWithID) {
    e := s.QueryX(req.ID)
    // 判定付费方式是否允许改变
    if *req.Payment != e.Payment {
        set := NewEnterpriseStatementWithModifier(s.modifier).Current(e)
        if set.Cost > 0 {
            snag.Panic("企业已产生费用, 无法修改支付方式")
        }
    }

    tx, err := ar.Ent.Tx(s.ctx)
    if err != nil {
        snag.Panic(err)
    }

    e, err = tx.Enterprise.ModifyOne(e, req.EnterpriseDetail).Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _, err = tx.EnterprisePrice.Delete().Where(enterpriseprice.EnterpriseID(e.ID)).Exec(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    _, err = tx.EnterpriseContract.Delete().Where(enterprisecontract.EnterpriseID(e.ID)).Exec(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    s.SaveEnterprise(tx, e, req.EnterpriseDetail)
    _ = tx.Commit()
}

// SaveEnterprise 保存企业信息
func (s *enterpriseService) SaveEnterprise(tx *ent.Tx, e *ent.Enterprise, req *model.EnterpriseDetail) {
    var err error
    // 存储价格信息
    cvm := make(map[string]struct{})
    for _, rp := range req.Prices {
        // 判断价格是否重复
        k := fmt.Sprintf("%d-%f", rp.CityID, rp.Voltage)
        if _, ok := cvm[k]; ok {
            snag.PanicCallbackX(tx.Rollback, "价格重复")
        }
        _, err = tx.EnterprisePrice.Create().SetPrice(rp.Price).SetCityID(rp.CityID).SetVoltage(rp.Voltage).SetEnterprise(e).Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)
        cvm[k] = struct{}{}
    }

    // 存储合同
    tt := tools.NewTime()
    var dates [][]int64
    for _, rc := range req.Contracts {
        rcs := tt.ParseDateStringX(rc.Start)
        rce := tt.ParseDateStringX(rc.End)
        for _, r := range dates {
            if rcs.Unix() <= r[0] && rce.Unix() >= r[1] {
                snag.PanicCallbackX(tx.Rollback, "日期重叠")
            }
        }
        _, err = tx.EnterpriseContract.Create().SetFile(rc.File).SetStart(rcs).SetEnd(rce).SetEnterprise(e).Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)
        dates = append(dates, []int64{rcs.Unix(), rce.Unix()})
    }
}

// DetailQuery 企业详情查询语句
func (s *enterpriseService) DetailQuery() *ent.EnterpriseQuery {
    return s.orm.QueryNotDeleted().WithStatements().WithCity().WithPrices(func(ep *ent.EnterprisePriceQuery) {
        ep.WithCity()
    }).WithContracts(func(ecq *ent.EnterpriseContractQuery) {
        ecq.Order(ent.Desc(enterprisecontract.FieldEnd))
    }).WithRiders(func(erq *ent.RiderQuery) {
        erq.Where(rider.DeletedAtIsNil())
    }).WithStatements(func(esq *ent.EnterpriseStatementQuery) {
        esq.Where(enterprisestatement.Date(carbon.Now().StartOfDay().Carbon2Time()))
    })
}

// List 列举企业
func (s *enterpriseService) List(req *model.EnterpriseListReq) *model.PaginationRes {
    tt := tools.NewTime()

    q := s.DetailQuery()
    if req.Name != nil {
        q.Where(enterprise.NameContainsFold(*req.Name))
    }
    if req.CityID != nil {
        q.Where(enterprise.CityID(*req.CityID))
    }
    if req.Status != nil {
        q.Where(enterprise.Status(*req.Status))
    }
    if req.Payment != nil {
        q.Where(enterprise.Payment(*req.Payment))
    }
    if req.ContactKeyword != nil {
        q.Where(enterprise.Or(
            enterprise.ContactNameContainsFold(*req.ContactKeyword),
            enterprise.ContactPhoneContainsFold(*req.ContactKeyword),
            enterprise.IdcardNumberContainsFold(*req.ContactKeyword),
        ))
    }
    if req.Start != nil {
        q.Where(enterprise.HasContractsWith(enterprisecontract.StartLTE(tt.ParseDateStringX(*req.Start))))
    }
    if req.End != nil {
        q.Where(enterprise.HasContractsWith(enterprisecontract.EndGTE(tt.ParseDateStringX(*req.End))))
    }
    return model.ParsePaginationResponse(
        q, req.PaginationReq,
        func(item *ent.Enterprise) model.EnterpriseRes {
            return s.Detail(item)
        },
    )
}

// GetDetail 获取企业详情
func (s *enterpriseService) GetDetail(req *model.IDParamReq) model.EnterpriseRes {
    item, _ := s.DetailQuery().Where(enterprise.ID(req.ID)).WithStatements(func(esq *ent.EnterpriseStatementQuery) {
        esq.Where(enterprisestatement.Date(carbon.Now().StartOfDay().Carbon2Time()))
    }).First(s.ctx)
    if item == nil {
        snag.Panic("未找到有效企业")
    }
    return s.Detail(item)
}

// Detail 企业详情
func (s *enterpriseService) Detail(item *ent.Enterprise) (res model.EnterpriseRes) {
    _ = copier.Copy(&res, item)
    res.Riders = len(item.Edges.Riders)
    res.City = model.City{
        ID:   item.Edges.City.ID,
        Name: item.Edges.City.Name,
    }
    contracts := item.Edges.Contracts
    if contracts != nil {
        res.Contracts = make([]model.EnterpriseContract, len(contracts))
        for i, ec := range contracts {
            res.Contracts[i] = model.EnterpriseContract{
                Start: ec.Start.Format(carbon.DateLayout),
                End:   ec.End.Format(carbon.DateLayout),
                File:  ec.File,
            }
        }
    }

    prices := item.Edges.Prices
    if prices != nil {
        res.Prices = make([]model.EnterprisePriceWithCity, len(prices))
        for i, ep := range prices {
            res.Prices[i] = model.EnterprisePriceWithCity{
                Voltage: ep.Voltage,
                Price:   ep.Price,
                City: model.City{
                    ID:   ep.Edges.City.ID,
                    Name: ep.Edges.City.Name,
                },
            }
        }
    }

    stas := item.Edges.Statements
    if item.Payment == model.EnterprisePaymentPostPay && stas != nil && len(stas) == 1 {
        res.Unsettlement = stas[0].Days
    }
    return
}

// QueryAllCollaborated 获取合作中的企业
func (s *enterpriseService) QueryAllCollaborated() []*ent.Enterprise {
    items, _ := s.orm.QueryNotDeleted().
        Where(enterprise.Status(model.EnterpriseStatusCollaborated)).
        WithPrices().
        All(s.ctx)
    return items
}

// QueryAllUsingSubscribe 获取所有待结算骑手团签订阅
func (s *enterpriseService) QueryAllUsingSubscribe(enterpriseID uint64, args ...any) []*ent.Subscribe {
    q := ar.Ent.Subscribe.QueryNotDeleted().
        Where(
            // 所属企业
            subscribe.EnterpriseID(enterpriseID),
            // 已开始使用
            subscribe.StartAtNotNil(),
        )

    end := time.Now()
    if len(args) > 0 {
        if d, ok := args[0].(time.Time); ok {
            end = d
        }
    }
    q.Where(
        // 获取未结算订阅: 上次结算日期为空或小于截止日期
        subscribe.Or(
            subscribe.LastBillDateIsNil(),
            subscribe.LastBillDateLTE(carbon.Time2Carbon(end).StartOfDay().Carbon2Time()),
        ),
    )

    items, _ := q.WithCity().All(s.ctx)
    return items
}

// GetPrices 获取企业价格表
func (s *enterpriseService) GetPrices(item *ent.Enterprise) (res map[string]float64) {
    var items []*ent.EnterprisePrice
    if item.Edges.Prices == nil {
        items, _ = ar.Ent.EnterprisePrice.QueryNotDeleted().Where(enterpriseprice.EnterpriseID(item.ID)).All(s.ctx)
    } else {
        items = item.Edges.Prices
    }

    res = make(map[string]float64)
    for _, ep := range items {
        res[fmt.Sprintf("%d-%.2f", ep.CityID, ep.Voltage)] = ep.Price
    }

    return res
}

func (s *enterpriseService) CalculateStatement(e *ent.Enterprise, end time.Time) (sta *ent.EnterpriseStatement, bills []model.StatementBillData) {
    tt := tools.NewTime()
    prices := s.GetPrices(e)
    sta = NewEnterpriseStatement().Current(e)

    // 获取所有骑手订阅
    subs := s.QueryAllUsingSubscribe(e.ID, end)
    bills = make([]model.StatementBillData, len(subs))
    for i, sub := range subs {
        // 计算使用天数
        var used int

        // 判定是否退订
        // 上次结算日期存在则从上次结算日开始计算
        from := carbon.Time2Carbon(*sub.StartAt).StartOfDay().Carbon2Time()
        to := end
        if sub.LastBillDate != nil {
            from = *sub.LastBillDate
        }

        // 是否已结束
        if sub.EndAt != nil && sub.EndAt.Before(end) {
            to = carbon.Time2Carbon(*sub.EndAt).StartOfDay().Carbon2Time()
        }
        used = tt.UsedDays(to, from)

        // 按城市/型号计算金额
        k := fmt.Sprintf("%d-%.2f", sub.CityID, sub.Voltage)
        p, ok := prices[k]
        if !ok {
            log.Errorf("%d [%d] 获取价格失败", sub.ID, sub.CityID)
        }

        cost, _ := decimal.NewFromFloat(p).Mul(decimal.NewFromInt(int64(used))).Float64()

        bills[i] = model.StatementBillData{
            EnterpriseID: *sub.EnterpriseID,
            RiderID:      sub.RiderID,
            SubscribeID:  sub.ID,
            City: model.City{
                ID:   sub.Edges.City.ID,
                Name: sub.Edges.City.Name,
            },
            StatementID: sta.ID,
            Days:        used,
            End:         end.Format(carbon.DateLayout),
            Start:       from.Format(carbon.DateLayout),
            Cost:        cost,
            Price:       p,
            Voltage:     sub.Voltage,
        }
    }

    return
}

// UpdateStatement 更新企业账单
func (s *enterpriseService) UpdateStatement(e *ent.Enterprise) {
    sta, bills := s.CalculateStatement(e, time.Now())

    // 总天数
    var days int
    // 总计费用
    var cost float64
    td := tools.NewDecimal()
    for _, bill := range bills {
        cost = td.Sum(cost, bill.Cost)
        days += bill.Days
    }

    // 企业付款方式
    var balance float64
    switch e.Payment {
    case model.EnterprisePaymentPrepay:
        // 预付费, 计算余额
        balance = tools.NewDecimal().Sub(sta.Balance, cost)
        break
    }

    _, err := e.Update().SetBalance(balance).Save(s.ctx)
    if err != nil {
        log.Errorf("[ENTERPRISE TASK] %d 更新失败: %v", e.ID, err)
    }

    now := carbon.Now().StartOfDay().Carbon2Time()
    _, err = sta.Update().
        SetRiderNumber(len(bills)).
        SetBalance(balance).
        SetDays(days).
        SetCost(cost).
        SetDate(now).
        Save(s.ctx)

    if err != nil {
        log.Errorf("[ENTERPRISE TASK] %d 更新失败: %v", e.ID, err)
    }

    log.Infof("[ENTERPRISE TASK] EntperirseID:[%d] 更新成功, 总使用人数: %d, 账期使用总天数: %d, 总费用: %.2f, 余额: %.2f, 出账日期: %s",
        e.ID,
        len(bills),
        days,
        cost,
        balance,
        now.Format(carbon.DateLayout),
    )
}

// Prepayment 预付费
func (s *enterpriseService) Prepayment(req *model.EnterprisePrepaymentReq) float64 {
    e := s.QueryX(req.ID)
    if e.Payment == model.EnterprisePaymentPostPay {
        snag.Panic("该企业支付方式为后付费")
    }

    set := NewEnterpriseStatementWithModifier(s.modifier).Current(e)

    before := e.Balance
    tx, _ := ar.Ent.Tx(s.ctx)

    // 创建预付费记录
    _, err := tx.EnterprisePrepayment.Create().SetEnterpriseID(e.ID).SetAmount(req.Amount).SetRemark(req.Remark).Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    td := tools.NewDecimal()

    // 更新余额
    // 账单表
    balance := td.Sum(e.Balance, req.Amount)
    _, err = tx.EnterpriseStatement.UpdateOne(set).SetBalance(balance).Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 更新企业表
    _, err = tx.Enterprise.UpdateOne(e).SetBalance(balance).Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 记录日志
    go logging.NewOperateLog().
        SetRef(e).
        SetModifier(s.modifier).
        SetOperate(model.OperateEnterprisePrepayment).
        SetDiff(fmt.Sprintf("余额%.2f元", before), fmt.Sprintf("余额%.2f元", balance)).
        SetRemark(req.Remark).
        Send()

    _ = tx.Commit()
    return balance
}

// Business 企业是否可以办理业务
func (s *enterpriseService) Business(e *ent.Enterprise) error {
    if e == nil {
        return errors.New("未找到企业信息")
    }

    if e.Status != model.EnterpriseStatusCollaborated {
        return errors.New("企业已终止合作")
    }

    if e.Payment == model.EnterprisePaymentPrepay && e.Balance < 0 {
        return errors.New("企业已欠费")
    }

    return nil
}
