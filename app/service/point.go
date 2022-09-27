// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-27
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
)

type pointService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
    orm          *ent.PointLogClient
}

func NewPoint() *pointService {
    return &pointService{
        ctx: context.Background(),
        orm: ent.Database.PointLog,
    }
}

func NewPointWithRider(r *ent.Rider) *pointService {
    s := NewPoint()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewPointWithModifier(m *model.Modifier) *pointService {
    s := NewPoint()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewPointWithEmployee(e *ent.Employee) *pointService {
    s := NewPoint()
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

// Modify 修改积分
func (s *pointService) Modify(req *model.PointModifyReq) {
    r := NewRider().Query(req.RiderID)
    after := r.Points + req.Points
    if after < 0 {
        snag.Panic("积分余额不能小于0")
    }
    creater := s.orm.Create().SetRiderID(req.RiderID).SetPoints(req.Points).SetReason(req.Reason).SetType(req.Type.Value()).SetAfter(after)
    err := creater.Exec(s.ctx)
    if err != nil {
        snag.Panic(err)
    }
}
