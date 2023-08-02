// Copyright (C) liasica. 2023-present.
//
// Created at 2023-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/LucaTheHacker/go-haversine"
	"github.com/auroraride/adapter"
	"github.com/auroraride/adapter/defs/cabdef"

	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/logging"
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

func (s *agentCabinetService) detail(ac *app.AgentContext, item *ent.Cabinet, lng *float64, lat *float64) *model.AgentCabinet {
	data := &model.AgentCabinet{
		ID:         item.ID,
		Serial:     item.Serial,
		Name:       item.Name,
		Status:     item.Status,
		Health:     item.Health,
		Permission: model.AgentCabinetPermissionView,
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

	ed := ac.Enterprise.Distance
	if lng != nil && lat != nil && ed >= 0 {
		distance := haversine.Distance(haversine.NewCoordinates(*lat, *lng), haversine.NewCoordinates(item.Lat, item.Lng)).Kilometers() * 1000.0
		data.Distance = silk.Pointer(distance)
		// 若无限制或当前距离小于控制距离上限则有全部权限
		// if ed == 0 || distance <= ed {
		// 	data.Permission = model.AgentCabinetPermissionAll
		// }
	}

	for i, bm := range item.Edges.Models {
		data.Models[i] = bm.Model
	}

	for i, b := range item.Bin {
		data.Bins[i] = &model.AgentCabinetBin{
			Ordinal:   b.Index + 1,
			Usable:    b.DoorHealth && !b.Deactivate,
			BatterySN: b.BatterySN,
			Soc:       b.Electricity,
		}
	}

	return data
}

// Detail 代理端查询电柜详情
func (s *agentCabinetService) Detail(ac *app.AgentContext, req *model.AgentCabinetDetailReq) *model.AgentCabinet {
	// 查找电柜
	q := ent.Database.Cabinet.QueryNotDeleted().WithStation().WithModels().Where(
		cabinet.Serial(req.Serial),
		cabinet.StatusNEQ(model.CabinetStatusPending.Value()),
		cabinet.StationIDIn(ac.StationIDs()...),
	)

	cab, _ := q.First(s.ctx)
	if cab == nil {
		snag.Panic("未找到有效电柜")
	}

	// 同步电柜并返回电柜详情
	NewCabinet().Sync(cab)

	return s.detail(ac, cab, req.Lng, req.Lat)
}

func (s *agentCabinetService) List(ac *app.AgentContext, req *model.AgentCabinetListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().Where(cabinet.StationIDIn(ac.StationIDs()...)).WithStation().WithModels()

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
			return s.detail(ac, item, req.Lng, req.Lat)
		},
		NewCabinet().SyncCabinets,
	)
}

func (s *agentCabinetService) Operable(ac *app.AgentContext, id uint64, lng, lat float64) *ent.Cabinet {
	cab, _ := s.orm.QueryNotDeleted().Where(
		cabinet.StationIDIn(ac.StationIDs()...),
		cabinet.StatusNEQ(model.CabinetStatusPending.Value()),
		cabinet.ID(id),
	).First(s.ctx)

	if cab == nil {
		snag.Panic("未找到有效电柜")
	}

	// 暂时屏蔽电柜控制权限
	snag.Panic("无权限操作")

	ed := ac.Enterprise.Distance
	if ed < 0 || (ed > 0 && haversine.Distance(haversine.NewCoordinates(lat, lng), haversine.NewCoordinates(cab.Lat, cab.Lng)).Kilometers()*1000.0 > ed) {
		snag.Panic("操作距离过远")
	}

	return cab
}

// Maintain 设置电柜操作维护
func (s *agentCabinetService) Maintain(ac *app.AgentContext, req *model.AgentMaintainReq) {
	cab := s.Operable(ac, req.ID, req.Lng, req.Lat)

	status := model.CabinetStatusNormal
	if *req.Maintain {
		status = model.CabinetStatusMaintenance
	}

	err := cab.Update().SetStatus(status.Value()).Exec(s.ctx)
	if err != nil {
		snag.Panic(err)
	}

	// 记录日志
	go logging.NewOperateLog().
		SetRef(cab).
		SetAgent(ac.Agent).
		SetOperate(model.OperateCabinetMaintain).
		SetDiff(model.CabinetStatus(cab.Status).String(), status.String()).
		Send()
}

// BinOpen 仓位操作
func (s *agentCabinetService) BinOpen(ac *app.AgentContext, req *model.AgentBinOperateReq, operate cabdef.Operate) []*cabdef.BinOperateResult {
	cab := s.Operable(ac, req.ID, req.Lng, req.Lat)

	payload := &cabdef.OperateBinRequest{
		Operate: operate,
		Ordinal: req.Ordinal,
		Serial:  cab.Serial,
		Remark:  "代理开仓",
	}

	results, err := adapter.Post[[]*cabdef.BinOperateResult](s.GetCabinetAdapterUrlX(cab, "/agent/operate/bin/open"), ac.Agent.AdapterUser(), payload)

	if err != nil {
		snag.Panic(err)
	}

	return results
}

// AllCabinet 查询代理所有智能或者非智能电柜
func (s *agentCabinetService) AllCabinet(ac *app.AgentContext, isIntelligent *bool) []*ent.Cabinet {
	q := s.orm.QueryNotDeleted().Where(
		cabinet.StationIDIn(ac.StationIDs()...),
		cabinet.EnterpriseID(ac.Enterprise.ID),
	).WithModels()
	if isIntelligent != nil {
		q.Where(cabinet.Intelligent(*isIntelligent))
	}
	return q.AllX(s.ctx)
}
