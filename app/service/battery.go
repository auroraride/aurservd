// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/batterymodel"
    "github.com/auroraride/aurservd/pkg/snag"
)

type batteryService struct {
    ctx      context.Context
    modifier *model.Modifier
    orm      *ent.BatteryModelClient
}

func NewBattery(modifier *model.Modifier) *batteryService {
    return &batteryService{
        ctx:      context.Background(),
        modifier: modifier,
        orm:      ar.Ent.BatteryModel,
    }
}

// ListModels 列举电池型号
func (s *batteryService) ListModels() (res *model.ItemListRes) {
    res = new(model.ItemListRes)
    var items []model.BatteryModel
    s.orm.Query().
        Order(ent.Desc(batterymodel.FieldCreatedAt)).
        Select(batterymodel.FieldCapacity, batterymodel.FieldID, batterymodel.FieldVoltage).
        ScanX(s.ctx, &items)

    model.SetItemListResItems[model.BatteryModel](res, items)
    return
}

// CreateModel 创建电池型号
func (s *batteryService) CreateModel(req *model.BatteryModelCreateReq) model.BatteryModel {
    // 查找同型号电池是否存在
    if s.orm.Query().
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
        SetLastModifier(s.modifier).
        SetCreator(s.modifier).
        SaveX(s.ctx)
    return model.BatteryModel{
        ID:       item.ID,
        Voltage:  item.Voltage,
        Capacity: item.Capacity,
    }
}
