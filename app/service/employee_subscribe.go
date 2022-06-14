// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "strconv"
    "strings"
    "time"
)

type employeeSubscribeService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
}

func NewEmployeeSubscribe() *employeeSubscribeService {
    return &employeeSubscribeService{
        ctx: context.Background(),
    }
}

func NewEmployeeSubscribeWithRider(r *ent.Rider) *employeeSubscribeService {
    s := NewEmployeeSubscribe()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewEmployeeSubscribeWithModifier(m *model.Modifier) *employeeSubscribeService {
    s := NewEmployeeSubscribe()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewEmployeeSubscribeWithEmployee(e *ent.Employee) *employeeSubscribeService {
    s := NewEmployeeSubscribe()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// Inactive 获取骑手待激活订阅详情
func (s *employeeSubscribeService) Inactive(qr string) *model.SubscribeActiveInfo {
    if s.employee.Edges.Store == nil {
        snag.Panic("你未上班")
    }
    if strings.HasPrefix(qr, "SUBSCRIBE:") {
        qr = strings.ReplaceAll(qr, "SUBSCRIBE:", "")
    }
    id, _ := strconv.ParseUint(strings.TrimSpace(qr), 10, 64)
    // 查询订单状态
    sub, _ := ar.Ent.Subscribe.QueryNotDeleted().
        Where(
            subscribe.ID(id),
            subscribe.RefundAtIsNil(),
            subscribe.StartAtIsNil(),
            subscribe.Or(
                subscribe.Type(0),
                subscribe.TypeIn(model.OrderTypeNewly, model.OrderTypeAgain),
            ),
            subscribe.Status(model.SubscribeStatusInactive),
        ).
        WithInitialOrder(func(oq *ent.OrderQuery) {
            oq.WithPlan().WithCommission()
        }).
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).
        WithEnterprise().
        Only(s.ctx)

    if sub == nil {
        snag.Panic("未找到待激活骑士卡")
    }

    r := sub.Edges.Rider
    if r == nil {
        snag.Panic("骑手信息获取失败")
    }

    p := sub.Edges.Rider.Edges.Person
    if p.Status != model.PersonAuthenticated.Raw() {
        snag.Panic("骑手还未认证")
    }

    res := &model.SubscribeActiveInfo{
        ID:           sub.ID,
        EnterpriseID: sub.EnterpriseID,
        Voltage:      sub.Voltage,
        CommissionID: nil,
        Rider: model.RiderBasic{
            ID:    r.ID,
            Phone: r.Phone,
            Name:  p.Name,
        },
    }

    // TODO 企业订单店员是否有提成?
    if sub.EnterpriseID == nil {
        o := sub.Edges.InitialOrder
        if o == nil || o.Status != model.OrderStatusPaid {
            snag.Panic("订单状态异常") // TODO 退款被拒绝的如何操作
        }
        res.Order = &model.SubscribeOrderInfo{
            ID:      o.ID,
            Status:  o.Status,
            PayAt:   o.CreatedAt.Format(carbon.DateTimeLayout),
            Payway:  o.Payway,
            Amount:  o.Amount,
            Deposit: o.Total - o.Amount,
            Total:   o.Total,
        }

        c := sub.Edges.InitialOrder.Edges.Commission
        if c != nil {
            res.CommissionID = &c.ID
        }
    } else {
        res.Enterprise = &model.EnterpriseBasic{
            ID:   sub.Edges.Enterprise.ID,
            Name: sub.Edges.Enterprise.Name,
        }
    }

    return res
}

// Active 激活订单
func (s *employeeSubscribeService) Active(req *model.QRPostReq) {
    info := s.Inactive(req.Qrcode)

    tx, _ := ar.Ent.Tx(s.ctx)

    // 激活
    _, err := tx.Subscribe.UpdateOneID(info.ID).
        SetStatus(model.SubscribeStatusUsing).
        SetStartAt(time.Now()).
        SetEmployeeID(s.employee.ID).
        SetStoreID(s.employee.Edges.Store.ID).
        Save(s.ctx)
    snag.PanicIfErrorX(err, tx.Rollback)

    // 团签账单
    if info.EnterpriseID != nil {
        sm := NewEnterpriseStatement().Current(*info.EnterpriseID)
        _, err = tx.EnterpriseStatement.UpdateOne(sm).AddRiderNumber(1).Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)
    }

    // 提成
    if info.CommissionID != nil {
        _, _ = tx.Commission.UpdateOneID(*info.CommissionID).SetEmployeeID(s.employee.ID).Save(s.ctx)
    }

    // 调出库存
    err = NewStockWithEmployee(s.employee).BatteryOutboundWithRider(
        tx.Stock.Create(),
        info.Rider.ID,
        s.employee.Edges.Store.ID,
        s.employee.ID,
        info.Voltage,
        model.StockTypeRiderObtain,
    )
    if err != nil {
        snag.PanicIfErrorX(err, tx.Rollback)
    }

    _ = tx.Commit()
}
