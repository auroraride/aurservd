// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-18
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/workwx"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/coupontemplate"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/questioncategory"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
)

type selectionService struct {
	ctx context.Context
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
		WithBrand().
		WithAgreement()

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
			children := make([]model.SelectOption, 0)

			for _, arr := range p.Complexes {
				for _, cl := range *arr {
					children = append(children, model.SelectOption{
						Value: cl.ID,
						Label: fmt.Sprintf("%s - %d天", cl.Model, cl.Days),
					})
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

// CommissionPlan 返佣方案筛选骑士卡
func (s *selectionService) CommissionPlan(req *model.CommissionPlanSelectionReq) (items []model.CascaderOptionLevel3) {
	items = make([]model.CascaderOptionLevel3, 0)

	now := time.Now()
	q := ent.Database.Plan.QueryNotDeleted().Where(
		plan.StartLTE(now),
		plan.EndGTE(now),
		plan.Enable(true),
		plan.ParentIDIsNil(),
	).WithComplexes(func(pq *ent.PlanQuery) {
		pq.Where(plan.DeletedAtIsNil())
	})

	if req.Keyword != nil {
		q.Where(
			plan.NameContainsFold(*req.Keyword),
		)
	}

	all, _ := q.All(s.ctx)
	if all == nil {
		return items
	}

	cmap := make(map[uint64]model.CascaderOptionLevel3)
	for _, r := range all {
		// 根据骑士卡分组
		if _, ok := cmap[r.ID]; !ok {
			cmap[r.ID] = model.CascaderOptionLevel3{
				SelectOption: model.SelectOption{
					Value: r.ID,
					Label: r.Name,
				},
				Children: silk.Pointer(make([]model.CascaderOptionLevel2, 0)),
			}
		}

		l2c := cmap[r.ID].Children

		p := NewPlan().PlanWithComplexes(r)
		children := make([]model.SelectOption, 0)

		for _, arr := range p.Complexes {
			for _, cl := range *arr {
				children = append(children, model.SelectOption{
					Value: cl.ID,
					Label: fmt.Sprintf(" %d天 - %.2f元", cl.Days, cl.Price),
				})
			}
		}

		*l2c = append(*l2c, model.CascaderOptionLevel2{
			SelectOption: model.SelectOption{
				Value: p.ID,
				Label: p.Model,
			},
			Children: children,
		})

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
func (s *selectionService) Cabinet(req *model.CabinetSelectionReq) (items []*model.CascaderOptionLevel2) {
	q := ent.Database.Cabinet.QueryNotDeleted().WithCity()

	if req.EnterpriseID > 0 {
		q.Where(cabinet.EnterpriseID(req.EnterpriseID))
	}

	if req.StationID > 0 {
		q.Where(cabinet.StationID(req.StationID))
	}

	res, _ := q.All(s.ctx)

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
	p, _ := ent.Database.Plan.QueryNotDeleted().Where(plan.ID(req.PlanID)).First(s.ctx)
	return []string{p.Model}
}

func (s *selectionService) CabinetModel(req *model.SelectionCabinetModelByCabinetReq) (items []string) {
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

func (s *selectionService) BatterySerialSearch(req *model.BatterySearchReq) (res []*model.Battery) {
	res = make([]*model.Battery, 0)
	if len(req.Serial) < 4 {
		return
	}

	q := ent.Database.Asset.Query().Where(asset.SnHasSuffix(req.Serial)).WithModel()

	if req.StationID != nil && *req.StationID > 0 {
		q.Where(asset.LocationsType(model.AssetLocationsTypeStation.Value()), asset.LocationsID(*req.StationID))
	}
	if req.EnterpriseID != nil && *req.EnterpriseID > 0 {
		es, _ := ent.Database.EnterpriseStation.QueryNotDeleted().Where(enterprisestation.EnterpriseID(*req.EnterpriseID)).All(s.ctx)
		if len(es) > 0 {
			ids := make([]uint64, 0)
			for _, e := range es {
				ids = append(ids, e.ID)
			}
			q.Where(asset.LocationsType(model.AssetLocationsTypeStation.Value()), asset.LocationsIDIn(ids...))
		}
	}

	items, _ := q.All(s.ctx)
	for _, item := range items {
		var m string
		if item.Edges.Model != nil {
			m = item.Edges.Model.Model
		}
		res = append(res, &model.Battery{
			ID:    item.ID,
			SN:    item.Sn,
			Model: m,
		})
	}

	return
}

// QuestionCategory 筛选问题分类
func (s *selectionService) QuestionCategory() (items []model.SelectOption) {
	res, _ := ent.Database.QuestionCategory.QueryNotDeleted().Order(ent.Desc(questioncategory.FieldSort)).All(s.ctx)
	items = make([]model.SelectOption, len(res))
	for i, r := range res {
		items[i] = model.SelectOption{
			Value: r.ID,
			Label: r.Name,
		}
	}
	items = append(items, model.SelectOption{
		Value: 0,
		Label: "其他",
	})
	return
}

// ModelByCity 首页电池型号
func (s *selectionService) ModelByCity(req *model.SelectionCabinetModelByCityReq) (res []string) {
	res = make([]string, 0)
	q := ent.Database.Cabinet.QueryNotDeleted().
		Where(cabinet.Status(model.CabinetStatusNormal.Value())).
		WithModels()
	if req.CityID != nil {
		q.Where(cabinet.CityID(*req.CityID))
	}
	list, _ := q.All(s.ctx)
	if len(list) == 0 {
		return res
	}
	// 并集电池型号
	uniqueModels := make(map[string]struct{})
	for _, v := range list {
		for _, m := range v.Edges.Models {
			uniqueModels[m.Model] = struct{}{}
		}
	}

	res = make([]string, 0, len(uniqueModels))
	for m := range uniqueModels {
		res = append(res, m)
	}
	return res
}

// EbikeBrandByCity 通过城市筛选品牌
func (s *selectionService) EbikeBrandByCity(req *model.SelectionBrandByCityReq) (items []model.SelectOption) {
	brands := NewEbikeBrand().All()
	if req.CityID != nil {
		brands = NewEbikeBrand().ListByCityAndPlan(*req.CityID)
	}

	items = make([]model.SelectOption, len(brands))
	for i, b := range brands {
		items[i] = model.SelectOption{
			Value: b.ID,
			Label: b.Name,
		}
	}
	return
}

// StoreGroup 筛选门店集合
func (s *selectionService) StoreGroup() (items []model.SelectOption) {
	items = make([]model.SelectOption, 0)
	list, _ := ent.Database.StoreGroup.QueryNotDeleted().All(context.Background())
	for _, v := range list {
		items = append(items, model.SelectOption{
			Value: v.ID,
			Label: v.Name,
		})
	}
	return
}

// Goods 筛选商品
func (s *selectionService) Goods() (items []model.SelectOptionGoods) {
	items = make([]model.SelectOptionGoods, 0)
	goods, _ := ent.Database.Goods.Query().All(s.ctx)
	for _, g := range goods {
		prices := make([]model.SelectOption, 0)
		for k, p := range g.PaymentPlans {
			prices = append(prices, model.SelectOption{
				Value: uint64(k),
				Label: fmt.Sprintf("%v期", len(p)),
			})
		}
		items = append(items, model.SelectOptionGoods{
			Value:  g.ID,
			Label:  g.Name,
			Prices: prices,
		})
	}
	return
}

// GoodsStore 筛选购车门店
func (s *selectionService) GoodsStore() (items []*model.CascaderOptionLevel2) {
	res, _ := ent.Database.Store.QueryNotDeleted().WithCity().
		Where(
			store.EbikeSale(true),
		).
		All(s.ctx)

	return cascaderLevel2IDName(res, func(r *ent.Store) model.IDName {
		return s.nilableCity(r.Edges.City)
	}, "未选择网点", true)
}
