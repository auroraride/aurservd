// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-16
// Based on aurservd by liasica, magicrolan@qq.com.

package ent

import (
	"context"
	"time"

	"entgo.io/ent"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent/stocksummary"
)

func (sc *StockCreate) Clone() (creator *StockCreate) {
	mutation := new(StockMutation)
	*mutation = *sc.mutation
	return &StockCreate{
		config:   sc.config,
		mutation: mutation,
		hooks:    sc.hooks,
		conflict: sc.conflict,
	}
}

type intr interface {
	Num() (r int, exists bool)
	StoreID() (r uint64, exists bool)
	GetType() (r uint8, exists bool)
	EnterpriseID() (r uint64, exists bool)
	StationID() (r uint64, exists bool)
	CabinetID() (r uint64, exists bool)
	EbikeID() (r uint64, exists bool)
	BatteryID() (r uint64, exists bool)
	RiderID() (r uint64, exists bool)
}

func NewStockHook() Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (value ent.Value, err error) {
			var (
				batteryID    *uint64
				eBikeID      *uint64
				riderID      *uint64
				cabinetID    *uint64
				stationID    *uint64
				enterpriseID *uint64
				storeID      *uint64
				num          int
			)

			if st, ok := m.(intr); ok {

				date := time.Now().Format("2006-01-02")
				q := Database.StockSummary.Query().Where(stocksummary.Date(date))

				if v, ok := st.EnterpriseID(); ok {
					enterpriseID = &v
				}

				if v, ok := st.StationID(); ok {
					stationID = &v
					q.Where(stocksummary.EnterpriseID(*enterpriseID), stocksummary.StationID(*stationID))
				}

				if v, ok := st.CabinetID(); ok {
					cabinetID = &v
					q.Where(stocksummary.CabinetID(*cabinetID))
				}

				if v, ok := st.StoreID(); ok {
					storeID = &v
					q.Where(stocksummary.StoreID(*storeID))
				}

				if v, ok := st.EbikeID(); ok {
					eBikeID = &v
				}
				if v, ok := st.BatteryID(); ok {
					batteryID = &v
				}
				if v, ok := st.RiderID(); ok {
					riderID = &v
				}

				if v, ok := st.Num(); ok {
					num = v
				}

				// 平台的出库不统计
				if v, ok := st.GetType(); ok && v == model.StockTypeTransfer {
					if enterpriseID == nil && stationID == nil && cabinetID == nil && riderID == nil && storeID == nil {
						return next.Mutate(ctx, m)
					}
				}

				stockSummary, _ := q.First(ctx)
				req := &model.StockSummaryReq{
					Num:          num,
					StoreID:      storeID,
					EnterpriseID: enterpriseID,
					StationID:    stationID,
					CabinetID:    cabinetID,
					EbikeID:      eBikeID,
					BatteryID:    batteryID,
					Date:         date,
				}
				// 查询统计表中是否有数据 有就更新 没有就创建
				if stockSummary == nil {
					// 创建统计数据
					_, err = Database.StockSummary.Create().Do(req).Save(ctx)
				} else {
					// 更新统计数据
					_, err = Database.StockSummary.UpdateOne(stockSummary).Do(req).Save(ctx)
				}
				if err != nil {
					return nil, err
				}
			}
			return next.Mutate(ctx, m)
		})
	}
}
