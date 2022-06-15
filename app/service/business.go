// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
)

// 门店业务处理专用
type businessService struct {
    ctx      context.Context
    employee *ent.Employee
}

func NewBusiness() *businessService {
    return &businessService{
        ctx: context.Background(),
    }
}

func NewBusinessWithEmployee(e *ent.Employee) *businessService {
    s := NewBusiness()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// CheckCity 检查城市
func (s *businessService) CheckCity(cityID uint64) {
    if s.employee.Edges.Store.CityID != cityID {
        snag.Panic("不能跨城市操作")
    }
}

// Detail 获取骑手订阅业务详情
func (s *businessService) Detail(id uint64) (res model.SubscribeBusiness) {
    r, err := NewRider().QueryForBusiness(id)
    if err != nil {
        snag.Panic(err)
    }
    // 获取最近的订阅
    sub := NewSubscribe().RecentDetail(r.ID)

    if sub == nil {
        snag.Panic("未找到有效订阅")
    }

    s.CheckCity(sub.City.ID)

    ic := r.Edges.Person.IDCardNumber
    res = model.SubscribeBusiness{
        ID:           r.ID,
        Status:       sub.Status,
        Name:         r.Edges.Person.Name,
        Phone:        r.Phone,
        IDCardNumber: ic[len(ic)-4:],
        Voltage:      sub.Voltage,
        SubscribeID:  sub.ID,
    }

    if sub.Enterprise != nil {
        res.EnterpriseName = sub.Enterprise.Name
    }

    if sub.Plan != nil {
        res.PlanName = sub.Plan.Name
    }

    res.Business = sub.Business
    return
}
