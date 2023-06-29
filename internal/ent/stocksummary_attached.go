package ent

import (
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/stocksummary"
)

// Do 物资更新统计数据
func (s *StockSummaryUpdateOne) Do(params *model.StockSummaryParams) *StockSummaryUpdateOne {
	// 更新当天的统计数据
	s.Where(stocksummary.Date(params.Date.Format(carbon.DateLayout)))
	// 物资型号
	s.Where(stocksummary.Model(params.Model))
	// 物资类型
	s.Where(stocksummary.MaterialEQ(stocksummary.Material(params.Material)))

	// 团签站点纬度修改统计数据
	if params.StationID != nil {
		s.Where(stocksummary.EnterpriseID(*params.EnterpriseID), stocksummary.StationID(*params.StationID))
	}

	// 门店纬度修改统计数据
	if params.StoreID != nil {
		s.Where(stocksummary.StoreID(*params.StoreID))
	}

	// 电柜纬度修改统计数据
	if params.CabinetID != nil {
		s.Where(stocksummary.CabinetID(*params.CabinetID))
	}

	// 修改电池统计 出/入库数量
	s.updateStockSummary(params)

	return s
}

// updateBattery 修改物资统计数据
func (s *StockSummaryUpdateOne) updateStockSummary(params *model.StockSummaryParams) {
	if params.Num > 0 {
		s.AddInboundNum(params.Num)
	} else {
		s.AddOutboundNum(params.Num)
	}

	s.AddNum(params.Num)

	s.AddTodayNum(params.Num)

	if params.RiderID != nil && params.StationID != nil {
		s.AddInRiderNum(-(params.Num))
	}

}

// Do 创建物资统计数据
func (s *StockSummaryCreate) Do(params *model.StockSummaryParams) *StockSummaryCreate {
	// 记录日期 以每天为单位统计
	s.SetDate(params.Date.Format(carbon.DateLayout))
	// 记录物资类型
	s.SetMaterial(stocksummary.Material(params.Material))
	// 记录物资型号
	s.SetModel(params.Model)

	// 团签站点纬度统计数量
	if params.StationID != nil && params.StoreID == nil {
		s.SetStationID(*params.StationID).SetEnterpriseID(*params.EnterpriseID)
	}

	// 门店纬度统计数量
	if params.StoreID != nil && params.StationID == nil {
		s.SetStoreID(*params.StoreID)
	}

	// 电柜纬度统计数量
	if params.CabinetID != nil {
		s.SetCabinetID(*params.CabinetID)
	}
	// 物资数量设置
	s.createStockSummary(params)

	return s
}

// 物资统计数据
func (s *StockSummaryCreate) createStockSummary(params *model.StockSummaryParams) {
	// 出入库数量
	if params.Num > 0 {
		s.SetInboundNum(params.Num)
	} else {
		s.SetOutboundNum(params.Num)
	}
	// 今天出入库总数
	s.SetTodayNum(params.Num)
	// 总数 总数加上今天的数量
	s.SetNum(params.StockNum + params.Num)

	// 如果有骑手ID 电池出库至骑手
	if params.RiderID != nil && params.StationID != nil {
		s.SetInRiderNum(-(params.Num))
	}
}
