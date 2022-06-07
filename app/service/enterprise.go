// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-05
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/enterprise"
    "github.com/auroraride/aurservd/internal/ent/enterprisecontract"
    "github.com/auroraride/aurservd/internal/ent/enterpriseprice"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "github.com/jinzhu/copier"
    "github.com/shopspring/decimal"
    log "github.com/sirupsen/logrus"
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

func (s *enterpriseService) Query(id uint64) *ent.Enterprise {
    e, _ := s.orm.QueryNotDeleted().Where(enterprise.ID(id)).Only(s.ctx)
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
    e := s.Query(req.ID)
    // 判定付费方式是否允许改变
    if *req.Payment != e.Payment {
        set := NewEnterpriseStatementWithModifier(s.modifier).Current(e.ID)
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

    tx.EnterprisePrice.Delete().Where(enterpriseprice.EnterpriseID(e.ID)).ExecX(s.ctx)
    tx.EnterpriseContract.Delete().Where(enterprisecontract.EnterpriseID(e.ID)).ExecX(s.ctx)

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
    item, _ := s.DetailQuery().Where(enterprise.ID(req.ID)).First(s.ctx)
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
func (s *enterpriseService) QueryAllUsingSubscribe(enterpriseID uint64) []*ent.Subscribe {
    items, _ := ar.Ent.Subscribe.QueryNotDeleted().Where(
        // 所属企业
        subscribe.EnterpriseID(enterpriseID),
        // 未结算
        subscribe.StatementIDIsNil(),
    ).All(s.ctx)
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
        res[fmt.Sprintf("%d-%f", ep.CityID, ep.Price)] = ep.Price
    }

    return res
}

// UpdateStatement 更新企业账单
func (s *enterpriseService) UpdateStatement(item *ent.Enterprise) {
    tt := tools.NewTime()
    prices := s.GetPrices(item)
    statement := NewEnterpriseStatement().Current(item.ID)

    // 获取所有骑手订阅
    var days int
    cost := decimal.NewFromFloat(0)
    subs := s.QueryAllUsingSubscribe(item.ID)
    for _, sub := range subs {
        // 计算使用天数
        var used int
        // 判定是否退订
        if sub.EndAt != nil {
            used = tt.DiffDaysOfStart(*sub.EndAt, *sub.StartAt)
        } else {
            // 排除当前时间当日0点
            used = tt.DiffDaysOfStartToNow(*sub.StartAt)
        }

        // 保存骑手费用详细
        // NewEnterpriseInvoice().Current(sub).Update().SetDays(used)

        // 总天数
        days += used

        // 按城市/型号计算金额
        k := fmt.Sprintf("%d-%f", sub.CityID, sub.Voltage)
        if p, ok := prices[k]; ok {
            cost = cost.Add(decimal.NewFromFloat(p))
        }
    }

    // 企业付款方式
    var balance float64
    switch item.Payment {
    case model.EnterprisePaymentPrepay:
        // 预付费, 计算余额
        balance, _ = decimal.NewFromFloat(statement.Amount).Sub(cost).Float64()
        break
    }

    _, err := item.Update().SetBalance(balance).Save(s.ctx)
    if err != nil {
        log.Errorf("[ENTERPRISE TASK] %d 更新失败: %v", item.ID, err)
    }

    cf, _ := cost.Float64()
    _, err = statement.Update().
        SetRiderNumber(len(subs)).
        SetBalance(balance).
        SetDays(days).
        SetCost(cf).
        SetBillTime(carbon.Now().StartOfDay().AddDays(-1).Carbon2Time()).
        Save(s.ctx)

    if err != nil {
        log.Errorf("[ENTERPRISE TASK] %d 更新失败: %v", item.ID, err)
    }

    log.Infof("[ENTERPRISE TASK] %d 更新成功, 账期使用总天数: %d, 总费用: %v", item.ID, days, cost)
}

// Prepayment 预付费
func (s *enterpriseService) Prepayment(req *model.EnterprisePrepaymentReq) float64 {
    e := s.Query(req.ID)
    if e.Payment == model.EnterprisePaymentPostPay {
        snag.Panic("该企业支付方式为后付费")
    }

    set := NewEnterpriseStatementWithModifier(s.modifier).Current(e.ID)

    before := e.Balance
    tx, _ := ar.Ent.Tx(s.ctx)

    // 创建预付费记录
    _, err := tx.EnterprisePrepayment.Create().SetEnterpriseID(e.ID).SetAmount(req.Amount).SetRemark(req.Remark).Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 更新余额
    b, _ := decimal.NewFromFloat(set.Balance).Add(decimal.NewFromFloat(req.Amount)).Float64()
    a, _ := decimal.NewFromFloat(set.Amount).Add(decimal.NewFromFloat(req.Amount)).Float64()

    // 更新账单表
    _, err = tx.EnterpriseStatement.UpdateOne(set).SetBalance(b).SetAmount(a).Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 更新企业表
    _, err = tx.Enterprise.UpdateOne(e).SetBalance(b).Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 记录日志
    go logging.NewOperateLog().
        SetRef(e).
        SetModifier(s.modifier).
        SetOperate(logging.OperateEnterprisePrepayment).
        SetDiff(fmt.Sprintf("余额%.2f元", before), fmt.Sprintf("余额%.2f元", b)).
        SetRemark(req.Remark).
        Send()

    _ = tx.Commit()
    return b
}
