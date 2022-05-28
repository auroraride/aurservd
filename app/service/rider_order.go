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
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "time"
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

func NewRiderOrderWithRider(rider *ent.Rider) *orderService {
    s := NewOrder()
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

// FindNotActived 查找用户未激活的新订单
func (s *riderOrderService) FindNotActived(riderID uint64) *ent.Order {
    o, _ := s.orm.QueryNotDeleted().Where(
        order.RiderID(riderID),
        order.Status(model.OrderStatusPaid),
        order.TypeIn(model.OrderTypeNewPlan, model.OrderTypeRenewal),
        order.StartAtIsNil(),
    ).Only(s.ctx)
    return o
}

// NotActivedDetail 获取未激活骑士卡详情
func (s *riderOrderService) NotActivedDetail(riderID uint64) *model.OrderNotActived {
    o, err := s.orm.QueryNotDeleted().
        Where(
            order.RiderID(riderID),
            order.Status(model.OrderStatusPaid),
            order.TypeIn(model.OrderTypeNewPlan, model.OrderTypeRenewal),
            order.StartAtIsNil(),
        ).
        WithCity().
        WithPlan(func(pq *ent.PlanQuery) {
            pq.WithPms()
        }).
        WithChildren(func(cq *ent.OrderQuery) {
            cq.Where(order.Type(model.OrderTypeDeposit))
        }).
        Only(s.ctx)
    if err != nil {
        return nil
    }
    item := &model.OrderNotActived{
        ID:     o.ID,
        Amount: o.Amount,
        Total:  o.Amount,
        Payway: o.Payway,
        Plan:   o.PlanDetail,
        City: model.City{
            ID:   o.Edges.City.ID,
            Name: o.Edges.City.Name,
        },
        Time: o.CreatedAt.Format(carbon.DateLayout),
    }
    for _, pm := range o.Edges.Plan.Edges.Pms {
        item.Models = append(item.Models, model.BatteryModel{
            ID:       pm.ID,
            Voltage:  pm.Voltage,
            Capacity: pm.Capacity,
        })
    }
    if len(o.Edges.Children) > 0 {
        item.Deposit = o.Edges.Children[0].Amount
    }
    item.Total += item.Deposit
    return item
}

// Recent 获取最近的订单
func (s *riderOrderService) Recent(riderID uint64) *model.RiderRecentOrder {
    now := time.Now()
    tt := tools.NewTime()
    o, _ := ar.Ent.Order.
        QueryNotDeleted().
        WithAlters().
        WithPauses().
        WithPlan(
            func(pq *ent.PlanQuery) {
                pq.WithPms()
            },
        ).
        Where(
            order.RiderID(riderID),
            order.Status(model.OrderStatusPaid),
            order.TypeIn(model.OrderRiderPlan...),
        ).
        Only(s.ctx)
    if o == nil {
        return nil
    }
    ro := &model.RiderRecentOrder{
        PlanID:   o.PlanID,
        PlanName: o.PlanDetail.Name,
        Days:     int(o.PlanDetail.Days),
        // Start:       "",
        // End:         "",
        Voltage: o.Edges.Plan.Edges.Pms[0].Voltage,
    }

    // 计算改动天数
    alter := 0
    for _, oa := range o.Edges.Alters {
        alter += oa.Days
    }

    // 计算暂停天数
    paused := 0
    for _, op := range o.Edges.Pauses {
        if op.EndAt.IsZero() {
            paused += tt.SubDays(now, op.StartAt)
        } else {
            paused += op.Days
        }
    }

    // 距离激活已过去多少天
    past := 0
    if o.StartAt != nil {
        ro.StartAt = o.StartAt.Format(carbon.DateLayout)
        past = tt.SubDays(now, *o.StartAt)
    }

    remaining := int(o.Days) + alter + paused - past
    ro.Remaining = remaining
    ro.PausedDays = paused
    ro.AlterDays = alter

    switch true {
    case o.StartAt == nil:
        // 未激活
        ro.Status = model.RiderOrderStatusPending
        break
    case o.PausedAt.IsZero() && o.StartAt != nil:
        if remaining >= 0 {
            // 计费中
            ro.Status = model.RiderOrderStatusNormal
        } else {
            ro.Status = model.RiderOrderStatusOverdue
        }
        break
    case !o.PausedAt.IsZero():
        // 暂停中
        ro.Status = model.RiderOrderStatusPaused
        break
    case o.EndAt != nil:
        // 已归还
        ro.Status = model.RiderOrderStatusRemanded
        ro.EndAt = o.EndAt.Format(carbon.DateLayout)
        break
    }

    return ro
}
