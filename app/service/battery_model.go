// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/batterymodel"
    "github.com/auroraride/aurservd/pkg/snag"
    "regexp"
    "strings"
)

type batteryModelService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
    orm          *ent.BatteryModelClient
}

func NewBatteryModel() *batteryModelService {
    return &batteryModelService{
        ctx: context.Background(),
        orm: ent.Database.BatteryModel,
    }
}

func NewBatteryModelWithRider(r *ent.Rider) *batteryModelService {
    s := NewBatteryModel()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewBatteryModelWithModifier(m *model.Modifier) *batteryModelService {
    s := NewBatteryModel()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewBatteryModelWithEmployee(e *ent.Employee) *batteryModelService {
    s := NewBatteryModel()
    if e != nil {
        s.employee = e
        s.employeeInfo = &model.Employee{
            ID:    e.ID,
            Name:  e.Name,
            Phone: e.Phone,
        }
        s.ctx = context.WithValue(s.ctx, "employee", s.employeeInfo)
    }
    return s
}

// List 列举电池型号
func (s *batteryModelService) List() (res *model.ItemListRes) {
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
        snag.Panic("参数错误")
    }
    return q.First(s.ctx)
}

func (s *batteryModelService) QueryX(t any) *ent.BatteryModel {
    bm, _ := s.Query(t)
    if bm == nil {
        snag.Panic("未找到电池型号")
    }
    return bm
}

// QueryIDs 根据ID查询电池型号
func (s *batteryModelService) QueryIDs(ids []uint64) []*ent.BatteryModel {
    return s.orm.Query().Where(batterymodel.IDIn(ids...)).AllX(s.ctx)
}

func (s *batteryModelService) QueryModelsX(models []string) []*ent.BatteryModel {
    items, _ := s.orm.Query().Where(batterymodel.ModelIn(models...)).All(s.ctx)
    if len(items) != len(models) {
        snag.Panic("电池型号查询失败")
    }
    return items
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

// Delete 删除电池型号
func (s *batteryModelService) Delete(req *model.BatteryModelReq) {
    s.orm.DeleteOne(s.QueryX(req.Model)).ExecX(s.ctx)
}
