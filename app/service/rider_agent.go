// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-03
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"time"

	"github.com/golang-module/carbon/v2"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/battery"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/ent/subscribealter"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type riderAgentService struct {
	ctx        context.Context
	modifier   *model.Modifier
	agent      *ent.Agent
	enterprise *ent.Enterprise
	orm        *ent.RiderClient
}

func NewRiderAgent() *riderAgentService {
	return &riderAgentService{
		ctx: context.Background(),
		orm: ent.Database.Rider,
	}
}

func NewRiderAgentWithAgent(ag *ent.Agent, en *ent.Enterprise) *riderAgentService {
	s := NewRiderAgent()
	s.agent = ag
	s.enterprise = en
	return s
}

func NewRiderAgentWithModifier(m *model.Modifier) *riderAgentService {
	s := NewRiderAgent()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

func (s *riderAgentService) detail(item *ent.Rider) model.AgentRider {
	today := carbon.Now().StartOfDay().Carbon2Time()
	res := model.AgentRider{
		ID:       item.ID,
		Phone:    item.Phone,
		Date:     item.CreatedAt.Format(carbon.DateLayout),
		Name:     item.Name,
		IsAuthed: NewRider().IsAuthed(item),
	}
	// 获取站点
	st := item.Edges.Station
	if st != nil {
		res.Station = st.Name
	}

	// 获取电池sn
	bat := item.Edges.Battery
	if bat != nil {
		res.BatterySN = bat.Sn
	}

	// 获取订阅信息
	subs := item.Edges.Subscribes
	if len(subs) > 0 {
		sub := subs[0]
		res.SubscribeID = sub.ID
		res.Model = sub.Model
		ci := sub.Edges.City
		res.City = &model.City{
			ID:   ci.ID,
			Name: ci.Name,
		}
		if sub.EndAt != nil {
			res.EndAt = sub.EndAt.Format(carbon.DateLayout)
		}
		if sub.AgentEndAt != nil {
			res.StopAt = sub.AgentEndAt.Format(carbon.DateLayout)
		}
		if sub.StartAt != nil {
			res.StartAt = sub.StartAt.Format(carbon.DateLayout)
			res.Used = tools.NewTime().UsedDaysToNow(*sub.StartAt)
		}
		switch sub.Status {
		case model.SubscribeStatusInactive:
			// 未激活
			res.Status = model.AgentRiderStatusInactive
		case model.SubscribeStatusUnSubscribed:
			// 已退租
			res.Status = model.AgentRiderStatusUnsubscribed
		case model.SubscribeStatusUsing:
			res.Status = model.AgentRiderStatusUsing
			// 计算剩余日期
			if sub.AgentEndAt != nil {
				res.Remaining = silk.Pointer(tools.NewTime().LastDays(*sub.AgentEndAt, today))
			}
		}
	} else {
		res.Status = model.AgentRiderStatusInactive
	}
	return res
}

// List 骑手列表查询
func (s *riderAgentService) List(enterpriseID uint64, req *model.AgentRiderListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().
		Where(rider.EnterpriseID(enterpriseID)).
		WithSubscribes(func(sq *ent.SubscribeQuery) {
			sq.Order(ent.Desc(subscribe.FieldCreatedAt)).WithCity()
		}).
		WithStation().
		WithBattery().
		WithPerson().
		Order(ent.Desc(rider.FieldCreatedAt))

	today := carbon.Now().StartOfDay().Carbon2Time()

	var subquery []predicate.Subscribe

	if req.CityID != 0 {
		subquery = append(subquery, subscribe.CityID(req.CityID))
	}

	if req.Keyword != "" {
		q.Where(
			rider.Or(
				rider.NameContainsFold(req.Keyword),
				rider.PhoneContainsFold(req.Keyword),
			),
		)
	}

	if req.StationID != 0 {
		q.Where(rider.StationID(req.StationID))
	}

	switch req.Status {
	case model.AgentRiderStatusInactive:
		// 未激活 也要查询未使用的
		subquery = append(subquery, subscribe.Status(model.SubscribeStatusInactive))
		q.Where(rider.Or(rider.HasSubscribesWith(subquery...), rider.Not(rider.HasSubscribes())))
	case model.AgentRiderStatusUsing:
		// 使用中
		subquery = append(subquery, subscribe.Status(model.SubscribeStatusUsing))
		q.Where(rider.HasSubscribesWith(subquery...))
	case model.AgentRiderStatusOverdue:
		// 已逾期
		// 代理商团签的逾期状态 = 使用中 并且 代理商处到期日期已过期
		subquery = append(
			subquery,
			subscribe.Status(model.SubscribeStatusUsing),
			subscribe.EndAtIsNil(),
			subscribe.AgentEndAtLT(today),
		)
		q.Where(rider.HasSubscribesWith(subquery...))
	case model.AgentRiderStatusUnsubscribed:
		// 已退订
		subquery = append(
			subquery,
			subscribe.EndAtNotNil(),
			subscribe.Status(model.SubscribeStatusUnSubscribed),
		)
		q.Where(rider.And(
			rider.HasSubscribesWith(subquery...),
			rider.Not(rider.HasSubscribesWith(subscribe.StatusIn(model.SubscribeNotUnSubscribed()...))),
		))
	case model.AgentRiderStatusWillOverdue:
		// 将逾期
		subquery = append(
			subquery,
			subscribe.Status(model.SubscribeStatusUsing),
			subscribe.AgentEndAtGTE(today),
			subscribe.AgentEndAtLTE(tools.NewTime().WillEnd(today, model.WillOverdueNum)),
		)
		q.Where(rider.HasSubscribesWith(subquery...))
	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Rider) model.AgentRider {
		return s.detail(item)
	})
}

func (s *riderAgentService) Detail(req *model.IDParamReq, enterpriseID uint64) model.AgentRider {
	item, _ := s.orm.QueryNotDeleted().
		Where(rider.EnterpriseID(enterpriseID), rider.ID(req.ID)).
		WithSubscribes(func(query *ent.SubscribeQuery) {
			query.Order(ent.Desc(subscribe.FieldCreatedAt)).WithCity()
		}).
		WithStation().
		WithBattery().
		WithPerson().
		First(s.ctx)
	if item == nil {
		snag.Panic("未找到骑手")
	}
	return s.detail(item)
}

func (s *riderAgentService) Log(req *model.AgentRiderLogReq) (items []model.AgentRiderLog) {
	logs, _ := ent.Database.SubscribeAlter.
		QueryNotDeleted().
		Where(subscribealter.RiderID(req.ID), subscribealter.EnterpriseID(req.EnterpriseID)).
		WithAgent().
		Order(ent.Desc(subscribealter.FieldCreatedAt)).
		All(s.ctx)
	items = make([]model.AgentRiderLog, len(logs))
	for i, log := range logs {
		items[i] = model.AgentRiderLog{
			Days: log.Days,
			Time: log.CreatedAt.Format(carbon.DateTimeLayout),
		}
		ag := log.Edges.Agent
		if ag != nil {
			items[i].Name = ag.Name
		} else {
			items[i].Name = "平台"
		}
	}
	return
}

// Active 激活骑手电池
func (s *enterpriseService) Active(req *model.AgentSubscribeActiveReq, enterpriseID uint64) {
	// 不是自己骑手的不能激活
	riderInfo, _ := ent.Database.Rider.
		QueryNotDeleted().
		Where(
			rider.ID(req.ID),
			rider.EnterpriseID(enterpriseID),
		).
		WithPerson().
		First(s.ctx)
	if riderInfo == nil {
		snag.Panic("未找到有效骑手")
	}

	if !NewRider().IsAuthed(riderInfo) {
		snag.Panic("骑手未实名认证")
	}

	// 查询电池是否存在
	batteryInfo := NewBattery(s.ctx).QueryIDX(req.BatteryID)
	if batteryInfo.CabinetID != nil {
		snag.Panic("电柜中的电池无法手动绑定骑手")
	}

	if batteryInfo.EnterpriseID == nil || batteryInfo.StationID == nil {
		snag.Panic("电池未绑定企业或者站点")
	}

	if *batteryInfo.EnterpriseID != enterpriseID || *batteryInfo.StationID != *riderInfo.StationID {
		snag.Panic("电池不属于当前企业或者骑手所在站点")
	}

	if batteryInfo.RiderID != nil || batteryInfo.SubscribeID != nil {
		snag.Panic("电池已被绑定")
	}
	// 查询骑手订阅信息
	sub, _ := NewSubscribe().QueryEffective(req.ID)
	if sub != nil && (sub.EnterpriseID == nil || (sub.EnterpriseID != nil && sub.Status != model.SubscribeStatusInactive)) {
		snag.Panic("已有正在使用订单,无法绑定电池")
	}
	subscribeInfo, _ := ent.Database.Subscribe.
		QueryNotDeleted().
		Where(subscribe.RiderID(req.ID), subscribe.Status(model.SubscribeStatusInactive)).
		First(s.ctx)
	if subscribeInfo == nil {
		snag.Panic("未找到未激活订单")
	}

	if batteryInfo.Model != sub.Model {
		snag.Panic("激活电池型号与订阅电池型号不符")
	}

	// 绑定电池
	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) error {
		// 激活订阅
		aendTime := tools.NewTime().WillEnd(time.Now(), subscribeInfo.InitialDays)
		if err := tx.Subscribe.UpdateOneID(subscribeInfo.ID).SetStatus(model.SubscribeStatusUsing).SetStartAt(time.Now()).
			SetNillableAgentEndAt(&aendTime).Exec(s.ctx); err != nil {
			zap.L().Error("激活订阅失败", zap.Error(err))
			snag.Panic("激活订阅失败")
		}
		// 删除原有信息
		if err := tx.Battery.Update().Where(battery.SubscribeID(subscribeInfo.ID)).ClearRiderID().ClearSubscribeID().Exec(s.ctx); err != nil {
			zap.L().Error("电池解绑失败", zap.Error(err))
			snag.Panic("电池解绑失败")
		}
		// 绑定电池
		if err := tx.Battery.Update().Where(battery.Sn(batteryInfo.Sn)).
			SetRiderID(req.ID).SetSubscribeID(subscribeInfo.ID).Exec(s.ctx); err != nil {
			zap.L().Error("电池绑定失败", zap.Error(err))
			snag.Panic("电池绑定失败")
		}
		// 更新电池流水记录表
		if err := tx.BatteryFlow.Create().SetSn(batteryInfo.Sn).SetRiderID(req.ID).
			SetBatteryID(batteryInfo.ID).
			SetNillableSubscribeID(&subscribeInfo.ID).Exec(s.ctx); err != nil {
			zap.L().Error("电池流水记录失败", zap.Error(err))
			snag.Panic("电池流水记录失败")
		}

		return nil
	})
}

func (s *enterpriseService) ExitEnterprise(r *ent.Rider) {
	// 查询订阅
	sub, _ := ent.Database.Subscribe.QueryNotDeleted().
		Where(subscribe.RiderID(r.ID),
			subscribe.StatusIn(model.SubscribeStatusInactive, model.SubscribeStatusUnSubscribed),
		).First(s.ctx)
	if sub == nil {
		snag.Panic("未找到订阅")
	}
	// 如果是未激活删除订阅
	if sub.Status == model.SubscribeStatusUnSubscribed {
		_, err := ent.Database.Subscribe.SoftDeleteOne(sub).Save(s.ctx)
		if err != nil {
			snag.Panic("订阅信息删除失败")
		}
	}
	//
	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) error {
		// 删除之前骑手信息
		_, err := tx.Rider.SoftDeleteOne(r).SetRemark("骑手退出团签").Save(s.ctx)
		if err != nil {
			return err
		}
		// 新增骑手信息
		_, err = tx.Rider.Create().SetPhone(r.Phone).Save(s.ctx)
		if err != nil {
			return err
		}
		return nil
	})
	// 清除团签信息
	// err := ent.Database.Rider.UpdateOne(r).ClearEnterpriseID().ClearStationID().Exec(s.ctx)
	// if err != nil {
	// 	snag.Panic("骑手团签信息清除失败")
	// }
}
