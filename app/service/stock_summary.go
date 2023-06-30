package service

import (
	"github.com/auroraride/aurservd/app"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/battery"
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
func (s *stockSummaryService) BatteryStockSummary(ac *app.AgentContext) model.BatteryStockSummaryRsp {
	rsp := model.BatteryStockSummaryRsp{}
	rsp.Overview = s.BatterySummary(ac)
	rsp.Group = s.BatteryGroup(ac)
	return rsp
}

// EbikeStockSummary 电动车物资统计
func (s *stockSummaryService) EbikeStockSummary(ac *app.AgentContext) model.EbikeStockSummaryRsp {
	rsp := model.EbikeStockSummaryRsp{}
	rsp.Overview = s.EbikeSummary(ac)
	rsp.Group = s.EbikeGroup(ac)
	return rsp
}

// BatterySummary 电池数据统计
func (s *stockSummaryService) BatterySummary(ac *app.AgentContext) model.BatterySummary {
	var v []model.BatterySummary
	s.orm.Query().
		Where(stocksummary.EnterpriseID(ac.Enterprise.ID), stocksummary.MaterialEQ(stocksummary.MaterialBattery)).
		Aggregate(
			ent.As(ent.Sum(stocksummary.FieldTodayNum), "stationBatteryTotal"),
			ent.As(ent.Sum(stocksummary.FieldInRiderNum), "riderBatteryTotal"),
		).ScanX(s.ctx, &v)
	// 在电柜非智能电池数量
	nonIntelligentNum, _ := s.NonIntelligentBatteryNum(ac)
	// 智能电池在电柜数量
	intelligentNum := ent.Database.Battery.Query().Where(battery.EnterpriseID(ac.Enterprise.ID), battery.CabinetIDNotNil()).CountX(s.ctx)

	// 电池总数
	rsp := v[0]
	rsp.CabinetBatteryTotal = nonIntelligentNum + intelligentNum
	rsp.BatteryTotal = rsp.StationBatteryTotal + rsp.CabinetBatteryTotal + rsp.RiderBatteryTotal
	return rsp
}

// EbikeSummary 电动车数据统计
func (s *stockSummaryService) EbikeSummary(ac *app.AgentContext) model.EbikeSummary {
	var v []model.EbikeSummary
	s.orm.Query().
		Where(stocksummary.EnterpriseID(ac.Enterprise.ID), stocksummary.MaterialEQ(stocksummary.MaterialEbike)).
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
func (s *stockSummaryService) BatteryGroup(ac *app.AgentContext) []*model.BatteryStockGroup {
	var rsp []*model.BatteryStockGroup

	s.orm.Query().
		Where(stocksummary.EnterpriseID(ac.Enterprise.ID), stocksummary.MaterialEQ(stocksummary.MaterialBattery)).
		Order(ent.Asc(stocksummary.FieldModel)).
		GroupBy(stocksummary.FieldModel).
		Aggregate(
			ent.As(ent.Sum(stocksummary.FieldTodayNum), "stationBatteryTotal"),
			ent.As(ent.Sum(stocksummary.FieldInRiderNum), "riderBatteryTotal"),
		).ScanX(s.ctx, &rsp)

	// 在电柜非智能电池数量
	_, nonIntelligent := s.NonIntelligentBatteryNum(ac)

	// 智能电池在电柜数量
	intelligent := s.IntelligentBatteryNum(ac)

	modelTotal := make(map[string]int)

	for k, v := range nonIntelligent {
		modelTotal[k] += v
	}

	for k, v := range intelligent {
		modelTotal[k] += v
	}

	processedModels := make(map[string]bool) // 创建一个用于跟踪已处理模型的映射

	// 处理已统计到的电池型号
	for i, v := range rsp {
		if total, ok := modelTotal[v.Model]; ok {
			rsp[i].CabinetBatteryTotal += total
			processedModels[v.Model] = true
		}
		// 电池总数
		rsp[i].BatteryTotal = rsp[i].StationBatteryTotal + rsp[i].CabinetBatteryTotal + rsp[i].RiderBatteryTotal
	}
	// 处理未统计到的电池型号
	for m, total := range modelTotal {
		if !processedModels[m] {
			newBatteryGroup := &model.BatteryStockGroup{
				Model: m,
				BatterySummary: model.BatterySummary{
					CabinetBatteryTotal: total,
					BatteryTotal:        total,
				},
			}
			rsp = append(rsp, newBatteryGroup)
		}
	}

	return rsp
}

// EbikeGroup 分组统计电动车物资
func (s *stockSummaryService) EbikeGroup(ac *app.AgentContext) []*model.EbikeStockGroup {
	var v []*model.EbikeStockGroup
	s.orm.Query().
		Where(stocksummary.EnterpriseID(ac.Enterprise.ID), stocksummary.MaterialEQ(stocksummary.MaterialEbike)).
		Order(ent.Asc(stocksummary.FieldModel)).
		GroupBy(stocksummary.FieldModel).
		Aggregate(
			ent.As(ent.Sum(stocksummary.FieldTodayNum), "stationEbikeTotal"),
			ent.As(ent.Sum(stocksummary.FieldInRiderNum), "riderEbikeTotal"),
		).ScanX(s.ctx, &v)

	if len(v) > 0 {
		for _, item := range v {
			item.EbikeTotal = item.StationEbikeTotal + item.RiderEbikeTotal
		}
	}
	return v
}

func (s *stockSummaryService) NonIntelligentBatteryNum(ac *app.AgentContext) (batteryNum int, modelsCounts map[string]int) {
	allCabinet := NewAgentCabinet().AllNonIntelligentCabinet(ac)
	NewCabinet().SyncCabinets(allCabinet)

	// 电池型号 非智能电池柜只有一种电池型号
	modelsCounts = make(map[string]int) // 电池型号 -> 数量

	// 计算电池总数
	for _, cab := range allCabinet {
		batteryNum += cab.BatteryNum

		// 统计电池型号 -> 数量
		if cab.Edges.Models != nil && len(cab.Edges.Models) > 0 {
			m := cab.Edges.Models[0].Model
			modelsCounts[m] += cab.BatteryNum
		}
	}

	return batteryNum, modelsCounts
}

func (s *stockSummaryService) IntelligentBatteryNum(ac *app.AgentContext) (modelsCounts map[string]int) {
	var intelligent []model.BatteryStockGroup
	ent.Database.Battery.Query().
		Where(
			battery.EnterpriseID(ac.Enterprise.ID),
			battery.CabinetIDNotNil(),
		).
		GroupBy(battery.FieldModel).Aggregate(ent.As(ent.Count(), "batteryTotal")).ScanX(s.ctx, &intelligent)
	modelsCounts = make(map[string]int)
	for _, v := range intelligent {
		modelsCounts[v.Model] = v.BatteryTotal
	}
	return modelsCounts
}
