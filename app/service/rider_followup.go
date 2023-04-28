// Copyright (C) liasica. 2022-present.
//
// Created at 2022-07-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/riderfollowup"
	"github.com/golang-module/carbon/v2"
)

type riderFollowupService struct {
	ctx      context.Context
	modifier *model.Modifier
	rider    *ent.Rider
	orm      *ent.RiderFollowUpClient
}

func NewRiderFollowup() *riderFollowupService {
	return &riderFollowupService{
		ctx: context.Background(),
		orm: ent.Database.RiderFollowUp,
	}
}

func NewRiderFollowupWithRider(r *ent.Rider) *riderFollowupService {
	s := NewRiderFollowup()
	s.ctx = context.WithValue(s.ctx, "rider", r)
	s.rider = r
	return s
}

func NewRiderFollowupWithModifier(m *model.Modifier) *riderFollowupService {
	s := NewRiderFollowup()
	s.ctx = context.WithValue(s.ctx, "modifier", m)
	s.modifier = m
	return s
}

func (s *riderFollowupService) Create(req *model.RiderFollowUpCreateReq) {
	s.orm.Create().
		SetRiderID(req.RiderID).
		SetRemark(req.Remark).
		SetManagerID(s.modifier.ID).
		SaveX(s.ctx)
}

func (s *riderFollowupService) List(req *model.RiderFollowUpListReq) *model.PaginationRes {
	return model.ParsePaginationResponse(
		s.orm.QueryNotDeleted().WithManager().Where(riderfollowup.RiderID(req.RiderID)),
		req.PaginationReq,
		func(item *ent.RiderFollowUp) model.RiderFollowUpListRes {
			m := item.Edges.Manager
			return model.RiderFollowUpListRes{
				ID: item.ID,
				Manager: model.Modifier{
					ID:    m.ID,
					Name:  m.Name,
					Phone: m.Phone,
				},
				Remark: item.Remark,
				Time:   item.CreatedAt.Format(carbon.DateTimeLayout),
			}
		},
	)
}
