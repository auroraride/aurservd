// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"errors"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/ebike"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/store"
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

func (s *ebikeBrandService) AllByCity(cityID uint64) []model.EbikeBrand {
	// 查出城市的门店
	stores, _ := ent.Database.Store.QueryNotDeleted().
		Where(
			store.CityID(cityID),
		).
		All(s.ctx)
	storeIds := make([]uint64, 0)
	for _, st := range stores {
		storeIds = append(storeIds, st.ID)
	}

	// 查询门店所有车辆及品牌信息
	ebikes, _ := ent.Database.Ebike.Query().
		Where(
			ebike.Enable(true),
			ebike.HasStoreWith(store.IDIn(storeIds...)),
		).
		WithBrand(func(query *ent.EbikeBrandQuery) {
			query.WithBrandAttribute()
		}).
		All(s.ctx)

	brands := make([]*ent.EbikeBrand, 0)
	brandIdMap := make(map[uint64]bool)
	for _, b := range ebikes {
		if b.Edges.Brand != nil && !brandIdMap[b.Edges.Brand.ID] {
			brands = append(brands, b.Edges.Brand)
			brandIdMap[b.Edges.Brand.ID] = true
		}
	}

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
