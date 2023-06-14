// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"fmt"
	"time"

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

// AlterReview 审批加时申请
func (s *agentSubscribeService) AlterReview(req *model.SubscribeAlterReviewReq) {
	// 查找申请记录
	q := ent.Database.SubscribeAlter.Query().
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
