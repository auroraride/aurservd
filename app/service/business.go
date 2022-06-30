// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/business"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
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
    r, err := NewRider().QueryForBusinessID(id)
    if err != nil {
        snag.Panic(err)
    }
    // 获取最近的订阅
    subd, _ := NewSubscribe().RecentDetail(r.ID)

    if subd == nil {
        snag.Panic("未找到有效订阅")
    }

    s.CheckCity(subd.City.ID)

    ic := r.Edges.Person.IDCardNumber
    res = model.SubscribeBusiness{
        ID:           r.ID,
        Status:       subd.Status,
        Name:         r.Edges.Person.Name,
        Phone:        r.Phone,
        IDCardNumber: ic[len(ic)-4:],
        Model:        subd.Model,
        SubscribeID:  subd.ID,
    }

    if subd.Enterprise != nil {
        res.EnterpriseName = subd.Enterprise.Name
    }

    if subd.Plan != nil {
        res.PlanName = subd.Plan.Name
    }

    res.Business = subd.Business
    return
}

// Plans 获取更换电池型号允许的套餐列表
// TODO 更换电池型号
func (s *businessService) Plans(subscribeID uint64) {
    sub := NewSubscribe().QueryEdgesX(subscribeID)

    s.CheckCity(sub.CityID)
    NewRider().CheckForBusiness(sub.Edges.Rider)

    if sub.Status != model.SubscribeStatusUsing {
        snag.Panic("当前为非使用中")
    }

    // 获取全部的电压列表
}

// listBasicQuery 列表基础查询语句
func (s *businessService) listBasicQuery(req *model.BusinessListReq) *ent.BusinessQuery {
    tt := tools.NewTime()

    q := ent.Database.Business.
        QueryNotDeleted().
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).
        WithEnterprise().
        WithPlan()

    if req.Type != nil {
        q.Where(business.TypeEQ(business.Type(*req.Type)))
    }

    if req.Start != nil {
        q.Where(business.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
    }

    if req.End != nil {
        q.Where(business.CreatedAtLTE(tt.ParseDateStringX(*req.End)))
    }

    switch req.Aimed {
    case model.BusinessAimedPersonal:
        q.Where(business.EnterpriseIDIsNil())
        break
    case model.BusinessAimedEnterprise:
        q.Where(business.EnterpriseIDNotNil())
        break
    }

    return q
}

// ListEmployee 业务列表 - 门店
func (s *businessService) ListEmployee(req *model.BusinessListReq) *model.PaginationRes {
    q := s.listBasicQuery(req).Where(business.EmployeeID(s.employee.ID))

    return model.ParsePaginationResponse(
        q,
        req.PaginationReq,
        func(item *ent.Business) (res model.BusinessEmployeeListRes) {
            res = model.BusinessEmployeeListRes{
                ID:    item.ID,
                Name:  item.Edges.Rider.Edges.Person.Name,
                Phone: item.Edges.Rider.Phone,
                Type:  item.Type.String(),
                Time:  item.CreatedAt.Format(carbon.DateTimeLayout),
            }
            p := item.Edges.Plan
            if p != nil {
                res.Plan = &model.Plan{
                    ID:   p.ID,
                    Name: p.Name,
                    Days: p.Days,
                }
            }

            e := item.Edges.Enterprise
            if e != nil {
                res.Enterprise = &model.EnterpriseBasic{
                    ID:   e.ID,
                    Name: e.Name,
                }
            }

            return
        },
    )
}
