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
    "github.com/auroraride/aurservd/app/task/reminder"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/contract"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/internal/ent/subscribesuspend"
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
        orm: ent.Database.Subscribe,
    }
}

func NewSubscribeWithRider(rider *ent.Rider) *subscribeService {
    s := NewSubscribe()
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

func (s *subscribeService) Query(id uint64) (*ent.Subscribe, error) {
    return s.orm.QueryNotDeleted().Where(subscribe.ID(id)).First(s.ctx)
}

func (s *subscribeService) QueryX(id uint64) *ent.Subscribe {
    sub, _ := s.Query(id)
    if sub == nil {
        snag.Panic("未找到订阅")
    }
    return sub
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
// params 查询实名认证的订阅
func (s *subscribeService) Recent(riderID uint64, params ...uint64) *ent.Subscribe {
    q := s.orm.QueryNotDeleted().
        Where(subscribe.StatusNotIn(model.SubscribeStatusCanceled)).
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
        WithEnterprise()
    if len(params) == 0 {
        q.Where(subscribe.RiderID(riderID))
    } else {
        ids, _ := ent.Database.Rider.Query().Where(rider.PersonID(params[0]), rider.IDNotIn(riderID)).IDs(s.ctx)
        ids = append(ids, riderID)
        q.Where(subscribe.RiderIDIn(ids...))
    }
    sub, _ := q.First(s.ctx)

    return sub
}

func (s *subscribeService) RecentX(riderID uint64) *ent.Subscribe {
    sub := s.Recent(riderID)
    if sub == nil {
        snag.Panic("未找到骑士卡")
    }
    return sub
}

// Detail 获取订阅详情
func (s *subscribeService) Detail(sub *ent.Subscribe) *model.Subscribe {
    if sub == nil {
        return nil
    }

    // 代理商模式计算剩余时间
    remaining := sub.Remaining
    if sub.AgentEndAt != nil {
        remaining = tools.NewTime().LastDays(*sub.AgentEndAt, time.Now())
    }

    res := &model.Subscribe{
        ID:          sub.ID,
        RiderID:     sub.RiderID,
        Status:      sub.Status,
        Model:       sub.Model,
        Days:        sub.InitialDays + sub.PauseDays + sub.AlterDays + sub.OverdueDays,
        PauseDays:   sub.PauseDays,
        AlterDays:   sub.AlterDays,
        InitialDays: sub.InitialDays,
        Remaining:   remaining,
        OverdueDays: sub.OverdueDays,
        Business:    model.SubscribeBusinessable(sub.Status),
        Suspend:     sub.SuspendAt != nil,
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
                ID:    pm.ID,
                Model: pm.Model,
            })
        }
    }

    if sub.StartAt != nil {
        res.StartAt = sub.StartAt.Format(carbon.DateLayout)
    } else {
        res.StartAt = "-"
    }

    if sub.EndAt != nil {
        res.EndAt = sub.EndAt.Format(carbon.DateLayout)
    } else {
        res.EndAt = "-"
    }

    // 已取消(退款)不显示到期日期
    if res.Status == model.SubscribeStatusCanceled {
        res.EndAt = "-"
    }

    // 如果是未激活则查询支付信息
    if res.Status == model.SubscribeStatusInactive {
        o := sub.Edges.InitialOrder
        if o != nil {
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
    }

    e := sub.Edges.Enterprise
    if e != nil {
        res.Enterprise = &model.Enterprise{
            ID:    e.ID,
            Name:  e.Name,
            Agent: e.Agent,
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
    return ent.Database.Subscribe.QueryNotDeleted().
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
    items, _ := ent.Database.Subscribe.QueryNotDeleted().
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
            spq.Order(ent.Desc(subscribepause.FieldStartAt))
        }).
        WithSuspends(func(spq *ent.SubscribeSuspendQuery) {
            spq.Order(ent.Desc(subscribesuspend.FieldStartAt))
        }).
        WithPlan().
        WithRider(func(query *ent.RiderQuery) {
            query.WithPerson()
        }).
        All(s.ctx)
    return items
}

// UpdateStatus 更新订阅状态
// notice 是否发送催费通知
func (s *subscribeService) UpdateStatus(item *ent.Subscribe, notice bool) error {
    if item.EnterpriseID != nil {
        return nil
    }

    tt := tools.NewTime()
    // 已用天数
    pastDays := tt.UsedDaysToNow(*item.StartAt)
    status := item.Status

    // 计算寄存
    pauses, suspends := item.AdditionalItems()

    // 当寄存中生效的时候暂停扣费会产生重复天数, 需要计算并扣除该部分重复的天数
    // 计算暂停计费和寄存之间的重复天数
    pause := ent.SubscribeAdditionalCalculate[*ent.SubscribePause](pauses)

    suspend := ent.SubscribeAdditionalCalculate[*ent.SubscribeSuspend](suspends)

    // 剩余天数
    remaining := item.InitialDays + item.AlterDays + item.OverdueDays + item.RenewalDays + pause.TotalDays + suspend.TotalDays - pastDays

    if remaining < 0 {
        status = model.SubscribeStatusOverdue
    }

    return ent.WithTx(s.ctx, func(tx *ent.Tx) error {
        up := tx.Subscribe.UpdateOne(item).SetPauseDays(pause.TotalDays).SetSuspendDays(suspend.TotalDays).SetRemaining(remaining)

        // 寄存中的如果欠费则自动退租: 超过寄存设置的最大时间继续计费, 直到欠费自动退租
        var unsub bool
        if pause.Current != nil {
            pup := tx.SubscribePause.UpdateOne(pause.Current).SetDays(pause.CurrentDays).SetOverdueDays(pause.CurrentOverdue).SetSuspendDays(pause.CurrentDuplicateDays)
            if remaining < 0 {
                status = model.SubscribeStatusUnSubscribed
                unsub = true
                reason := "寄存超期自动退租"

                pup.SetEndAt(time.Now()).SetRemark(reason).SetPauseOverdue(true)
                up.SetPauseOverdue(true).ClearPausedAt().SetUnsubscribeReason(reason)
                log.Infof("[SUBSCRIBE TASK PAUSE] %d 寄存超期自动退租", item.ID)
            }
            _, err := pup.Save(s.ctx)
            if err != nil {
                log.Errorf("[SUBSCRIBE TASK PAUSE] %d 更新失败: %v", pause.Current.ID, err)
                return err
            }
        }

        if suspend.Current != nil {
            _, err := tx.SubscribeSuspend.UpdateOne(suspend.Current).SetDays(suspend.CurrentDays).Save(s.ctx)
            if err != nil {
                log.Errorf("[SUBSCRIBE TASK SUSPEND] %d 更新失败: %v", suspend.Current.ID, err)
                return err
            }
        }

        // 更新
        sub, err := up.
            SetStatus(status).
            Save(context.Background())
        if err != nil {
            log.Errorf("[SUBSCRIBE TASK] %d 更新失败: %v", item.ID, err)
            return err
        }
        sub.Edges = item.Edges

        *item = *sub
        log.Infof("[SUBSCRIBE TASK] %d 更新成功, 状态: %d, 剩余天数: %d", item.ID, status, remaining)

        if unsub {
            // 标记需要签约
            _, _ = tx.Rider.UpdateOneID(sub.RiderID).SetContractual(false).Save(s.ctx)

            // 查询并标记用户合同为失效
            _, _ = tx.Contract.Update().Where(contract.RiderID(sub.RiderID)).SetEffective(false).Save(s.ctx)
        }

        if notice {
            reminder.Subscribe(item)
        }
        return nil
    })
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
    // 2022-07-07 和博文沟通后把时间限制又给取消了
    // if req.Days+sub.Remaining < 0 {
    //     snag.Panic("不能将剩余时间调整为负值")
    // }

    before := sub.Remaining
    after := sub.Remaining + req.Days
    status := sub.Status

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
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
            snag.Panic("时间修改失败")
        }

        // 更新订阅
        if after > 0 && status == model.SubscribeStatusOverdue {
            status = model.SubscribeStatusUsing
        }
        if after < 0 {
            status = model.SubscribeStatusOverdue
        }

        sub, err = tx.Subscribe.
            UpdateOneID(sub.ID).
            AddAlterDays(req.Days).
            AddRemaining(req.Days).
            SetStatus(status).
            Save(s.ctx)
        if err != nil {
            log.Error(err)
            snag.Panic("时间修改失败")
        }
        return
    })

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
        Model:     sub.Model,
    }
}

// OverdueFee 计算逾期费用
func (s *subscribeService) OverdueFee(riderID uint64, sub *ent.Subscribe) (fee float64, formula string, o *ent.Order) {
    remaining := sub.Remaining
    if remaining > 0 {
        return
    }

    o, _ = NewOrder().RencentSubscribeOrder(riderID)
    var p *ent.Plan
    if o != nil {
        p = o.Edges.Plan
    } else if sub.PlanID != nil {
        p, _ = ent.Database.Plan.Query().Where(plan.ID(*sub.PlanID)).First(s.ctx)
    }
    if p == nil {
        snag.Panic("上次购买骑士卡获取失败")
    }

    fee, formula = p.OverdueFee(remaining)
    return
}

func (s *subscribeService) OverdueFeeFormula(price float64, days uint, remaining int) (fee float64, formula string) {
    fee, _ = decimal.NewFromFloat(price).Div(decimal.NewFromInt(int64(days))).Mul(decimal.NewFromInt(int64(remaining)).Neg()).Mul(decimal.NewFromFloat(1.24)).Float64()
    fee = math.Round(fee*100) / 100

    formula = fmt.Sprintf("(上次购买骑士卡价格 %.2f元 ÷ 天数 %d天) × 逾期天数 %d天 × 1.24 = 逾期费用 %.2f元", price, days, remaining, fee)
    return
}
