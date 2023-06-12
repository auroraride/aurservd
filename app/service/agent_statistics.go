// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"strconv"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/ent/subscribealter"
	"github.com/auroraride/aurservd/pkg/tools"
)

type agentStatisticsService struct {
	*BaseService
}

func NewAgentStatistics(params ...any) *agentStatisticsService {
	return &agentStatisticsService{
		BaseService: newService(params...),
	}
}

// Overview 代理小程序数据统计概览
func (s *agentStatisticsService) Overview(en *ent.Enterprise) *model.AgentStatisticsOverviewRes {
	start := carbon.Now().StartOfDay().Carbon2Time()
	endtime := tools.NewTime().WillEnd(start, model.WillOverdueNum)
	var v []model.AgentStatisticsOverviewRes
	ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.EnterpriseID(en.ID)).Modify(
		func(s *sql.Selector) {
			// 统计骑手数量
			s.Select(
				// // 新签骑手
				// sql.As("COUNT(CASE WHEN status IN ("+
				// 	strconv.FormatUint(uint64(model.SubscribeStatusUsing), 10)+","+
				// 	strconv.FormatUint(uint64(model.SubscribeStatusUnSubscribed), 10)+
				// 	") THEN rider_id END)", "new_rental_count"),
				// 退签骑手
				sql.As("COUNT(CASE WHEN status="+
					strconv.FormatUint(uint64(model.SubscribeStatusUnSubscribed), 10)+
					" THEN rider_id END)", "quitRiderTotal"),
				// 计费中骑手
				sql.As("COUNT(CASE WHEN status = "+
					strconv.FormatUint(uint64(model.SubscribeStatusUsing), 10)+
					"  THEN rider_id END)", "billingRiderTotal"),
				// 临期
				sql.As("COUNT(CASE WHEN agent_end_at <= '"+endtime.Format(carbon.DateLayout)+"' AND agent_end_at >= '"+start.Format(carbon.DateLayout)+"' THEN rider_id END)", "expiringRiderTotal"),
				// 总数 =计费中+未激活
				sql.As("COUNT(CASE WHEN status IN ("+
					strconv.FormatUint(uint64(model.SubscribeStatusUsing), 10)+","+
					strconv.FormatUint(uint64(model.SubscribeStatusInactive), 10)+
					") THEN rider_id END)", "riderTotal"),
			)
		},
	).ScanX(s.ctx, &v)
	if len(v) == 0 {
		return &model.AgentStatisticsOverviewRes{
			AgentStatisticsOverviewAmount:  model.AgentStatisticsOverviewAmount{},
			AgentStatisticsOverviewRider:   model.AgentStatisticsOverviewRider{},
			AgentStatisticsOverviewBattery: model.AgentStatisticsOverviewBattery{},
		}
	}

	v[0].Balance = en.Balance
	// 新签=临期+计费中
	v[0].NewRiderTotal = v[0].BillingRiderTotal + v[0].ExpiringRiderTotal
	// 骑手加时审核数
	v[0].OverTimeRiderTotal = ent.Database.SubscribeAlter.Query().Where(subscribealter.EnterpriseID(en.ID), subscribealter.StatusEQ(model.SubscribeAlterStatusPending)).CountX(s.ctx)
	return &v[0]
}
