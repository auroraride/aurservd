// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"fmt"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/ent/subscribealter"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type subscribeAlterService struct {
	*BaseService
}

func NewSubscribeAlter(params ...any) *subscribeAlterService {
	return &subscribeAlterService{
		BaseService: newService(params...),
	}
}

// AlterDays 申请加时
func (s *subscribeAlterService) AlterDays(r *ent.Rider, req *model.SubscribeAlterRiderReq) {

	// 查询骑手团签
	sub, _ := NewSubscribe().QueryEffective(r.ID)
	if sub == nil {
		snag.Panic("订阅状态异常")
	}

	// 查询骑手申请是否有未审批的
	q := ent.Database.SubscribeAlter.Query().Where(
		subscribealter.RiderID(r.ID),
		subscribealter.Status(model.SubscribeAlterStatusPending),
		subscribealter.SubscribeID(sub.ID),
	)

	exists, _ := q.Exist(s.ctx)
	if exists {
		snag.Panic("您有正在审核中的加时申请,不能重复提交")
	}

	// 增加记录
	_, err := ent.Database.SubscribeAlter.Create().
		SetRiderID(r.ID).
		SetEnterpriseID(*r.EnterpriseID).
		SetSubscribeID(sub.ID).
		SetDays(req.Days).
		SetStatus(model.SubscribeAlterStatusPending).
		SetSubscribeEndAt(carbon.CreateFromStdTime(*sub.AgentEndAt).EndOfDay().ToStdTime()).
		Save(s.ctx)
	if err != nil {
		snag.Panic("申请失败")
	}
}

// List 加时申请列表
func (s *subscribeAlterService) List(req *model.SubscribeAlterListReq) *model.PaginationRes {
	q := ent.Database.SubscribeAlter.Query().
		Where(
			subscribealter.HasRiderWith(rider.DeletedAtIsNil()),
			subscribealter.HasSubscribeWith(subscribe.StatusNotIn(model.SubscribeStatusUnSubscribed)),
		).
		Order(ent.Desc(subscribealter.FieldCreatedAt)).
		WithRider().
		WithSubscribe().
		WithAgent()

	tt := tools.NewTime()
	if req.Start != nil {
		q.Where(subscribealter.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
	}
	if req.End != nil {
		q.Where(subscribealter.CreatedAtLT(tt.ParseNextDateStringX(*req.End)))
	}

	if req.Status != nil {
		q.Where(subscribealter.Status(*req.Status))
	}

	if req.Keyword != nil {
		q.Where(subscribealter.HasRiderWith(
			rider.Or(
				rider.NameContainsFold(*req.Keyword),
				rider.PhoneContainsFold(*req.Keyword),
			),
		))
	}

	if req.RiderID != nil {
		q.Where(subscribealter.RiderID(*req.RiderID))
	}

	if req.EnterpriseID > 0 {
		q.Where(subscribealter.EnterpriseID(req.EnterpriseID))
	}

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.SubscribeAlter) model.SubscribeAlterApplyListRes {
			res := model.SubscribeAlterApplyListRes{
				ID:        item.ID,
				Days:      item.Days,
				ApplyTime: item.CreatedAt.Format(carbon.DateTimeLayout),
				Status:    item.Status,
			}
			if item.SubscribeEndAt != nil {
				res.SubscribeEndAt = item.SubscribeEndAt.Format(carbon.DateLayout)
			}
			if item.ReviewTime != nil {
				res.ReviewTime = item.ReviewTime.Format(carbon.DateTimeLayout)
			}
			if item.Edges.Rider != nil {
				res.Rider = &model.Rider{
					ID:    item.Edges.Rider.ID,
					Phone: item.Edges.Rider.Phone,
					Name:  item.Edges.Rider.Name,
				}
			}
			// 操作人
			if item.Edges.Agent != nil {
				res.Operator = fmt.Sprintf("代理 - %s", item.Edges.Agent.Name)
			}
			if item.Creator != nil {
				res.Operator = fmt.Sprintf("后台 - %s", item.Creator.Name)
			}
			return res
		},
	)
}
