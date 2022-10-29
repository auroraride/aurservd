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
    "github.com/golang-module/carbon/v2"
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

func NewCouponTemplateWithModifier(m *model.Modifier) *couponTemplateService {
    s := NewCouponTemplate()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func (s *couponTemplateService) Query(id uint64, enable bool) (*ent.CouponTemplate, error) {
    return ent.Database.CouponTemplate.Query().Where(coupontemplate.Enable(enable), coupontemplate.ID(id)).First(s.ctx)
}

func (s *couponTemplateService) QueryEnableX(id uint64) *ent.CouponTemplate {
    ct, _ := s.Query(id, true)
    if ct == nil {
        snag.Panic("未找到启用的模板")
    }
    return ct
}

func (s *couponTemplateService) QueryDisableX(id uint64) *ent.CouponTemplate {
    ct, _ := s.Query(id, false)
    if ct == nil {
        snag.Panic("未找到禁用的模板")
    }
    return ct
}

// Create 创建优惠券模板
func (s *couponTemplateService) Create(req *model.CouponTemplateCreateReq) {
    meta := &model.CouponTemplateMeta{
        CouponTemplate: req.CouponTemplate,
    }

    s.CityAndPlan(meta, req.CityIDs, req.PlanIDs)

    // 保存数据
    _, err := s.orm.Create().SetName(req.Name).SetMeta(meta).SetRemark(req.Remark).Save(s.ctx)
    if err != nil {
        snag.Panic(err)
    }
}

func (s *couponTemplateService) Status(req *model.CouponTemplateStatusReq) {
    var updater *ent.CouponTemplateUpdateOne
    if req.Enable {
        updater = s.orm.UpdateOne(s.QueryDisableX(req.ID))
    } else {
        updater = s.orm.UpdateOne(s.QueryEnableX(req.ID))
    }
    updater.SetEnable(req.Enable).SaveX(s.ctx)
}

func (s *couponTemplateService) CityAndPlan(meta *model.CouponTemplateMeta, cities, plans *[]uint64) {
    // 查找城市
    if cities != nil {
        cs, _ := ent.Database.City.Query().Where(city.IDIn(*cities...)).All(s.ctx)
        if len(cs) != len(*cities) {
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
    if plans != nil {
        ps, _ := ent.Database.Plan.Query().Where(plan.IDIn(*plans...)).All(s.ctx)
        if len(ps) != len(*plans) {
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
}

func (s *couponTemplateService) List(req *model.CouponTemplateListReq) *model.PaginationRes {
    enable := true
    if req.Enable != nil {
        enable = *req.Enable
    }
    q := s.orm.Query().Where(coupontemplate.Enable(enable)).Order(ent.Desc(coupontemplate.FieldUpdatedAt))
    return model.ParsePaginationResponse(q.WithCoupons(), req.PaginationReq, func(item *ent.CouponTemplate) (res model.CouponTemplateListRes) {
        cps := item.Edges.Coupons
        res = model.CouponTemplateListRes{
            ID:     item.ID,
            Total:  len(cps),
            Name:   item.Name,
            Enable: item.Enable,
            Time:   item.UpdatedAt.Format(carbon.DateTimeLayout),
            Remark: item.Remark,
            CouponTemplateMeta: model.CouponTemplateMeta{
                CouponTemplate: model.CouponTemplate{
                    Rule:           item.Meta.Rule,
                    CouponDuration: item.Meta.CouponDuration,
                    Multiple:       item.Meta.Multiple,
                },
                Cities: item.Meta.Cities,
                Plans:  item.Meta.Plans,
            },
        }
        for _, c := range cps {
            switch NewCoupon().Status(c) {
            case model.CouponStatusExpired:
                res.Expired += 1
            case model.CouponStatusInStock:
                res.InStock += 1
            case model.CouponStatusUnused:
                res.Unused += 1
            case model.CouponStatusUsed:
                res.Used += 1
            }
        }
        return
    })
}
