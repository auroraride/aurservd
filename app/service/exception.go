// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-17
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
)

type exceptionService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
    orm      *ent.ExceptionClient
}

func NewException() *exceptionService {
    return &exceptionService{
        ctx: context.Background(),
        orm: ar.Ent.Exception,
    }
}

func NewExceptionWithRider(r *ent.Rider) *exceptionService {
    s := NewException()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewExceptionWithModifier(m *model.Modifier) *exceptionService {
    s := NewException()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewExceptionWithEmployee(e *ent.Employee) *exceptionService {
    s := NewException()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func (s *exceptionService) Setting() model.ExceptionEmployeeSetting {
    set := NewSetting().GetSetting(model.SettingException)
    return model.ExceptionEmployeeSetting{
        Items: NewInventory().ListInventory(model.InventoryListReq{
            Count:    true,
            Transfer: true,
            Purchase: true,
        }),
        Reasons: set.([]interface{}),
    }
}

func (s *exceptionService) Create(req *model.ExceptionEmployeeReq) {
    if (req.Voltage == 0 && req.Name == "") || (req.Voltage != 0 && req.Name != "") {
        snag.Panic("请求参数错误")
    }

    if s.employee.Edges.Store.CityID == nil {
        snag.Panic("门店所在城市异常")
    }

    ec := s.orm.Create().
        SetEmployee(s.employee).
        SetStore(s.employee.Edges.Store).
        SetCityID(*s.employee.Edges.Store.CityID).
        SetNum(req.Num).
        SetReason(req.Reason).
        SetDescription(req.Description).
        SetAttachments(req.Attachments)

    if req.Name != "" {
        ec.SetName(req.Name)
    }
    if req.Voltage != 0 {
        ec.SetName(NewBattery().VoltageName(req.Voltage)).SetVoltage(req.Voltage)
    }

    ec.SaveX(s.ctx)
}
