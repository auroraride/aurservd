// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
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
	brands, _ := s.orm.QueryNotDeleted().All(s.ctx)
	items := make([]model.EbikeBrand, len(brands))
	for i, b := range brands {
		items[i] = model.EbikeBrand{
			ID:    b.ID,
			Name:  b.Name,
			Cover: b.Cover,
		}
	}
	return items
}

func (s *ebikeBrandService) Create(req *model.EbikeBrandCreateReq) {
	s.orm.Create().SetName(req.Name).SetCover(req.Cover).ExecX(s.ctx)
}

func (s *ebikeBrandService) Modify(req *model.EbikeBrandModifyReq) {
	updater := s.orm.Update().Where(ebikebrand.ID(req.ID))
	if req.Name != "" {
		updater.SetName(req.Name)
	}
	if req.Cover != "" {
		updater.SetCover(req.Cover)
	}
	updater.ExecX(s.ctx)
}
