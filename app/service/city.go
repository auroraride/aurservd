// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/pkg/snag"
)

type cityService struct {
	orm      *ent.CityClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewCity() *cityService {
	return &cityService{
		orm: ent.Database.City,
		ctx: context.Background(),
	}
}

func NewCityWithModifier(m *model.Modifier) *cityService {
	s := NewCity()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

func (s *cityService) Query(id uint64) *ent.City {
	item, err := s.orm.QueryNotDeleted().Where(city.ID(id)).First(s.ctx)
	if err != nil {
		snag.Panic("未找到城市")
	}
	return item
}

// List 获取城市列表
func (s *cityService) List(req *model.CityListReq) (items []*model.CityItem) {
	fields := []string{
		city.FieldID, city.FieldName, city.FieldParentID, city.FieldOpen,
	}
	q := s.orm.QueryNotDeleted().
		Where(city.ParentIDIsNil()).
		WithChildren(func(query *ent.CityQuery) {
			query.Select(fields...).Order(ent.Asc(city.FieldID))
			if req.Status > 0 {
				query.Where(city.Open(req.Status == model.CityStatusOpen))
			}
		}).
		Order(ent.Asc(city.FieldID))
	if req.Status > 0 {
		q.Where(city.HasChildrenWith(city.Open(req.Status == model.CityStatusOpen)))
	}
	cities := q.Select(fields...).AllX(s.ctx)
	items = make([]*model.CityItem, len(cities))
	for i, c := range cities {
		item := &model.CityItem{
			ID:       c.ID,
			Name:     c.Name,
			Children: make([]model.CityItem, len(c.Edges.Children)),
		}
		for n, child := range c.Edges.Children {
			item.Children[n] = model.CityItem{
				ID:   child.ID,
				Open: child.Open,
				Name: child.Name,
			}
		}
		items[i] = item
	}
	return
}

// Modify 修改城市
func (s *cityService) Modify(req *model.CityModifyReq) *bool {
	if exists, _ := s.orm.QueryNotDeleted().Where(city.ID(req.ID), city.ParentIDNotNil()).Exist(context.Background()); !exists {
		snag.Panic("城市ID错误")
	}
	c := s.orm.UpdateOneID(req.ID).SetOpen(*req.Open).SaveX(s.ctx)
	return c.Open
}

// OpenedCities 获取已开通城市列表
func (s *cityService) OpenedCities() []model.CityWithLocation {
	items := s.orm.QueryNotDeleted().
		Where(city.Open(true)).
		Where(city.ParentIDNotNil()).
		AllX(s.ctx)
	res := make([]model.CityWithLocation, len(items))
	for i, item := range items {
		res[i] = model.CityWithLocation{
			ID:   item.ID,
			Name: item.Name,
			Lng:  item.Lng,
			Lat:  item.Lat,
		}
	}
	return res
}

func (s *cityService) NameFromID(id uint64) string {
	p, _ := ent.Database.City.QueryNotDeleted().Where(city.ID(id)).First(s.ctx)
	if p == nil {
		return "-"
	}
	return p.Name
}

func (s *cityService) QueryIDs(ids []uint64) ([]*ent.City, error) {
	return s.orm.QueryNotDeleted().Where(city.IDIn(ids...)).All(s.ctx)
}
