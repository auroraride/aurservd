// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/ebikebrandattribute"
)

type ebikeBrandAttributeService struct {
	*BaseService
	orm *ent.EbikeBrandAttributeClient
}

func NewEbikeBrandAttribute(params ...any) *ebikeBrandAttributeService {
	return &ebikeBrandAttributeService{
		BaseService: newService(params...),
		orm:         ent.Database.EbikeBrandAttribute,
	}
}

// Create 创建电动车品牌属性
func (s *ebikeBrandAttributeService) Create(req *model.EbikeBrandAttributeCreateReq) error {
	if len(req.BrandAttribute) == 0 {
		return nil
	}

	bulk := make([]*ent.EbikeBrandAttributeCreate, 0)
	for _, v := range req.BrandAttribute {
		b := s.orm.Create().
			SetBrandID(req.BrandID).
			SetValue(v.Value).
			SetName(v.Name)
		bulk = append(bulk, b)
	}
	err := s.orm.CreateBulk(bulk...).Exec(s.ctx)
	if err != nil {
		return err
	}
	return nil
}

// Update 更新电动车品牌属性
func (s *ebikeBrandAttributeService) Update(req *model.EbikeBrandAttributeUpdateReq) error {
	if len(req.BrandAttribute) == 0 {
		return nil
	}

	// 删除原有属性
	_, err := s.orm.Delete().Where(ebikebrandattribute.BrandID(req.BrandID)).Exec(s.ctx)
	if err != nil {
		return err
	}

	// 创建新属性
	err = s.Create(&model.EbikeBrandAttributeCreateReq{
		BrandID:        req.BrandID,
		BrandAttribute: req.BrandAttribute,
	})
	if err != nil {
		return err
	}

	return nil
}
