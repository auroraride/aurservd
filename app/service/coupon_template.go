// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-26
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/coupontemplate"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/pkg/snag"
    log "github.com/sirupsen/logrus"
)

type couponTemplateService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    employeeInfo *model.Employee
    orm          *ent.CouponTemplateClient
}

func NewCouponTemplate() *couponTemplateService {
    return &couponTemplateService{
        ctx: context.Background(),
        orm: ent.Database.CouponTemplate,
    }
}

func NewCouponTemplateWithRider(r *ent.Rider) *couponTemplateService {
    s := NewCouponTemplate()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewCouponTemplateWithModifier(m *model.Modifier) *couponTemplateService {
    s := NewCouponTemplate()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewCouponTemplateWithEmployee(e *ent.Employee) *couponTemplateService {
    s := NewCouponTemplate()
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

// Create 创建优惠券模板
func (s *couponTemplateService) Create(req *model.CouponTemplateCreateReq) {
    meta := &model.CouponTemplateMeta{
        CouponTemplate: req.CouponTemplate,
    }

    // 查找城市
    if req.CityIDs != nil {
        cs, _ := ent.Database.City.QueryNotDeleted().Where(city.IDIn(*req.CityIDs...)).All(s.ctx)
        if len(cs) != len(*req.CityIDs) {
            snag.Panic("城市有错")
        }
        meta.Cities = make([]model.City, len(cs))
        for i, c := range cs {
            meta.Cities[i] = model.City{
                ID:   c.ID,
                Name: c.Name,
            }
        }
    }

    // 查找骑士卡
    if req.PlanIDs != nil {
        ps, _ := ent.Database.Plan.QueryNotDeleted().Where(plan.IDIn(*req.PlanIDs...)).All(s.ctx)
        if len(ps) != len(*req.PlanIDs) {
            snag.Panic("骑士卡有错")
        }
        meta.Plans = make([]model.Plan, len(ps))
        for i, p := range ps {
            meta.Plans[i] = model.Plan{
                ID:   p.ID,
                Name: p.Name,
                Days: p.Days,
            }
        }
    }

    // 保存数据
    _, err := s.orm.Create().SetName(req.Name).SetMeta(meta).Save(s.ctx)
    if err != nil {
        snag.Panic(err)
    }
}

func (s *couponTemplateService) List(req *model.CouponTemplateListReq) *model.PaginationRes {
    enable := true
    if req.Enable != nil {
        enable = *req.Enable
    }
    q := s.orm.Query().Where(coupontemplate.Enable(enable)).Order(ent.Desc(coupontemplate.FieldCreatedAt))
    log.Println(q)
    return nil
}
