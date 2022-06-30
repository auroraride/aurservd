// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-16
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/business"
    "github.com/auroraride/aurservd/pkg/snag"
)

type businessLogService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
    orm      *ent.BusinessClient
    creator  *ent.BusinessCreate
}

func NewBusinessLog(sub *ent.Subscribe) *businessLogService {
    s := &businessLogService{
        ctx:     context.Background(),
        orm:     ent.Database.Business,
        creator: ent.Database.Business.Create(),
    }
    s.setSubscribe(sub)
    return s
}

func NewBusinessLogWithModifier(m *model.Modifier, sub *ent.Subscribe) *businessLogService {
    s := NewBusinessLog(sub)
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    s.creator.SetLastModifier(m)
    return s
}

func NewBusinessLogWithEmployee(e *ent.Employee, sub *ent.Subscribe) *businessLogService {
    s := NewBusinessLog(sub)
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    s.creator.SetEmployee(e).SetStoreID(e.Edges.Store.ID)
    return s
}

func (s *businessLogService) setSubscribe(sub *ent.Subscribe) {
    s.creator.SetRiderID(sub.RiderID).
        SetSubscribeID(sub.ID).
        SetCityID(sub.CityID).
        SetNillableEnterpriseID(sub.EnterpriseID).
        SetNillableStationID(sub.StationID).
        SetNillablePlanID(sub.PlanID)
}

func (s *businessLogService) Save(typ business.Type) (*ent.Business, error) {
    return s.creator.SetType(typ).Save(s.ctx)
}

func (s *businessLogService) SaveAsync(typ business.Type) {
    go func() {
        _, _ = s.Save(typ)
    }()
}

func (s *businessLogService) SaveX(typ business.Type) *ent.Business {
    biz, err := s.Save(typ)
    if err != nil {
        snag.Panic("日志保存失败")
    }
    return biz
}
