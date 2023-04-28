// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/planintroduce"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type planIntroduceService struct {
	*BaseService
	orm *ent.PlanIntroduceClient
}

func NewPlanIntroduce(params ...any) *planIntroduceService {
	return &planIntroduceService{
		BaseService: newService(params...),
		orm:         ent.Database.PlanIntroduce,
	}
}

func (s *planIntroduceService) Query(id uint64) (*ent.PlanIntroduce, error) {
	q := s.orm.Query().Where(planintroduce.ID(id))
	return q.First(s.ctx)
}

func (s *planIntroduceService) QueryX(id uint64) *ent.PlanIntroduce {
	intro, _ := s.Query(id)
	if intro == nil {
		snag.Panic("未找到介绍")
	}
	return intro
}

func (s *planIntroduceService) QueryModelBrand(req model.PlanIntroduceQuery) (*ent.PlanIntroduce, error) {
	q := s.orm.Query().Where(planintroduce.Model(req.Model))
	if req.EbikeBrandID != nil {
		q.Where(planintroduce.BrandID(*req.EbikeBrandID))
	}
	return q.First(s.ctx)
}

func (s *planIntroduceService) QueryModelBrandX(req model.PlanIntroduceQuery) *ent.PlanIntroduce {
	intro, _ := s.QueryModelBrand(req)
	if intro == nil {
		snag.Panic("未找到介绍")
	}
	return intro
}

// Key 获取简介Key
func (s *planIntroduceService) Key(model string, brandID *uint64) string {
	k := model
	if brandID != nil {
		k += fmt.Sprintf("-%d", *brandID)
	}
	return k
}

// Notset 获取未设定介绍的车电型号
func (s *planIntroduceService) Notset() (res []model.PlanIntroduceOption) {
	res = make([]model.PlanIntroduceOption, 0)

	m := make(ar.Map)
	models := NewBatteryModel().Models()
	brands := NewEbikeBrand().All()
	for _, bm := range models {
		m[bm] = bm
		for _, brand := range brands {
			m[s.Key(bm, silk.UInt64(brand.ID))] = model.PlanIntroduceEbike{
				Model: bm,
				Name:  brand.Name,
				ID:    brand.ID,
			}
		}
	}

	items, _ := s.orm.Query().All(s.ctx)
	for _, item := range items {
		k := s.Key(item.Model, item.BrandID)
		if _, ok := m[k]; ok {
			delete(m, k)
		}
	}

	r := make(map[string]*model.PlanIntroduceOption)

	for _, v := range m {
		var k string
		var b *model.PlanIntroduceEbike

		switch o := v.(type) {
		case string:
			k = o
		case model.PlanIntroduceEbike:
			k = o.Model
			b = &o
		}
		rv, ok := r[k]
		if !ok {
			_, notSet := m[k]
			rv = &model.PlanIntroduceOption{
				Model:       k,
				EbikeBrands: make([]model.EbikeBrand, 0),
				ModelSet:    !notSet,
			}
			r[k] = rv
		}

		if b != nil {
			rv.EbikeBrands = append(rv.EbikeBrands, model.EbikeBrand{
				ID:   b.ID,
				Name: b.Name,
			})
		}
	}

	for _, o := range r {
		res = append(res, *o)
	}

	sort.Slice(res, func(i, j int) bool {
		return strings.Compare(res[i].Model, res[j].Model) < 0
	})

	return res
}

func (s *planIntroduceService) List() (res []model.PlanIntroduce) {
	items, _ := s.orm.Query().WithBrand().All(s.ctx)
	res = make([]model.PlanIntroduce, len(items))
	for i, item := range items {
		res[i] = model.PlanIntroduce{
			ID:    item.ID,
			Model: item.Model,
			Image: item.Image,
		}
		b := item.Edges.Brand
		if b != nil {
			res[i].EbikeBrand = &model.EbikeBrand{
				ID:   b.ID,
				Name: b.Name,
			}
		}
	}
	return
}

func (s *planIntroduceService) Create(req *model.PlanIntroduceCreateReq) {
	// 查找是否重复
	b, _ := s.QueryModelBrand(model.PlanIntroduceQuery{
		EbikeBrandID: req.EbikeBrandID,
		Model:        req.Model,
	})
	if b != nil {
		snag.Panic("此条已设置过")
	}
	s.orm.Create().SetNillableBrandID(req.EbikeBrandID).SetImage(req.Image).SetModel(req.Model).ExecX(s.ctx)
}

func (s *planIntroduceService) Modify(req *model.PlanIntroduceModifyReq) {
	s.QueryX(req.ID).Update().SetImage(req.Image).ExecX(s.ctx)
}

func (s *planIntroduceService) QueryMap() (res map[string]string) {
	res = make(map[string]string)
	items, _ := s.orm.Query().All(s.ctx)
	for _, item := range items {
		res[s.Key(item.Model, item.BrandID)] = item.Image
	}
	return
}
