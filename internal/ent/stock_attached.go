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
	"github.com/auroraride/aurservd/internal/ent/stock"
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

type StockInfo interface {
	Num() (r int, exists bool)
	StoreID() (r uint64, exists bool)
	GetType() (r uint8, exists bool)
	EnterpriseID() (r uint64, exists bool)
	StationID() (r uint64, exists bool)
	CabinetID() (r uint64, exists bool)
	RiderID() (r uint64, exists bool)
	Model() (r string, exists bool)
	Material() (r stock.Material, exists bool)
	Name() (r string, exists bool)
}

func NewStockHook() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (value ent.Value, err error) {

			if st, ok := m.(StockInfo); ok {

				date := time.Now()
				q := Database.StockSummary.Query().Where(stocksummary.Date(date.Format("2006-01-02")))

				params := &model.StockSummaryParams{Date: date}

				// 获取参数
				GetParams(st, params)

				// 设置查询条件
				// 型号
				if params.Model != "" {
					q.Where(stocksummary.Model(params.Model))
				}

				if params.EnterpriseID != nil && params.StationID != nil {
					q.Where(
						stocksummary.EnterpriseID(*params.EnterpriseID),
						stocksummary.StationID(*params.StationID),
					)
				}
				if params.StoreID != nil {
					q.Where(stocksummary.StoreID(*params.StoreID))
				}
				if params.CabinetID != nil {
					q.Where(stocksummary.CabinetID(*params.CabinetID))
				}

				// 平台的出库不统计
				if params.Type == model.StockTypeTransfer {
					if params.EnterpriseID == nil && params.StationID == nil && params.CabinetID == nil && params.RiderID == nil && params.StoreID == nil {
						return next.Mutate(ctx, m)
					}
				}

				stockSummary, _ := q.First(ctx)
				// 查询统计表中是否有数据 有就更新 没有就创建
				if stockSummary == nil {
					// 创建数据时获取物资总数
					GetStockNum(ctx, params)
					// 创建统计数据
					_, err = Database.StockSummary.Create().Do(params).Save(ctx)
				} else {
					// 更新统计数据
					_, err = Database.StockSummary.UpdateOne(stockSummary).Do(params).Save(ctx)
				}
				if err != nil {
					return nil, err
				}
			}
			return next.Mutate(ctx, m)
		})
	}
}

// GetParams 获取参数
func GetParams(st StockInfo, params *model.StockSummaryParams) {
	if v, ok := st.EnterpriseID(); ok {
		params.EnterpriseID = &v
	}
	if v, ok := st.StationID(); ok {
		params.StationID = &v
	}
	if v, ok := st.CabinetID(); ok {
		params.CabinetID = &v
	}
	if v, ok := st.StoreID(); ok {
		params.StoreID = &v
	}
	if v, ok := st.RiderID(); ok {
		params.RiderID = &v
	}

	if v, ok := st.Model(); ok {
		params.Model = v
	}
	// 电车型号 其他物资名称
	if v, ok := st.Name(); ok && params.Material != stock.MaterialBattery.String() {
		params.Model = v
	}

	if v, ok := st.Material(); ok {
		params.Material = v.String()
	}

	if v, ok := st.Num(); ok {
		params.Num = v
	}

	if v, ok := st.GetType(); ok {
		params.Type = v
	}
}

// GetStockNum 获取物资总数
func GetStockNum(ctx context.Context, Params *model.StockSummaryParams) {
	// 创建数据时获取总数
	var ag []struct {
		StockNum int `json:"stockNum"` // 总数量
	}
	q := Database.StockSummary.Query()

	if Params.Model != "" {
		q.Where(stocksummary.Model(Params.Model))
	}

	switch Params.Material {
	case stock.MaterialBattery.String(): // 电池物资
		q.Where(stocksummary.MaterialEQ(stocksummary.Material(stock.MaterialBattery)))
	case stock.MaterialEbike.String(): // 电车物资
		q.Where(stocksummary.MaterialEQ(stocksummary.Material(stock.MaterialEbike)))
	case stock.MaterialOthers.String(): // 其他物资
		q.Where(stocksummary.MaterialEQ(stocksummary.Material(stock.MaterialOthers)))
	}

	//  站点查询
	if Params.StationID != nil {
		q.Where(stocksummary.StationID(*Params.StationID))
	}
	// 门店查询
	if Params.StoreID != nil {
		q.Where(stocksummary.StoreID(*Params.StoreID))
	}
	// 电柜查询
	if Params.CabinetID != nil {
		q.Where(stocksummary.CabinetID(*Params.CabinetID))
	}

	q.Aggregate(
		As(Sum(stocksummary.FieldTodayNum), "stockNum"),
	).ScanX(ctx, &ag)
	Params.StockNum = ag[0].StockNum
}
