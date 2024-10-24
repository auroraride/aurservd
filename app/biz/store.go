package biz

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/storegoods"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
)

type storeBiz struct {
	orm *ent.StoreClient
	ctx context.Context
}

func NewStore() *storeBiz {
	return &storeBiz{
		orm: ent.Database.Store,
		ctx: context.Background(),
	}
}

// List 门店列表
func (b *storeBiz) List(req *definition.StoreListReq) (res []*definition.StoreDetail, err error) {
	res = make([]*definition.StoreDetail, 0)
	maxDistance := definition.DefaultMaxDistance
	if req.Distance != nil && *req.Distance < maxDistance {
		maxDistance = *req.Distance
	}
	q := b.orm.QueryNotDeleted().
		WithCity().
		WithEmployee().
		WithAsset().
		Modify(func(sel *sql.Selector) {
			sel.
				AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(ST_GeographyFromText('SRID=4326;POINT(' || "store"."lng" || ' ' || "store"."lat" || ')'),ST_GeographyFromText('SRID=4326;POINT(%f  %f)'))`, req.Lng, req.Lat)), "distance").
				OrderBy(sql.Asc("distance")).
				Where(sql.P(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`ST_DWithin(ST_GeographyFromText('SRID=4326;POINT(' || "store"."lng" || ' ' || "store"."lat" || ')'),ST_GeographyFromText('SRID=4326;POINT(%f  %f)'), %f)`, req.Lng, req.Lat, maxDistance))
				}))
		})

	if req.CityID != nil {
		q.Where(store.CityID(*req.CityID))
	}

	// 门店只需要查询营业和休息中两种状态
	q.Where(store.StatusIn(model.StoreStatusOpen.Value(), model.StoreStatusClose.Value()))

	if req.Status != nil {
		q.Where(store.Status(req.Status.Value()))
	}

	if req.BusinessType != nil {
		switch *req.BusinessType {
		case model.StoreBusinessTypeObtain:
			q.Where(store.EbikeObtain(true))
		case model.StoreBusinessTypeRepair:
			q.Where(store.EbikeRepair(true))
		case model.StoreBusinessTypeSale:
			q.Where(store.EbikeSale(true))
		case model.StoreBusinessTypeRest:
			q.Where(store.Rest(true))
		}
	}

	if req.Keyword != nil {
		q.Where(store.NameContains(*req.Keyword))
	}

	if req.GoodsID != nil {
		q.Where(store.HasGoodsWith(storegoods.GoodsID(*req.GoodsID)))
	}

	list, _ := q.All(b.ctx)
	if len(list) == 0 {
		return res, nil
	}

	var pl *ent.Plan
	if req.PlanID != nil {
		pl, _ = ent.Database.Plan.QueryNotDeleted().Where(plan.ID(*req.PlanID)).First(b.ctx)
		if pl == nil {
			return nil, errors.New("未找到有效套餐")
		}

		if pl.Type == model.PlanTypeBattery.Value() {
			return nil, errors.New("套餐类型错误")
		}

	}
	for _, v := range list {
		// 查询门店库存
		if req.PlanID != nil && pl != nil {
			ebikeNum, batteryNum := b.QueryStocks(v, pl)
			// 电车或电池库存为0则不展示
			if ebikeNum <= 0 || batteryNum <= 0 {
				continue
			}
		}
		res = append(res, b.detail(v))
	}
	return
}

// Detail 门店详情
func (b *storeBiz) Detail(req *definition.StoreDetailReq) (res *definition.StoreDetail) {
	q, _ := b.orm.QueryNotDeleted().
		Where(store.ID(req.ID)).
		WithCity().
		WithEmployee().
		WithAsset().
		Modify(func(sel *sql.Selector) {
			sel.AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(ST_GeographyFromText('SRID=4326;POINT(' || "store"."lng" || ' ' || "store"."lat" || ')'),ST_GeographyFromText('SRID=4326;POINT(%f  %f)'))`, req.Lng, req.Lat)), "distance").
				OrderBy(sql.Asc("distance"))
		}).
		First(b.ctx)
	if q == nil {
		return nil
	}
	return b.detailForStock(q)
}

func (b *storeBiz) detail(item *ent.Store) (res *definition.StoreDetail) {
	res = &definition.StoreDetail{
		ID:            item.ID,
		Name:          item.Name,
		Status:        model.StoreStatus(item.Status),
		Lng:           item.Lng,
		Lat:           item.Lat,
		Address:       item.Address,
		EbikeRepair:   item.EbikeRepair,
		EbikeObtain:   item.EbikeObtain,
		EbikeSale:     item.EbikeSale,
		BusinessHours: item.BusinessHours,
		Rest:          item.Rest,
		Photos:        item.Photos,
		Phone:         item.Phone,
		HeadPic:       item.HeadPic,
	}
	if item.Edges.Employee != nil {
		res.Employee = &model.Employee{
			ID:    item.Edges.Employee.ID,
			Name:  item.Edges.Employee.Name,
			Phone: item.Edges.Employee.Phone,
		}
	}
	if item.Edges.City != nil {
		res.City = model.City{
			ID:   item.Edges.City.ID,
			Name: item.Edges.City.Name,
		}
	}

	distance, err := item.Value("distance")
	if distance != nil || err == nil {
		distanceFloat, ok := distance.(float64)
		if ok {
			res.Distance = distanceFloat
		}
	}
	return

}

// QueryStocks 查询门店电池 车库存
func (b *storeBiz) QueryStocks(item *ent.Store, pl *ent.Plan) (ebikeNum, batteryNum int) {
	for _, st := range item.Edges.Asset {
		switch model.AssetType(st.Type) {
		case model.AssetTypeEbike:
			if st.Status == model.AssetStatusStock.Value() && st.BrandID != nil && *st.BrandID == *pl.BrandID {
				ebikeNum += 1
			}
		case model.AssetTypeSmartBattery, model.AssetTypeNonSmartBattery:
			m, _ := st.QueryModel().First(b.ctx)
			if st.Status == model.AssetStatusStock.Value() && m != nil && m.Model == pl.Model {
				batteryNum += 1
			}
		default:
		}
	}
	return
}

// StoreBySubscribe 根据订阅查询门店
func (b *storeBiz) StoreBySubscribe(r *ent.Rider, req *definition.StoreDetailReq) (res *definition.StoreDetail, err error) {
	q, _ := ent.Database.Subscribe.QueryNotDeleted().
		Where(subscribe.ID(req.ID), subscribe.RiderID(r.ID)).
		WithStore(
			func(query *ent.StoreQuery) {
				query.WithCity().WithEmployee()
			}).Modify(func(sel *sql.Selector) {
		t := sql.Table(store.Table).As("store")
		sel.LeftJoin(t).On(t.C(store.FieldID), sel.C(subscribe.FieldStoreID))
		sel.
			AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(ST_GeographyFromText('SRID=4326;POINT(' || "store"."lng" || ' ' || "store"."lat" || ')'),ST_GeographyFromText('SRID=4326;POINT(%f  %f)'))`, req.Lng, req.Lat)), "distance").
			OrderBy(sql.Asc("distance"))
	}).First(b.ctx)
	if q == nil || q.Edges.Store == nil {
		return nil, errors.New("未找到门店")
	}

	if q.Edges.Store != nil {
		item := q.Edges.Store
		res = &definition.StoreDetail{
			ID:            item.ID,
			Name:          item.Name,
			Status:        model.StoreStatus(item.Status),
			Lng:           item.Lng,
			Lat:           item.Lat,
			Address:       item.Address,
			EbikeRepair:   item.EbikeRepair,
			EbikeObtain:   item.EbikeObtain,
			EbikeSale:     item.EbikeSale,
			BusinessHours: item.BusinessHours,
		}
		if item.Edges.Employee != nil {
			res.Employee = &model.Employee{
				ID:    item.Edges.Employee.ID,
				Name:  item.Edges.Employee.Name,
				Phone: item.Edges.Employee.Phone,
			}
		}
		if item.Edges.City != nil {
			res.City = model.City{
				ID:   item.Edges.City.ID,
				Name: item.Edges.City.Name,
			}
		}
	}
	var distance ent.Value
	distance, err = q.Value("distance")
	if distance != nil || err == nil {
		distanceFloat, ok := distance.(float64)
		if ok {
			res.Distance = distanceFloat
		}
	}
	return res, nil
}

// queryStocksByStore 查询门店电车是否有库存
func (b *storeBiz) queryStocksByStore(item *ent.Store, brandIds []uint64) (eBrandIds []uint64) {
	brandIdMap := make(map[uint64]bool)
	for _, brandId := range brandIds {
		brandIdMap[brandId] = true
	}

	for _, st := range item.Edges.Asset {
		if st.Type == model.AssetTypeEbike.Value() && st.Status == model.AssetStatusStock.Value() &&
			st.BrandID != nil && brandIdMap[*st.BrandID] {
			eBrandIds = append(eBrandIds, *st.BrandID)
		}
	}
	return
}

// detailForStock 插叙门店详情及商品、电车套餐列表
func (b *storeBiz) detailForStock(item *ent.Store) (res *definition.StoreDetail) {
	res = &definition.StoreDetail{
		ID:            item.ID,
		Name:          item.Name,
		Status:        model.StoreStatus(item.Status),
		Lng:           item.Lng,
		Lat:           item.Lat,
		Address:       item.Address,
		EbikeRepair:   item.EbikeRepair,
		EbikeObtain:   item.EbikeObtain,
		EbikeSale:     item.EbikeSale,
		BusinessHours: item.BusinessHours,
		Rest:          item.Rest,
		Photos:        item.Photos,
		Phone:         item.Phone,
	}
	if item.Edges.Employee != nil {
		res.Employee = &model.Employee{
			ID:    item.Edges.Employee.ID,
			Name:  item.Edges.Employee.Name,
			Phone: item.Edges.Employee.Phone,
		}
	}
	if item.Edges.City != nil {
		res.City = model.City{
			ID:   item.Edges.City.ID,
			Name: item.Edges.City.Name,
		}
	}

	distance, err := item.Value("distance")
	if distance != nil || err == nil {
		distanceFloat, ok := distance.(float64)
		if ok {
			res.Distance = distanceFloat
		}
	}

	if item.Edges.Asset != nil {
		var brandIds []uint64
		for _, st := range item.Edges.Asset {
			if st.BrandID != nil {
				brandIds = append(brandIds, *st.BrandID)
			}
		}

		// 新租车列表
		res.StoreBrands = NewPlanBiz().ListByStoreById(item.ID)

	}

	// 查询门店商品
	res.SaleGoods = NewGoods().ListByStoreId(item.ID)

	return

}

// SelectionList 门店列表筛选项
func (b *storeBiz) SelectionList() (res []*model.SelectOption) {
	res = make([]*model.SelectOption, 0)
	list, _ := b.orm.QueryNotDeleted().Order(ent.Asc(store.FieldID)).All(b.ctx)

	for _, item := range list {
		res = append(res, &model.SelectOption{
			Value: item.ID,
			Label: item.Name,
		})
	}
	return
}

// ListByEmployee 店员已配置门店列表
func (b *storeBiz) ListByEmployee(ep *ent.Employee) (res []*model.CascaderOptionLevel2) {
	res = make([]*model.CascaderOptionLevel2, 0)
	if ep == nil {
		return
	}

	// 查询人员配置的仓库门店信息
	eep, err := ent.Database.Employee.QueryNotDeleted().
		Where(employee.ID(ep.ID)).
		WithStores(func(query *ent.StoreQuery) {
			query.Where(store.DeletedAtIsNil()).WithCity()
		}).First(b.ctx)
	if err != nil || eep == nil {
		return
	}

	// 数据组合
	stList := eep.Edges.Stores
	cityIds := make([]uint64, 0)
	cityIdMap := make(map[uint64]*ent.City)
	cityIdListMap := make(map[uint64][]model.SelectOption)
	for _, st := range stList {
		if st.Edges.City != nil {
			cId := st.Edges.City.ID
			if cityIdMap[cId] == nil {
				cityIds = append(cityIds, cId)
				cityIdMap[cId] = st.Edges.City
			}

			cityIdListMap[cId] = append(cityIdListMap[cId], model.SelectOption{
				Label: st.Name,
				Value: st.ID,
			})
		}
	}

	for _, cityId := range cityIds {
		if cityIdMap[cityId] != nil && len(cityIdListMap[cityId]) != 0 {
			res = append(res, &model.CascaderOptionLevel2{
				SelectOption: model.SelectOption{
					Value: cityIdMap[cityId].ID,
					Label: cityIdMap[cityId].Name,
				},
				Children: cityIdListMap[cityId],
			})
		}
	}

	return
}
