// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-12, by aurb

package biz

import (
	"context"
	"errors"
	"fmt"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
)

type batteryModelBiz struct {
	orm      *ent.BatteryModelClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewBatteryModel() *batteryModelBiz {
	return &batteryModelBiz{
		orm: ent.Database.BatteryModel,
		ctx: context.Background(),
	}
}

func NewBatteryModelWithModifier(m *model.Modifier) *batteryModelBiz {
	s := NewBatteryModel()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// List 列表
func (b *batteryModelBiz) List(req *definition.BatteryModelListReq) (res *model.PaginationRes, err error) {
	q := b.orm.Query()

	if req.Type != nil {
		q.Where(batterymodel.Type(req.Type.Value()))
	}

	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.BatteryModel) (result *definition.BatteryModelDetail) {
		return b.detail(item)
	})
	return
}

// detail 详情数据
func (b *batteryModelBiz) detail(item *ent.BatteryModel) (res *definition.BatteryModelDetail) {
	res = &definition.BatteryModelDetail{
		ID:       item.ID,
		Model:    fmt.Sprintf("%dV%dAH", item.Voltage, item.Capacity),
		Type:     definition.BatteryModelType(item.Type),
		Voltage:  item.Voltage,
		Capacity: item.Capacity,
	}
	return
}

// queryById 通过ID查询结果
func (b *batteryModelBiz) queryById(id uint64) (item *ent.BatteryModel, err error) {
	return b.orm.Query().Where(batterymodel.ID(id)).First(b.ctx)
}

// queryByModel 通过型号查询结果
func (b *batteryModelBiz) queryByModel(model string) (item *ent.BatteryModel, err error) {
	return b.orm.Query().Where(batterymodel.Model(model)).First(b.ctx)
}

// queryByModelNotSelf 查询非自身同型号结果
func (b *batteryModelBiz) queryByModelNotSelf(model string, id uint64) (item *ent.BatteryModel, err error) {
	return b.orm.Query().Where(batterymodel.Model(model), batterymodel.IDNotIn(id)).First(b.ctx)
}

// Detail 详情
func (b *batteryModelBiz) Detail(id uint64) (*definition.BatteryModelDetail, error) {
	g, _ := b.queryById(id)
	if g == nil {
		return nil, errors.New("数据不存在")
	}
	return b.detail(g), nil
}

// Create 创建
func (b *batteryModelBiz) Create(req *definition.BatteryModelCreateReq) (err error) {
	batModel := fmt.Sprintf("%dV%dAH", req.Voltage, req.Capacity)
	bm, _ := b.queryByModel(batModel)
	if bm != nil {
		return errors.New("电池型号已存在")
	}
	_, err = b.orm.Create().
		SetType(uint8(req.Type)).
		SetVoltage(req.Voltage).
		SetCapacity(req.Capacity).
		SetModel(batModel).
		Save(b.ctx)
	if err != nil {
		return err
	}
	return
}

// Modify 编辑
func (b *batteryModelBiz) Modify(req *definition.BatteryModelModifyReq) (err error) {
	bm, _ := b.queryById(req.ID)
	if bm == nil {
		return errors.New("数据不存在")
	}

	batModel := fmt.Sprintf("%dV%dAH", req.Voltage, req.Capacity)

	sbm, _ := b.queryByModelNotSelf(batModel, req.ID)
	if sbm != nil {
		return errors.New("电池型号已存在")
	}

	_, err = b.orm.UpdateOneID(req.ID).
		SetType(uint8(req.Type)).
		SetVoltage(req.Voltage).
		SetCapacity(req.Capacity).
		SetModel(batModel).
		Save(b.ctx)
	if err != nil {
		return err
	}

	return
}

// Delete 删除仓库
func (b *batteryModelBiz) Delete(id uint64) (err error) {
	bm, _ := b.queryById(id)
	if bm == nil {
		return errors.New("数据不存在")
	}
	err = b.orm.DeleteOne(bm).Exec(b.ctx)
	if err != nil {
		return err
	}
	return
}
