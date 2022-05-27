// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/batterymodel"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "github.com/jinzhu/copier"
    "time"
)

type planService struct {
    ctx      context.Context
    orm      *ent.PlanClient
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewPlan() *planService {
    return &planService{
        ctx: context.Background(),
        orm: ar.Ent.Plan,
    }
}

func NewPlanWithRider(rider *ent.Rider) *planService {
    s := NewPlan()
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

func NewPlanWithModifier(m *model.Modifier) *planService {
    s := NewPlan()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// Query 查找骑士卡
func (s *planService) Query(id uint64) *ent.Plan {
    item, err := s.orm.QueryNotDeleted().Where(plan.ID(id)).Only(s.ctx)
    if err != nil || item == nil {
        snag.Panic("未找到有效的骑士卡")
    }
    return item
}

// QueryEffectiveWithID 获取当前生效的骑行卡
func (s *planService) QueryEffectiveWithID(id uint64) *ent.Plan {
    now := time.Now()
    item, err := s.orm.QueryNotDeleted().
        Where(
            plan.Enable(true),
            plan.ID(id),
            plan.StartLTE(now),
            plan.EndGTE(now),
        ).
        Only(s.ctx)
    if err != nil || item == nil {
        snag.Panic("未找到有效的骑士卡")
    }
    return item
}

// Duplicated 查询骑士卡是否冲突
func (s *planService) Duplicated(cities, models []uint64, start, end time.Time, days uint) bool {
    for _, cityID := range cities {
        for _, modelID := range models {
            if exists, _ := s.orm.QueryNotDeleted().
                Where(
                    plan.Enable(true),
                    plan.Days(days),
                    plan.HasCitiesWith(city.ID(cityID)),
                    plan.HasPmsWith(batterymodel.ID(modelID)),
                    plan.Or(
                        plan.StartLTE(start),
                        plan.EndGTE(end),
                    ),
                ).Exist(s.ctx); exists {
                snag.Panic("骑士卡冲突")
            }
        }
    }
    return false
}

// CreatePlan 创建骑士卡
func (s *planService) CreatePlan(req *model.PlanCreateReq) {
    start, _ := time.Parse(carbon.DateLayout, req.Start)
    end, _ := time.Parse(carbon.DateLayout, req.End)
    // 查询是否重复
    if s.Duplicated(req.Cities, req.Models, start, end, req.Days) {
        snag.Panic("骑士卡冲突")
    }
    original := req.Original
    if original == 0 {
        original = req.Price
    }
    s.orm.Create().
        AddCityIDs(req.Cities...).
        AddPmIDs(req.Models...).
        SetStart(start).
        SetEnd(end).
        SetDays(req.Days).
        SetEnable(req.Enable).
        SetName(req.Name).
        SetPrice(req.Price).
        SetCommission(req.Commission).
        SetOriginal(original).
        SetDesc(req.Desc).
        SaveX(s.ctx)
}

// UpdateEnable 修改骑士卡状态
func (s *planService) UpdateEnable(req *model.PlanEnableModifyReq) {
    item := s.Query(req.ID)
    s.orm.UpdateOne(item).SetEnable(*req.Enable).SaveX(s.ctx)
}

// Delete 软删除骑士卡
func (s *planService) Delete(req *model.IDParamReq) {
    s.orm.SoftDeleteOne(s.Query(req.ID)).SaveX(s.ctx)
}

// List 列举骑士卡
func (s *planService) List(req *model.PlanListReq) *model.PaginationRes {
    res := new(model.PaginationRes)
    q := s.orm.QueryNotDeleted().WithCities().WithPms()
    if req.CityID != nil {
        q.Where(plan.HasCitiesWith(city.ID(*req.CityID)))
    }
    if req.Name != nil {
        q.Where(plan.NameContainsFold(*req.Name))
    }
    if req.Enable != nil {
        q.Where(plan.Enable(*req.Enable))
    }
    res.Pagination = q.PaginationResult(req.PaginationReq)
    items := q.AllX(s.ctx)
    out := make([]model.PlanItemRes, len(items))
    for i, item := range items {
        _ = copier.CopyWithOption(&out[i], item, copier.Option{Converters: []copier.TypeConverter{
            {
                SrcType: time.Time{},
                DstType: copier.String,
                Fn: func(src interface{}) (interface{}, error) {
                    t, ok := src.(time.Time)
                    if !ok {
                        return "", nil
                    }
                    return t.Format(carbon.DateLayout), nil
                },
            },
        }})
        cities := make([]model.City, len(item.Edges.Cities))
        for ci, c := range item.Edges.Cities {
            cities[ci] = model.City{
                ID:   c.ID,
                Name: c.Name,
            }
        }
        out[i].Cities = cities

        models := make([]model.BatteryModel, len(item.Edges.Pms))
        for mi, pm := range item.Edges.Pms {
            models[mi] = model.BatteryModel{
                ID:       pm.ID,
                Voltage:  pm.Voltage,
                Capacity: pm.Capacity,
            }
        }
        out[i].Models = models
    }
    res.Items = out
    return res
}

// RiderList 获取骑士卡列表
func (s *planService) RiderList(req *model.PlanListRiderReq) (res []model.RiderPlanItem) {
    now := time.Now()
    items := s.orm.QueryNotDeleted().
        Where(
            plan.Enable(true),
            plan.StartLTE(now),
            plan.EndGTE(now),
            plan.HasPmsWith(batterymodel.Voltage(*req.Voltage)),
        ).
        Order(ent.Asc(plan.FieldDays)).
        AllX(s.ctx)
    res = make([]model.RiderPlanItem, len(items))
    for i, item := range items {
        res[i] = model.RiderPlanItem{
            ID:       item.ID,
            Name:     item.Name,
            Price:    item.Price,
            Days:     item.Days,
            Original: item.Original,
            Desc:     item.Desc,
        }
    }
    return
}
