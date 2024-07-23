// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-15, by aurb

package biz

import (
	"context"
	"errors"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/material"
)

type materialBiz struct {
	orm      *ent.MaterialClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewMaterial() *materialBiz {
	return &materialBiz{
		orm: ent.Database.Material,
		ctx: context.Background(),
	}
}

func NewMaterialWithModifier(m *model.Modifier) *materialBiz {
	s := NewMaterial()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// List 列表
func (b *materialBiz) List(req *definition.MaterialListReq) (res []*definition.MaterialDetail, err error) {
	res = make([]*definition.MaterialDetail, 0)

	q := b.orm.QueryNotDeleted()

	if req.Keyword != nil {
		q.Where(material.NameContains(*req.Keyword))
	}

	if req.Type != nil {
		q.Where(material.Type(req.Type.Value()))
	}

	list, _ := q.All(b.ctx)

	if len(list) == 0 {
		return res, nil
	}

	for _, v := range list {
		res = append(res, b.detail(v))
	}
	return
}

// detail 详情数据
func (b *materialBiz) detail(item *ent.Material) (res *definition.MaterialDetail) {
	res = &definition.MaterialDetail{
		ID:        item.ID,
		Name:      item.Name,
		Type:      definition.MaterialType(item.Type),
		Statement: item.Statement,
		Allot:     item.Allot,
	}

	return
}

// Detail 详情
func (b *materialBiz) Detail(id uint64) (*definition.MaterialDetail, error) {
	g, _ := b.queryById(id)
	if g == nil {
		return nil, errors.New("数据不存在")
	}
	return b.detail(g), nil
}

// Create 创建
func (b *materialBiz) Create(req *definition.MaterialCreateReq) (err error) {
	_, err = b.orm.Create().
		SetName(req.Name).
		SetType(req.Type.Value()).
		SetStatement(req.Statement).
		SetAllot(req.Allot).
		Save(b.ctx)
	if err != nil {
		return err
	}
	return
}

// queryById 通过ID查询
func (b *materialBiz) queryById(id uint64) (item *ent.Material, err error) {
	return b.orm.QueryNotDeleted().Where(material.ID(id)).First(b.ctx)
}

// Modify 编辑
func (b *materialBiz) Modify(req *definition.MaterialModifyReq) (err error) {
	g, _ := b.queryById(req.ID)
	if g == nil {
		return errors.New("数据不存在")
	}

	_, err = b.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetType(req.Type.Value()).
		SetStatement(req.Statement).
		SetAllot(req.Allot).
		Save(b.ctx)
	if err != nil {
		return err
	}

	return
}

// Delete 删除
func (b *materialBiz) Delete(id uint64) (err error) {
	g, _ := b.queryById(id)
	if g == nil {
		return errors.New("数据不存在")
	}
	_, err = b.orm.SoftDeleteOne(g).Save(b.ctx)
	if err != nil {
		return err
	}
	return
}
