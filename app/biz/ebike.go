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
	"github.com/auroraride/aurservd/internal/ent/ebike"
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
func (b *ebikeBiz) EbikeBrandDetail(req *definition.EbikeDetailReq) (*definition.EbikeDetailRes, error) {
	d, _ := ent.Database.EbikeBrand.QueryNotDeleted().
		Where(ebikebrand.ID(req.ID)).
		WithBrandAttribute().
		First(b.ctx)
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
		First(b.ctx)

	pl, _ := ent.Database.Plan.QueryNotDeleted().Where(plan.ID(req.PlanID)).WithComplexes().WithParent(func(query *ent.PlanQuery) {
		query.WithComplexes()
	}).First(b.ctx)
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

// DeleteBrand 删除车电品牌
func (b *ebikeBiz) DeleteBrand(req *definition.EbikeBrandDeleteReq) error {
	// 查询车电品牌下是否有车辆在使用
	i, err := ent.Database.Ebike.Query().Where(ebike.BrandID(req.ID)).Count(b.ctx)
	if err != nil {
		return err
	}
	if i > 0 {
		return errors.New("车电品牌下有车辆在使用")
	}

	err = ent.Database.EbikeBrand.SoftDelete().Where(ebikebrand.ID(req.ID)).Exec(b.ctx)
	if err != nil {
		return err
	}
	return nil
}

// BatchModify 批量修改电车型号
func (b *ebikeBiz) BatchModify(req *definition.EbikeBatchModifyReq) []error {
	errs := make([]error, 0)
	for _, v := range req.SN {
		e, _ := ent.Database.Ebike.Query().Where(ebike.Sn(v), ebike.RiderIDIsNil()).First(b.ctx)
		if e == nil {
			errs = append(errs, fmt.Errorf("sn:%s, %s", v, "未找到车电或已分配骑手,无法修改"))
			continue
		}
		err := e.Update().SetBrandID(req.BrandID).Exec(b.ctx)
		if err != nil {
			errs = append(errs, fmt.Errorf("sn:%s, %s", v, err.Error()))
		}
	}
	return errs
}
