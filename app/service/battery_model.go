// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/batterymodel"
	"github.com/auroraride/aurservd/pkg/snag"
)

type batteryModelService struct {
	ctx      context.Context
	modifier *model.Modifier
	orm      *ent.BatteryModelClient
}

func NewBatteryModel() *batteryModelService {
	return &batteryModelService{
		ctx: context.Background(),
		orm: ent.Database.BatteryModel,
	}
}

func NewBatteryModelWithModifier(m *model.Modifier) *batteryModelService {
	s := NewBatteryModel()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

// ListModel 列举电池型号
func (s *batteryModelService) ListModel() (res *model.ItemListRes) {
	res = new(model.ItemListRes)
	var items []model.BatteryModel
	s.orm.Query().
		Select(batterymodel.FieldID, batterymodel.FieldModel).
		ScanX(s.ctx, &items)
	model.SetItemListResItems[model.BatteryModel](res, items)
	return
}

// CreateModel 创建电池型号
func (s *batteryModelService) CreateModel(req *model.BatteryModelReq) model.BatteryModel {
	// 判断电池型号是否合法
	if match, _ := regexp.Match(`^[0-9]+(\.?[0-9]+)?V[0-9]+(\.?[0-9]+)?AH$`, []byte(req.Model)); !match {
		snag.Panic("电池型号名称校验失败")
	}
	// 查找同型号电池是否存在
	if s.orm.Query().
		Where(batterymodel.Model(req.Model)).
		ExistX(s.ctx) {
		snag.Panic("电池型号已存在")
	}
	// 创建电池型号
	item := s.orm.Create().
		SetModel(req.Model).
		SaveX(s.ctx)
	return model.BatteryModel{
		ID:    item.ID,
		Model: item.Model,
	}
}

// Query 查询电池型号
func (s *batteryModelService) Query(t any) (*ent.BatteryModel, error) {
	q := s.orm.Query()
	switch u := t.(type) {
	case uint64:
		q.Where(batterymodel.ID(u))
	case string:
		q.Where(batterymodel.Model(strings.ToUpper(u)))
	default:
		return nil, errors.New("参数类型错误")
	}
	return q.First(s.ctx)
}

// QueryIDs 根据ID查询电池型号
func (s *batteryModelService) QueryIDs(ids []uint64) ([]*ent.BatteryModel, error) {
	return s.orm.Query().Where(batterymodel.IDIn(ids...)).All(s.ctx)
}

// Models 列出所有电池型号
func (s *batteryModelService) Models() []string {
	items, _ := s.orm.Query().All(s.ctx)
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = item.Model
	}
	return out
}

// List 列表
func (s *batteryModelService) List(req *model.BatteryModelListReq) (res *model.PaginationRes, err error) {
	q := s.orm.Query()
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.BatteryModel) (result *model.BatteryModelDetail) {
		return s.detail(item)
	}), nil
}

// detail 详情数据
func (s *batteryModelService) detail(item *ent.BatteryModel) (res *model.BatteryModelDetail) {
	res = &model.BatteryModelDetail{
		ID:       item.ID,
		Model:    fmt.Sprintf("%dV%dAH", item.Voltage, item.Capacity),
		Voltage:  item.Voltage,
		Capacity: item.Capacity,
	}
	return
}

// queryById 通过ID查询结果
func (s *batteryModelService) queryById(id uint64) (item *ent.BatteryModel, err error) {
	return s.orm.Query().Where(batterymodel.ID(id)).First(s.ctx)
}

// queryByModel 通过型号查询结果
func (s *batteryModelService) queryByModel(model string) (item *ent.BatteryModel, err error) {
	return s.orm.Query().Where(batterymodel.Model(model)).First(s.ctx)
}

// queryByModelNotSelf 查询非自身同型号结果
func (s *batteryModelService) queryByModelNotSelf(model string, id uint64) (item *ent.BatteryModel, err error) {
	return s.orm.Query().Where(batterymodel.Model(model), batterymodel.IDNotIn(id)).First(s.ctx)
}

// Detail 详情
func (s *batteryModelService) Detail(id uint64) (*model.BatteryModelDetail, error) {
	g, _ := s.queryById(id)
	if g == nil {
		return nil, errors.New("数据不存在")
	}
	return s.detail(g), nil
}

// Create 创建
func (s *batteryModelService) Create(req *model.BatteryModelCreateReq) (err error) {
	batModel := fmt.Sprintf("%dV%dAH", req.Voltage, req.Capacity)
	bm, _ := s.queryByModel(batModel)
	if bm != nil {
		return errors.New("电池型号已存在")
	}
	_, err = s.orm.Create().
		SetVoltage(req.Voltage).
		SetCapacity(req.Capacity).
		SetModel(batModel).
		Save(s.ctx)
	if err != nil {
		return err
	}
	return
}

// Modify 编辑
func (s *batteryModelService) Modify(req *model.BatteryModelModifyReq) (err error) {
	bm, _ := s.queryById(req.ID)
	if bm == nil {
		return errors.New("数据不存在")
	}

	batModel := fmt.Sprintf("%dV%dAH", req.Voltage, req.Capacity)

	sbm, _ := s.queryByModelNotSelf(batModel, req.ID)
	if sbm != nil {
		return errors.New("电池型号已存在")
	}

	_, err = s.orm.UpdateOneID(req.ID).
		SetVoltage(req.Voltage).
		SetCapacity(req.Capacity).
		SetModel(batModel).
		Save(s.ctx)
	if err != nil {
		return err
	}

	return
}

// Delete 删除
func (s *batteryModelService) Delete(id uint64) (err error) {
	bm, _ := s.queryById(id)
	if bm == nil {
		return errors.New("数据不存在")
	}
	err = s.orm.DeleteOne(bm).Exec(s.ctx)
	if err != nil {
		return err
	}
	return
}

// SelectionModels 电池型号筛选项
func (s *batteryModelService) SelectionModels() (res []model.SelectOption) {
	items, _ := s.orm.Query().All(s.ctx)
	res = make([]model.SelectOption, 0)
	for _, item := range items {
		res = append(res, model.SelectOption{
			Value: item.ID,
			Label: item.Model,
		})
	}
	return
}

func (s *batteryModelService) QueryModels(m []string) ([]*ent.BatteryModel, error) {
	return ent.Database.BatteryModel.
		Query().
		Where(
			batterymodel.ModelIn(m...),
		).
		All(context.Background())
}
