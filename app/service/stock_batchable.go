// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-03
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
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

	if result == nil || len(result) < 0 {
		return 0
	}
	return result[0].Sum
}
