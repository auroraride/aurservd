// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-31, by aurb

package biz

import (
	"context"
	"errors"

	"github.com/auroraride/aurservd/internal/ent/material"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
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

// List 其他物资列表
func (b *materialBiz) List(req *definition.MaterialListReq) (res *model.PaginationRes, err error) {
	q := b.orm.QueryNotDeleted()

	if req.Keyword != nil {
		q.Where(material.NameContains(*req.Keyword))
	}

	if req.Type != nil {
		q.Where(material.Type(req.Type.Value()))
	}

	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Material) (result *definition.MaterialDetail) {
		return b.otherDetail(item)
	})
	return
}

// detail 详情数据
func (b *materialBiz) otherDetail(item *ent.Material) (res *definition.MaterialDetail) {
	res = &definition.MaterialDetail{
		ID:        item.ID,
		Name:      item.Name,
		Type:      model.AssetType(item.Type),
		Statement: item.Statement,
	}
	return
}

// Create 其他物资创建
func (b *materialBiz) Create(req *definition.MaterialCreateReq) (err error) {
	if ex, _ := b.orm.QueryNotDeleted().Where(material.Type(req.Type.Value()), material.Name(req.Name)).First(b.ctx); ex != nil {
		return errors.New("该类型配件已存在")
	}

	_, err = b.orm.Create().
		SetType(req.Type.Value()).
		SetName(req.Name).
		SetStatement(req.Statement).
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

// Modify 其他物资编辑
func (b *materialBiz) Modify(req *definition.MaterialModifyReq) (err error) {
	g, _ := b.queryById(req.ID)
	if g == nil {
		return errors.New("数据不存在")
	}

	if ex, _ := b.orm.QueryNotDeleted().Where(material.IDNotIn(g.ID), material.Type(g.Type), material.Name(req.Name)).First(b.ctx); ex != nil {
		return errors.New("该类型配件已存在")
	}

	_, err = b.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetStatement(req.Statement).
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
