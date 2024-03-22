package biz

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/store"
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
func (s *storeBiz) List(req *definition.StoreListReq) (res []*definition.Store) {
	res = make([]*definition.Store, 0)
	q := s.orm.QueryNotDeleted().
		WithCity().
		WithEmployee().
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
		return res
	}
	for _, v := range list {
		res = append(res, s.detail(v))
	}

	return
}

// Detail 门店详情
func (s *storeBiz) Detail(id uint64) (res *definition.Store) {
	q, _ := s.orm.QueryNotDeleted().
		Where(store.ID(id)).
		WithCity().
		WithEmployee().
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
