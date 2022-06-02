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
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "time"
)

type subscribeService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    orm      *ent.SubscribeClient
}

func NewSubscribe() *subscribeService {
    return &subscribeService{
        ctx: context.Background(),
        orm: ar.Ent.Subscribe,
    }
}

func NewSubscribeWithRider(rider *ent.Rider) *orderService {
    s := NewOrder()
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

func NewSubscribeWithModifier(m *model.Modifier) *subscribeService {
    s := NewSubscribe()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// QueryRecentOnly 获取骑手最近的一次订阅
func (s *subscribeService) QueryRecentOnly(riderID uint64) *ent.Subscribe {
    sub, _ := s.orm.QueryNotDeleted().
        Where(subscribe.RiderID(riderID)).
        Order(ent.Desc(subscribe.FieldCreatedAt)).
        Only(s.ctx)
    return sub
}

// Recent 获取骑手最近订阅详情
func (s *subscribeService) Recent(riderID uint64) *model.Subscribe {
    sub, err := s.orm.QueryNotDeleted().
        Where(subscribe.RiderID(riderID)).
        Order(ent.Desc(subscribe.FieldCreatedAt)).
        WithPlan(func(pq *ent.PlanQuery) {
            pq.WithPms()
        }).
        WithCity().
        // 查询初始订单和对应的押金订单
        WithInitialOrder(func(ioq *ent.OrderQuery) {
            ioq.WithChildren(func(doq *ent.OrderQuery) {
                doq.Where(order.Type(model.OrderTypeDeposit))
            })
        }).
        Limit(1).
        First(s.ctx)

    if err != nil {
        log.Error(err)
    }

    if sub == nil {
        return nil
    }

    res := &model.Subscribe{
        ID:        sub.ID,
        Status:    sub.Status,
        Voltage:   sub.Voltage,
        Days:      sub.Days,
        PauseDays: sub.PauseDays,
        AlterDays: sub.AlterDays,
        Remaining: sub.Remaining,
        City: &model.City{
            ID:   sub.Edges.City.ID,
            Name: sub.Edges.City.Name,
        },
        Plan: &model.Plan{
            ID:   sub.Edges.Plan.ID,
            Name: sub.Edges.Plan.Name,
            Days: sub.Edges.Plan.Days,
        },
    }

    startAt := time.Now()
    if sub.StartAt != nil {
        startAt = *sub.StartAt
        res.StartAt = sub.StartAt.Format(carbon.DateLayout)
    } else {
        res.StartAt = ""
    }
    if sub.EndAt != nil {
        res.EndAt = sub.EndAt.Format(carbon.DateLayout)
    } else {
        res.EndAt = startAt.AddDate(0, 0, sub.Remaining).Format(carbon.DateLayout)
    }

    for _, pm := range sub.Edges.Plan.Edges.Pms {
        res.Models = append(res.Models, model.BatteryModel{
            ID:       pm.ID,
            Voltage:  pm.Voltage,
            Capacity: pm.Capacity,
        })
    }

    // 已取消(退款)不显示到期日期
    if res.Status == model.SubscribeStatusCanceled {
        res.EndAt = "-"
    }

    // 如果是未激活则查询支付信息
    if res.Status == model.SubscribeStatusInactive {
        o := sub.Edges.InitialOrder
        res.Order = &model.SubscribeOrderInfo{
            ID:     o.ID,
            Status: o.Status,
            PayAt:  o.CreatedAt.Format(carbon.DateTimeLayout),
            Payway: o.Payway,
            Amount: o.Amount,
            Total:  o.Total,
        }
        if len(o.Edges.Children) > 0 {
            res.Order.Deposit = o.Edges.Children[0].Amount
        }
    }

    return res
}

// QueryAllEffective 获取所有生效的订阅
func (s *subscribeService) QueryAllEffective() []*ent.Subscribe {
    items, _ := ar.Ent.Subscribe.Query().
        Where(
            subscribe.RefundAtIsNil(),
            subscribe.EndAtIsNil(),
            subscribe.StartAtNotNil(),
        ).
        WithPauses(func(spq *ent.SubscribePauseQuery) {
            spq.Where(subscribepause.EndAtIsNil())
        }).
        All(s.ctx)
    return items
}

// UpdateStatus 更新订阅状态
func (s *subscribeService) UpdateStatus(item *ent.Subscribe) {
    tt := tools.NewTime()
    pauseDays := item.PauseDays
    alterDays := item.AlterDays
    days := item.Days
    pastDays := tt.SubDaysToNow(*item.StartAt)
    status := model.SubscribeStatusUsing
    // 寄存中
    if item.PausedAt != nil && item.Edges.Pauses != nil {
        p := item.Edges.Pauses[0]
        pauseDays += tt.SubDaysToNow(p.StartAt)
    }
    // 剩余天数
    remaining := days + alterDays + pauseDays - pastDays
    if remaining < 0 {
        status = model.SubscribeStatusOverdue
    }
    _, err := item.Update().
        SetPauseDays(pauseDays).
        SetStatus(status).
        SetRemaining(remaining).
        Save(context.Background())
    if err != nil {
        log.Errorf("[SUBSCRIBE TASK] %d 更新失败: %v", item.ID, err)
    }
    log.Infof("[SUBSCRIBE TASK] %d 更新成功, 状态: %d, 剩余天数: %d", item.ID, status, remaining)
}
