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
    "github.com/golang-module/carbon/v2"
)

// import (
//     "context"
//     "github.com/auroraride/aurservd/app/model"
//     "github.com/auroraride/aurservd/internal/ar"
//     "github.com/auroraride/aurservd/internal/ent"
//     "github.com/auroraride/aurservd/internal/ent/order"
//     "github.com/auroraride/aurservd/pkg/tools"
//     "github.com/golang-module/carbon/v2"
//     "time"
// )
//
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

// List 分页获取骑手订单
func (s *riderOrderService) List(riderID uint64, req *model.PaginationReq) *model.PaginationRes {
    if req.PageSize > 10 {
        req.PageSize = 10
    }
    q := s.orm.QueryNotDeleted().
        WithCity().
        WithPlan(func(pq *ent.PlanQuery) {
            pq.WithPms()
        }).
        Order(ent.Desc(order.FieldCreatedAt)).
        Where(
            order.RiderID(riderID),
            order.TypeIn(model.OrderSubscribeTypes...),
        )
    return model.ParsePaginationResponse[model.RiderOrder, ent.Order](
        q,
        *req,
        s.Detail,
    )
}

// Detail 订单详情
func (s *riderOrderService) Detail(item *ent.Order) model.RiderOrder {
    rc := item.Edges.City
    res := model.RiderOrder{
        ID:         item.ID,
        Type:       item.Type,
        Status:     item.Status,
        Payway:     item.Payway,
        PayAt:      item.CreatedAt.Format(carbon.DateTimeLayout),
        Amount:     item.Amount,
        OutTradeNo: item.OutTradeNo,
        City: model.City{
            ID:   rc.ID,
            Name: rc.Name,
        },
    }
    // 骑士卡订阅订单
    rp := item.Edges.Plan
    if rp != nil {
        res.Plan = &model.Plan{
            ID:   rp.ID,
            Name: rp.Name,
            Days: rp.Days,
        }
        for _, pm := range rp.Edges.Pms {
            res.Models = append(res.Models, model.BatteryModel{
                ID:       pm.ID,
                Voltage:  pm.Voltage,
                Capacity: pm.Capacity,
            })
        }
    }
    return res
}
