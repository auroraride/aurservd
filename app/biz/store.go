package biz

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"entgo.io/ent/dialect/sql"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/store"
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
func (s *storeBiz) List(req *definition.StoreListReq) (res []*definition.StoreDetail, err error) {
	res = make([]*definition.StoreDetail, 0)
	q := s.orm.QueryNotDeleted().
		WithCity().
		WithEmployee().
		WithStocks().
		Modify(func(sel *sql.Selector) {
			sel.
				AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(ST_GeographyFromText('SRID=4326;POINT(' || "store"."lng" || ' ' || "store"."lat" || ')'),ST_GeographyFromText('SRID=4326;POINT(%f  %f)'))`, req.Lng, req.Lat)), "distance").
				OrderBy(sql.Asc("distance"))
			if req.Distance != nil {
				if *req.Distance > 100000 {
					*req.Distance = 100000
				}
				sel.Where(sql.P(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`ST_DWithin(ST_GeographyFromText('SRID=4326;POINT(' || "store"."lng" || ' ' || "store"."lat" || ')'),ST_GeographyFromText('SRID=4326;POINT(%f  %f)'), %f)`, req.Lng, req.Lat, *req.Distance))
				}))
			}
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

	list, _ := q.All(s.ctx)
	if len(list) == 0 {
		return res, nil
	}

	var pl *ent.Plan
	if req.PlanID != nil {
		pl, _ = ent.Database.Plan.QueryNotDeleted().Where(plan.ID(*req.PlanID)).First(s.ctx)
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
			ebikeNum, batteryNum := s.queryStocks(v, pl)
			// 电车或电池库存为0则不展示
			if ebikeNum <= 0 || batteryNum <= 0 {
				continue
			}
		}
		res = append(res, s.detail(v))
	}
	return
}

// Detail 门店详情
func (s *storeBiz) Detail(req *definition.StoreDetailReq) (res *definition.StoreDetail) {
	q, _ := s.orm.QueryNotDeleted().
		Where(store.ID(req.ID)).
		WithCity().
		WithEmployee().
		WithStocks().
		Modify(func(sel *sql.Selector) {
			sel.AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(ST_GeographyFromText('SRID=4326;POINT(' || "store"."lng" || ' ' || "store"."lat" || ')'),ST_GeographyFromText('SRID=4326;POINT(%f  %f)'))`, req.Lng, req.Lat)), "distance").
				OrderBy(sql.Asc("distance"))
		}).
		First(s.ctx)
	if q == nil {
		return nil
	}
	return s.detailForStock(q)
}

func (s *storeBiz) detail(item *ent.Store) (res *definition.StoreDetail) {
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

// 查询门店电车库存
func (s *storeBiz) queryStocks(item *ent.Store, pl *ent.Plan) (ebikeNum, batteryNum int) {
	bikes := make(map[string]*model.StockMaterial)
	batteries := make(map[string]*model.StockMaterial)
	for _, st := range item.Edges.Stocks {
		switch true {
		case st.BrandID != nil && *st.BrandID == *pl.BrandID:
			// 电车
			service.NewStock().Calculate(bikes, st)
		case st.Model != nil && *st.Model == pl.Model:
			// 电池
			service.NewStock().Calculate(batteries, st)
		}
	}
	for _, bike := range bikes {
		ebikeNum += bike.Surplus
	}

	for _, battery := range batteries {
		batteryNum += battery.Surplus
	}
	return
}

// StoreBySubscribe 根据订阅查询门店
func (s *storeBiz) StoreBySubscribe(r *ent.Rider, req *definition.StoreDetailReq) (res *definition.StoreDetail, err error) {
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
	}).First(s.ctx)
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
	bikes := make(map[string]*model.StockMaterial)

	brandIdMap := make(map[uint64]bool)
	for _, brandId := range brandIds {
		brandIdMap[brandId] = true
	}

	for _, st := range item.Edges.Stocks {
		switch {
		case st.BrandID != nil && brandIdMap[*st.BrandID]:
			service.NewStock().Calculate(bikes, st)
		}
	}
	for bId, bike := range bikes {
		brandId, _ := strconv.Atoi(bId)
		if bike.Surplus > 0 {
			eBrandIds = append(eBrandIds, uint64(brandId))
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

	if item.Edges.Stocks != nil {
		var brandIds []uint64
		for _, st := range item.Edges.Stocks {
			if st.BrandID != nil {
				brandIds = append(brandIds, *st.BrandID)
			}
		}

		// 查询门店存在库存的电车数据
		sBids := b.queryStocksByStore(item, brandIds)
		ebikeRes := NewPlanBiz().EbikeList(sBids)
		if ebikeRes.Brands != nil {
			res.Brands = ebikeRes.Brands
			res.RtoBrands = ebikeRes.RtoBrands
		}
	}

	// 查询门店商品
	res.SaleGoods = NewGoods().ListByStoreId(item.ID)

	return

}
