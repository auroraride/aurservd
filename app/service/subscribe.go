// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/logging"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "github.com/shopspring/decimal"
    log "github.com/sirupsen/logrus"
    "math"
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

func (s *subscribeService) QueryEdges(id uint64) (*ent.Subscribe, error) {
    return s.orm.QueryNotDeleted().
        Where(subscribe.ID(id)).
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
        WithEnterprise().
        First(s.ctx)
}

func (s *subscribeService) QueryEdgesX(id uint64) *ent.Subscribe {
    item, _ := s.QueryEdges(id)
    if item == nil {
        snag.Panic("未找到有效订阅")
    }
    return item
}

// QueryAndDetailX 根据订阅ID查询订阅以及计算详情
func (s *subscribeService) QueryAndDetailX(id uint64) (sub *ent.Subscribe, detail *model.Subscribe) {
    sub = s.QueryEdgesX(id)
    detail = s.Detail(sub)
    return
}

// Recent 查询骑手最近的订阅
func (s *subscribeService) Recent(riderID uint64) *ent.Subscribe {
    sub, _ := s.orm.QueryNotDeleted().
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
        WithEnterprise().
        First(s.ctx)

    return sub
}

// Detail 获取订阅详情
func (s *subscribeService) Detail(sub *ent.Subscribe) *model.Subscribe {
    if sub == nil {
        return nil
    }

    res := &model.Subscribe{
        ID:          sub.ID,
        RiderID:     sub.RiderID,
        Status:      sub.Status,
        Voltage:     sub.Voltage,
        Days:        sub.InitialDays + sub.PauseDays + sub.AlterDays + sub.OverdueDays,
        PauseDays:   sub.PauseDays,
        AlterDays:   sub.AlterDays,
        InitialDays: sub.InitialDays,
        Remaining:   sub.Remaining,
        OverdueDays: sub.OverdueDays,
        Business:    model.SubscribeBusinessable(sub.Status),
        City: &model.City{
            ID:   sub.Edges.City.ID,
            Name: sub.Edges.City.Name,
        },
    }

    if sub.Edges.Plan != nil {
        p := sub.Edges.Plan
        res.Plan = &model.Plan{
            ID:   p.ID,
            Name: p.Name,
            Days: p.Days,
        }

        for _, pm := range sub.Edges.Plan.Edges.Pms {
            res.Models = append(res.Models, model.BatteryModel{
                ID:       pm.ID,
                Voltage:  pm.Voltage,
                Capacity: pm.Capacity,
            })
        }
    }

    startAt := time.Now()
    if sub.StartAt != nil {
        startAt = *sub.StartAt
        res.StartAt = sub.StartAt.Format(carbon.DateLayout)
    } else {
        res.StartAt = "-"
    }

    if sub.EndAt != nil {
        res.EndAt = sub.EndAt.Format(carbon.DateLayout)
    } else {
        res.EndAt = startAt.AddDate(0, 0, sub.Remaining).Format(carbon.DateLayout)
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

    e := sub.Edges.Enterprise
    if e != nil {
        res.Enterprise = &model.EnterpriseBasic{
            ID:   e.ID,
            Name: e.Name,
        }
        res.Business = e.Status == model.EnterpriseStatusCollaborated
    }

    return res
}

// RecentDetail 获取骑手最近订阅详情
func (s *subscribeService) RecentDetail(riderID uint64) (*model.Subscribe, *ent.Subscribe) {
    sub := s.Recent(riderID)

    if sub == nil {
        return nil, nil
    }

    return s.Detail(sub), sub
}

// QueryEffective 获取骑手当前生效中的订阅
func (s *subscribeService) QueryEffective(riderID uint64) (*ent.Subscribe, error) {
    return ar.Ent.Subscribe.QueryNotDeleted().
        Where(
            subscribe.RiderID(riderID),
            subscribe.StatusIn(
                model.SubscribeStatusInactive,
                model.SubscribeStatusUsing,
                model.SubscribeStatusPaused,
                model.SubscribeStatusOverdue,
            ),
        ).First(s.ctx)
}

// QueryAllRidersEffective 获取所有骑手生效中的订阅
func (s *subscribeService) QueryAllRidersEffective() []*ent.Subscribe {
    items, _ := ar.Ent.Subscribe.Query().
        Where(
            // 未退款
            subscribe.RefundAtIsNil(),
            // 未结束
            subscribe.EndAtIsNil(),
            // 已开始
            subscribe.StartAtNotNil(),
            // 非企业
            subscribe.EnterpriseIDIsNil(),
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
    // 已用天数
    pastDays := tt.UsedDaysToNow(*item.StartAt)
    status := model.SubscribeStatusUsing
    // 寄存中
    if item.PausedAt != nil && item.Edges.Pauses != nil {
        p := item.Edges.Pauses[0]
        // 寄存已过时间需要尽可能的少算
        diff := s.PausedDays(p.StartAt, time.Now())
        pauseDays += diff
    }
    // 剩余天数
    remaining := item.InitialDays + item.AlterDays + item.OverdueDays + item.RenewalDays + pauseDays - pastDays
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

// AlterDays 修改骑手时间
func (s *subscribeService) AlterDays(req *model.SubscribeAlter) (res model.RiderItemSubscribe) {
    sub, _ := s.orm.QueryNotDeleted().Where(subscribe.ID(req.ID)).WithRider().Only(s.ctx)
    // 团签用户禁止修改
    if sub.EnterpriseID != nil {
        snag.Panic("团签用户无法修改")
    }

    u := sub.Edges.Rider
    if sub == nil {
        snag.Panic("订阅不存在")
    }
    if req.Days+sub.Remaining < 0 {
        snag.Panic("不能将剩余时间调整为负值")
    }

    before := sub.Remaining

    tx, err := ar.Ent.Tx(s.ctx)
    if err != nil {
        log.Error(err)
        snag.Panic("时间修改失败")
    }

    // 插入时间修改
    _, err = tx.SubscribeAlter.
        Create().
        SetRiderID(sub.RiderID).
        SetManagerID(s.modifier.ID).
        SetSubscribeID(sub.ID).
        SetDays(req.Days).
        SetRemark(req.Reason).
        Save(s.ctx)
    if err != nil {
        log.Error(err)
        _ = tx.Rollback()
        snag.Panic("时间修改失败")
    }

    // 更新订阅
    sub, err = tx.Subscribe.
        UpdateOneID(sub.ID).
        AddAlterDays(req.Days).
        AddRemaining(req.Days).
        Save(s.ctx)
    if err != nil {
        log.Error(err)
        _ = tx.Rollback()
        snag.Panic("时间修改失败")
    }

    _ = tx.Commit()

    // 记录日志
    go logging.NewOperateLog().
        SetRef(u).
        SetModifier(s.modifier).
        SetOperate(model.OperateSubscribeAlter).
        SetDiff(fmt.Sprintf("剩余%d天", before), fmt.Sprintf("剩余%d天", sub.Remaining)).
        SetRemark(req.Reason).
        Send()

    return model.RiderItemSubscribe{
        ID:        sub.ID,
        Status:    sub.Status,
        Remaining: sub.Remaining,
        Voltage:   sub.Voltage,
    }
}

// PausedDays 计算寄存天数
// 寄存天数 = 结束寄存当天0点 - 寄存当日24点(第二天0点)
func (s *subscribeService) PausedDays(start time.Time, end time.Time) int {
    return int(math.Abs(float64(carbon.Time2Carbon(start).StartOfDay().AddDay().DiffInDays(carbon.Time2Carbon(end).StartOfDay()))))
}

// OverdueFee 计算逾期费用
func (s *subscribeService) OverdueFee(riderID uint64, remaining int) (fee float64, formula string, o *ent.Order) {
    if remaining > 0 {
        return
    }

    o, _ = NewOrder().RencentSubscribeOrder(riderID)
    p := o.Edges.Plan
    if p == nil {
        snag.Panic("上次购买骑士卡获取失败")
    }

    price := p.Price
    days := p.Days
    fee, _ = decimal.NewFromFloat(price).Div(decimal.NewFromInt(int64(days))).Mul(decimal.NewFromInt(int64(remaining)).Neg()).Float64()

    formula = fmt.Sprintf("(上次购买骑士卡价格 %.2f元 ÷ 天数 %d天) × 逾期天数 %d天 = 逾期费用 %.2f元", price, days, remaining, fee)
    return
}
