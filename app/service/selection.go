// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-18
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/workwx"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/coupontemplate"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "sort"
    "strconv"
    "time"
)

type selectionService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewSelection() *selectionService {
    return &selectionService{
        ctx: context.Background(),
    }
}

// Plan 筛选骑行卡
func (s *selectionService) Plan(req *model.PlanSelectionReq) (items []model.CascaderOptionLevel3) {
    q := ent.Database.Plan.QueryNotDeleted().
        Where(plan.ParentIDIsNil()).
        WithComplexes(func(pq *ent.PlanQuery) {
            pq.Where(plan.DeletedAtIsNil())
        }).
        WithCities().
        WithModels()

    if req.Effect != nil && *req.Effect != 0 {
        now := time.Now()
        if *req.Effect == 1 {
            q.Where(
                plan.StartLTE(now),
                plan.EndGTE(now),
            )
        } else {
            q.Where(
                plan.Or(
                    plan.StartGTE(now),
                    plan.EndLTE(now),
                ),
            )
        }
    }

    if req.Status != nil && *req.Status != 0 {
        enable := *req.Status == 1
        q.Where(plan.Enable(enable))
    }

    res, _ := q.All(s.ctx)

    cmap := make(map[uint64]model.CascaderOptionLevel3)

    for _, r := range res {
        cs := r.Edges.Cities
        for _, c := range cs {
            if _, ok := cmap[c.ID]; !ok {
                cmap[c.ID] = model.CascaderOptionLevel3{
                    SelectOption: model.SelectOption{
                        Value: c.ID,
                        Label: c.Name,
                    },
                    Children: silk.Pointer(make([]model.CascaderOptionLevel2, 0)),
                }
            }

            l2c := cmap[c.ID].Children

            p := NewPlan().PlanWithComplexes(r)
            children := make([]model.SelectOption, len(p.Complexes))

            for k, cl := range p.Complexes {
                children[k] = model.SelectOption{
                    Value: cl.ID,
                    Label: strconv.Itoa(int(cl.Days)),
                }
            }

            *l2c = append(*l2c, model.CascaderOptionLevel2{
                SelectOption: model.SelectOption{
                    Value: p.ID,
                    Label: p.Name,
                },
                Children: children,
            })
        }
    }

    items = make([]model.CascaderOptionLevel3, 0)
    for _, m := range cmap {
        items = append(items, m)
    }

    return
}

// Rider 筛选骑手
func (s *selectionService) Rider(req *model.RiderSelectionReq) (items []model.SelectOption) {
    q := ent.Database.Rider.QueryNotDeleted()
    if req.Keyword != nil {
        q.Where(
            rider.Or(
                rider.PhoneContainsFold(*req.Keyword),
                rider.NameContainsFold(*req.Keyword),
            ),
        )
    }
    res, _ := q.All(s.ctx)
    items = make([]model.SelectOption, len(res))

    for i, r := range res {
        name := "[未认证]"
        if r.Name != "" {
            name = r.Name
        }
        items[i] = model.SelectOption{
            Value: r.ID,
            Label: fmt.Sprintf("%s - %s", r.Phone, name),
        }
    }

    return
}

// Role 筛选角色
func (s *selectionService) Role() (items []model.SelectOption) {
    roles, _ := ent.Database.Role.Query().All(s.ctx)
    for _, role := range roles {
        items = append(items, model.SelectOption{
            Value: role.ID,
            Label: role.Name,
        })
    }
    return
}

// City 筛选城市
func (s *selectionService) City() (items []*model.CascaderOptionLevel2) {
    res, _ := ent.Database.City.QueryNotDeleted().WithChildren(func(cq *ent.CityQuery) {
        cq.Where(city.Open(true))
    }).Where(
        city.ParentIDIsNil(),
        city.HasChildrenWith(
            city.Open(true),
        ),
    ).All(s.ctx)

    items = make([]*model.CascaderOptionLevel2, len(res))

    for i, r := range res {
        items[i] = &model.CascaderOptionLevel2{
            SelectOption: model.SelectOption{
                Value: r.ID,
                Label: r.Name,
            },
            Children: make([]model.SelectOption, len(r.Edges.Children)),
        }

        for k, child := range r.Edges.Children {
            items[i].Children[k] = model.SelectOption{
                Value: child.ID,
                Label: child.Name,
            }
        }
    }

    return
}

type cascader[T any] func(data T) (parent model.SelectOption, item model.SelectOption)

func cascaderLevel2[T any](res []T, cb cascader[T], params ...any) (items []*model.CascaderOptionLevel2) {
    smap := make(map[uint64]*model.CascaderOptionLevel2)
    for _, r := range res {
        p, c := cb(r)

        ol, ok := smap[p.Value]
        if !ok {
            ol = &model.CascaderOptionLevel2{
                SelectOption: p,
                Children:     make([]model.SelectOption, 0),
            }
            smap[p.Value] = ol
        }

        ol.Children = append(ol.Children, c)
    }

    items = make([]*model.CascaderOptionLevel2, 0)
    for _, m := range smap {
        items = append(items, m)
    }

    if len(params) > 0 && params[0].(bool) {
        sort.Slice(items, func(i, j int) bool {
            return items[i].Value < items[j].Value
        })
    }

    return
}

type cascaderX[T any] func(data T) (parent *model.CascaderOption, item *model.CascaderOption)

func cascaderParser[T any](res []T, cb cascaderX[T]) (items []*model.CascaderOption) {
    smap := make(map[any]*model.CascaderOption)
    for _, r := range res {
        p, c := cb(r)

        ol, ok := smap[p.Value]
        if !ok {
            ol = &model.CascaderOption{
                Value:    p.Value,
                Label:    p.Label,
                Children: silk.Pointer(make([]*model.CascaderOption, 0)),
            }
            smap[p.Value] = ol
        }

        *ol.Children = append(*ol.Children, c)
    }

    items = make([]*model.CascaderOption, 0)
    for _, m := range smap {
        items = append(items, m)
    }

    return
}

func selectOptionIDName[T model.IDName, K model.IDName](r T, pb func(r T) K, message string) (p model.SelectOption, c model.SelectOption) {
    parent := pb(r)
    if parent.GetID() == 0 {
        p = model.SelectOption{
            Value: 0,
            Label: message,
        }
    } else {
        p = model.SelectOption{
            Value: parent.GetID(),
            Label: parent.GetName(),
        }
    }
    c = model.SelectOption{
        Value: r.GetID(),
        Label: r.GetName(),
    }
    return
}

func cascaderLevel2IDName[T model.IDName, K model.IDName](res []T, pb func(r T) K, message string, params ...any) (items []*model.CascaderOptionLevel2) {
    cb := func(r T) (model.SelectOption, model.SelectOption) {
        return selectOptionIDName(r, pb, message)
    }
    return cascaderLevel2(res, cb, params...)
}

func (s *selectionService) nilableCity(item *ent.City) model.IDName {
    if item == nil {
        return new(model.NilIDName)
    }
    return item
}

// Store 筛选门店
func (s *selectionService) Store() (items []*model.CascaderOptionLevel2) {
    res, _ := ent.Database.Store.QueryNotDeleted().WithCity().All(s.ctx)

    return cascaderLevel2IDName(res, func(r *ent.Store) model.IDName {
        return s.nilableCity(r.Edges.City)
    }, "未选择网点", true)
}

// Employee 筛选店员
func (s *selectionService) Employee() (items []*model.CascaderOptionLevel2) {
    res, _ := ent.Database.Employee.QueryNotDeleted().WithCity().All(s.ctx)

    return cascaderLevel2IDName(res, func(r *ent.Employee) model.IDName {
        return s.nilableCity(r.Edges.City)
    }, "未选择城市", true)
}

// Branch 筛选网点
func (s *selectionService) Branch() (items []*model.CascaderOptionLevel2) {
    res, _ := ent.Database.Branch.QueryNotDeleted().WithCity().All(s.ctx)

    return cascaderLevel2IDName(res, func(r *ent.Branch) model.IDName {
        return s.nilableCity(r.Edges.City)
    }, "未选择网点", true)
}

// Enterprise 筛选企业
func (s *selectionService) Enterprise() (items []*model.CascaderOptionLevel2) {
    res, _ := ent.Database.Enterprise.QueryNotDeleted().WithCity().All(s.ctx)

    return cascaderLevel2IDName(res, func(r *ent.Enterprise) model.IDName {
        return s.nilableCity(r.Edges.City)
    }, "未选择城市", true)
}

// Cabinet 筛选电柜
func (s *selectionService) Cabinet() (items []*model.CascaderOptionLevel2) {
    res, _ := ent.Database.Cabinet.QueryNotDeleted().WithCity().All(s.ctx)

    return cascaderLevel2IDName(res, func(r *ent.Cabinet) model.IDName {
        return s.nilableCity(r.Edges.City)
    }, "未选择网点", true)
}

// WorkwxEmployee 筛选企业微信成员
func (s *selectionService) WorkwxEmployee() (items []any) {
    wx := workwx.New()
    userlist, err := wx.UserSimpleList(1)

    if err != nil {
        snag.Panic(err)
    }

    items = make([]any, len(userlist.Userlist))

    for i, u := range userlist.Userlist {
        items[i] = ar.Map{
            "value": u.Userid,
            "label": u.Name,
        }
    }
    return
}

// PlanModel 筛选骑行卡电池型号
func (s *selectionService) PlanModel(req *model.SelectionPlanModelReq) []string {
    p, _ := ent.Database.Plan.QueryNotDeleted().Where(plan.ID(req.PlanID)).WithModels().First(s.ctx)
    items := make([]string, len(p.Edges.Models))
    for i, pm := range p.Edges.Models {
        items[i] = pm.Model
    }
    return items
}

func (s *selectionService) CabinetModel(req *model.SelectionCabinetModelReq) (items []string) {
    cab, _ := ent.Database.Cabinet.QueryNotDeleted().
        WithModels().
        Where(cabinet.ID(req.CabinetID)).
        First(s.ctx)
    items = make([]string, 0)
    if cab == nil {
        return
    }
    for _, bm := range cab.Edges.Models {
        items = append(items, bm.Model)
    }
    return
}

// CabinetModelX 筛选电柜型号
func (s *selectionService) CabinetModelX() (items []model.CascaderOption) {
    res, _ := ent.Database.Cabinet.QueryNotDeleted().
        WithCity().
        WithModels().
        All(s.ctx)

    cmap := make(map[uint64]model.CascaderOption)

    for _, r := range res {
        c := r.Edges.City
        if c == nil {
            continue
        }
        if _, ok := cmap[c.ID]; !ok {
            cmap[c.ID] = model.CascaderOption{
                Value:    c.ID,
                Label:    c.Name,
                Children: silk.Pointer(make([]*model.CascaderOption, 0)),
            }
        }

        l2c := cmap[c.ID].Children

        children := make([]*model.CascaderOption, len(r.Edges.Models))

        for k, b := range r.Edges.Models {
            children[k] = &model.CascaderOption{
                Value: b.Model,
                Label: b.Model,
            }
        }

        *l2c = append(*l2c, &model.CascaderOption{
            Value:    r.ID,
            Label:    r.Name,
            Children: &children,
        })
    }

    items = make([]model.CascaderOption, 0)
    for _, m := range cmap {
        items = append(items, m)
    }

    return
}

func (s *selectionService) Models() []string {
    return NewBatteryModel().Models()
}

// CouponTemplate 选择优惠券模板
func (s *selectionService) CouponTemplate() (items []model.SelectOptionGroup) {
    cts, _ := ent.Database.CouponTemplate.Query().Order(ent.Desc(coupontemplate.FieldUpdatedAt)).All(s.ctx)
    var enable, disable []model.SelectOption
    for _, ct := range cts {
        t := model.SelectOption{
            Value: ct.ID,
            Label: ct.Name,
            Desc:  ct.Remark,
        }
        if ct.Enable {
            enable = append(enable, t)
        } else {
            disable = append(disable, t)
        }
    }
    if len(enable) > 0 {
        items = append(items, model.SelectOptionGroup{
            Label:   "已启用",
            Options: enable,
        })
    }
    if len(disable) > 0 {
        items = append(items, model.SelectOptionGroup{
            Label:   "已禁用",
            Options: disable,
        })
    }
    return
}

func (s *selectionService) EbikeBrand() (items []model.SelectOption) {
    brands := NewEbikeBrand().All()
    items = make([]model.SelectOption, len(brands))
    for i, b := range brands {
        items[i] = model.SelectOption{
            Value: b.ID,
            Label: b.Name,
        }
    }
    return
}
