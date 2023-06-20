// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-27
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/order"
	"github.com/auroraride/aurservd/pkg/snag"
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
		orm: ent.Database.Order,
	}
}

func NewRiderOrderWithRider(rider *ent.Rider) *riderOrderService {
	s := NewRiderOrder()
	s.ctx = context.WithValue(s.ctx, model.CtxRiderKey{}, rider)
	s.rider = rider
	return s
}

func NewRiderOrderWithModifier(m *model.Modifier) *riderOrderService {
	s := NewRiderOrder()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

// List 分页获取骑手订单
func (s *riderOrderService) List(riderID uint64, req model.PaginationReq) *model.PaginationRes {
	if req.PageSize > 10 {
		req.PageSize = 10
	}
	q := s.orm.QueryNotDeleted().
		WithCity().
		WithPlan().
		Order(ent.Desc(order.FieldCreatedAt)).
		Where(
			order.RiderID(riderID),
			order.TypeIn(model.OrderSubscribeTypes...),
		).
		WithRider().
		WithSubscribe(func(sq *ent.SubscribeQuery) {
			sq.WithEmployee().WithStore()
		})
	return model.ParsePaginationResponse[model.Order, ent.Order](
		q,
		req,
		NewOrder().Detail,
	)
}

// Query 查询订单
func (s *riderOrderService) Query(riderID, orderID uint64) *ent.Order {
	item, _ := s.orm.QueryNotDeleted().
		Where(order.RiderID(riderID), order.ID(orderID)).
		WithCity().
		WithPlan().
		WithRider().
		WithSubscribe(func(sq *ent.SubscribeQuery) {
			sq.WithEmployee().WithStore()
		}).
		WithRefund().
		WithCoupons().
		First(s.ctx)
	if item == nil {
		snag.Panic("未找到订单")
	}
	return item
}
