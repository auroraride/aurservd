// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"fmt"
	"time"

	"github.com/golang-module/carbon/v2"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/ent/subscribealter"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type agentSubscribeService struct {
	*BaseService
}

func NewAgentSubscribe(params ...any) *agentSubscribeService {
	return &agentSubscribeService{
		BaseService: newService(params...),
	}
}

// AlterList 加时申请列表
func (s *agentSubscribeService) AlterList(enterpriseId uint64, req *model.SubscribeAlterApplyReq) *model.PaginationRes {
	q := ent.Database.SubscribeAlter.QueryNotDeleted().
		Where(
			subscribealter.EnterpriseID(enterpriseId),
			subscribealter.HasRiderWith(rider.DeletedAtIsNil()),
			subscribealter.HasSubscribeWith(subscribe.StatusNotIn(model.SubscribeStatusUnSubscribed)),
		).
		Order(ent.Desc(subscribealter.FieldCreatedAt)).
		WithRider().
		WithSubscribe()

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
		q.Where(subscribealter.HasRiderWith(rider.Or(rider.NameContainsFold(*req.Keyword),
			rider.PhoneContainsFold(*req.Keyword))))
	}
	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.SubscribeAlter) model.SubscribeAlterApplyListRsp {
			rsp := model.SubscribeAlterApplyListRsp{
				ID:   item.ID,
				Days: item.Days,
				// 申请时间
				ApplyTime: item.CreatedAt.Format(carbon.DateTimeLayout),
				// 审批状态
				Status: item.Status,
			}
			if item.ExpireTime != nil {
				rsp.ExpireTime = item.ExpireTime.Format(carbon.DateTimeLayout)
			}
			if item.ReviewTime != nil {
				rsp.ReviewTime = item.ReviewTime.Format(carbon.DateTimeLayout)
			}
			if item.Edges.Rider != nil {
				// 骑手姓名
				rsp.RiderName = item.Edges.Rider.Name
				// 骑手手机号
				rsp.RiderPhone = item.Edges.Rider.Phone
			}
			return rsp
		})
}

// AlterReview 审批加时申请
func (s *agentSubscribeService) AlterReview(req *model.SubscribeAlterReviewReq) {
	// 查找申请记录
	q := ent.Database.SubscribeAlter.QueryNotDeleted().
		Where(
			subscribealter.IDIn(req.Ids...),
			subscribealter.HasRiderWith(rider.DeletedAtIsNil()),
			subscribealter.HasSubscribeWith(subscribe.StatusNotIn(model.SubscribeStatusUnSubscribed)),
			subscribealter.Status(model.SubscribeAlterStatusPending),
		).
		WithRider().
		WithSubscribe()

	alters, _ := q.All(s.ctx)
	if len(alters) == 0 {
		snag.Panic("申请记录不存在")
	}

	for _, v := range alters {
		// 事务
		ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
			// 查询订阅信息
			sub := v.Edges.Subscribe

			err = tx.SubscribeAlter.UpdateOne(v).SetStatus(req.Status).SetReviewTime(time.Now()).Exec(s.ctx)
			if err != nil {
				zap.L().Log(zap.ErrorLevel, "审批加时申请失败", zap.Error(err))
				return
			}
			// 审批不通过不继续
			if req.Status == model.SubscribeAlterStatusRefuse {
				return
			}

			// 加时前剩余天数
			before := tools.NewTime().LastDaysToNow(*sub.AgentEndAt)
			// 加时后剩余天数
			after := before + v.Days

			// 更新订阅时间
			if err = tx.Subscribe.UpdateOne(sub).AddAlterDays(v.Days).
				SetAgentEndAt(tools.NewTime().WillEnd(*sub.AgentEndAt, v.Days, true)).
				Exec(s.ctx); err != nil {
				zap.L().Log(zap.ErrorLevel, "更新订阅时间失败", zap.Error(err))
				return
			}

			// 记录日志
			go logging.NewOperateLog().
				SetRef(v.Edges.Rider).
				SetModifier(s.modifier).
				SetAgent(s.agent).
				SetOperate(model.OperateAgentSubscribeAlter).
				SetDiff(fmt.Sprintf("剩余%d天", before), fmt.Sprintf("剩余%d天", after)).
				Send()

			return
		})
	}
}
