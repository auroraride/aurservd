// Copyright (C) liasica. 2023-present.
//
// Created at 2023-06-21
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"entgo.io/ent/dialect/sql"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/enterpriseprice"
	"github.com/auroraride/aurservd/pkg/snag"
)

type enterprisePriceService struct {
	*BaseService

	orm *ent.EnterprisePriceClient
}

func NewEnterprisePrice(params ...any) *enterprisePriceService {
	return &enterprisePriceService{
		BaseService: newService(params...),
		orm:         ent.Database.EnterprisePrice,
	}
}

// CityList 团签开通城市列表
func (s *enterprisePriceService) CityList(enterpriseID uint64) (res []model.City) {
	_ = ent.Database.City.QueryNotDeleted().Modify(func(sel *sql.Selector) {
		t := sql.Table(enterpriseprice.Table)
		sel.LeftJoin(t).On(sel.C(city.FieldID), t.C(enterpriseprice.FieldCityID))
		sel.Where(sql.EQ(t.C(enterpriseprice.FieldEnterpriseID), enterpriseID))
		sel.Where(sql.IsNull(t.C(enterpriseprice.FieldDeletedAt)))
		sel.Select(sel.C(city.FieldID), sel.C(city.FieldName))
		sel.GroupBy(sel.C(city.FieldID))
	}).Scan(s.ctx, &res)
	if len(res) == 0 {
		res = make([]model.City, 0)
	}
	return
}

// PriceList 团签价格列表
func (s *enterprisePriceService) PriceList(enterpriseId uint64) (res []model.EnterprisePriceWithCity) {
	pr, _ := ent.Database.EnterprisePrice.QueryNotDeleted().WithCity().WithBrand().Where(enterpriseprice.EnterpriseID(enterpriseId)).All(s.ctx)
	if len(pr) == 0 {
		snag.Panic("团签价格不存在")
	}
	res = make([]model.EnterprisePriceWithCity, 0)
	for _, ep := range pr {
		data := model.EnterprisePriceWithCity{
			ID:          ep.ID,
			Model:       ep.Model,
			Price:       ep.Price,
			Intelligent: ep.Intelligent,
			City: model.City{
				ID:   ep.Edges.City.ID,
				Name: ep.Edges.City.Name,
			},
		}
		if ep.Edges.Brand != nil {
			data.EbikeBrand = &model.EbikeBrand{
				ID:    ep.Edges.Brand.ID,
				Name:  ep.Edges.Brand.Name,
				Cover: ep.Edges.Brand.Cover,
			}
		}
		res = append(res, data)
	}
	return
}
