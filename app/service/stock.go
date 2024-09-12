// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-12
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/battery"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/ebike"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/stock"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/tools"
)

type stockService struct {
	ctx          context.Context
	modifier     *model.Modifier
	employee     *ent.Employee
	orm          *ent.StockClient
	employeeInfo *model.Employee
	*BaseService
}

func NewStock(params ...any) *stockService {
	return &stockService{
		ctx:         context.Background(),
		orm:         ent.Database.Stock,
		BaseService: newService(params...),
	}
}

func NewStockWithModifier(m *model.Modifier) *stockService {
	s := NewStock()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

func NewStockWithEmployee(e *ent.Employee) *stockService {
	s := NewStock()
	if e != nil {
		s.employee = e
		s.employeeInfo = &model.Employee{
			ID:    e.ID,
			Name:  e.Name,
			Phone: e.Phone,
		}
		s.ctx = context.WithValue(s.ctx, model.CtxEmployeeKey{}, s.employeeInfo)
	}
	return s
}

// StoreList 门店物资
// func (s *stockService) StoreList(req *model.StockListReq) *model.PaginationRes {
// 	q := ent.Database.Store.QueryNotDeleted().
// 		Where(store.HasStocks()).
// 		WithCity().
// 		WithStocks().
// 		WithExceptions(func(eq *ent.ExceptionQuery) {
// 			eq.Where(exception.Status(model.ExceptionStatusPending))
// 		})
// 	if req.Name != nil {
// 		q.Where(
// 			store.NameContainsFold(*req.Name),
// 		)
// 	}
// 	if req.CityID != nil {
// 		q.Where(
// 			store.HasBranchWith(branch.CityID(*req.CityID)),
// 		)
// 	}
// 	if req.StoreID != nil {
// 		q.Where(
// 			store.ID(*req.StoreID),
// 		)
// 	}
//
// 	if req.Start != nil {
// 		start := carbon.Parse(*req.Start).StartOfDay().StdTime()
// 		if start.IsZero() {
// 			snag.Panic("开始时间错误")
// 		}
// 		q.Where(
// 			store.Or(
// 				store.HasStocksWith(stock.CreatedAtGTE(start)),
// 			),
// 		)
// 	}
// 	if req.End != nil {
// 		end := carbon.Parse(*req.End).EndOfDay().StdTime()
// 		if end.IsZero() {
// 			snag.Panic("结束时间错误")
// 		}
// 		q.Where(
// 			store.Or(
// 				store.HasStocksWith(stock.CreatedAtLTE(end)),
// 			),
// 		)
// 	}
// 	if req.EbikeBrandID != 0 {
// 		q.Where(store.HasStocksWith(stock.BrandID(req.EbikeBrandID)))
// 	}
// 	if req.Model != "" {
// 		q.Where(store.HasStocksWith(stock.Model(req.Model)))
// 	}
// 	if req.Keyword != "" {
// 		q.Where(store.HasStocksWith(stock.NameContainsFold(req.Keyword)))
// 	}
// 	return model.ParsePaginationResponse(
// 		q,
// 		req.PaginationReq,
// 		func(item *ent.Store) model.StockListRes {
// 			res := model.StockListRes{
// 				Store: model.Store{
// 					ID:   item.ID,
// 					Name: item.Name,
// 				},
// 				City: model.City{
// 					ID:   item.Edges.City.ID,
// 					Name: item.Edges.City.Name,
// 				},
// 				BatteryTotal: 0,
// 				Batteries:    make([]*model.StockMaterial, 0),
// 				Materials:    make([]*model.StockMaterial, 0),
// 				Ebikes:       make([]*model.StockMaterial, 0),
// 			}
//
// 			// 计算所有物资
// 			batteries := make(map[string]*model.StockMaterial)
// 			materials := make(map[string]*model.StockMaterial)
// 			bikes := make(map[string]*model.StockMaterial)
//
// 			// 出入库
// 			for _, st := range item.Edges.Stocks {
// 				switch true {
// 				case st.Model != nil:
// 					// 电池
// 					s.Calculate(batteries, st)
// 				case st.BrandID != nil:
// 					// 电车
// 					s.Calculate(bikes, st)
// 				default:
// 					// 其他物资
// 					s.Calculate(materials, st)
// 				}
// 			}
//
// 			for _, ex := range item.Edges.Exceptions {
// 				if ex.Model != nil {
// 					s.CalculateException(batteries, ex)
// 				} else {
// 					s.CalculateException(materials, ex)
// 				}
// 			}
//
// 			for _, battery := range batteries {
// 				res.Batteries = append(res.Batteries, battery)
// 				res.BatteryTotal += battery.Surplus
// 			}
//
// 			for _, bike := range bikes {
// 				res.Ebikes = append(res.Ebikes, bike)
// 				res.EbikeTotal += bike.Surplus
// 			}
//
// 			for _, material := range materials {
// 				res.Materials = append(res.Materials, material)
// 			}
//
// 			// 排序
// 			sort.Slice(res.Batteries, func(i, j int) bool {
// 				return strings.Compare(res.Batteries[i].Name, res.Batteries[j].Name) < 0
// 			})
// 			sort.Slice(res.Materials, func(i, j int) bool {
// 				return strings.Compare(res.Materials[i].Name, res.Materials[j].Name) < 0
// 			})
//
// 			return res
// 		},
// 	)
// }
//
// // EnterpriseList 团签物资
// func (s *stockService) EnterpriseList(req *model.StockListReq) *model.PaginationRes {
// 	q := ent.Database.EnterpriseStation.QueryNotDeleted().WithEnterprise(func(query *ent.EnterpriseQuery) {
// 		query.WithStations().WithCity()
// 	}).WithStocks()
//
// 	if req.CityID != nil {
// 		q.Where(
// 			enterprisestation.HasEnterpriseWith(enterprise.HasCityWith(city.ID(*req.CityID))),
// 		)
// 	}
//
// 	if req.EnterpriseID != nil {
// 		q.Where(
// 			enterprisestation.HasEnterpriseWith(enterprise.ID(*req.EnterpriseID)),
// 		)
// 	}
//
// 	if req.StationID != nil {
// 		q.Where(enterprisestation.ID(*req.StationID))
// 	}
//
// 	if req.Start != nil {
// 		q.Where(enterprisestation.HasStocksWith(
// 			stock.CreatedAtGTE(tools.NewTime().ParseDateStringX(*req.Start)),
// 		))
// 	}
//
// 	if req.End != nil {
// 		q.Where(enterprisestation.HasStocksWith(
// 			stock.CreatedAtLT(tools.NewTime().ParseNextDateStringX(*req.Start)),
// 		))
// 	}
//
// 	if req.Model != "" {
// 		q.Where(enterprisestation.HasEnterpriseWith(enterprise.HasStocksWith(stock.Model(req.Model))))
// 	}
//
// 	return model.ParsePaginationResponse(
// 		q,
// 		req.PaginationReq,
// 		func(item *ent.EnterpriseStation) model.StockListRes {
// 			rsp := model.StockListRes{
// 				StationName: item.Name,
// 				StationID:   item.ID,
// 			}
// 			if item.Edges.Enterprise != nil {
// 				rsp.EnterpriseName = item.Edges.Enterprise.Name
// 				rsp.EnterpriseID = item.Edges.Enterprise.ID
// 				if item.Edges.Enterprise.Edges.City != nil {
// 					rsp.City = model.City{
// 						ID:   item.Edges.Enterprise.Edges.City.ID,
// 						Name: item.Edges.Enterprise.Edges.City.Name,
// 					}
// 				}
// 			}
// 			// 电池总数
// 			rsp.BatteryTotal = 0
// 			// 计算电池物资
// 			batteries := make(map[string]*model.StockMaterial)
// 			// 计算电车物资
// 			bikes := make(map[string]*model.StockMaterial)
// 			// 出入库
// 			for _, st := range item.Edges.Stocks {
// 				switch true {
// 				case st.Model != nil:
// 					// 电池
// 					s.Calculate(batteries, st)
// 				case st.BrandID != nil:
// 					// 电车
// 					s.Calculate(bikes, st)
// 				}
// 			}
// 			for _, bat := range batteries {
// 				rsp.Batteries = append(rsp.Batteries, bat)
// 				rsp.BatteryTotal += bat.Surplus
// 			}
// 			for _, bike := range bikes {
// 				rsp.Ebikes = append(rsp.Ebikes, bike)
// 				rsp.EbikeTotal += bike.Surplus
// 			}
//
// 			return rsp
// 		})
// }

// TODO 统计故障电车
// func (s *stockService) CalculateException(items map[string]*model.StockMaterial, ex *ent.Exception) {
// 	name := ex.Name
// 	if _, ok := items[name]; !ok {
// 		items[name] = &model.StockMaterial{
// 			Name:     name,
// 			Outbound: 0,
// 			Inbound:  0,
// 			Surplus:  0,
// 		}
// 	}
// 	items[name].Exception += ex.Num
// }

// func (s *stockService) getKey(st *ent.Stock) string {
// 	if st.BrandID != nil {
// 		return fmt.Sprintf("%d", *st.BrandID)
// 	}
// 	return st.Name
// }

// func (s *stockService) Calculate(items map[string]*model.StockMaterial, st *ent.Stock) {
// 	key := s.getKey(st)
// 	if _, ok := items[key]; !ok {
// 		items[key] = &model.StockMaterial{
// 			Name:     st.Name,
// 			Outbound: 0,
// 			Inbound:  0,
// 			Surplus:  0,
// 		}
// 	}
// 	if st.Num > 0 {
// 		items[key].Inbound += st.Num
// 	} else {
// 		items[key].Outbound += int(math.Abs(float64(st.Num)))
// 	}
// 	items[key].Surplus += st.Num
// }
//
// func (s *stockService) BatteryOverview(req *model.StockOverviewReq) (items []model.StockBatteryOverviewRes) {
// 	var extends []string
//
// 	switch req.Goal {
// 	case model.StockGoalStore:
// 		if req.StoreID != 0 {
// 			extends = append(extends, fmt.Sprintf("AND (%s = %d)", stock.FieldStoreID, req.StoreID))
// 		} else {
// 			extends = append(extends, fmt.Sprintf("AND (%s IS NOT NULL OR (%s IS NULL AND %s IS NULL AND %s > 0))", stock.FieldStoreID, stock.FieldStoreID, stock.FieldCabinetID, stock.FieldType))
// 		}
// 	case model.StockGoalCabinet:
// 		if req.CabinetID != 0 {
// 			extends = append(extends, fmt.Sprintf("AND (%s = %d)", stock.FieldCabinetID, req.StoreID))
// 		} else {
// 			extends = append(extends, fmt.Sprintf("AND (%s IS NOT NULL)", stock.FieldCabinetID))
// 		}
// 	default:
// 		extends = append(extends, fmt.Sprintf("AND (%s IS NOT NULL OR %s IS NOT NULL OR %s IS NOT NULL)", stock.FieldStoreID, stock.FieldCabinetID, stock.FieldRiderID))
// 	}
//
// 	if req.CityID != 0 {
// 		extends = append(extends, fmt.Sprintf("AND (%s = %d)", stock.FieldCityID, req.CityID))
// 	}
//
// 	if req.Start != "" && req.End != "" {
// 		start := tools.NewTime().ParseDateStringX(req.Start).Format(carbon.DateTimeLayout)
// 		end := tools.NewTime().ParseNextDateStringX(req.End).Format(carbon.DateTimeLayout)
// 		extends = append(extends, fmt.Sprintf("AND (%s >= '%s'::timestamp AND %s < '%s'::timestamp)", stock.FieldCreatedAt, start, stock.FieldCreatedAt, end))
// 	}
//
// 	extend := fmt.Sprintf("WHERE model IS NOT NULL %s", strings.Join(extends, " "))
// 	query := strings.Replace(assets.SQLStockOverview, "WHERE model IS NOT NULL", extend, 1)
//
// 	rows, err := ent.Database.QueryContext(s.ctx, query)
//
// 	if err != nil {
// 		snag.Panic("请求失败")
// 	}
//
// 	defer func(rows *stdsql.Rows) {
// 		_ = rows.Close()
// 	}(rows)
//
// 	for rows.Next() {
// 		var b []byte
// 		err = rows.Scan(&b)
// 		if err != nil {
// 			break
// 		}
// 		var item model.StockBatteryOverviewRes
// 		_ = jsoniter.Unmarshal(b, &item)
// 		items = append(items, item)
// 	}
//
// 	return
// }

// RiderBusiness 骑手业务 电池 / 电车 出入库
// 此方法操作库存
func (s *stockService) RiderBusiness(req *model.StockBusinessReq) (err error) {
	if req.StoreID == nil && req.EbikeStoreID == nil && req.BatStoreID == nil && req.CabinetID == nil && req.EnterpriseID == nil && req.StationID == nil {
		err = errors.New("参数校验错误")
		return
	}

	// 如果是骑手自己操作 激活和取消寄存拿走电池 会有电柜任务已经执行调拨
	if s.operator.Type == model.OperatorTypeRider && req.CabinetID != nil && (req.AssetTransferType == model.AssetTransferTypeActive || req.AssetTransferType == model.AssetTransferTypeContinue) {
		return
	}

	// 查询资产
	var ebikeInfo *ent.Asset
	var batteryInfo *ent.Asset
	var fromLocationType model.AssetLocationsType
	var fromLocationID uint64
	var toLocationType model.AssetLocationsType
	var toLocationID uint64
	details := make([]model.AssetTransferCreateDetail, 0)
	ebikeDetails := make([]model.AssetTransferCreateDetail, 0)
	assetType := model.AssetTypeSmartBattery

	if req.Ebike != nil {
		ebikeInfo, _ = NewAsset().QuerySn(req.Ebike.Sn)
		if ebikeInfo == nil {
			err = errors.New("电车不存在")
			return
		}
	}

	if req.Battery != nil {
		if req.Battery.SN == "" {
			batteryInfo, _ = NewAsset().QueryID(req.Battery.ID)
		} else {
			batteryInfo, _ = NewAsset().QuerySn(req.Battery.SN)
		}
		if batteryInfo == nil {
			err = errors.New("电池不存在")
			return
		}
	}
	var storeID uint64
	if req.StoreID != nil {
		storeID = *req.StoreID
	}
	if req.BatStoreID != nil {
		storeID = *req.BatStoreID
	}

	// 激活和取消寄存 需要判定非智能库存
	if req.Battery == nil && (req.AssetTransferType == model.AssetTransferTypeActive || req.AssetTransferType == model.AssetTransferTypeContinue) {
		assetType = model.AssetTypeNonSmartBattery
		if req.StoreID != nil || req.BatStoreID != nil {
			batteryInfo, _ = s.CheckBusinessBattery(req, model.AssetLocationsTypeStore, storeID)
		}
		// 团签业务
		if req.CabinetID == nil && req.EnterpriseID != nil && req.StationID != nil {
			batteryInfo, _ = s.CheckBusinessBattery(req, model.AssetLocationsTypeStation, *req.StationID)
		}

		// 电柜业务
		if req.CabinetID != nil {
			batteryInfo, _ = s.CheckBusinessBattery(req, model.AssetLocationsTypeStation, *req.CabinetID)
		}
	}

	switch req.AssetTransferType {
	case model.AssetTransferTypeActive, model.AssetTransferTypeContinue:
		// 激活和取消寄存 某个位置的库存调拨到骑手
		toLocationType = model.AssetLocationsTypeRider
		toLocationID = req.RiderID
		if batteryInfo != nil {
			fromLocationType = model.AssetLocationsType(batteryInfo.LocationsType)
			fromLocationID = batteryInfo.LocationsID
		}
		if ebikeInfo != nil {
			fromLocationType = model.AssetLocationsType(ebikeInfo.LocationsType)
			fromLocationID = ebikeInfo.LocationsID
		}
	case model.AssetTransferTypePause, model.AssetTransferTypeUnSubscribe:
		// 寄存和退租 骑手的库存调拨到某个位置
		fromLocationType = model.AssetLocationsTypeRider
		fromLocationID = req.RiderID
		if req.BatStoreID != nil || req.StoreID != nil {
			toLocationType = model.AssetLocationsTypeStore
			toLocationID = storeID
		}
		if req.CabinetID != nil {
			toLocationType = model.AssetLocationsTypeCabinet
			toLocationID = *req.CabinetID
		}
		if req.EnterpriseID != nil && req.StationID != nil {
			toLocationType = model.AssetLocationsTypeStation
			toLocationID = *req.StationID
		}
	default:
		return nil, nil, errors.New("业务类型错误")
	}

	if ebikeInfo != nil {
		ebikeDetails = append(ebikeDetails, model.AssetTransferCreateDetail{
			AssetType: model.AssetTypeEbike,
			SN:        silk.String(ebikeInfo.Sn),
		})
	}
	if batteryInfo != nil {
		if assetType == model.AssetTypeSmartBattery {
			details = append(details, model.AssetTransferCreateDetail{
				AssetType: assetType,
				SN:        silk.String(batteryInfo.Sn),
			})
		} else {
			details = append(details, model.AssetTransferCreateDetail{
				AssetType: assetType,
				Num:       silk.UInt(1),
				ModelID:   batteryInfo.ModelID,
			})
		}
	}

	autoIn := true
	if s.operator.Type == model.OperatorTypeRider {
		// 骑手操作不自动入库
		autoIn = false
	}

	if len(details) != 0 {
		// 创建调拨单
		_, failed, err := NewAssetTransfer().Transfer(s.ctx, &model.AssetTransferCreateReq{
			FromLocationType:  &fromLocationType,
			FromLocationID:    &fromLocationID,
			ToLocationType:    toLocationType,
			ToLocationID:      toLocationID,
			Details:           details,
			Reason:            req.AssetTransferType.String() + "骑手业务",
			AssetTransferType: req.AssetTransferType,
			OperatorID:        s.operator.ID,
			OperatorType:      s.operator.Type,
			AutoIn:            autoIn,
			SkipLimit:         true,
		}, &model.Modifier{
			ID:    s.operator.ID,
			Name:  s.operator.Name,
			Phone: s.operator.Phone,
		})
		if err != nil {
			return nil, nil, err
		}
		if len(failed) > 0 {
			return nil, nil, errors.New(failed[0])
		}
	}

	if len(ebikeDetails) != 0 {
		// 电车创建调拨单
		_, failed, err := NewAssetTransfer().Transfer(s.ctx, &model.AssetTransferCreateReq{
			FromLocationType:  &fromLocationType,
			FromLocationID:    &fromLocationID,
			ToLocationType:    toLocationType,
			ToLocationID:      toLocationID,
			Details:           ebikeDetails,
			Reason:            req.AssetTransferType.String() + "骑手业务",
			AssetTransferType: req.AssetTransferType,
			OperatorID:        s.operator.ID,
			OperatorType:      s.operator.Type,
			AutoIn:            autoIn,
			SkipLimit:         true,
		}, &model.Modifier{
			ID:    s.operator.ID,
			Name:  s.operator.Name,
			Phone: s.operator.Phone,
		})
		if err != nil {
			return nil, nil, err
		}
		if len(failed) > 0 {
			return nil, nil, errors.New(failed[0])
		}
	}
	return
}

// CheckBusinessBattery 激活和取消寄存非智能电池数量判定
func (s *stockService) CheckBusinessBattery(req *model.StockBusinessReq, locationsType model.AssetLocationsType, locationsID uint64) (batteryInfo *ent.Asset, err error) {
	// 判定非智能电池库存
	batteryInfo, _ = NewAsset().CheckAsset(locationsType, locationsID, req.Model)
	if batteryInfo == nil {
		err = errors.New("电池库存不足")
		return
	}
	if batteryInfo.Edges.Model == nil {
		err = errors.New("电池型号不存在")
		return
	}
	return batteryInfo, nil
}

// func (s *stockService) RiderBusiness(tx *ent.Tx, req *model.StockBusinessReq) (batSk, ebikeSk *ent.Stock, err error) {
// 	num := model.StockNumberOfRiderBusiness(req.StockType)
//
// 	if req.StoreID == nil && req.EbikeStoreID == nil && req.BatStoreID == nil && req.CabinetID == nil && req.EnterpriseID == nil && req.StationID == nil {
// 		err = errors.New("参数校验错误")
// 		return
// 	}
//
// 	creator := tx.Stock.Create()
//
// 	// TODO 平台管理员可操作性时处理出入库逻辑
// 	switch {
// 	case req.StoreID != nil:
// 		creator.SetStoreID(*req.StoreID)
// 		// 判定非智能电池库存
// 		if req.Battery == nil && num < 0 && !s.CheckStore(*req.StoreID, req.Model, int(math.Round(math.Abs(float64(num))))) {
// 			err = errors.New("电池库存不足")
// 			return
// 		}
// 	case req.BatStoreID != nil:
// 		// 退租电池返回门店（退租时电池是查询出来的一定不会为空，所以此处无需上面的判定非智能电池库存）
// 		creator.SetStoreID(*req.BatStoreID)
// 		if req.Battery == nil && num < 0 && !s.CheckStore(*req.StoreID, req.Model, int(math.Round(math.Abs(float64(num))))) {
// 			err = errors.New("电池库存不足")
// 			return
// 		}
// 	}
//
// 	// 判定团签库存
// 	// 未使用电柜激活的时候才需要判定团签站点库存
// 	if req.CabinetID == nil && req.EnterpriseID != nil && req.StationID != nil {
// 		creator.SetEnterpriseID(*req.EnterpriseID).SetStationID(*req.StationID)
// 		// 判断团签非智能电池库存是否足够
// 		if req.Battery == nil && num < 0 && !s.CheckStation(*req.StationID, req.Model, int(math.Round(math.Abs(float64(num))))) {
// 			err = errors.New("电池库存不足")
// 			return
// 		}
// 	}
//
// 	// 主出入库
// 	creator.SetNillableEmployeeID(req.EmployeeID).
// 		SetRiderID(req.RiderID).
// 		SetType(req.StockType).
// 		SetNum(num).
// 		SetCityID(req.CityID).
// 		SetNillableSubscribeID(req.SubscribeID).
// 		SetNillableAgentID(req.AgentID)
//
// 	son := creator.Clone()
//
// 	// 判定电柜库存
// 	if req.CabinetID != nil {
// 		creator.SetCabinetID(*req.CabinetID)
// 		if num < 0 && !s.CheckCabinet(*req.CabinetID, req.Model, int(math.Round(math.Abs(float64(num))))) {
// 			err = errors.New("电池库存不足")
// 			return
// 		}
// 	}
//
// 	var batID *uint64
// 	if req.Battery != nil {
// 		batID = silk.UInt64(req.Battery.ID)
// 	}
// 	batSk, err = creator.SetName(req.Model).
// 		SetModel(req.Model).
// 		SetMaterial(stock.MaterialBattery).
// 		SetSn(tools.NewUnique().NewSN()).
// 		SetNillableBatteryID(batID). // 设置智能电池ID
// 		Save(s.ctx)
// 	if err != nil {
// 		return
// 	}
//
// 	time.Sleep(10 * time.Millisecond)
//
// 	// 当电车需要参与出入库时
// 	if req.Ebike != nil && !req.Rto {
// 		ebq := son.SetParent(batSk).
// 			SetEbikeID(req.Ebike.ID).
// 			SetName(req.Ebike.BrandName).
// 			SetBrandID(req.Ebike.BrandID).
// 			SetMaterial(stock.MaterialEbike).
// 			SetSn(tools.NewUnique().NewSN())
//
// 		if req.EbikeStoreID != nil {
// 			ebq.SetStoreID(*req.EbikeStoreID)
// 		}
//
// 		ebikeSk, err = ebq.Save(s.ctx)
// 		if err != nil {
// 			return
// 		}
// 	}
//
// 	return
// }

// EmployeeOverview 店员物资概览
//
//	func (s *stockService) EmployeeOverview() (res model.StockEmployeeOverview) {
//		st := s.employee.Edges.Store
//		if st == nil {
//			snag.Panic("未上班")
//		}
//
//		start := carbon.Now().StartOfDay().Timestamp()
//
//		res = model.StockEmployeeOverview{
//			Batteries: make([]*model.StockEmployeeOverviewBattery, 0),
//			Materials: make([]*model.StockEmployeeOverviewMaterial, 0),
//		}
//
//		// 计算所有物资
//		batteries := make(map[string]*model.StockEmployeeOverviewBattery)
//		materials := make(map[string]*model.StockEmployeeOverviewMaterial)
//
//		items, _ := s.orm.Query().Where(stock.StoreID(st.ID)).All(s.ctx)
//		for _, item := range items {
//			name := st.Name
//			if item.Model != nil {
//				if _, ok := batteries[name]; !ok {
//					batteries[name] = &model.StockEmployeeOverviewBattery{
//						Outbound: 0,
//						Inbound:  0,
//						Surplus:  0,
//						Model:    *item.Model,
//					}
//				}
//				// 判断是否今日
//				if item.CreatedAt.Unix() > start {
//					if item.Num > 0 {
//						batteries[name].Inbound += item.Num
//					} else {
//						batteries[name].Outbound += int(math.Abs(float64(item.Num)))
//					}
//				}
//				batteries[name].Surplus += item.Num
//			} else {
//				if _, ok := materials[name]; !ok {
//					materials[name] = &model.StockEmployeeOverviewMaterial{
//						Name:    name,
//						Surplus: 0,
//					}
//				}
//				materials[name].Surplus += item.Num
//			}
//		}
//
//		for _, battery := range batteries {
//			res.Batteries = append(res.Batteries, battery)
//		}
//
//		for _, material := range materials {
//			res.Materials = append(res.Materials, material)
//		}
//
//		// 排序
//		sort.Slice(res.Batteries, func(i, j int) bool {
//			return strings.Compare(res.Batteries[i].Model, res.Batteries[j].Model) < 0
//		})
//		sort.Slice(res.Materials, func(i, j int) bool {
//			return strings.Compare(res.Materials[i].Name, res.Materials[j].Name) < 0
//		})
//
//		return
//	}
//
// // listBasicQuery 列表基础查询语句
//
//	func (s *stockService) listBasicQuery(req *model.StockEmployeeListReq) *ent.StockQuery {
//		tt := tools.NewTime()
//
//		q := s.orm.Query().WithRider()
//
//		if req.Outbound {
//			q.Where(stock.NumLT(0))
//		} else {
//			q.Where(stock.NumGT(0))
//		}
//
//		if req.Start != nil {
//			q.Where(stock.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
//		}
//
//		if req.End != nil {
//			q.Where(stock.CreatedAtLTE(tt.ParseNextDateStringX(*req.End)))
//		}
//
//		return q
//	}
//
//	func (s *stockService) EmployeeList(req *model.StockEmployeeListReq) model.StockEmployeeListRes {
//		st := s.employee.Edges.Store
//		if st == nil {
//			snag.Panic("未上班")
//		}
//		q := s.listBasicQuery(req).Where(stock.StoreID(st.ID))
//
//		res := model.ParsePaginationResponse(
//			q,
//			req.PaginationReq,
//			func(item *ent.Stock) model.StockEmployeeListResItem {
//				r := item.Edges.Rider
//				res := model.StockEmployeeListResItem{
//					ID:    item.ID,
//					Type:  item.Type,
//					Model: item.Model,
//					Num:   item.Num,
//					Time:  item.CreatedAt.Format(carbon.DateTimeLayout),
//				}
//				if r != nil {
//					res.Phone = r.Phone
//					res.Name = r.Name
//				}
//				return res
//			},
//		)
//
//		var today *int
//
//		if res.Pagination.Current == 1 {
//			today = new(int)
//			// 获取今日数量
//			var result []struct {
//				ID  uint64 `json:"id"`
//				Sum int    `json:"sum"`
//			}
//			cq := s.orm.Query().Where(
//				stock.CreatedAtGTE(carbon.Now().StartOfDay().StdTime()),
//				stock.StoreID(st.ID),
//			)
//			if req.Outbound {
//				cq.Where(stock.NumLT(0))
//			} else {
//				cq.Where(stock.NumGT(0))
//			}
//			_ = cq.Modify().GroupBy(stock.FieldID).Aggregate(ent.Sum(stock.FieldNum)).Scan(s.ctx, &result)
//
//			if len(result) > 0 {
//				today = &result[0].Sum
//			}
//		}
//
//		return model.StockEmployeeListRes{
//			Today:         today,
//			PaginationRes: res,
//		}
//	}
//
// // CabinetList 电柜物资列表
//
//	func (s *stockService) CabinetList(req *model.StockCabinetListReq) *model.PaginationRes {
//		q := ent.Database.Cabinet.QueryNotDeleted().
//			Where(cabinet.HasStocks()).
//			WithCity().
//			WithStocks()
//
//		if req.CabinetID != 0 {
//			q.Where(cabinet.ID(req.CabinetID))
//		}
//		if req.Serial != "" {
//			q.Where(cabinet.Serial(req.Serial))
//		}
//		if req.CityID != 0 {
//			q.Where(cabinet.CityID(req.CityID))
//		}
//		if req.Start != "" {
//			q.Where(cabinet.CreatedAtGTE(tools.NewTime().ParseDateStringX(req.Start)))
//		}
//		if req.End != "" {
//			q.Where(cabinet.CreatedAtLT(tools.NewTime().ParseNextDateStringX(req.End)))
//		}
//
//		return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Cabinet) model.StockCabinetListRes {
//			res := model.StockCabinetListRes{
//				ID:        item.ID,
//				Serial:    item.Serial,
//				Name:      item.Name,
//				Batteries: make([]*model.StockMaterial, 0),
//			}
//			c := item.Edges.City
//			if c != nil {
//				res.City = model.City{
//					ID:   c.ID,
//					Name: c.Name,
//				}
//			}
//			batteries := make(map[string]*model.StockMaterial)
//
//			// 出入库
//			for _, st := range item.Edges.Stocks {
//				s.Calculate(batteries, st)
//			}
//
//			for _, battery := range batteries {
//				res.Batteries = append(res.Batteries, battery)
//			}
//
//			return res
//		})
//	}
func (s *stockService) listFilter(req model.StockDetailFilter) (q *ent.StockQuery, info ar.Map) {
	info = make(ar.Map)

	q = s.orm.Query().
		Modify(func(sel *sql.Selector) {
			// 不做对象查询时需要去重排除配偶
			if req.Goal == 0 && req.Serial == "" && req.StoreID == 0 && req.CabinetID == 0 && req.EnterpriseID == 0 && req.StationID == 0 {
				sel.FromExpr(sql.Raw("(SELECT DISTINCT ON (id + COALESCE(stock_spouse, 0)) * FROM stock) stock"))
			}
		}).
		WithCabinet().
		WithStore().
		WithSpouse(func(sq *ent.StockQuery) {
			sq.WithStore().WithCabinet().WithRider().WithStation().WithBattery().WithEnterprise()
		}).
		WithRider().
		WithEmployee().
		WithCity().
		WithEbike().
		WithStation().
		WithBattery().
		WithEnterprise().
		WithAgent()
	// 排序
	if req.Positive {
		q.Order(ent.Asc(stock.FieldSn))
	} else {
		q.Order(ent.Desc(stock.FieldSn))
	}

	if req.Start != "" {
		info["开始时间"] = req.Start
		q.Where(stock.CreatedAtGTE(tools.NewTime().ParseDateStringX(req.Start)))
	}

	if req.End != "" {
		info["结束时间"] = req.Start
		q.Where(stock.CreatedAtLT(tools.NewTime().ParseNextDateStringX(req.End)))
	}

	if req.CityID != 0 {
		info["城市"] = ent.NewExportInfo(req.CityID, city.Table)
		q.Where(stock.CityID(req.CityID))
	}

	if req.Serial != "" {
		info["电柜编号"] = req.Serial
		q.Where(stock.HasCabinetWith(cabinet.Serial(req.Serial)))
	}

	if req.EbikeSN != "" {
		info["车架号"] = req.EbikeSN
		q.Where(stock.HasEbikeWith(ebike.Sn(req.EbikeSN)))
	}

	if req.BatterySN != "" {
		info["电池编码"] = req.BatterySN
		q.Where(stock.HasBatteryWith(battery.Sn(req.BatterySN)))
	}

	switch req.Goal {
	case model.StockGoalStore:
		// 门店物资
		info["查询目标"] = "门店"
		q.Where(
			stock.StoreIDNotNil(),
			stock.CabinetIDIsNil(),
			stock.StationIDIsNil(),
		)
	case model.StockGoalCabinet:
		// 电柜物资
		info["查询目标"] = "电柜"
		q.Where(
			stock.StoreIDIsNil(),
			stock.CabinetIDNotNil(),
			stock.StationIDIsNil(),
		)
	case model.StockGoalStation:
		// 站点
		info["查询目标"] = "站点"
		q.Where(
			stock.StoreIDIsNil(),
			stock.CabinetIDIsNil(),
			stock.StationIDNotNil(),
		)
	default:
		// 门店或电柜物资或站点
		info["查询目标"] = "电柜或门店或站点"
		q.Where(
			stock.Or(
				stock.StoreIDNotNil(),
				stock.CabinetIDNotNil(),
				stock.StationIDNotNil(),
			),
		)
	}

	// 筛选物资类别
	if req.Materials == "" {
		req.Materials = fmt.Sprintf("%s,%s", stock.MaterialBattery, stock.MaterialEbike)
	} else {
		req.Materials = strings.ReplaceAll(req.Materials, " ", "")
	}
	materials := strings.Split(req.Materials, ",")

	if len(materials) > 0 {
		var mtext []string
		var predicates []predicate.Stock
		for _, material := range materials {
			switch stock.Material(material) {
			case stock.MaterialBattery:
				mtext = append(mtext, "电池")
				predicates = append(predicates, stock.ModelNotNil())
			case stock.MaterialEbike:
				mtext = append(mtext, "电车")
				predicates = append(predicates, stock.EbikeIDNotNil())
			case stock.MaterialOthers:
				mtext = append(mtext, "其他")
				predicates = append(predicates, stock.ModelIsNil())
			}
		}
		info["物资"] = strings.Join(mtext, ",")
		q.Where(stock.Or(predicates...))
	}

	if req.Type != 0 {
		info["类型"] = model.StockTypesText[req.Type]
		q.Where(stock.Type(req.Type))
	}

	if req.StoreID != 0 {
		info["门店"] = ent.NewExportInfo(req.StoreID, store.Table)
		q.Where(stock.StoreID(req.StoreID))
	}

	if req.CabinetID != 0 {
		info["电柜"] = ent.NewExportInfo(req.CabinetID, cabinet.Table)
		q.Where(stock.CabinetID(req.CabinetID))
	}

	if req.EnterpriseID != 0 {
		info["团签"] = ent.NewExportInfo(req.EnterpriseID, enterprise.Table)
		q.Where(stock.EnterpriseID(req.EnterpriseID))
	}
	if req.StationID != 0 {
		info["站点"] = ent.NewExportInfo(req.StationID, enterprisestation.Table)
		q.Where(stock.StationID(req.StationID))
	}

	return
}

// detailInfo 库存出入明细信息
func (s *stockService) detailInfo(item *ent.Stock) model.StockDetailRes {
	res := model.StockDetailRes{
		ID:     item.ID,
		Sn:     item.Sn,
		Name:   item.Name,
		Num:    int(math.Abs(float64(item.Num))),
		Time:   item.CreatedAt.Format(carbon.DateTimeLayout),
		Remark: item.Remark,
	}

	// 城市
	c := item.Edges.City
	if c != nil {
		res.City = c.Name
	}

	// 电车
	bike := item.Edges.Ebike
	if bike != nil {
		res.Name = fmt.Sprintf("[%s] %s", item.Name, bike.Sn)
	}

	em := item.Creator
	er := item.Edges.Rider
	ee := item.Edges.Employee
	es := item.Edges.Store
	ec := item.Edges.Cabinet
	en := item.Edges.Enterprise
	st := item.Edges.Station
	ba := item.Edges.Battery
	ag := item.Edges.Agent

	// 站点调拨电池
	if ba != nil {
		res.Name = fmt.Sprintf("[%s] %s", *item.Model, ba.Sn)
	}
	if item.Type == model.StockTypeTransfer {
		// 平台调拨记录
		res.Type = "平台调拨"
		res.Operator = "后台"
		if em != nil {
			res.Operator = fmt.Sprintf("后台 - %s", em.Name)
		}
		var ses *ent.Store
		var sec *ent.Cabinet
		var sst *ent.EnterpriseStation
		var sen *ent.Enterprise
		sp := item.Edges.Spouse
		if sp != nil {
			ses = sp.Edges.Store
			sec = sp.Edges.Cabinet
			sst = sp.Edges.Station
			sen = sp.Edges.Enterprise
		}

		// 出入库对象判定
		if item.Num > 0 {
			res.Inbound = s.target(es, ec, st, en)
			res.Outbound = s.target(ses, sec, sst, sen)
		} else {
			res.Inbound = s.target(ses, sec, sst, sen)
			res.Outbound = s.target(es, ec, st, en)
		}
	} else {
		// 业务调拨记录
		var riderName string
		var agentName string

		if er != nil {
			riderName = er.Name
			res.Rider = fmt.Sprintf("%s - %s", riderName, er.Phone)
		}

		if ag != nil {
			agentName = ag.Name
		}

		tm := map[uint8]string{
			model.StockTypeRiderActive:      "新签",
			model.StockTypeRiderPause:       "寄存",
			model.StockTypeRiderContinue:    "取消寄存",
			model.StockTypeRiderUnSubscribe: "退租",
		}

		var tmr string
		if ec != nil {
			res.Operator = fmt.Sprintf("骑手 - %s", riderName)
			tmr = "电柜"
		} else {
			if ee != nil {
				tmr = "门店"
				res.Operator = fmt.Sprintf("店员 - %s", ee.Name)
			} else if item.Creator != nil {
				tmr = "后台"
				res.Operator = fmt.Sprintf("后台 - %s", item.Creator.Name)
			} else if st != nil {
				tmr = "站点"
				res.Operator = fmt.Sprintf("代理 - %s", agentName)
			}
		}

		res.Type = tmr + tm[item.Type]

		// 出入库对象
		target := fmt.Sprintf("[骑手] %s", riderName)
		switch item.Type {
		case model.StockTypeRiderActive, model.StockTypeRiderContinue:
			res.Inbound = target
			res.Outbound = s.target(es, ec, st, en)
		case model.StockTypeRiderPause, model.StockTypeRiderUnSubscribe:
			res.Inbound = s.target(es, ec, st, en)
			res.Outbound = target
		}
	}

	return res
}

// target 出入库对象
func (s *stockService) target(es *ent.Store, ec *ent.Cabinet, st *ent.EnterpriseStation, en *ent.Enterprise) (target string) {
	target = "平台"
	if es != nil {
		target = fmt.Sprintf("[门店] %s", es.Name)
	}
	if ec != nil {
		target = fmt.Sprintf("[电柜] %s - %s", ec.Name, ec.Serial)
	}
	if st != nil {
		target = fmt.Sprintf("[代理商站点] %s - %s", en.Name, st.Name)
	}
	return
}

// Detail 出入库明细
func (s *stockService) Detail(req *model.StockDetailReq) *model.PaginationRes {
	q, _ := s.listFilter(req.StockDetailFilter)

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Stock) model.StockDetailRes {
		return s.detailInfo(item)
	})
}

//
// // StoreCurrent 列出当前门店所有库存物资
// func (s *stockService) StoreCurrent(id uint64) []model.InventoryNum {
// 	ins := make([]model.InventoryNum, 0)
// 	_ = s.orm.Query().
// 		Where(stock.StoreID(id)).
// 		Modify(func(sel *sql.Selector) {
// 			sel.GroupBy(stock.FieldName, stock.FieldModel).
// 				Select(stock.FieldName, stock.FieldModel).
// 				AppendSelectExprAs(sql.Raw(fmt.Sprintf("%s IS NOT NULL", stock.FieldModel)), "battery").
// 				AppendSelectExprAs(sql.Raw(fmt.Sprintf("SUM(%s)", stock.FieldNum)), "num")
// 		}).
// 		Scan(s.ctx, &ins)
//
// 	return ins
// }
//
// func (s *stockService) StoreCurrentMap(id uint64) map[string]model.InventoryNum {
// 	inm := make(map[string]model.InventoryNum)
// 	for _, in := range s.StoreCurrent(id) {
// 		inm[in.Name] = in
// 	}
// 	return inm
// }
//
// // CurrentBatteryNum 获取当前电池库存总数
// func (s *stockService) CurrentBatteryNum(ids []uint64, field string) map[uint64]int {
// 	var result []struct {
// 		TargetID uint64 `json:"target_id"`
// 		Sum      int    `json:"sum"`
// 	}
// 	v := make([]interface{}, len(ids))
// 	for i := range v {
// 		v[i] = ids[i]
// 	}
// 	_ = s.orm.Query().
// 		Modify(func(sel *sql.Selector) {
// 			sel.Where(sql.In(sel.C(field), v...)).
// 				Select(
// 					sql.As(sel.C(field), "target_id"),
// 					sql.As(sql.Sum(stock.FieldNum), "sum"),
// 				).
// 				GroupBy(field)
// 		}).
// 		Scan(s.ctx, &result)
// 	m := make(map[uint64]int)
// 	for _, r := range result {
// 		m[r.TargetID] = r.Sum
// 	}
// 	return m
// }
//
// func (s *stockService) CurrentBattery(id uint64, field string) int {
// 	return s.CurrentBatteryNum([]uint64{id}, field)[id]
// }
//
// // Inventory 查询所有物资
// func (s *stockService) Inventory(req *model.StockInventoryReq) (items []model.StockInventory) {
// 	_ = s.orm.Query().
// 		Modify(func(sel *sql.Selector) {
// 			sel.Select(stock.FieldCabinetID, stock.FieldStoreID, stock.FieldName, stock.FieldMaterial).
// 				AppendSelectExprAs(sql.Raw(fmt.Sprintf("SUM(%s)", stock.FieldNum)), "num").
// 				GroupBy(stock.FieldCabinetID, stock.FieldStoreID, stock.FieldName, stock.FieldMaterial)
// 			// 如果请求参数为空则查询全部门店和电柜的全部物资
// 			if req == nil {
// 				sel.Where(sql.Or(
// 					sql.NotNull(stock.FieldCabinetID),
// 					sql.NotNull(stock.FieldStoreID),
// 				))
// 			} else {
// 				if req.Material != "" {
// 					sel.Where(sql.EQ(stock.FieldMaterial, req.Material))
// 				}
// 				if req.Goal != model.StockGoalAll {
// 					col := stock.FieldStoreID
// 					if req.Goal == model.StockGoalCabinet {
// 						col = stock.FieldCabinetID
// 					}
// 					if len(req.IDs) > 0 {
// 						ids := make([]any, len(req.IDs))
// 						for i, d := range req.IDs {
// 							ids[i] = d
// 						}
// 						sel.Where(sql.In(col, ids...))
// 					} else {
// 						sel.Where(sql.NotNull(col))
// 					}
// 					if req.Name != "" {
// 						sel.Where(sql.EQ(stock.FieldName, req.Name))
// 					}
// 				}
// 			}
//
// 		}).
// 		Scan(s.ctx, &items)
// 	return
// }
//
// func (s *stockService) InventoryMap(req *model.StockInventoryReq) (data model.StockInventoryMapData) {
// 	data = make(map[uint64]map[string]map[string]model.StockInventory)
// 	for _, item := range s.Inventory(req) {
// 		id := item.CabinetID
// 		if id == 0 {
// 			id = item.StoreID
// 		}
// 		if _, ok := data[id]; !ok {
// 			data[id] = make(map[string]map[string]model.StockInventory)
// 		}
// 		if _, ok := data[id][item.Material]; !ok {
// 			data[id][item.Material] = make(map[string]model.StockInventory)
// 		}
// 		data[id][item.Material][item.Name] = item
// 	}
// 	return
// }
//
// // Export 出入库明细导出
// func (s *stockService) Export(req *model.StockDetailExportReq) model.ExportRes {
// 	q, info := s.listFilter(req.StockDetailFilter)
//
// 	return NewExportWithModifier(s.modifier).Start("出入库明细", req.StockDetailFilter, info, req.Remark, func(path string) {
// 		items, _ := q.All(s.ctx)
// 		var rows tools.ExcelItems
// 		title := []any{
// 			"编号",
// 			"城市",
// 			"调出",
// 			"调入",
// 			"物资",
// 			"数量",
// 			"类型",
// 			"操作人",
// 			"骑手",
// 			"备注",
// 			"操作时间",
// 		}
// 		rows = append(rows, title)
//
// 		for _, item := range items {
// 			detail := s.detailInfo(item)
// 			rows = append(rows, []any{
// 				detail.Sn,
// 				detail.City,
// 				detail.Outbound,
// 				detail.Inbound,
// 				detail.Name,
// 				detail.Num,
// 				detail.Type,
// 				detail.Operator,
// 				detail.Rider,
// 				detail.Remark,
// 				detail.Time,
// 			})
// 		}
//
// 		tools.NewExcel(path).AddValues(rows).Done()
// 	})
// }
//
// // Transfer 调拨物资
// // 此方法操作库存
// func (s *stockService) Transfer(req *model.StockTransferReq) (failed []string) {
// 	err := req.Validate()
// 	if err != nil {
// 		snag.Panic(err)
// 	}
//
// 	failed = make([]string, 0)
//
// 	var cityID uint64
//
// 	// 查询电柜
// 	var cab *ent.Cabinet
// 	var cabID uint64
//
// 	// 检查电柜是否初始化调拨过
// 	if req.InboundTarget == model.StockTargetCabinet {
// 		cabID = req.InboundID
// 	}
// 	if req.OutboundTarget == model.StockTargetCabinet {
// 		cabID = req.OutboundID
// 	}
// 	if cabID > 0 {
// 		cab = NewCabinet().QueryOne(cabID)
// 		if !cab.Transferred && !req.Force {
// 			snag.Panic("电柜未初始化调拨")
// 		}
// 		if cab.CityID != nil {
// 			cityID = *cab.CityID
// 		}
// 	}
//
// 	// 查询门店
// 	var st *ent.Store
// 	var stID uint64
// 	if req.InboundTarget == model.StockTargetStore {
// 		stID = req.InboundID
// 	}
// 	if req.OutboundTarget == model.StockTargetStore {
// 		stID = req.OutboundID
// 	}
// 	if stID > 0 {
// 		st = NewStore().Query(stID)
// 		cityID = st.CityID
// 	}
// 	if cab != nil && cab.CityID != nil && st != nil && st.CityID != *cab.CityID {
// 		snag.Panic("不同城市电柜和门店无法调拨")
// 	}
//
// 	var stationID uint64
// 	var stationInfo *ent.EnterpriseStation
// 	var enterpriseId uint64
// 	// 团签调拨
// 	if req.InboundTarget == model.StockTargetStation {
// 		stationID = req.InboundID
// 	}
// 	if req.OutboundTarget == model.StockTargetStation {
// 		stationID = req.OutboundID
// 	}
// 	if stationID > 0 {
// 		stationInfo, _ = ent.Database.EnterpriseStation.QueryNotDeleted().WithEnterprise().Where(enterprisestation.ID(stationID)).First(s.ctx)
// 		if stationInfo == nil {
// 			snag.Panic("电池调拨站点不存在")
// 		}
// 		if stationInfo.Edges.Enterprise != nil {
// 			enterpriseId = stationInfo.Edges.Enterprise.ID
// 			cityID = stationInfo.Edges.Enterprise.CityID
// 		}
// 	}
//
// 	in := &req.InboundID
// 	if req.InboundID == 0 {
// 		in = nil
// 	}
//
// 	out := &req.OutboundID
// 	if req.OutboundID == 0 {
// 		out = nil
// 	}
//
// 	batteryModel := &req.Model
// 	if req.Model == "" {
// 		batteryModel = nil
// 	}
//
// 	num := req.RealNumber()
// 	name := req.RealName()
// 	batchable := req.Batchable()
//
// 	// 批量调拨, 调出检查
// 	// 跳过电车
// 	if req.OutboundID > 0 && len(req.Ebikes) == 0 && len(req.BatterySn) == 0 && NewStockBatchable().Fetch(req.OutboundTarget, req.OutboundID, name) < req.Num {
// 		snag.Panic("操作失败, 调出物资大于库存物资")
// 	}
//
// 	material := func(req *model.StockTransferReq) stock.Material {
// 		switch true {
// 		case len(req.Ebikes) > 0:
// 			return stock.MaterialEbike
// 		case req.Model != "" || len(req.BatterySn) > 0:
// 			return stock.MaterialBattery
// 		}
// 		return stock.MaterialOthers
// 	}(req)
//
// 	outCreator := s.orm.Create().
// 		SetNum(-num).
// 		SetCityID(cityID).
// 		SetType(model.StockTypeTransfer).
// 		SetMaterial(material).
// 		SetRemark(req.Remark).
// 		SetSn(tools.NewUnique().NewSN())
// 	switch req.OutboundTarget {
// 	case model.StockTargetStore:
// 		outCreator.SetNillableStoreID(out)
// 	case model.StockTargetCabinet:
// 		outCreator.SetNillableCabinetID(out)
// 	case model.StockTargetStation:
// 		outCreator.SetNillableStationID(out)
// 		outCreator.SetEnterpriseID(enterpriseId)
// 	}
//
// 	inCreator := s.orm.Create().
// 		SetNum(num).
// 		SetCityID(cityID).
// 		SetType(model.StockTypeTransfer).
// 		SetMaterial(material).
// 		SetRemark(req.Remark).
// 		SetSn(tools.NewUnique().NewSN())
// 	switch req.InboundTarget {
// 	case model.StockTargetStore:
// 		inCreator.SetNillableStoreID(in)
// 	case model.StockTargetCabinet:
// 		inCreator.SetNillableCabinetID(in)
// 	case model.StockTargetStation:
// 		inCreator.SetNillableStationID(in)
// 		inCreator.SetEnterpriseID(enterpriseId)
// 	}
//
// 	var looppers []model.StockTransferLoopper
//
// 	switch true {
// 	case len(req.Ebikes) > 0:
// 		// failed = append(failed, NewStockEbike(s.modifier, s.employee, s.rider).Transfer(cityID, in, out, req)...)
// 		looppers, failed = NewStockEbike().Loopers(req, enterpriseId)
// 	case len(req.BatterySn) > 0:
// 		looppers, failed = NewStockBatchable().Loopers(req, enterpriseId)
// 	default:
// 		looppers = []model.StockTransferLoopper{
// 			{BatteryModel: batteryModel},
// 		}
// 	}
//
// 	for _, l := range looppers {
// 		err = ent.WithTx(s.ctx, func(tx *ent.Tx) (err error) {
// 			// 判定名称
// 			if l.BrandName != nil {
// 				name = *l.BrandName
// 			}
//
// 			// 调出
// 			var spouse *ent.Stock
// 			spouse, err = outCreator.
// 				SetName(name).
// 				SetNillableEbikeID(l.EbikeID).
// 				SetNillableBrandID(l.BrandID).
// 				SetNillableModel(l.BatteryModel).
// 				SetNillableBatteryID(l.BatteryID).
// 				Save(s.ctx)
// 			if err != nil {
// 				if batchable {
// 					return
// 				}
// 				return fmt.Errorf("出库失败: %s", l.Message)
// 			}
//
// 			// 调入
// 			_, err = inCreator.
// 				SetName(name).
// 				SetNillableEbikeID(l.EbikeID).
// 				SetNillableBrandID(l.BrandID).
// 				SetSpouse(spouse).
// 				SetNillableModel(l.BatteryModel).
// 				SetNillableBatteryID(l.BatteryID).
// 				Save(s.ctx)
// 			if err != nil {
// 				if batchable {
// 					return
// 				}
// 				return fmt.Errorf("入库失败: %s", l.Message)
// 			}
//
// 			// 电车调拨完成更新所属
// 			if l.EbikeID != nil {
// 				// 是否可调拨检查
// 				if exists, _ := NewEbike().AllocatableBaseFilter().Where(ebike.ID(*l.EbikeID)).Exist(s.ctx); !exists {
// 					return fmt.Errorf("电车无法调拨: %s", l.Message)
// 				}
//
// 				updater := tx.Ebike.UpdateOneID(*l.EbikeID)
// 				// 调拨到门店
// 				if req.IsToStore() {
// 					updater.SetNillableStoreID(in)
// 				}
// 				// 调拨到平台 清空所有关联
// 				if req.IsToPlaform() {
// 					updater.ClearStoreID().ClearStationID().ClearEnterpriseID()
// 				}
// 				// 调拨到站点
// 				if req.IsToStation() {
// 					updater.SetNillableStationID(in)
// 					updater.SetEnterpriseID(enterpriseId)
// 				}
//
// 				if updater.Exec(s.ctx) != nil {
// 					return fmt.Errorf("电车更新失败: %s", l.Message)
// 				}
// 			}
//
// 			if l.BatteryID != nil {
// 				updater := tx.Battery.UpdateOneID(*l.BatteryID)
// 				switch req.InboundTarget {
// 				case model.StockTargetStation:
// 					// 调拨到站点
// 					updater.SetNillableStationID(in).SetEnterpriseID(enterpriseId)
// 				case model.StockTargetPlaform:
// 					// 调拨到平台
// 					updater.ClearStationID().ClearEnterpriseID()
// 				}
// 				if updater.Exec(s.ctx) != nil {
// 					return fmt.Errorf("电池更新失败: %s", *l.BatterySN)
// 				}
// 			}
//
// 			return
// 		})
//
// 		if req.Batchable() && err != nil {
// 			failed = append(failed, err.Error())
// 		}
// 	}
//
// 	return
// }
