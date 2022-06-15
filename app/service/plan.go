// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/batterymodel"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "math"
    "sort"
    "time"
)

type planService struct {
    ctx      context.Context
    orm      *ent.PlanClient
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewPlan() *planService {
    return &planService{
        ctx: context.Background(),
        orm: ar.Ent.Plan,
    }
}

func NewPlanWithRider(rider *ent.Rider) *planService {
    s := NewPlan()
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

func NewPlanWithModifier(m *model.Modifier) *planService {
    s := NewPlan()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// Query 查找骑士卡
func (s *planService) Query(id uint64) *ent.Plan {
    item, err := s.orm.QueryNotDeleted().Where(plan.ID(id)).Only(s.ctx)
    if err != nil || item == nil {
        snag.Panic("未找到有效的骑士卡")
    }
    return item
}

// QueryEffectiveWithID 获取当前生效的骑行卡
func (s *planService) QueryEffectiveWithID(id uint64) *ent.Plan {
    now := time.Now()
    item, err := s.orm.QueryNotDeleted().
        Where(
            plan.Enable(true),
            plan.ID(id),
            plan.StartLTE(now),
            plan.EndGTE(now),
        ).
        WithPms().
        Only(s.ctx)
    if err != nil || item == nil {
        snag.Panic("未找到有效的骑士卡")
    }
    return item
}

// checkDuplicate 查询骑士卡是否冲突
func (s *planService) checkDuplicate(cities, models []uint64, start, end time.Time, parentID ...uint64) {
    for _, cityID := range cities {
        for _, modelID := range models {
            q := s.orm.QueryNotDeleted().
                Where(
                    plan.Enable(true),
                    plan.HasCitiesWith(city.ID(cityID)),
                    plan.HasPmsWith(batterymodel.ID(modelID)),
                    plan.StartLTE(end),
                    plan.EndGTE(start),
                )
            if len(parentID) > 0 {
                id := parentID[0]
                q.Where(
                    plan.IDNEQ(id),
                    plan.ParentIDNEQ(id),
                )
            }
            if exists, _ := q.Exist(s.ctx); exists {
                snag.Panic("骑士卡冲突")
            }
        }
    }
}

func (s *planService) cloneCreator(creator *ent.PlanCreate) *ent.PlanCreate {
    c := *creator
    return &c
}

func (s *planService) getCitiesAndModels(reqCities, reqModels []uint64) (cities ent.Cities, pms ent.BatteryModels) {
    var err error
    cities, err = ar.Ent.City.QueryNotDeleted().Where(city.IDIn(reqCities...)).All(s.ctx)
    if err != nil {
        snag.Panic("城市参数错误")
    }
    pms, err = ar.Ent.BatteryModel.QueryNotDeleted().Where(batterymodel.IDIn(reqModels...)).All(s.ctx)
    if err != nil {
        snag.Panic("电池型号错误")
    }
    return
}

// Create 创建骑士卡
func (s *planService) Create(req *model.PlanCreateReq) model.PlanWithComplexes {
    cities, pms := s.getCitiesAndModels(req.Cities, req.Models)

    start := carbon.ParseByLayout(req.Start, carbon.DateLayout).Carbon2Time()
    end := carbon.ParseByLayout(req.End, carbon.DateLayout).Carbon2Time()

    // 查询是否重复
    s.checkDuplicate(req.Cities, req.Models, start, end)

    // 排序
    sort.Slice(req.Complexes, func(i, j int) bool {
        return req.Complexes[i].Days < req.Complexes[j].Days
    })

    // 开始创建
    tx, _ := ar.Ent.Tx(s.ctx)
    creator := tx.Plan.Create().
        SetName(req.Name).
        SetEnable(req.Enable).
        AddCityIDs(req.Cities...).
        AddPmIDs(req.Models...).
        SetStart(start).
        SetEnd(end)

    var parent *ent.Plan
    for i, cl := range req.Complexes {
        c := s.cloneCreator(creator).
            SetPrice(cl.Price).
            SetOriginal(cl.Original).
            SetCommission(cl.Commission).
            SetDesc(cl.Desc).
            SetDays(cl.Days)
        if i > 0 {
            c.SetParent(parent)
        }
        r, err := c.Save(s.ctx)
        snag.PanicIfErrorX(err, tx.Rollback)
        if i == 0 {
            parent = r
            parent.Edges.Cities = cities
            parent.Edges.Pms = pms
            parent.Edges.Complexes = make([]*ent.Plan, len(req.Complexes))
            parent.Edges.Complexes[i] = r
        } else {
            parent.Edges.Complexes[i] = r
        }
    }

    _ = tx.Commit()

    return s.PlanWithComplexes(parent)
}

// // Modify 修改骑士卡 TODO: 修改太麻烦了, 情况贼多, 暂时不做?
// func (s *planService) Modify(req *model.PlanModifyReq) model.PlanWithComplexes {
//     old, err := s.orm.QueryNotDeleted().Where(plan.ID(req.ID)).WithPms().WithCities().WithComplexes().First(s.ctx)
//     if err != nil {
//         snag.Panic("未找到骑士卡")
//     }
//     if old.ParentID != nil {
//         snag.Panic("骑士卡子项无法单独修改, 请携带原始骑士卡ID")
//     }
//
//     start := carbon.ParseByLayout(req.Start, carbon.DateLayout).Carbon2Time()
//     end := carbon.ParseByLayout(req.End, carbon.DateLayout).Carbon2Time()
//
//     // 查询是否重复
//     s.checkDuplicate(req.Cities, req.Models, start, end, req.ID)
//
//     // 排序
//     sort.Slice(req.Complexes, func(i, j int) bool {
//         return req.Complexes[i].Days < req.Complexes[j].Days
//     })
//
//     tx, _ := ar.Ent.Tx(s.ctx)
//
//     var parent *ent.Plan
//
//     // 判定父级骑士卡是否改变
//     first := req.Complexes[0]
//     if first.ID != old.ID {}
//
//     _ = tx.Commit()
//     return model.PlanWithComplexes{}
// }

// UpdateEnable 修改骑士卡状态
func (s *planService) UpdateEnable(req *model.PlanEnableModifyReq) {
    item := s.Query(req.ID)
    if item.ParentID != nil {
        snag.Panic("子项不能单独操作")
    }
    s.orm.Update().
        Where(plan.Or(
            plan.ID(req.ID),
            plan.ParentID(req.ID),
        )).
        SetEnable(*req.Enable).
        SaveX(s.ctx)
}

// Delete 软删除骑士卡
func (s *planService) Delete(req *model.IDParamReq) {
    item := s.Query(req.ID)
    if item.ParentID != nil {
        snag.Panic("子项不能单独操作")
    }
    s.orm.SoftDelete().Where(plan.Or(
        plan.ID(req.ID),
        plan.ParentID(req.ID),
    )).SaveX(s.ctx)
}

// PlanWithComplexes 骑士卡详情
func (s *planService) PlanWithComplexes(item *ent.Plan) (res model.PlanWithComplexes) {
    sort.Slice(item.Edges.Complexes, func(i, j int) bool {
        return item.Edges.Complexes[i].Days < item.Edges.Complexes[j].Days
    })

    res = model.PlanWithComplexes{
        ID:        item.ID,
        Name:      item.Name,
        Enable:    item.Enable,
        Start:     item.Start.Format(carbon.DateLayout),
        End:       item.End.Format(carbon.DateLayout),
        Cities:    make([]model.City, len(item.Edges.Cities)),
        Models:    make([]model.BatteryModel, len(item.Edges.Pms)),
        Complexes: make([]model.PlanComplex, len(item.Edges.Complexes)+1),
    }

    for i, c := range item.Edges.Cities {
        res.Cities[i] = model.City{
            ID:   c.ID,
            Name: c.Name,
        }
    }

    for i, pm := range item.Edges.Pms {
        res.Models[i] = model.BatteryModel{
            ID:       pm.ID,
            Voltage:  pm.Voltage,
            Capacity: pm.Capacity,
        }
    }

    res.Complexes[0] = model.PlanComplex{
        ID:         item.ID,
        Price:      item.Price,
        Days:       item.Days,
        Original:   item.Original,
        Desc:       item.Desc,
        Commission: item.Commission,
    }

    for i, child := range item.Edges.Complexes {
        res.Complexes[i+1] = model.PlanComplex{
            ID:         child.ID,
            Price:      child.Price,
            Days:       child.Days,
            Original:   child.Original,
            Desc:       child.Desc,
            Commission: child.Commission,
        }
    }
    return
}

// List 列举骑士卡
func (s *planService) List(req *model.PlanListReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().
        Where(plan.ParentIDIsNil()).
        WithComplexes(func(pq *ent.PlanQuery) {
            pq.Where(plan.DeletedAtIsNil())
        }).
        WithCities().
        WithPms().
        Order(ent.Desc(plan.FieldStart), ent.Asc(plan.FieldEnd))

    if req.CityID != nil {
        q.Where(plan.HasCitiesWith(city.ID(*req.CityID)))
    }
    if req.Name != nil {
        q.Where(plan.NameContainsFold(*req.Name))
    }
    if req.Enable != nil {
        q.Where(plan.Enable(*req.Enable))
    }

    return model.ParsePaginationResponse(
        q,
        req.PaginationReq,
        func(item *ent.Plan) model.PlanWithComplexes {
            return s.PlanWithComplexes(item)
        },
    )
}

// RiderList 获取骑士卡列表
func (s *planService) RiderList(req *model.PlanListRiderReq) (res []model.RiderPlanItem) {
    now := time.Now()
    items := s.orm.QueryNotDeleted().
        Where(
            plan.Enable(true),
            plan.StartLTE(now),
            plan.EndGTE(now),
            plan.HasPmsWith(batterymodel.Voltage(req.Voltage)),
        ).
        Order(ent.Asc(plan.FieldDays)).
        AllX(s.ctx)
    res = make([]model.RiderPlanItem, len(items))
    for i, item := range items {
        res[i] = model.RiderPlanItem{
            ID:       item.ID,
            Name:     item.Name,
            Price:    item.Price,
            Days:     item.Days,
            Original: item.Original,
            Desc:     item.Desc,
        }
    }
    return
}

// RiderListNewly 获取新购骑士卡列表
func (s *planService) RiderListNewly(req *model.PlanListRiderReq) []model.RiderPlanItem {
    if sub, _ := NewSubscribe().QueryEffective(s.rider.ID); sub != nil {
        snag.Panic("骑手有生效中的订阅, 无法新购")
    }

    return s.RiderList(req)
}

// RiderListRenewal 获取续费骑士卡列表
func (s *planService) RiderListRenewal() model.RiderPlanRenewalRes {
    sub, _ := NewSubscribe().QueryEffective(s.rider.ID)
    if sub == nil {
        snag.Panic("骑手无生效中的订阅, 无法续费")
    }

    var fee float64
    var formula string
    if sub.Remaining < 0 {
        fee, formula = NewSubscribe().OverdueFee(s.rider.ID, math.Abs(float64(sub.Remaining)))
    }

    return model.RiderPlanRenewalRes{
        Overdue: sub.Remaining < 0,
        Fee:     fee,
        Formula: formula,
        Items: s.RiderList(&model.PlanListRiderReq{
            CityID:  sub.CityID,
            Voltage: sub.Voltage,
        }),
    }
}
