// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-05-29, by aurb

package biz

import (
	"context"
	"errors"

	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/goods"
	"github.com/auroraride/aurservd/internal/ent/storegoods"
	"github.com/auroraride/aurservd/pkg/tools"
)

type goodsBiz struct {
	orm      *ent.GoodsClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewGoods() *goodsBiz {
	return &goodsBiz{
		orm: ent.Database.Goods,
		ctx: context.Background(),
	}
}

func NewGoodsWithModifierBiz(m *model.Modifier) *goodsBiz {
	s := NewGoods()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// List 获取商品列表
func (b *goodsBiz) List(req *definition.GoodsListReq) *model.PaginationRes {
	query := b.orm.QueryNotDeleted().Order(ent.Desc(goods.FieldWeight)).
		WithStores(
			func(q *ent.StoreGoodsQuery) {
				q.Where(storegoods.DeletedAtIsNil()).WithStore()
			},
		)
	if req.Keyword != nil {
		query.Where(
			goods.Or(
				goods.SnContains(*req.Keyword),
				goods.NameContains(*req.Keyword),
			),
		)
	}
	if req.Status != nil {
		query.Where(goods.Status(req.Status.Value()))
	}
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		query.Where(
			goods.CreatedAtGTE(start),
			goods.CreatedAtLTE(end),
		)
	}
	return model.ParsePaginationResponse(query, req.PaginationReq, func(item *ent.Goods) *definition.GoodsDetail {
		return toGoodsDetail(item)
	})
}

// toGoodsDetail 商品详情数据拼接
func toGoodsDetail(item *ent.Goods) *definition.GoodsDetail {
	// 查询配置的门店信息
	var storeIds []uint64
	var stores []model.Store
	for _, es := range item.Edges.Stores {
		storeIds = append(storeIds, es.Edges.Store.ID)
		stores = append(stores, model.Store{
			ID:   es.Edges.Store.ID,
			Name: es.Edges.Store.Name,
		})
	}
	return &definition.GoodsDetail{
		ID: item.ID,
		Goods: definition.Goods{
			Sn:        item.Sn,
			Name:      item.Name,
			Type:      definition.GoodsType(item.Type),
			Lables:    item.Lables,
			Price:     item.Price,
			Weight:    item.Weight,
			HeadPic:   item.HeadPic,
			Photos:    item.Photos,
			Intro:     item.Intro,
			Stores:    stores,
			CreatedAt: item.CreatedAt.Format(carbon.DateTimeLayout),
			Status:    definition.GoodsStatus(item.Status),
			Remark:    item.Remark,
			StoreIds:  storeIds,
		},
	}
}

// Create 创建商品
func (b *goodsBiz) Create(req *definition.GoodsCreateReq) (err error) {
	sn := tools.NewUnique().NewSN()
	var item *ent.Goods
	item, err = b.orm.Create().
		SetSn(sn).
		SetName(req.Name).
		SetType(definition.GoodsTypeEbike.Value()).
		SetLables(req.Lables).
		SetPrice(req.Price).
		SetWeight(req.Weight).
		SetHeadPic(req.HeadPic).
		SetPhotos(req.Photos).
		SetIntro(req.Intro).
		SetRemark(req.Remark).
		SetStatus(definition.GoodsStatusOnline.Value()).
		Save(b.ctx)
	if err != nil {
		return err
	}

	bulk := make([]*ent.StoreGoodsCreate, len(req.StoreIds))
	for i, storeId := range req.StoreIds {
		bulk[i] = ent.Database.StoreGoods.Create().SetGoodsID(item.ID).SetStoreID(storeId)
	}

	return ent.Database.StoreGoods.CreateBulk(bulk...).Exec(b.ctx)
}

// queryById 通过ID查询商品
func (b *goodsBiz) queryById(id uint64) (item *ent.Goods, err error) {
	return b.orm.QueryNotDeleted().Where(goods.ID(id)).First(b.ctx)
}

// Delete 删除商品
func (b *goodsBiz) Delete(id uint64) (err error) {
	g, _ := b.queryById(id)
	if g == nil {
		return errors.New("商品不存在")
	}
	_, err = b.orm.SoftDeleteOne(g).Save(b.ctx)
	if err != nil {
		return err
	}
	return
}

// Modify 编辑商品
func (b *goodsBiz) Modify(req *definition.GoodsModifyReq) (err error) {
	g, _ := b.queryById(req.ID)
	if g == nil {
		return errors.New("商品不存在")
	}

	_, err = b.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetType(definition.GoodsTypeEbike.Value()).
		SetLables(req.Lables).
		SetPrice(req.Price).
		SetWeight(req.Weight).
		SetHeadPic(req.HeadPic).
		SetPhotos(req.Photos).
		SetIntro(req.Intro).
		SetRemark(req.Remark).
		SetStatus(definition.GoodsStatusOnline.Value()).
		Save(b.ctx)
	if err != nil {
		return err
	}

	// 直接先删除已配置的门店
	_, _ = ent.Database.StoreGoods.Delete().Where(storegoods.GoodsID(g.ID)).Exec(b.ctx)

	for _, storeId := range req.StoreIds {
		_, _ = ent.Database.StoreGoods.Create().SetGoodsID(g.ID).SetStoreID(storeId).Save(b.ctx)
	}

	return
}

// Detail 查询商品详情
func (b *goodsBiz) Detail(id uint64) (*definition.GoodsDetail, error) {
	item, err := b.orm.Query().Where(goods.ID(id)).First(b.ctx)
	if err != nil {
		return nil, err
	}
	return toGoodsDetail(item), nil
}

// UpdateStatus 更新商品上下架状态
func (b *goodsBiz) UpdateStatus(req *definition.GoodsUpdateStatusReq) (err error) {
	g, _ := b.queryById(req.ID)
	if g == nil {
		return errors.New("商品不存在")
	}

	if g.Status == req.Status.Value() {
		return errors.New("商品状态已存在")
	}

	_, err = b.orm.UpdateOneID(req.ID).
		SetStatus(req.Status.Value()).
		Save(b.ctx)
	if err != nil {
		return err
	}
	return
}

// ListByStoreId 通过storeId查询商品数据
func (b *goodsBiz) ListByStoreId(storeId uint64) (res []*definition.GoodsDetail) {
	items, _ := b.orm.QueryNotDeleted().Order(ent.Desc(goods.FieldWeight)).
		Where(
			goods.Status(definition.GoodsStatusOnline.Value()),
			goods.HasStoresWith(storegoods.StoreID(storeId), storegoods.DeletedAtIsNil()),
		).All(b.ctx)

	for _, item := range items {
		res = append(res, toGoodsDetail(item))
	}
	return
}
