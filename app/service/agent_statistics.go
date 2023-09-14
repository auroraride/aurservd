// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/rider"
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
func (s *agentStatisticsService) Overview(ac *app.AgentContext) *model.AgentStatisticsOverviewRes {
	var rsp model.AgentStatisticsOverviewRes
	rsp.AgentStatisticsOverviewRider = s.RiderSummary(ac.Enterprise)
	// 加时审核数
	rsp.SubscribeAlterTotal = ent.Database.SubscribeAlter.Query().Where(
		subscribealter.EnterpriseID(ac.Enterprise.ID),
		subscribealter.StatusEQ(model.SubscribeAlterStatusPending),
		subscribealter.HasRiderWith(rider.DeletedAtIsNil()),
		subscribealter.HasSubscribeWith(subscribe.StatusNEQ(model.SubscribeStatusUnSubscribed)),
	).CountX(s.ctx)

	ns := NewStockSummary()
	// 电池数据统计
	rsp.BatterySummary = ns.BatterySummary(ac)
	// 车+电骑手数
	rsp.EbikeSummary = ns.EbikeSummary(ac)
	return &rsp
}

// RiderSummary 骑手数据统计
func (s *agentStatisticsService) RiderSummary(en *ent.Enterprise) model.AgentStatisticsOverviewRider {
	var v []model.AgentStatisticsOverviewRider
	start := carbon.Now().StartOfDay().ToStdTime()
	endtime := tools.NewTime().WillEnd(start, 4, true)
	ent.Database.Rider.Query().Where(rider.EnterpriseID(en.ID)).Modify(
		func(s *sql.Selector) {
			t := sql.Table(subscribe.Table)
			s.LeftJoin(t).
				On(s.C(rider.FieldID), t.C(subscribe.FieldRiderID)).
				Select(
					// 骑手总数
					sql.As("COUNT(DISTINCT CASE WHEN rider.deleted_at IS NULL THEN rider.id END)", "riderTotal"),
					// 计费中骑手
					sql.As(fmt.Sprintf("COUNT(DISTINCT CASE WHEN rider.deleted_at IS NULL AND t1.status = %d THEN rider.id END)", model.SubscribeStatusUsing), "billingRiderTotal"),
					// 临期骑手
					sql.As(fmt.Sprintf("COUNT(DISTINCT CASE WHEN rider.deleted_at IS NULL AND t1.agent_end_at < '%s' AND t1.agent_end_at >= '%s' AND t1.status = %d THEN rider.id END)",
						endtime.Format(carbon.DateLayout), start.Format(carbon.DateLayout), model.SubscribeStatusUsing), "expiringRiderTotal"),
					// 累计退订骑手
					sql.As(fmt.Sprintf("COUNT(DISTINCT CASE WHEN t1.status = %d AND t1.enterprise_id = %d THEN t1.rider_id END)", model.SubscribeStatusUnSubscribed, en.ID), "unSubscribeTotal"),
					// 单电骑手
					sql.As(fmt.Sprintf("COUNT(DISTINCT CASE WHEN rider.deleted_at IS NULL AND t1.enterprise_id = %d AND t1.brand_id IS NULL THEN t1.rider_id END)", en.ID), "riderOnlyBatteryTotal"),
					// 车+电骑手
					sql.As(fmt.Sprintf("COUNT(DISTINCT CASE WHEN rider.deleted_at IS NULL AND t1.enterprise_id = %d AND t1.brand_id IS NOT NULL THEN t1.rider_id END)", en.ID), "riderBatteryAndEbikeTotal"),
				)
		},
	).ScanX(s.ctx, &v)
	rsp := v[0]
	// 累计激活骑手 = 计费中骑手 + 累计退订骑手
	rsp.ActivationTotal = rsp.UnSubscribeTotal + rsp.BillingRiderTotal
	return rsp
}
