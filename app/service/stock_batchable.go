// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-03
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"fmt"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/battery"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/stock"
	"github.com/auroraride/aurservd/pkg/snag"
)

type stockBatchableService struct {
	*BaseService
	orm *ent.StockClient
}

func NewStockBatchable(params ...any) *stockBatchableService {
	return &stockBatchableService{
		BaseService: newService(params...),
		orm:         ent.Database.Stock,
	}
}

// Fetch 获取可批量调拨的物资对应库存
func (s *stockBatchableService) Fetch(target uint8, id uint64, name string) int {
	var result []struct {
		Sum       int    `json:"sum"`
		StoreID   uint64 `json:"store_id"`
		CabinetID uint64 `json:"cabinet_id"`
		StationID uint64 `json:"station_id"`
	}
	q := s.orm.Query()
	var idw predicate.Stock
	switch target {
	case model.StockTargetStore:
		idw = stock.StoreID(id)
	case model.StockTargetCabinet:
		idw = stock.CabinetID(id)
	case model.StockTargetStation:
		idw = stock.StationID(id)
		// 对于非智能电池的查询 非智能电池没有电池id
		q.Where(stock.BatteryIDIsNil())
	}

	var err error
	q.Where(stock.Name(name), idw)
	if target == model.StockTargetStation {
		err = q.GroupBy(stock.FieldStationID).Aggregate(ent.Sum(stock.FieldNum)).
			Scan(s.ctx, &result)
	} else {
		q.GroupBy(stock.FieldStoreID, stock.FieldCabinetID)
		err = q.Aggregate(ent.Sum(stock.FieldNum)).
			Scan(s.ctx, &result)
	}

	if err != nil {
		snag.Panic("物资数量获取失败")
	}
	if len(result) == 0 {
		return 0
	}
	return result[0].Sum
}

func (s *stockBatchableService) Loopers(req *model.StockTransferReq, enterpriseId uint64) ([]model.StockTransferLoopper, []string) {
	failed := make([]string, 0)
	// 查询电池信息
	q := ent.Database.Battery.Query().Where(battery.SnIn(req.BatterySn...), battery.RiderIDIsNil())

	// 平台往站点调拨 需要判断当前电池有没有被使用
	if req.OutboundTarget == model.StockTargetPlaform && req.InboundTarget == model.StockTargetStation {
		q.Where(
			battery.EnterpriseIDIsNil(),
			battery.StationIDIsNil(),
			battery.CabinetIDIsNil(),
			battery.SubscribeIDIsNil(),
		)
	}

	all, _ := q.All(s.ctx)
	if len(all) == 0 {
		snag.Panic("电池信息获取失败或电池已被使用")
	}
	Loopers := make([]model.StockTransferLoopper, 0)
	for _, bat := range all {
		// 站点调拨到站点 只能同一团签
		if req.InboundTarget == model.StockTargetStation && req.OutboundTarget == model.StockTargetStation && *bat.EnterpriseID != enterpriseId {
			failed = append(failed, fmt.Sprintf("电池调拨失败，电池[%s]不属于当前团签", bat.Sn))
			continue
		}
		// 调出到平台  不是自己站点的电池不允许调拨
		if req.OutboundTarget == model.StockTargetStation && req.InboundTarget == model.StockTargetPlaform &&
			bat.StationID != nil && *bat.StationID != req.OutboundID {
			failed = append(failed, fmt.Sprintf("电池调拨失败，[%s]不属于当前站点", bat.Sn))
			continue
		}

		Loopers = append(Loopers, model.StockTransferLoopper{
			BatterySN:    &bat.Sn,
			BatteryID:    &bat.ID,
			BatteryModel: &bat.Model,
			BrandName:    &bat.Model,
		})
	}
	return Loopers, failed
}
