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
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "github.com/jinzhu/copier"
    "time"
)

type planService struct {
    ctx context.Context
    orm *ent.PlanClient
}

func NewPlan() *planService {
    return &planService{
        ctx: context.Background(),
        orm: ar.Ent.Plan,
    }
}

// Query 查找骑士卡
func (s *planService) Query(id uint64) *ent.Plan {
    item, err := s.orm.QueryNotDeleted().Where(plan.ID(id)).Only(s.ctx)
    if err != nil || item == nil {
        snag.Panic("未找到有效的骑士卡")
    }
    return item
}

// CreatePlan 创建骑士卡
func (s *planService) CreatePlan(m *model.Modifier, req *model.PlanCreateReq) {
    start, _ := time.Parse(carbon.DateLayout, req.Start)
    end, _ := time.Parse(carbon.DateLayout, req.End)
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
        SetLastModifier(m).
        SaveX(s.ctx)
}

// UpdateEnable 修改骑士卡状态
func (s *planService) UpdateEnable(m *model.Modifier, req *model.PlanEnableModifyReq) {
    item := s.Query(req.ID)
    s.orm.UpdateOne(item).SetEnable(*req.Enable).SetLastModifier(m).SaveX(s.ctx)
}

// Delete 软删除骑士卡
func (s *planService) Delete(m *model.Modifier, req *model.IDParamReq) {
    item := s.Query(req.ID)
    s.orm.UpdateOne(item).SetDeletedAt(time.Now()).SetLastModifier(m).SaveX(s.ctx)
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
        _ = copier.Copy(&out[i], item)
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
