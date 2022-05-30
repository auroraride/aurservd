// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/order"
    "github.com/auroraride/aurservd/pkg/snag"
    "strconv"
    "strings"
    "time"
)

type employeeOrderService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    orm      *ent.OrderClient
    employee *model.Employee
}

func NewEmployeeOrder() *employeeOrderService {
    return &employeeOrderService{
        ctx: context.Background(),
        orm: ar.Ent.Order,
    }
}

func NewEmployeeOrderWithRider(rider *ent.Rider) *employeeOrderService {
    s := NewEmployeeOrder()
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

func NewEmployeeOrderWithModifier(m *model.Modifier) *employeeOrderService {
    s := NewEmployeeOrder()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewEmployeeOrderWithEmployee(e *model.Employee) *employeeOrderService {
    s := NewEmployeeOrder()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// Active 激活订单
// TODO 后续逻辑
func (s *employeeOrderService) Active(req *model.QRPostReq) {
    id, _ := strconv.ParseUint(strings.TrimSpace(strings.ReplaceAll(req.Qrcode, "ORDER:", "")), 10, 64)
    // 查询订单状态
    o, _ := s.orm.QueryNotDeleted().
        Where(
            order.ID(id),
            order.Status(model.OrderStatusPaid),
            order.RefundAtIsNil(),
            order.TypeIn(model.OrderTypeNewly, model.OrderTypeAgain),
            order.StartAtIsNil(),
        ).
        Only(s.ctx)
    if o == nil {
        snag.Panic("未找到订单")
    }
    o.Update().SetStartAt(time.Now()).SetEmployee(s.employee).SaveX(s.ctx)
}
