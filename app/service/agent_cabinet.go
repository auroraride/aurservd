// Copyright (C) liasica. 2023-present.
//
// Created at 2023-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/pkg/snag"
)

type agentCabinetService struct {
	*BaseService
	orm *ent.CabinetClient
}

func NewAgentCabinet(params ...any) *agentCabinetService {
	return &agentCabinetService{
		BaseService: newService(params...),
	}
}

// Detail 代理端查询电柜详情
func (s *agentCabinetService) Detail(serial string, ag *ent.Agent, sts ent.EnterpriseStations) *model.AgentCabinetDetailRes {
	// 查找电柜
	q := s.orm.QueryNotDeleted().WithStation().Where(
		cabinet.EnterpriseID(ag.EnterpriseID),
		cabinet.Serial(serial),
		cabinet.StatusNEQ(model.CabinetStatusPending.Value()),
	)

	// 如果站点不为空, 则只查询站点有权限的电柜
	if sts != nil {
		ids := make([]uint64, len(sts))
		for i, st := range sts {
			ids[i] = st.ID
		}
		q.Where(cabinet.StationIDIn(ids...))
	}

	// 查询唯一电柜
	cab, _ := q.First(s.ctx)
	if cab != nil {
		snag.Panic("未找到有效电柜")
	}

	// 同步电柜并返回电柜详情
	NewCabinet().Sync(cab)
	res := &model.AgentCabinetDetailRes{
		Serial: cab.Serial,
		Name:   cab.Name,
		Status: cab.Status,
		Health: cab.Health,
		Bins:   make([]*model.AgentCabinetBin, len(cab.Bin)),
	}

	if cab.Edges.Station != nil {
		res.Station = cab.Edges.Station.Name
	}

	for i, cb := range cab.Bin {
		res.Bins[i] = &model.AgentCabinetBin{
			Ordinal:   cb.Index + 1,
			BatterySN: cb.BatterySN,
			Soc:       cb.Electricity,
			Usable:    cb.IsBizUsable(),
		}
	}

	return res
}
