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
	"github.com/auroraride/aurservd/internal/ent/store"
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
	q := b.orm.QueryNotDeleted().
		WithStores(
			func(query *ent.StoreGoodsQuery) {
				query.Where(storegoods.DeletedAtIsNil()).WithStore()
			},
		).
		Order(ent.Desc(goods.FieldWeight))

	b.listFilter(req, q)

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Goods) *definition.GoodsDetail {
		return detail(item)
	})
}

// listFilter 条件筛选
func (b *goodsBiz) listFilter(req *definition.GoodsListReq, q *ent.GoodsQuery) {
	if req.Keyword != nil {
		q.Where(
			goods.Or(
				goods.SnContains(*req.Keyword),
				goods.NameContains(*req.Keyword),
			),
		)
	}
	if req.Status != nil {
		q.Where(goods.Status(req.Status.Value()))
	}
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(
			goods.CreatedAtGTE(start),
			goods.CreatedAtLTE(end),
		)
	}
}

// detail 商品详情数据拼接
func detail(item *ent.Goods) *definition.GoodsDetail {
	// 查询配置的门店信息
	storeIds := make([]uint64, 0)
	stores := make([]model.Store, 0)
	for _, es := range item.Edges.Stores {
		if es.Edges.Store != nil {
			storeIds = append(storeIds, es.Edges.Store.ID)
			stores = append(stores, model.Store{
				ID:   es.Edges.Store.ID,
				Name: es.Edges.Store.Name,
			})
		}
	}

	// 解析付款方案数据
	payPlans := make([][]float64, 0)
	for _, p := range item.PaymentPlans {
		payPlan := make([]float64, 0)
		for _, o := range p {
			payPlan = append(payPlan, o.Amount)
		}
		payPlans = append(payPlans, payPlan)
	}

	return &definition.GoodsDetail{
		ID: item.ID,
		Goods: definition.Goods{
			Sn:           item.Sn,
			Name:         item.Name,
			Type:         definition.GoodsType(item.Type),
			Lables:       item.Lables,
			Price:        item.Price,
			Weight:       item.Weight,
			HeadPic:      item.HeadPic,
			Photos:       item.Photos,
			Intro:        item.Intro,
			Stores:       stores,
			CreatedAt:    item.CreatedAt.Format(carbon.DateTimeLayout),
			Status:       definition.GoodsStatus(item.Status),
			Remark:       item.Remark,
			StoreIds:     storeIds,
			PaymentPlans: payPlans,
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
		SetPaymentPlans(req.ParsePaymentPlans()).
		SetStatus(definition.GoodsStatusOnline.Value()).
		Save(b.ctx)
	if err != nil {
		return err
	}

	// 创建商品配置门店的对应数据
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
		SetPaymentPlans(req.ParsePaymentPlans()).
		SetStatus(definition.GoodsStatusOnline.Value()).
		Save(b.ctx)
	if err != nil {
		return err
	}

	// 直接先删除已配置的门店
	_, _ = ent.Database.StoreGoods.Delete().Where(storegoods.GoodsID(g.ID)).Exec(b.ctx)
	// 创建商品配置门店的对应数据
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
	return detail(item), nil
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
		res = append(res, detail(item))
	}
	return
}

// ListForRider App获取商品列表
func (b *goodsBiz) ListForRider(req *definition.GoodsListForRiderReq) []*definition.GoodsDetail {
	items, _ := b.orm.QueryNotDeleted().Order(ent.Desc(goods.FieldWeight)).
		Where(
			goods.HasStoresWith(
				storegoods.DeletedAtIsNil(),
				storegoods.HasStoreWith(store.CityID(req.CityID)),
			),
		).
		WithStores(
			func(q *ent.StoreGoodsQuery) {
				q.Where(storegoods.DeletedAtIsNil())
			},
		).All(b.ctx)
	res := make([]*definition.GoodsDetail, 0, len(items))
	for _, v := range items {
		res = append(res, detail(v))
	}
	return res
}
