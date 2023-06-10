// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

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
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type subscribeService struct {
	ctx        context.Context
	modifier   *model.Modifier
	rider      *ent.Rider
	orm        *ent.SubscribeClient
	agent      *ent.Agent
	enterprise *ent.Enterprise
}

func NewSubscribe() *subscribeService {
	return &subscribeService{
		ctx: context.Background(),
		orm: ent.Database.Subscribe,
	}
}

func NewSubscribeWithRider(rider *ent.Rider) *subscribeService {
	s := NewSubscribe()
	s.ctx = context.WithValue(s.ctx, model.CtxRiderKey{}, rider)
	s.rider = rider
	return s
}

func NewSubscribeWithModifier(m *model.Modifier) *subscribeService {
	s := NewSubscribe()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}
func NewSubscribeWithAgent(ag *ent.Agent, en *ent.Enterprise) *subscribeService {
	s := NewSubscribe()
	s.agent = ag
	s.enterprise = en
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
		WithPlan().
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
		WithPlan().
		WithCity().
		// 查询初始订单和对应的押金订单
		WithInitialOrder(func(ioq *ent.OrderQuery) {
			ioq.WithChildren(func(doq *ent.OrderQuery) {
				doq.Where(order.Type(model.OrderTypeDeposit))
			})
		}).
		WithEnterprise().
		WithBrand().
		WithEbike()
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
		sub.EndAt = sub.AgentEndAt
	}

	res := &model.Subscribe{
		ID:           sub.ID,
		RiderID:      sub.RiderID,
		Status:       sub.Status,
		Model:        sub.Model,
		Days:         sub.InitialDays + sub.PauseDays + sub.AlterDays + sub.OverdueDays,
		PauseDays:    sub.PauseDays,
		AlterDays:    sub.AlterDays,
		InitialDays:  sub.InitialDays,
		Remaining:    remaining,
		OverdueDays:  sub.OverdueDays,
		Business:     model.SubscribeBusinessable(sub.Status),
		Suspend:      sub.SuspendAt != nil,
		NeedContract: sub.NeedContract,
		Intelligent:  sub.Intelligent,
		City: &model.City{
			ID:   sub.Edges.City.ID,
			Name: sub.Edges.City.Name,
		},
		Ebike: NewEbike().Detail(sub.Edges.Ebike, sub.Edges.Brand),
	}

	if sub.Edges.Plan != nil {
		p := sub.Edges.Plan
		res.Plan = p.BasicInfo()

		res.Models = []model.BatteryModel{
			{Model: p.Model},
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
func (s *subscribeService) QueryEffective(riderID uint64, edges ...ent.SubscribeQueryWith) (*ent.Subscribe, error) {
	return ent.Database.Subscribe.QueryNotDeleted().
		Where(
			subscribe.RiderID(riderID),
			subscribe.StatusIn(
				model.SubscribeStatusInactive,
				model.SubscribeStatusUsing,
				model.SubscribeStatusPaused,
				model.SubscribeStatusOverdue,
			),
		).
		With(edges...).
		First(s.ctx)
}

func (s *subscribeService) QueryEffectiveX(riderID uint64, edges ...ent.SubscribeQueryWith) *ent.Subscribe {
	sub, _ := s.QueryEffective(riderID, edges...)
	if sub == nil {
		snag.Panic("未找到生效中的订阅")
	}
	return sub
}

func (s *subscribeService) QueryEffectiveIntelligentX(riderID uint64, edges ...ent.SubscribeQueryWith) *ent.Subscribe {
	sub, _ := s.QueryEffective(riderID, edges...)
	if sub == nil {
		snag.Panic("未找到生效中的订阅")
	}
	if !sub.Intelligent {
		snag.Panic("骑手当前为非智能柜订阅")
	}
	return sub
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
		WithRider().
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

	// 计算寄存和暂停
	pr, sr := item.GetAdditionalItems()

	// 剩余天数
	remaining := item.InitialDays + item.AlterDays + item.OverdueDays + item.RenewalDays + pr.Days + sr.Days - pastDays
	formula := fmt.Sprintf(
		"剩余时间(%d) = 初始天数(%d) + 调整天数(%d) + 已缴滞纳金天数(%d) + 续费天数(%d) + 寄存天数(%d) + 暂停天数(%d) - 已过天数(%d)",
		remaining,
		item.InitialDays,
		item.AlterDays,
		item.OverdueDays,
		item.RenewalDays,
		pr.Days,
		sr.Days,
		pastDays,
	)

	if remaining < 0 {
		status = model.SubscribeStatusOverdue
	}

	return ent.WithTx(s.ctx, func(tx *ent.Tx) error {
		up := tx.Subscribe.UpdateOne(item).SetPauseDays(pr.Days).SetSuspendDays(sr.Days).SetRemaining(remaining).SetFormula(formula)

		// 寄存中的如果欠费则自动退租: 超过寄存设置的最大时间继续计费, 直到欠费自动退租
		var unsub bool
		if pr.Current != nil {
			pup := tx.SubscribePause.UpdateOne(pr.Current).SetDays(pr.CurrentDays).SetOverdueDays(pr.CurrentOverdueDays).SetSuspendDays(pr.CurrentDuplicateDays)
			// 查询当前寄存是否超期
			if remaining < 0 && pr.CurrentOverdueDays > 0 {
				status = model.SubscribeStatusUnSubscribed
				unsub = true
				reason := "寄存超期自动退租"

				pup.SetEndAt(time.Now()).SetRemark(reason).SetPauseOverdue(true)
				up.SetPauseOverdue(true).ClearPausedAt().SetUnsubscribeReason(reason).SetEndAt(time.Now())
				zap.L().Info("寄存超期自动退租: " + strconv.FormatUint(item.ID, 10))
			}
			_, err := pup.Save(s.ctx)
			if err != nil {
				zap.L().Error("寄存更新失败: "+strconv.FormatUint(pr.Current.ID, 10), zap.Error(err))
				return err
			}
		}

		if sr.Current != nil {
			_, err := tx.SubscribeSuspend.UpdateOne(sr.Current).SetDays(sr.CurrentDays).Save(s.ctx)
			if err != nil {
				zap.L().Error("暂停更新失败: "+strconv.FormatUint(sr.Current.ID, 10), zap.Error(err))
				return err
			}
		}

		// 更新
		sub, err := up.SetStatus(status).Save(context.Background())
		if err != nil {
			zap.L().Error("订阅更新失败: "+strconv.FormatUint(sub.ID, 10), zap.Error(err))
			return err
		}
		sub.Edges = item.Edges

		*item = *sub
		zap.L().Info("订阅更新成功: " + strconv.FormatUint(sub.ID, 10) +
			", 状态: " + model.SubscribeStatusText(status) +
			", 剩余天数: " + strconv.Itoa(remaining))

		if unsub {
			// 标记需要签约
			_, _ = tx.Rider.UpdateOneID(sub.RiderID).Save(s.ctx)

			// 查询并标记用户合同为失效
			_, _ = tx.Contract.Update().Where(contract.RiderID(sub.RiderID)).SetEffective(false).Save(s.ctx)
		}

		if !notice {
			return nil
		}

		// 提醒
		var (
			fee        *float64
			feeFormula *string
		)
		if sub.Remaining < 0 {
			f, fl := s.OverdueFee(sub)
			feeFormula = silk.Pointer(fl)
			fee = silk.Pointer(f)
		}
		reminder.Subscribe(item, fee, feeFormula)

		return nil
	})
}

// AlterDays 修改骑手时间
func (s *subscribeService) AlterDays(req *model.SubscribeAlter) (res model.RiderItemSubscribe) {
	sq := s.orm.QueryNotDeleted().WithEnterprise().WithRider().Where(subscribe.ID(req.ID), subscribe.StatusIn(model.SubscribeNotUnSubscribed()...))
	if s.agent != nil && s.enterprise != nil {
		sq.Where(subscribe.EnterpriseID(s.enterprise.ID))
	}

	sub, _ := sq.First(s.ctx)

	if sub == nil {
		snag.Panic("订阅不存在")
	}

	u := sub.Edges.Rider

	// 团签(代理)骑手禁止修改
	if sub.Edges.Enterprise != nil {
		snag.Panic("团签用户无法修改")
	}

	// 剩余天数
	before := sub.Remaining

	// if se != nil && se.Agent {
	// 	if sub.AgentEndAt == nil {
	// 		snag.Panic("骑手订阅异常")
	// 	}
	// 	// 计算剩余天数
	// 	before = tools.NewTime().LastDaysToNow(*sub.AgentEndAt)
	// }
	after := before + req.Days
	status := sub.Status

	// 2022-07-07 和博文沟通后把时间限制又给取消了
	// if req.Days+sub.Remaining < 0 {
	//     snag.Panic("不能将剩余时间调整为负值")
	// }

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		// 插入时间修改
		tsa := tx.SubscribeAlter.
			Create().
			SetRiderID(sub.RiderID).
			SetSubscribeID(sub.ID).
			SetDays(req.Days).
			SetRemark(req.Reason).
			SetStatus(model.SubscribeAlterStatusAgree)
		if s.agent != nil {
			tsa.SetAgentID(s.agent.ID).SetEnterpriseID(s.agent.EnterpriseID)
		}
		if s.modifier != nil {
			tsa.SetManagerID(s.modifier.ID)
		}
		_, err = tsa.Save(s.ctx)
		if err != nil {
			snag.Panic("时间修改失败")
		}

		// 更新订阅
		if after > 0 && status == model.SubscribeStatusOverdue {
			status = model.SubscribeStatusUsing
		}
		if after < 0 {
			status = model.SubscribeStatusOverdue
		}

		ts := tx.Subscribe.
			UpdateOneID(sub.ID).
			AddAlterDays(req.Days).
			SetStatus(status)

		// 计算代理商处到期日期
		// if se != nil && se.Agent {
		// 	ts.SetAgentEndAt(tools.NewTime().WillEnd(*sub.AgentEndAt, req.Days, true))
		// } else {
		// 	ts.AddRemaining(req.Days)
		// }
		ts.AddRemaining(req.Days)

		sub, err = ts.Save(s.ctx)
		if err != nil {
			snag.Panic("时间修改失败")
		}
		return
	})

	// 记录日志
	go logging.NewOperateLog().
		SetRef(u).
		SetModifier(s.modifier).
		SetAgent(s.agent).
		SetOperate(model.OperateSubscribeAlter).
		SetDiff(fmt.Sprintf("剩余%d天", before), fmt.Sprintf("剩余%d天", after)).
		SetRemark(req.Reason).
		Send()

	out := model.RiderItemSubscribe{
		ID:        sub.ID,
		Status:    sub.Status,
		Remaining: after,
		Model:     sub.Model,
		Suspend:   sub.SuspendAt != nil,
	}

	if sub.AgentEndAt != nil {
		out.AgentEndAt = sub.AgentEndAt.Format(carbon.DateLayout)
	}

	if sub.EnterpriseID == nil {
		go func() {
			_ = s.UpdateStatus(sub, false)
		}()
	}
	return out
}

func (s *subscribeService) OverdueFee(sub *ent.Subscribe) (fee float64, formula string) {
	if sub.Remaining > 0 {
		return
	}

	remaining := -sub.Remaining

	_, dr := NewSetting().DailyRent(nil, sub.CityID, sub.Model, sub.BrandID)
	fee, _ = decimal.NewFromFloat(dr).Mul(decimal.NewFromInt(int64(remaining))).Mul(decimal.NewFromFloat(1.24)).Float64()
	fee = math.Round(fee*100) / 100

	formula = fmt.Sprintf("该骑士卡日租价格 %.2f元 × 逾期天数 %d天 × 1.24 = 逾期费用 %.2f元", dr, remaining, fee)
	return
}

// CalculateOverdueFee 计算逾期费用
func (s *subscribeService) CalculateOverdueFee(sub *ent.Subscribe) (fee float64, formula string, o *ent.Order) {
	if sub.Remaining > 0 {
		return
	}

	o, _ = NewOrder().RencentSubscribeOrder(sub.RiderID)
	var p *ent.Plan
	if o != nil {
		p = o.Edges.Plan
	} else if sub.PlanID != nil {
		p, _ = ent.Database.Plan.Query().Where(plan.ID(*sub.PlanID)).First(s.ctx)
	}
	if p == nil {
		snag.Panic("上次购买骑士卡获取失败")
	}

	fee, formula = s.OverdueFee(sub)
	return
}

// NeedContract 查询订阅是否需要签约
func (s *subscribeService) NeedContract(sub *ent.Subscribe) bool {
	if !sub.NeedContract {
		return false
	}
	exists, _ := ent.Database.Contract.Query().Where(
		contract.Status(model.ContractStatusSuccess.Value()),
		contract.Effective(true),
		contract.SubscribeID(sub.ID),
	).Exist(s.ctx)
	return !exists
}

func (s *subscribeService) Signed(riderID, subscribeID uint64) (res model.SubscribeSigned) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	start := time.Now()
	for ; true; <-ticker.C {

		c, _ := ent.Database.Contract.Query().Where(
			contract.RiderID(riderID),
			contract.SubscribeID(subscribeID),
		).First(s.ctx)

		// 未找到签约信息
		if c == nil {
			return
		}

		switch model.ContractStatus(c.Status) {
		case model.ContractStatusSigning:
			res.Signed = 1
		case model.ContractStatusSuccess:
			res.Signed = 2
			return
		default:
			res.Signed = 0
			return
		}

		if time.Since(start).Seconds() > 50 {
			return
		}
	}

	return
}
