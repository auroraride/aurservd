package ent

import (
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/stocksummary"
)

// Do 物资更新统计数据
func (s *StockSummaryUpdateOne) Do(req *model.StockSummaryReq) *StockSummaryUpdateOne {

	// 更新当天的统计数据
	s.Where(stocksummary.Date(req.Date))

	// 记录团签物资出入库数量
	if req.EnterpriseID != nil && req.StationID != nil {
		s.Where(stocksummary.EnterpriseID(*req.EnterpriseID), stocksummary.StationID(*req.StationID))
	}
	// 记录门店物资出入库数量
	if req.StoreID != nil {
		s.Where(stocksummary.StoreID(*req.StoreID))
	}
	// 电柜
	if req.CabinetID != nil {
		s.Where(stocksummary.CabinetID(*req.CabinetID))
	}

	// 记录电池 出/入库数量
	if req.Num > 0 && req.BatteryID != nil {
		s.AddBatteryInboundNum(req.Num)
	} else if req.Num < 0 && req.BatteryID != nil {
		s.AddBatteryOutboundNum(req.Num)
	}

	// 记录电车 出/入库数量
	if req.Num > 0 && req.EbikeID != nil {
		s.AddBikeInboundNum(req.Num)
	} else if req.Num < 0 && req.EbikeID != nil {
		s.AddBikeOutboundNum(req.Num)
	}

	// 电池总数增/减
	if req.BatteryID != nil {
		s.AddBatteryNum(req.Num)
	}
	// 电车总数增/减
	if req.EbikeID != nil {
		s.AddBikeNum(req.Num)
	}

	return s
}

// Do 物资创建统计数据
func (s *StockSummaryCreate) Do(req *model.StockSummaryReq) *StockSummaryCreate {
	// 记录日期
	s.SetDate(req.Date)
	// 记录团签物资出入库数量
	if req.EnterpriseID != nil && req.StationID != nil {
		s.SetStationID(*req.StationID).SetEnterpriseID(*req.EnterpriseID)
	}
	// 记录门店物资出入库数量
	if req.StoreID != nil {
		s.SetStoreID(*req.StoreID)
	}
	// 电柜
	if req.CabinetID != nil {
		s.SetCabinetID(*req.CabinetID)
	}

	// 记录电池 出/入库数量
	if req.Num > 0 && req.BatteryID != nil {
		s.SetBatteryInboundNum(req.Num)
	} else if req.Num < 0 && req.BatteryID != nil {
		s.SetBatteryOutboundNum(req.Num)
	}

	// 记录电车 出/入库数量
	if req.Num > 0 && req.EbikeID != nil {
		s.SetBikeInboundNum(req.Num)
	} else if req.Num < 0 && req.EbikeID != nil {
		s.SetBikeOutboundNum(req.Num)
	}

	// 电池总数增/减
	if req.BatteryID != nil {
		s.SetBatteryNum(req.Num)
	}
	// 电车总数增/减
	if req.EbikeID != nil {
		s.SetBikeNum(req.Num)
	}

	return s
}
