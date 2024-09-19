// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-22, by aurb

package biz

import (
	"context"
	"errors"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/storegroup"
)

type storeGroupBiz struct {
	orm      *ent.StoreGroupClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewStoreGroup() *storeGroupBiz {
	return &storeGroupBiz{
		orm: ent.Database.StoreGroup,
		ctx: context.Background(),
	}
}

func NewStoreGroupWithModifier(m *model.Modifier) *storeGroupBiz {
	b := NewStoreGroup()
	b.ctx = context.WithValue(b.ctx, model.CtxModifierKey{}, m)
	b.modifier = m
	return b
}

// Create 创建
func (b *storeGroupBiz) Create(req *definition.StoreGroupCreateRep) (err error) {
	sg, _ := b.orm.QueryNotDeleted().Where(storegroup.Name(req.Name)).First(b.ctx)
	if sg != nil {
		return errors.New("门店集合已存在")
	}

	_, err = b.orm.Create().SetName(req.Name).Save(b.ctx)

	return
}

// List 列表
func (b *storeGroupBiz) List() (res []*definition.StoreGroupListRes) {
	res = make([]*definition.StoreGroupListRes, 0)

	sgs, _ := b.orm.QueryNotDeleted().All(b.ctx)
	for _, sg := range sgs {
		res = append(res, &definition.StoreGroupListRes{
			ID:   sg.ID,
			Name: sg.Name,
		})
	}
	return
}

// Delete 删除
func (b *storeGroupBiz) Delete(id uint64) (err error) {
	_, err = b.orm.SoftDeleteOneID(id).Save(b.ctx)
	return
}
