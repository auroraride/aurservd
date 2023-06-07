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
	}

	var idw predicate.Stock
	switch target {
	case model.StockTargetStore:
		idw = stock.StoreID(id)
	case model.StockTargetCabinet:
		idw = stock.CabinetID(id)
	}
	q := s.orm.QueryNotDeleted().
		Where(stock.Name(name), idw).
		GroupBy(stock.FieldStoreID, stock.FieldCabinetID).
		Aggregate(ent.Sum(stock.FieldNum))
	err := q.Scan(s.ctx, &result)
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

	all, _ := q.All(s.ctx)
	if len(all) == 0 {
		snag.Panic("电池信息获取失败或电池已被使用")
	}
	Loopers := make([]model.StockTransferLoopper, 0)
	for _, bat := range all {
		if bat.EnterpriseID == nil || bat.StationID == nil { // 电池未绑定团签，归属于平台不能调拨
			failed = append(failed, fmt.Sprintf("电池调拨失败，电池[%s]未绑定团签，归属于平台不能调拨", bat.Sn))
			continue
		}
		// 站点调拨到站点 只能同一团签
		if req.InboundTarget == model.StockTargetStation && req.OutboundTarget == model.StockTargetStation && *bat.EnterpriseID != enterpriseId {
			failed = append(failed, fmt.Sprintf("电池调拨失败，电池[%s]不属于当前团签", bat.Sn))
			continue
		}
		brandName := bat.Brand.String()
		Loopers = append(Loopers, model.StockTransferLoopper{
			BatterySN:    &bat.Sn,
			BatteryID:    &bat.ID,
			BatteryModel: bat.Model,
			BrandName:    &brandName,
		})
	}
	return Loopers, failed
}
