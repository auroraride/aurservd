// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"errors"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/pkg/snag"
)

type ebikeBrandService struct {
	*BaseService
	orm *ent.EbikeBrandClient
}

func NewEbikeBrand(params ...any) *ebikeBrandService {
	return &ebikeBrandService{
		BaseService: newService(params...),
		orm:         ent.Database.EbikeBrand,
	}
}

func (s *ebikeBrandService) Query(id uint64) (*ent.EbikeBrand, error) {
	return s.orm.QueryNotDeleted().Where(ebikebrand.ID(id)).First(s.ctx)
}

func (s *ebikeBrandService) QueryX(id uint64) *ent.EbikeBrand {
	b, _ := s.Query(id)
	if b == nil {
		snag.Panic("未找到电车型号")
	}
	return b
}

func (s *ebikeBrandService) All() []model.EbikeBrand {
	brands, _ := s.orm.QueryNotDeleted().WithBrandAttribute().All(s.ctx)
	items := make([]model.EbikeBrand, len(brands))
	for i, b := range brands {
		brandAttribute := make([]*model.EbikeBrandAttribute, 0)
		if b.Edges.BrandAttribute != nil {
			for _, ba := range b.Edges.BrandAttribute {
				brandAttribute = append(brandAttribute, &model.EbikeBrandAttribute{
					Name:  ba.Name,
					Value: ba.Value,
				})
			}
		}
		items[i] = model.EbikeBrand{
			ID:             b.ID,
			Name:           b.Name,
			Cover:          b.Cover,
			MainPic:        b.MainPic,
			BrandAttribute: brandAttribute,
		}
	}
	return items
}

// Create 创建电车品牌
func (s *ebikeBrandService) Create(req *model.EbikeBrandCreateReq) error {
	br, err := s.orm.Create().SetName(req.Name).SetCover(req.Cover).SetMainPic(req.MainPic).Save(s.ctx)
	if err != nil && ent.IsConstraintError(err) {
		return errors.New("电车品牌已存在")
	}

	if br == nil {
		return errors.New("创建电车品牌失败")
	}

	err = NewEbikeBrandAttribute().Create(&model.EbikeBrandAttributeCreateReq{
		BrandAttribute: req.BrandAttribute,
		BrandID:        br.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *ebikeBrandService) Modify(req *model.EbikeBrandModifyReq) error {
	updater := s.orm.Update().Where(ebikebrand.ID(req.ID))
	if req.Name != "" {
		updater.SetName(req.Name)
	}
	if req.Cover != "" {
		updater.SetCover(req.Cover)
	}
	if len(req.MainPic) > 0 {
		updater.SetMainPic(req.MainPic)
	}
	err := updater.Exec(s.ctx)
	if err != nil {
		return err
	}

	err = NewEbikeBrandAttribute().Update(&model.EbikeBrandAttributeUpdateReq{
		BrandAttribute: req.BrandAttribute,
		BrandID:        req.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

// ListByCityAndPlan 查询城市套餐电车品牌
func (s *ebikeBrandService) ListByCityAndPlan(cityID uint64) []model.EbikeBrand {
	brands, _ := s.orm.QueryNotDeleted().
		Where(
			ebikebrand.HasPlansWith(
				plan.HasCitiesWith(city.ID(cityID)),
				plan.Enable(true),
				plan.DeletedAtIsNil(),
			),
		).
		All(s.ctx)

	items := make([]model.EbikeBrand, len(brands))
	for i, b := range brands {
		attrs := make([]*model.EbikeBrandAttribute, 0)
		if b.Edges.BrandAttribute != nil {
			for _, ba := range b.Edges.BrandAttribute {
				attrs = append(attrs, &model.EbikeBrandAttribute{
					Name:  ba.Name,
					Value: ba.Value,
				})
			}
		}
		items[i] = model.EbikeBrand{
			ID:             b.ID,
			Name:           b.Name,
			Cover:          b.Cover,
			MainPic:        b.MainPic,
			BrandAttribute: attrs,
		}
	}
	return items
}
