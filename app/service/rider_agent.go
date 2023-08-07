// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-03
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
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
	isAuthed := NewRider().IsAuthed(item)
	res := model.AgentRider{
		ID:       item.ID,
		Phone:    item.Phone,
		Date:     item.CreatedAt.Format(carbon.DateLayout),
		Name:     item.Name,
		IsAuthed: isAuthed,
	}

	// 已实名认证的显示实名姓名
	if isAuthed && item.Edges.Person != nil {
		res.Name = item.Edges.Person.Name
	}

	// 获取站点
	st := item.Edges.Station
	if st != nil {
		res.Station = &model.EnterpriseStation{
			ID:   st.ID,
			Name: st.Name,
		}
	}

	// 获取电池sn
	bat := item.Edges.Battery
	if bat != nil {
		res.BatterySN = bat.Sn
	}

	// 加入团签时间
	res.JoinEnterpriseAt = item.CreatedAt.Format(carbon.DateTimeLayout)

	// 获取订阅信息
	subs := item.Edges.Subscribes
	if len(subs) > 0 {
		sub := subs[0]
		res.SubscribeID = sub.ID
		res.Model = sub.Model
		res.Intelligent = sub.Intelligent
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
				// 判定当前状态
				if sub.AgentEndAt.After(today) && *res.Remaining > model.WillOverdueNum {
					// 若代理商处到期日期晚于今天, 则是使用中
					res.Status = model.AgentRiderStatusUsing
				} else if *res.Remaining <= model.WillOverdueNum && *res.Remaining >= 0 { // 即将到期暂定3天
					res.Status = model.AgentRiderStatusWillOverdue
				} else {
					// 否则是已逾期
					res.Status = model.AgentRiderStatusOverdue
				}
			}
		}

		// 查找电车信息
		if sub.Edges.Brand != nil {
			res.Ebike = &model.Ebike{Brand: &model.EbikeBrand{
				ID:    sub.Edges.Brand.ID,
				Name:  sub.Edges.Brand.Name,
				Cover: sub.Edges.Brand.Cover,
			}}

			if sub.Edges.Ebike != nil {
				bike := sub.Edges.Ebike
				res.Ebike.ID = bike.ID
				res.Ebike.SN = bike.Sn
				res.Ebike.ExFactory = bike.ExFactory
				res.Ebike.Plate = bike.Plate
				res.Ebike.Color = bike.Color
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
			subscribe.AgentEndAtLTE(carbon.Time2Carbon(tools.NewTime().WillEnd(today, model.WillOverdueNum, true)).EndOfDay().Carbon2Time()),
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
			query.Order(ent.Desc(subscribe.FieldCreatedAt)).WithCity().WithBrand().WithEbike()
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
