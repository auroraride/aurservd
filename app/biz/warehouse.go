// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-10, by Jorjan

package biz

import (
	"context"
	"errors"
	"fmt"

	"github.com/lithammer/shortuuid/v4"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
)

type warehouseBiz struct {
	orm      *ent.WarehouseClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewWarehouse() *warehouseBiz {
	return &warehouseBiz{
		orm: ent.Database.Warehouse,
		ctx: context.Background(),
	}
}

func NewWarehouseWithModifier(m *model.Modifier) *warehouseBiz {
	s := NewWarehouse()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// List 仓库列表
func (b *warehouseBiz) List(req *definition.WareHouseListReq) (res []*definition.WarehouseDetail, err error) {
	res = make([]*definition.WarehouseDetail, 0)

	q := b.orm.QueryNotDeleted().WithCity()

	if req.CityID != nil {
		q.Where(warehouse.CityID(*req.CityID))
	}

	if req.Keyword != nil {
		q.Where(warehouse.NameContains(*req.Keyword))
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
func (b *warehouseBiz) detail(item *ent.Warehouse) (res *definition.WarehouseDetail) {
	res = &definition.WarehouseDetail{
		ID:      item.ID,
		Name:    item.Name,
		Lng:     item.Lng,
		Lat:     item.Lat,
		Address: item.Address,
		QRCode:  fmt.Sprintf("WAREHOUSE:%s", item.Sn),
	}

	if item.Edges.City != nil {
		res.City = model.City{
			ID:   item.Edges.City.ID,
			Name: item.Edges.City.Name,
		}
	}
	return
}

// Detail 查询仓库详情
func (b *warehouseBiz) Detail(id uint64) (*definition.WarehouseDetail, error) {
	g, _ := b.queryById(id)
	if g == nil {
		return nil, errors.New("仓库不存在")
	}
	return b.detail(g), nil
}

// Create 创建仓库
func (b *warehouseBiz) Create(req *definition.WarehouseCreateReq) (err error) {
	_, err = b.orm.Create().
		SetName(req.Name).
		SetCityID(req.CityID).
		SetAddress(req.Address).
		SetRemark(req.Remark).
		SetLat(req.Lat).
		SetLng(req.Lng).
		SetSn(shortuuid.New()).
		Save(b.ctx)
	if err != nil {
		return err
	}
	return
}

// queryById 通过ID查询仓库
func (b *warehouseBiz) queryById(id uint64) (item *ent.Warehouse, err error) {
	return b.orm.QueryNotDeleted().Where(warehouse.ID(id)).First(b.ctx)
}

// Modify 编辑仓库
func (b *warehouseBiz) Modify(req *definition.WarehouseModifyReq) (err error) {
	g, _ := b.queryById(req.ID)
	if g == nil {
		return errors.New("仓库不存在")
	}

	_, err = b.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetName(req.Name).
		SetCityID(req.CityID).
		SetAddress(req.Address).
		SetRemark(req.Remark).
		SetLat(req.Lat).
		SetLng(req.Lng).
		Save(b.ctx)
	if err != nil {
		return err
	}

	return
}

// Delete 删除仓库
func (b *warehouseBiz) Delete(id uint64) (err error) {
	g, _ := b.queryById(id)
	if g == nil {
		return errors.New("仓库不存在")
	}
	_, err = b.orm.SoftDeleteOne(g).Save(b.ctx)
	if err != nil {
		return err
	}
	return
}
