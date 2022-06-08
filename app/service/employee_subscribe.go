// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-30
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/snag"
    "strconv"
    "strings"
    "time"
)

type employeeSubscribeService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *model.Employee
}

func NewEmployeeSubscribe() *employeeSubscribeService {
    return &employeeSubscribeService{
        ctx: context.Background(),
    }
}

func NewEmployeeSubscribeWithRider(r *ent.Rider) *employeeSubscribeService {
    s := NewEmployeeSubscribe()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewEmployeeSubscribeWithModifier(m *model.Modifier) *employeeSubscribeService {
    s := NewEmployeeSubscribe()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewEmployeeSubscribeWithEmployee(e *model.Employee) *employeeSubscribeService {
    s := NewEmployeeSubscribe()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// Active 激活订单
// TODO 后续逻辑
func (s *employeeSubscribeService) Active(req *model.QRPostReq) {
    id, _ := strconv.ParseUint(strings.TrimSpace(strings.ReplaceAll(req.Qrcode, "SUBSCRIBE:", "")), 10, 64)
    // 查询订单状态
    sub, _ := ar.Ent.Subscribe.QueryNotDeleted().
        Where(
            subscribe.ID(id),
            subscribe.RefundAtIsNil(),
            subscribe.StartAtIsNil(),
            subscribe.Or(
                subscribe.Type(0),
                subscribe.TypeIn(model.OrderTypeNewly, model.OrderTypeAgain),
            ),
        ).
        WithInitialOrder().
        Only(s.ctx)
    if sub == nil {
        snag.Panic("未找到骑士卡")
    }
    if sub.EnterpriseID == nil && sub.Edges.InitialOrder.Status == model.OrderStatusRefundPending {
        snag.Panic("骑士卡已申请退款")
    }
    sub.Update().
        SetStatus(model.SubscribeStatusUsing).
        SetStartAt(time.Now()).
        SetEmployeeID(s.employee.ID).
        SaveX(s.ctx)
    if sub.EnterpriseID != nil {
        _, _ = NewEnterpriseStatement().Current(*sub.EnterpriseID).Update().AddRiderNumber(1).Save(s.ctx)
    }
}
