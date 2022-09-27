// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "strconv"
    "strings"
)

type businessEmployeeService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
}

func NewBusinessEmployee() *businessEmployeeService {
    return &businessEmployeeService{
        ctx: context.Background(),
    }
}

func NewBusinessEmployeeWithEmployee(e *ent.Employee) *businessEmployeeService {
    s := NewBusinessEmployee()
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

// Inactive 获取骑手待激活订阅详情
func (s *businessEmployeeService) Inactive(qr string) (*model.SubscribeActiveInfo, *ent.Subscribe) {
    if strings.HasPrefix(qr, "SUBSCRIBE:") {
        qr = strings.ReplaceAll(qr, "SUBSCRIBE:", "")
    }
    id, _ := strconv.ParseUint(strings.TrimSpace(qr), 10, 64)
    return NewBusinessRiderWithEmployee(s.employee).Inactive(id)
}

// Active 激活订阅
func (s *businessEmployeeService) Active(req *model.QRPostReq) {
    NewBusinessRiderWithEmployee(s.employee).Active(s.Inactive(req.Qrcode))
}
