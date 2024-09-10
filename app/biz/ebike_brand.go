// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-21, by aurb

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
)

type ebikeBrandBiz struct {
	orm      *ent.EbikeBrandClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewEbikeBrand() *ebikeBrandBiz {
	return &ebikeBrandBiz{
		orm: ent.Database.EbikeBrand,
		ctx: context.Background(),
	}
}

func NewEbikeBrandWithModifier(m *model.Modifier) *ebikeBrandBiz {
	b := NewEbikeBrand()
	if m != nil {
		b.ctx = context.WithValue(b.ctx, model.CtxModifierKey{}, m)
		b.modifier = m
	}
	return b
}

// List 列表
func (b *ebikeBrandBiz) List(req *definition.EbikeBrandListReq) *model.PaginationRes {
	q := b.orm.QueryNotDeleted().WithBrandAttribute()

	if req.Keyword != nil {
		q.Where(
			ebikebrand.NameContains(*req.Keyword),
		)
	}
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.EbikeBrand) model.EbikeBrand {
		return b.detail(item)
	},
	)
}

// detail 数据拼接
func (b *ebikeBrandBiz) detail(eb *ent.EbikeBrand) model.EbikeBrand {
	brandAttribute := make([]*model.EbikeBrandAttribute, 0)
	if eb.Edges.BrandAttribute != nil {
		for _, ba := range eb.Edges.BrandAttribute {
			brandAttribute = append(brandAttribute, &model.EbikeBrandAttribute{
				Name:  ba.Name,
				Value: ba.Value,
			})
		}
	}
	return model.EbikeBrand{
		ID:             eb.ID,
		Name:           eb.Name,
		Cover:          eb.Cover,
		MainPic:        eb.MainPic,
		BrandAttribute: brandAttribute,
	}
}

// Create 创建电车品牌
func (b *ebikeBrandBiz) Create(req *model.EbikeBrandCreateReq) error {
	return service.NewEbikeBrand().Create(req)
}

// Modify 编辑电车品牌
func (b *ebikeBrandBiz) Modify(req *model.EbikeBrandModifyReq) error {
	return service.NewEbikeBrand().Modify(req)
}

// Delete 删除车电品牌
func (b *ebikeBrandBiz) Delete(req *definition.EbikeBrandDeleteReq) error {
	return NewEbikeBiz().DeleteBrand(req)
}
