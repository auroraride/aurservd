package biz

import (
	"context"
	"errors"
	"fmt"

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
func (s *storeBiz) List(req *definition.StoreListReq) (res []*definition.Store, err error) {
	res = make([]*definition.Store, 0)
	q := s.orm.QueryNotDeleted().
		WithCity().
		WithEmployee().
		WithStocks().
		Modify(func(sel *sql.Selector) {
			sel.
				AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance( ST_SetSRID(ST_MakePoint(%s, %s), 4326), ST_SetSRID(ST_MakePoint(%f, %f), 4326))`, sel.C("lng"), sel.C("lat"), req.Lng, req.Lat)), "distance").
				OrderBy(sql.Asc("distance"))
			if req.Distance != nil {
				if *req.Distance > 100000 {
					*req.Distance = 100000
				}
				sel.Where(sql.P(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`ST_DWithin(ST_SetSRID(ST_MakePoint(%s, %s), 4326), ST_SetSRID(ST_MakePoint(%f, %f), 4326), %f)`, sel.C("lng"), sel.C("lat"), req.Lng, req.Lat, *req.Distance))
				}))
			}
		})

	if req.CityID != nil {
		q.Where(store.CityID(*req.CityID))
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
func (s *storeBiz) Detail(req *definition.StoreDetail) (res *definition.Store) {
	q, _ := s.orm.QueryNotDeleted().
		Where(store.ID(req.ID)).
		WithCity().
		WithEmployee().
		Modify(func(sel *sql.Selector) {
			sel.AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance( ST_SetSRID(ST_MakePoint(%s, %s), 4326), ST_SetSRID(ST_MakePoint(%f, %f), 4326))`, sel.C("lng"), sel.C("lat"), req.Lng, req.Lat)), "distance").
				OrderBy(sql.Asc("distance"))

		}).
		First(s.ctx)
	if q == nil {
		return nil
	}
	return s.detail(q)
}

func (s *storeBiz) detail(item *ent.Store) (res *definition.Store) {
	res = &definition.Store{
		ID:            item.ID,
		Name:          item.Name,
		Status:        item.Status,
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
	value, err := item.Value("distance")
	if err == nil {
		res.Distance = value.(float64)
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
func (s *storeBiz) StoreBySubscribe(r *ent.Rider, req *definition.StoreBySubscribe) (res *definition.Store, err error) {
	q, _ := ent.Database.Subscribe.QueryNotDeleted().
		Where(subscribe.ID(req.ID), subscribe.RiderID(r.ID)).
		WithStore(
			func(query *ent.StoreQuery) {
				query.WithCity().WithEmployee()
			}).First(s.ctx)
	if q == nil || q.Edges.Store == nil {
		return nil, errors.New("未找到门店")
	}
	return s.detail(q.Edges.Store), nil
}
