// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
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
        orm: ent.Database.Order,
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
func (s *riderOrderService) List(riderID uint64, req model.PaginationReq) *model.PaginationRes {
    if req.PageSize > 10 {
        req.PageSize = 10
    }
    q := s.orm.QueryNotDeleted().
        WithCity().
        WithPlan().
        Order(ent.Desc(order.FieldCreatedAt)).
        Where(
            order.RiderID(riderID),
            order.TypeIn(model.OrderSubscribeTypes...),
        ).
        WithRider().
        WithSubscribe(func(sq *ent.SubscribeQuery) {
            sq.WithEmployee().WithStore()
        })
    return model.ParsePaginationResponse[model.Order, ent.Order](
        q,
        req,
        s.Detail,
    )
}

// Detail 订单详情
func (s *riderOrderService) Detail(item *ent.Order) model.Order {
    rc := item.Edges.City
    no := item.TradeNo
    if item.Payway == model.OrderPaywayManual {
        no = ""
    }
    res := model.Order{
        ID:         item.ID,
        Type:       item.Type,
        Status:     item.Status,
        Payway:     item.Payway,
        PayAt:      item.CreatedAt.Format(carbon.DateTimeLayout),
        Amount:     item.Amount,
        OutTradeNo: item.OutTradeNo,
        TradeNo:    no,
        City: model.City{
            ID:   rc.ID,
            Name: rc.Name,
        },
        PointAmount:   tools.NewDecimal().Mul(float64(item.Points), item.PointRatio),
        DiscountNewly: item.DiscountNewly,
        CouponAmount:  item.CouponAmount,
    }
    if len(item.Edges.Coupons) > 0 {
        res.Coupons = make([]model.CouponRider, len(item.Edges.Coupons))
        for i, c := range item.Edges.Coupons {
            res.Coupons[i] = NewCoupon().RiderDetail(c)
        }
    }
    // 骑士卡订阅订单
    op := item.Edges.Plan
    if op != nil {
        res.Plan = &model.Plan{
            ID:   op.ID,
            Name: op.Name,
            Days: op.Days,
        }
    }

    // 骑手信息
    or := item.Edges.Rider
    if or != nil {
        res.Rider = model.Rider{
            ID:    or.ID,
            Phone: or.Phone,
            Name:  or.Name,
        }
    }

    // store
    osub := item.Edges.Subscribe
    if osub != nil {
        res.Model = osub.Model
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

    // ebike
    bike := item.Edges.Ebike
    if bike != nil {
        res.Ebike = NewEbike().Detail(bike, bike.Edges.Brand)
    }
    return res
}

// Query 查询订单
func (s *riderOrderService) Query(riderID, orderID uint64) *ent.Order {
    item, _ := s.orm.QueryNotDeleted().
        Where(order.RiderID(riderID), order.ID(orderID)).
        WithCity().
        WithPlan().
        WithRider().
        WithSubscribe(func(sq *ent.SubscribeQuery) {
            sq.WithEmployee().WithStore()
        }).
        WithRefund().
        WithCoupons().
        First(s.ctx)
    if item == nil {
        snag.Panic("未找到订单")
    }
    return item
}
