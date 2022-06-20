// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/batterymodel"
    "github.com/auroraride/aurservd/pkg/snag"
)

type batteryService struct {
    ctx context.Context
    orm *ent.BatteryModelClient
}

func NewBattery() *batteryService {
    return &batteryService{
        ctx: context.Background(),
        orm: ar.Ent.BatteryModel,
    }
}

func NewBatteryWithModifier(m *model.Modifier) *batteryService {
    s := NewBattery()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    return s
}

// ListModels 列举电池型号
func (s *batteryService) ListModels() (res *model.ItemListRes) {
    res = new(model.ItemListRes)
    var items []model.BatteryModel
    s.orm.QueryNotDeleted().
        Order(ent.Desc(batterymodel.FieldCreatedAt)).
        Select(batterymodel.FieldCapacity, batterymodel.FieldID, batterymodel.FieldVoltage).
        ScanX(s.ctx, &items)

    model.SetItemListResItems[model.BatteryModel](res, items)
    return
}

// CreateModel 创建电池型号
func (s *batteryService) CreateModel(req *model.BatteryModelCreateReq) model.BatteryModel {
    // 查找同型号电池是否存在
    if s.orm.QueryNotDeleted().
        Where(batterymodel.Capacity(req.Capacity)).
        Where(batterymodel.Voltage(req.Voltage)).
        Where(batterymodel.DeletedAtIsNil()).
        ExistX(s.ctx) {
        snag.Panic("电池型号已存在")
    }
    // 创建电池型号
    item := s.orm.Create().
        SetVoltage(req.Voltage).
        SetCapacity(req.Capacity).
        SaveX(s.ctx)
    return model.BatteryModel{
        ID:       item.ID,
        Voltage:  item.Voltage,
        Capacity: item.Capacity,
    }
}

// QueryIDs 根据ID查询电池型号
func (s *batteryService) QueryIDs(ids []uint64) []*ent.BatteryModel {
    return s.orm.QueryNotDeleted().Where(batterymodel.IDIn(ids...)).AllX(s.ctx)
}

// ListVoltages 列出所有型号电压
func (s *batteryService) ListVoltages(excludes ...float64) []float64 {
    var items []float64
    q := s.orm.QueryNotDeleted().
        Select(batterymodel.FieldVoltage)
    if len(excludes) > 0 {
        q.Where(batterymodel.VoltageNotIn(excludes...))
    }
    q.GroupBy(batterymodel.FieldVoltage).ScanX(s.ctx, &items)
    return items
}

// VoltageName 获取电池电压型号名称
func (s *batteryService) VoltageName(voltage float64) string {
    return fmt.Sprintf("%.2gV电池", voltage)
}
