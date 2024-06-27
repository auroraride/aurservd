package biz

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"entgo.io/ent/dialect/sql"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/store"
)

type ebikeBiz struct {
	orm *ent.EbikeClient
	ctx context.Context
}

func NewEbikeBiz() *ebikeBiz {
	return &ebikeBiz{
		orm: ent.Database.Ebike,
		ctx: context.Background(),
	}
}

// EbikeBrandDetail 车电品牌详情
func (s *ebikeBiz) EbikeBrandDetail(req *definition.EbikeDetailReq) (*definition.EbikeDetailRes, error) {
	d, _ := ent.Database.EbikeBrand.QueryNotDeleted().
		Where(ebikebrand.ID(req.ID)).
		WithBrandAttribute().
		First(s.ctx)
	if d == nil {
		return nil, nil
	}

	storeItem, _ := ent.Database.Store.QueryNotDeleted().
		WithStocks().
		Where(store.ID(req.StoreID)).
		Modify(func(sel *sql.Selector) {
			sel.
				AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(ST_GeographyFromText('SRID=4326;POINT(' || "store"."lng" || ' ' || "store"."lat" || ')'),ST_GeographyFromText('SRID=4326;POINT(%f  %f)'))`, req.Lng, req.Lat)), "distance").
				OrderBy(sql.Asc("distance"))
		}).
		First(s.ctx)

	pl, _ := ent.Database.Plan.QueryNotDeleted().Where(plan.ID(req.PlanID)).WithComplexes().WithParent(func(query *ent.PlanQuery) {
		query.WithComplexes()
	}).First(s.ctx)
	if pl == nil {
		return nil, errors.New("套餐不存在")
	}

	res := &definition.EbikeDetailRes{

		Brand: model.EbikeBrand{
			ID:             d.ID,
			Name:           d.Name,
			Cover:          d.Cover,
			MainPic:        d.MainPic,
			BrandAttribute: make([]*model.EbikeBrandAttribute, 0),
		},
	}

	slPlan := NewPlanBiz().GetPlanItems(pl)
	if len(slPlan) > 0 {
		// 排序 价格升序
		sort.Slice(slPlan, func(i, j int) bool {
			return slPlan[i].Price < slPlan[j].Price
		})
		res.Plan = model.Plan{
			ID:          slPlan[0].ID,
			Name:        slPlan[0].Name,
			Price:       slPlan[0].Price,
			Days:        slPlan[0].Days,
			Intelligent: slPlan[0].Intelligent,
			Type:        model.PlanType(slPlan[0].Type),
		}
	}

	if d.Edges.BrandAttribute != nil {
		for _, v := range d.Edges.BrandAttribute {
			res.Brand.BrandAttribute = append(res.Brand.BrandAttribute, &model.EbikeBrandAttribute{
				Name:  v.Name,
				Value: v.Value,
			})
		}
	}

	if storeItem != nil {
		res.Store = struct {
			model.StoreLngLat
			Address  string  `json:"address"`  // 地址
			Distance float64 `json:"distance"` // 距离
		}{
			StoreLngLat: model.StoreLngLat{
				Store: model.Store{
					ID:   storeItem.ID,
					Name: storeItem.Name,
				},
				Lat: storeItem.Lat,
				Lng: storeItem.Lng,
			},
			Address:  storeItem.Address,
			Distance: 0,
		}
		distance, err := storeItem.Value("distance")
		if distance != nil || err == nil {
			distanceFloat, ok := distance.(float64)
			if ok {
				res.Store.Distance = distanceFloat
			}
		}
	}

	return res, nil
}
