// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-03
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/predicate"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/internal/ent/subscribealter"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
)

type riderAgentService struct {
    ctx        context.Context
    modifier   *model.Modifier
    rider      *ent.Rider
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
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *riderAgentService) detail(item *ent.Rider) model.AgentRider {
    today := carbon.Now().StartOfDay().Carbon2Time()
    res := model.AgentRider{
        ID:    item.ID,
        Phone: item.Phone,
        Date:  item.CreatedAt.Format(carbon.DateLayout),
    }
    // 获取站点
    st := item.Edges.Station
    if st != nil {
        res.Station = st.Name
    }
    // 获取实名信息
    p := item.Edges.Person
    if p != nil {
        res.Name = p.Name
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
            // 计算剩余日期
            if sub.AgentEndAt != nil {
                res.Remaining = tools.Pointer(tools.NewTime().LastDays(*sub.AgentEndAt, today))
                // 判定当前状态
                if sub.AgentEndAt.After(today) {
                    // 若代理商处到期日期晚于今天, 则是使用中
                    res.Status = model.AgentRiderStatusUsing
                } else {
                    // 否则是已超期
                    res.Status = model.AgentRiderStatusOverdue
                }
            }
        }
    }
    return res
}

// List 骑手列表查询
func (s *riderAgentService) List(enterpriseID uint64, req *model.AgentRiderListReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().
        Where(rider.EnterpriseID(enterpriseID)).
        WithSubscribes(func(query *ent.SubscribeQuery) {
            query.Order(ent.Desc(subscribe.FieldCreatedAt)).WithCity()
        }).
        WithStation().
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
                rider.HasPersonWith(
                    person.NameContainsFold(req.Keyword),
                ),
                rider.PhoneContainsFold(req.Keyword),
            ),
        )
    }

    if req.StationID != 0 {
        q.Where(rider.StationID(req.StationID))
    }

    switch req.Status {
    case model.AgentRiderStatusInactive:
        // 未激活
        subquery = append(subquery, subscribe.Status(model.SubscribeStatusInactive))
        q.Where(rider.HasSubscribesWith(subquery...))
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
    }

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Rider) model.AgentRider {
        return s.detail(item)
    })
}

func (s *riderAgentService) Detail(req *model.IDParamReq, enterpriseID uint64) model.AgentRiderDetail {
    item, _ := s.orm.QueryNotDeleted().
        Where(rider.EnterpriseID(enterpriseID), rider.ID(req.ID)).
        WithSubscribes(func(query *ent.SubscribeQuery) {
            query.Order(ent.Desc(subscribe.FieldCreatedAt)).WithCity()
        }).
        WithStation().
        WithPerson().
        First(s.ctx)
    if item == nil {
        snag.Panic("未找到骑手")
    }
    return model.AgentRiderDetail{
        AgentRider: s.detail(item),
        Logs: s.Log(&model.AgentRiderLogReq{
            ID:           req.ID,
            EnterpriseID: enterpriseID,
        }),
    }
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