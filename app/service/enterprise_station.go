// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/enterprisestation"
    "github.com/auroraride/aurservd/pkg/snag"
)

type enterpriseStationService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
    orm      *ent.EnterpriseStationClient
}

func NewEnterpriseStation() *enterpriseStationService {
    return &enterpriseStationService{
        ctx: context.Background(),
        orm: ent.Database.EnterpriseStation,
    }
}

func NewEnterpriseStationWithRider(r *ent.Rider) *enterpriseStationService {
    s := NewEnterpriseStation()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewEnterpriseStationWithModifier(m *model.Modifier) *enterpriseStationService {
    s := NewEnterpriseStation()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewEnterpriseStationWithEmployee(e *ent.Employee) *enterpriseStationService {
    s := NewEnterpriseStation()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func (s *enterpriseStationService) Query(id uint64) *ent.EnterpriseStation {
    item, _ := s.orm.QueryNotDeleted().Where(enterprisestation.ID(id)).First(s.ctx)
    if item == nil {
        snag.Panic("未找到站点")
    }
    return item
}

// Create 创建站点
func (s *enterpriseStationService) Create(req *model.EnterpriseStationCreateReq) uint64 {
    es := s.orm.Create().SetEnterpriseID(req.EnterpriseID).SetName(req.Name).SaveX(s.ctx)
    return es.ID
}

// Modify 修改站点
func (s *enterpriseStationService) Modify(req *model.EnterpriseStationModifyReq) {
    s.Query(req.ID).Update().SetName(req.Name).SaveX(s.ctx)
}

func (s *enterpriseStationService) List(req *model.EnterpriseStationListReq) (res []model.EnterpriseStation) {
    res = make([]model.EnterpriseStation, 0)
    items, _ := s.orm.QueryNotDeleted().Where(enterprisestation.EnterpriseID(req.EnterpriseID)).All(s.ctx)
    for _, item := range items {
        res = append(res, model.EnterpriseStation{
            ID:   item.ID,
            Name: item.Name,
        })
    }
    return
}
