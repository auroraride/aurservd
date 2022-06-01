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

// GetStatus 获取订阅状态
func (s *subscribeService) GetStatus(sub *ent.Subscribe) uint8 {
    // 已退款
    if sub.RefundAt != nil {
        return model.SubscribeStatusCanceled
    }

    // 已退订
    if sub.EndAt != nil {
        return model.SubscribeStatusUnSubscribed
    }

    // 寄存中
    if sub.PausedAt != nil {
        return model.SubscribeStatusPaused
    }

    // 未激活
    if sub.StartAt == nil {
        return model.SubscribeStatusInactive
    }

    return model.SubscribeStatusUsing
}

// PausedDays 计算暂停天数
func (s *subscribeService) PausedDays(sub *ent.Subscribe) int {
    pause, _ := ar.Ent.SubscribePause.Query().Where(
        subscribepause.EndAtIsNil(),
        subscribepause.SubscribeID(sub.ID),
    ).Only(s.ctx)
    var dd int
    if pause != nil {
        dd = tools.NewTime().SubDaysToNow(pause.StartAt)
    }
    return int(sub.PauseDays) + dd
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
        Status:    s.GetStatus(sub),
        Voltage:   sub.Voltage,
        Days:      int(sub.Days),
        PauseDays: s.PausedDays(sub),
        AlterDays: int(sub.AlterDays),
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

    for _, pm := range sub.Edges.Plan.Edges.Pms {
        res.Models = append(res.Models, model.BatteryModel{
            ID:       pm.ID,
            Voltage:  pm.Voltage,
            Capacity: pm.Capacity,
        })
    }

    start := sub.StartAt
    if res.Status == model.SubscribeStatusInactive || res.Status == model.SubscribeStatusCanceled {
        // 骑士卡未激活时, 剩余时间等于骑士卡初始时间
        res.Remaining = res.Days
        start = tools.NewPointer().Time(time.Now())
    } else {
        // 骑士卡已激活剩余时间
        res.Remaining = res.Days + res.PauseDays + res.AlterDays - tools.NewTime().SubDaysToNow(*sub.StartAt)
        res.StartAt = start.Format(carbon.DateLayout)
    }

    // 若已退订或已取消
    if res.Status == model.SubscribeStatusUnSubscribed || res.Status == model.SubscribeStatusCanceled {
        res.Remaining = 0
    }

    // 结束日期
    if sub.EndAt == nil {
        res.EndAt = start.AddDate(0, 0, res.Remaining).Format(carbon.DateLayout)
    } else {
        res.EndAt = sub.EndAt.Format(carbon.DateLayout)
    }

    // 已取消(退款)不显示到期日期
    if res.Status == model.SubscribeStatusCanceled {
        res.EndAt = "-"
    }

    // 若剩余天数小于0 则为逾期状态
    if res.Remaining < 0 {
        res.Status = model.SubscribeStatusOverdue
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
