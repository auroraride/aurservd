// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-18
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/plan"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/pkg/tools"
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

func (s *selectionService) Plan(req *model.PlanSelectionReq) (items []model.CascaderOptionLevel3) {
    q := ar.Ent.Plan.QueryNotDeleted().
        Where(plan.ParentIDIsNil()).
        WithComplexes(func(pq *ent.PlanQuery) {
            pq.Where(plan.DeletedAtIsNil())
        }).
        WithCities().
        WithPms()

    if req.Effect != nil {
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

    if req.Status != nil {
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
                    Children: tools.NewPointerInterface(make([]model.CascaderOptionLevel2, 0)),
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

func (s *selectionService) Rider(req *model.RiderSelectionReq) (items []model.SelectOption) {
    q := ar.Ent.Rider.QueryNotDeleted().WithPerson()
    if req.Keyword != nil {
        q.Where(
            rider.Or(
                rider.PhoneContainsFold(*req.Keyword),
                rider.HasPersonWith(person.NameContainsFold(*req.Keyword)),
            ),
        )
    }
    res, _ := q.All(s.ctx)
    items = make([]model.SelectOption, len(res))

    for i, r := range res {
        name := "[未认证]"
        if r.Edges.Person != nil {
            name = r.Edges.Person.Name
        }
        items[i] = model.SelectOption{
            Value: r.ID,
            Label: fmt.Sprintf("%s - %s", r.Phone, name),
        }
    }

    return
}

func (s *selectionService) Store() (items []*model.CascaderOptionLevel2) {
    res, _ := ar.Ent.Store.QueryNotDeleted().WithCity().All(s.ctx)

    smap := make(map[uint64]*model.CascaderOptionLevel2)

    for _, r := range res {
        c := r.Edges.City
        var cid uint64
        var cname string
        if c == nil {
            cid = 99999999
            cname = "未选择网点"
        } else {
            cid = c.ID
            cname = c.Name
        }

        ol, ok := smap[cid]
        if !ok {
            ol = &model.CascaderOptionLevel2{
                SelectOption: model.SelectOption{
                    Value: cid,
                    Label: cname,
                },
                Children: make([]model.SelectOption, 0),
            }
            smap[cid] = ol
        }

        ol.Children = append(ol.Children, model.SelectOption{
            Value: r.ID,
            Label: r.Name,
        })
    }

    items = make([]*model.CascaderOptionLevel2, 0)
    for _, m := range smap {
        items = append(items, m)
    }

    sort.Slice(items, func(i, j int) bool {
        return items[i].Value < items[j].Value
    })
    return
}

func (s *selectionService) Employee() (items []*model.CascaderOptionLevel2) {
    res, _ := ar.Ent.Employee.QueryNotDeleted().WithCity().All(s.ctx)

    smap := make(map[uint64]*model.CascaderOptionLevel2)

    for _, r := range res {
        c := r.Edges.City
        cid := c.ID
        cname := c.Name

        ol, ok := smap[cid]
        if !ok {
            ol = &model.CascaderOptionLevel2{
                SelectOption: model.SelectOption{
                    Value: cid,
                    Label: cname,
                },
                Children: make([]model.SelectOption, 0),
            }
            smap[cid] = ol
        }

        ol.Children = append(ol.Children, model.SelectOption{
            Value: r.ID,
            Label: fmt.Sprintf("%s - %s", r.Name, r.Phone),
        })
    }

    items = make([]*model.CascaderOptionLevel2, 0)
    for _, m := range smap {
        items = append(items, m)
    }

    sort.Slice(items, func(i, j int) bool {
        return items[i].Value < items[j].Value
    })

    return
}

func (s *selectionService) City() (items []*model.CascaderOptionLevel2) {
    res, _ := ar.Ent.City.QueryNotDeleted().WithChildren(func(cq *ent.CityQuery) {
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
