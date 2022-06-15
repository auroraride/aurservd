// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
)

type riderOrderService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    orm      *ent.OrderClient
}

func NewRiderOrder() *riderOrderService {
    return &riderOrderService{
        ctx: context.Background(),
        orm: ar.Ent.Order,
    }
}

func NewRiderOrderWithRider(rider *ent.Rider) *riderOrderService {
    s := NewRiderOrder()
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

func NewRiderOrderWithModifier(m *model.Modifier) *riderOrderService {
    s := NewRiderOrder()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// List 分页获取骑手订单
func (s *riderOrderService) List(riderID uint64, req *model.PaginationReq) *model.PaginationRes {
    if req.PageSize > 10 {
        req.PageSize = 10
    }
    q := s.orm.QueryNotDeleted().
        WithCity().
        WithPlan(func(pq *ent.PlanQuery) {
            pq.WithPms()
        }).
        Order(ent.Desc(order.FieldCreatedAt)).
        Where(
            order.RiderID(riderID),
            order.TypeIn(model.OrderSubscribeTypes...),
        ).
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).
        WithSubscribe(func(sq *ent.SubscribeQuery) {
            sq.WithEmployee().WithStore()
        })
    return model.ParsePaginationResponse[model.RiderOrder, ent.Order](
        q,
        *req,
        s.Detail,
    )
}

// Detail 订单详情
func (s *riderOrderService) Detail(item *ent.Order) model.RiderOrder {
    rc := item.Edges.City
    res := model.RiderOrder{
        ID:         item.ID,
        Type:       item.Type,
        Status:     item.Status,
        Payway:     item.Payway,
        PayAt:      item.CreatedAt.Format(carbon.DateTimeLayout),
        Amount:     item.Amount,
        OutTradeNo: item.OutTradeNo,
        City: model.City{
            ID:   rc.ID,
            Name: rc.Name,
        },
    }
    // 骑士卡订阅订单
    op := item.Edges.Plan
    if op != nil {
        res.Plan = &model.Plan{
            ID:   op.ID,
            Name: op.Name,
            Days: op.Days,
        }
        for _, pm := range op.Edges.Pms {
            res.Models = append(res.Models, model.BatteryModel{
                ID:       pm.ID,
                Voltage:  pm.Voltage,
                Capacity: pm.Capacity,
            })
        }
    }

    // rider
    or := item.Edges.Rider
    if or != nil {
        res.Rider = model.RiderBasic{
            ID:    or.ID,
            Phone: or.Phone,
        }
        if or.Edges.Person != nil {
            res.Rider.Name = or.Edges.Person.Name
        }
    }

    // store
    osub := item.Edges.Subscribe
    if osub != nil {
        os := osub.Edges.Store
        if os != nil {
            res.Store = &model.Store{
                ID:   os.ID,
                Name: os.Name,
            }
        }

        oe := osub.Edges.Employee
        if oe != nil {
            res.Employee = &model.Employee{
                ID:    oe.ID,
                Name:  oe.Name,
                Phone: oe.Phone,
            }
        }
    }

    // refund
    rf := item.Edges.Refund
    if rf != nil {
        res.Refund = &model.Refund{
            Status:      rf.Status,
            Amount:      rf.Amount,
            OutRefundNo: rf.OutRefundNo,
            Reason:      rf.Reason,
            CreatedAt:   rf.CreatedAt.Format(carbon.DateTimeLayout),
            Remark:      rf.Remark,
            Modifier:    rf.LastModifier,
        }
        if rf.RefundAt != nil {
            res.Refund.RefundAt = rf.RefundAt.Format(carbon.DateTimeLayout)
        }
    }
    return res
}

// Query 查询订单
func (s *riderOrderService) Query(riderID, orderID uint64) *ent.Order {
    item, _ := s.orm.QueryNotDeleted().
        Where(order.RiderID(riderID), order.ID(orderID)).
        WithCity().
        WithPlan(func(pq *ent.PlanQuery) {
            pq.WithPms()
        }).
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).
        WithSubscribe(func(sq *ent.SubscribeQuery) {
            sq.WithEmployee().WithStore()
        }).
        WithRefund().
        First(s.ctx)
    if item == nil {
        snag.Panic("未找到订单")
    }
    return item
}
