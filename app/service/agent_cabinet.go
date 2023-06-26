// Copyright (C) liasica. 2023-present.
//
// Created at 2023-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/LucaTheHacker/go-haversine"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type agentCabinetService struct {
	*BaseService
	orm *ent.CabinetClient
}

func NewAgentCabinet(params ...any) *agentCabinetService {
	return &agentCabinetService{
		BaseService: newService(params...),
		orm:         ent.Database.Cabinet,
	}
}

func (s *agentCabinetService) detail(item *ent.Cabinet, lng *float64, lat *float64) *model.AgentCabinet {
	data := &model.AgentCabinet{
		Serial:     item.Serial,
		Name:       item.Name,
		Status:     item.Status,
		Health:     item.Health,
		Permission: model.AgentCabinetPermissionAll,
		Address:    item.Address,
		Lng:        item.Lng,
		Lat:        item.Lat,
		Station: model.EnterpriseStation{
			ID:   item.Edges.Station.ID,
			Name: item.Edges.Station.Name,
		},
		Bins:   make([]*model.AgentCabinetBin, len(item.Bin)),
		Models: make([]string, len(item.Edges.Models)),
	}

	if lng != nil && lat != nil {
		data.Distance = silk.Pointer(haversine.Distance(haversine.NewCoordinates(*lat, *lng), haversine.NewCoordinates(item.Lat, item.Lng)).Kilometers() * 1000.0)
	}

	for i, bm := range item.Edges.Models {
		data.Models[i] = bm.Model
	}

	for i, b := range item.Bin {
		data.Bins[i] = &model.AgentCabinetBin{
			Ordinal:   b.Index + 1,
			Usable:    b.CanUse(),
			BatterySN: b.BatterySN,
			Soc:       b.Electricity,
		}
	}

	return data
}

// Detail 代理端查询电柜详情
func (s *agentCabinetService) Detail(stations []uint64, req *model.AgentCabinetDetailReq) *model.AgentCabinet {
	// 查找电柜
	q := ent.Database.Cabinet.QueryNotDeleted().WithStation().WithModels().Where(
		cabinet.Serial(req.Serial),
		cabinet.StatusNEQ(model.CabinetStatusPending.Value()),
		cabinet.StationIDIn(stations...),
	)

	cab, _ := q.First(s.ctx)
	if cab == nil {
		snag.Panic("未找到有效电柜")
	}

	// 同步电柜并返回电柜详情
	NewCabinet().Sync(cab)

	return s.detail(cab, req.Lng, req.Lat)
}

func (s *agentCabinetService) List(stations []uint64, req *model.AgentCabinetListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().Where(cabinet.StationIDIn(stations...)).WithStation().WithModels()

	// 筛选站点
	if req.StationID != nil {
		q.Where(cabinet.StationID(*req.StationID))
	}

	// 筛选电池型号
	if req.Model != nil {
		q.Where(cabinet.HasModelsWith(batterymodel.Model(*req.Model)))
	}

	// 筛选电柜编码
	if req.Serial != nil {
		q.Where(cabinet.SerialContainsFold(*req.Serial))
	}

	// 按距离排序
	if req.Lng != nil && req.Lat != nil {
		point := (&model.Geometry{Lng: *req.Lng, Lat: *req.Lat}).String()
		q.Modify(func(sel *sql.Selector) {
			sel.AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(%s, ST_GeographyFromText('%s'))`, sel.C(cabinet.FieldGeom), point)), "distance")
			sel.OrderBy("distance")
		})
	}

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.Cabinet) *model.AgentCabinet {
			return s.detail(item, req.Lng, req.Lat)
		},
		NewCabinet().SyncCabinets,
	)
}

func (s *agentCabinetService) Open() {

}
