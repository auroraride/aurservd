// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-05-29, by aurb

package biz

import (
	"context"
	"errors"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/goods"
	"github.com/auroraride/aurservd/internal/ent/store"
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

func (s *goodsBiz) List(req *definition.GoodsListReq) *model.PaginationRes {
	query := s.orm.Query().Order(ent.Desc(goods.FieldWeight))
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

func toGoodsDetail(item *ent.Goods) *definition.GoodsDetail {
	// 查询配置的门店信息
	var stores []model.Store
	sis, _ := ent.Database.Store.Query().Where(store.IDIn(item.StoreIds...)).All(context.Background())
	for _, si := range sis {
		stores = append(stores, model.Store{
			ID:   si.ID,
			Name: si.Name,
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
			StoreIds:  item.StoreIds,
			Stores:    stores,
			CreatedAt: item.CreatedAt,
			Status:    definition.GoodsStatus(item.Status),
			Remark:    item.Remark,
		},
	}
}

func (s *goodsBiz) Create(req *definition.GoodsCreateReq) (err error) {
	sn := tools.NewUnique().NewSN()
	_, err = s.orm.Create().
		SetSn(sn).
		SetName(req.Name).
		SetType(definition.GoodsTypeEbike.Value()).
		SetPrice(req.Price).
		SetWeight(req.Weight).
		SetHeadPic(req.HeadPic).
		SetPhotos(req.Photos).
		SetIntro(req.Intro).
		SetStoreIds(req.StoreIds).
		SetRemark(req.Remark).
		SetStatus(definition.GoodsStatusOnline.Value()).
		Save(s.ctx)
	return
}

func (s *goodsBiz) Delete(id uint64) (err error) {
	g, _ := s.orm.QueryNotDeleted().Where(goods.ID(id)).First(s.ctx)
	if g == nil {
		return errors.New("商品不存在")
	}
	_, err = s.orm.SoftDeleteOneID(id).Save(s.ctx)
	if err != nil {
		return err
	}
	return
}

func (s *goodsBiz) Modify(req *definition.GoodsModifyReq) (err error) {
	g, _ := s.orm.QueryNotDeleted().Where(goods.ID(req.ID)).First(s.ctx)
	if g == nil {
		return errors.New("商品不存在")
	}

	_, err = s.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetType(definition.GoodsTypeEbike.Value()).
		SetPrice(req.Price).
		SetWeight(req.Weight).
		SetHeadPic(req.HeadPic).
		SetPhotos(req.Photos).
		SetIntro(req.Intro).
		SetStoreIds(req.StoreIds).
		SetRemark(req.Remark).
		SetStatus(definition.GoodsStatusOnline.Value()).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return
}

func (s *goodsBiz) Detail(id uint64) (*definition.GoodsDetail, error) {
	item, err := s.orm.Query().Where(goods.ID(id)).First(s.ctx)
	if err != nil {
		return nil, err
	}
	return toGoodsDetail(item), nil
}

func (s *goodsBiz) UpdateStatus(req *definition.GoodsUpdateStatusReq) (err error) {
	g, _ := s.orm.QueryNotDeleted().Where(goods.ID(req.ID)).First(s.ctx)
	if g == nil {
		return errors.New("商品不存在")
	}

	if g.Status == req.Status.Value() {
		return errors.New("商品状态已存在")
	}

	_, err = s.orm.UpdateOneID(req.ID).
		SetStatus(req.Status.Value()).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return
}
