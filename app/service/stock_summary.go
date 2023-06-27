package service

import (
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/stocksummary"
)

type stockSummaryService struct {
	*BaseService
	orm *ent.StockSummaryClient
}

func NewStockSummary(params ...any) *stockSummaryService {
	return &stockSummaryService{
		BaseService: newService(params...),
		orm:         ent.Database.StockSummary,
	}
}

// BatteryStockSummary 电池物资统计
func (s *stockSummaryService) BatteryStockSummary(req *model.StockSummaryReq) model.BatteryStockSummaryRsp {
	rsp := model.BatteryStockSummaryRsp{}
	rsp.Overview = s.BatterySummary(req.EnterpriseID)
	rsp.Group = s.BatteryGroup(req.EnterpriseID)
	return rsp
}

// EbikeStockSummary 电动车物资统计
func (s *stockSummaryService) EbikeStockSummary(req *model.StockSummaryReq) model.EbikeStockSummaryRsp {
	rsp := model.EbikeStockSummaryRsp{}
	rsp.Overview = s.EbikeSummary(req.EnterpriseID)
	rsp.Group = s.EbikeGroup(req.EnterpriseID)
	return rsp
}

// BatterySummary 电池数据统计
func (s *stockSummaryService) BatterySummary(enterpriseID uint64) model.BatterySummary {
	var v []model.BatterySummary
	s.orm.Query().
		Where(stocksummary.EnterpriseID(enterpriseID), stocksummary.MaterialEQ(stocksummary.MaterialBattery)).
		Aggregate(
			ent.As(ent.Sum(stocksummary.FieldTodayNum), "stationBatteryTotal"),
			ent.As(ent.Sum(stocksummary.FieldInRiderNum), "riderBatteryTotal"),
			ent.As(ent.Sum(stocksummary.FieldInCabinetNum), "cabinetBatteryTotal"),
		).ScanX(s.ctx, &v)

	// 电池总数
	rsp := v[0]
	rsp.BatteryTotal = rsp.StationBatteryTotal + rsp.CabinetBatteryTotal + rsp.RiderBatteryTotal
	return rsp
}

// EbikeSummary 电动车数据统计
func (s *stockSummaryService) EbikeSummary(enterpriseID uint64) model.EbikeSummary {
	var v []model.EbikeSummary
	s.orm.Query().
		Where(stocksummary.EnterpriseID(enterpriseID), stocksummary.MaterialEQ(stocksummary.MaterialEbike)).
		Aggregate(
			ent.As(ent.Sum(stocksummary.FieldTodayNum), "stationEbikeTotal"),
			ent.As(ent.Sum(stocksummary.FieldInRiderNum), "riderEbikeTotal"),
		).ScanX(s.ctx, &v)
	// 电动车总数
	rsp := v[0]
	rsp.EbikeTotal = rsp.StationEbikeTotal + rsp.RiderEbikeTotal
	return rsp
}

// BatteryGroup 分组统计电池物资
func (s *stockSummaryService) BatteryGroup(enterpriseID uint64) []model.BatteryStockGroup {
	var v []model.BatteryStockGroup

	s.orm.Query().
		Where(stocksummary.EnterpriseID(enterpriseID), stocksummary.MaterialEQ(stocksummary.MaterialBattery)).
		GroupBy(stocksummary.FieldModel).
		Aggregate(
			ent.As(ent.Sum(stocksummary.FieldTodayNum), "stationBatteryTotal"),
			ent.As(ent.Sum(stocksummary.FieldInCabinetNum), "cabinetBatteryTotal"),
			ent.As(ent.Sum(stocksummary.FieldInRiderNum), "riderBatteryTotal"),
		).ScanX(s.ctx, &v)
	if len(v) > 0 {
		for _, item := range v {
			item.BatteryTotal = item.StationBatteryTotal + item.CabinetBatteryTotal + item.RiderBatteryTotal
		}
	}
	return v
}

// EbikeGroup 分组统计电动车物资
func (s *stockSummaryService) EbikeGroup(enterpriseID uint64) []model.EbikeStockGroup {
	var v []model.EbikeStockGroup
	s.orm.Query().
		Where(stocksummary.EnterpriseID(enterpriseID), stocksummary.MaterialEQ(stocksummary.MaterialEbike)).
		GroupBy(stocksummary.FieldModel).
		Aggregate(
			ent.As(ent.Sum(stocksummary.FieldTodayNum), "stationEbikeTotal"),
			ent.As(ent.Sum(stocksummary.FieldInRiderNum), "riderEbikeTotal"),
		).ScanX(s.ctx, &v)
	rsp := make([]model.EbikeStockGroup, 0)
	if len(v) > 0 {
		for _, item := range v {
			item.EbikeTotal = item.StationEbikeTotal + item.RiderEbikeTotal
		}
	}
	return rsp
}
