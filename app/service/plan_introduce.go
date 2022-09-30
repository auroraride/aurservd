// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/tools"
)

type planIntroduceService struct {
    *BaseService
    orm *ent.PlanIntroduceClient
}

func NewPlanIntroduce(params ...any) *planIntroduceService {
    return &planIntroduceService{
        BaseService: NewService(params...),
        orm:         ent.Database.PlanIntroduce,
    }
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
            m[s.Key(bm, tools.NewPointer().UInt64(brand.ID))] = model.PlanIntroduceEbike{
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
            rv = &model.PlanIntroduceOption{
                Model:  k,
                Brands: make([]model.EbikeBrand, 0),
            }
            r[k] = rv
        }
        if b != nil {
            rv.Brands = append(rv.Brands, model.EbikeBrand{
                ID:   b.ID,
                Name: b.Name,
            })
        }
    }

    for _, o := range r {
        res = append(res, *o)
    }

    return res
}
