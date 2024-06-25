// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/exp/slices"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/model/promotion"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agreement"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/promotioncommission"
	"github.com/auroraride/aurservd/internal/ent/promotioncommissionplan"
	"github.com/auroraride/aurservd/internal/ent/setting"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
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
		orm: ent.Database.Plan,
	}
}

func NewPlanWithRider(rider *ent.Rider) *planService {
	s := NewPlan()
	s.ctx = context.WithValue(s.ctx, model.CtxRiderKey{}, rider)
	s.rider = rider
	return s
}

func NewPlanWithModifier(m *model.Modifier) *planService {
	s := NewPlan()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

// Query 查找骑士卡
func (s *planService) Query(id uint64) *ent.Plan {
	item, err := s.orm.QueryNotDeleted().Where(plan.ID(id)).First(s.ctx)
	if err != nil || item == nil {
		snag.Panic("未找到有效的骑士卡")
	}
	return item
}

// QueryEffectiveWithID 获取当前生效的骑行卡
func (s *planService) QueryEffectiveWithID(id uint64) *ent.Plan {
	today := carbon.Now().StartOfDay().StdTime()
	item, err := s.orm.QueryNotDeleted().
		Where(
			plan.Enable(true),
			plan.ID(id),
			plan.StartLTE(today),
			plan.EndGTE(today),
		).
		WithBrand().
		Only(s.ctx)
	if err != nil || item == nil {
		snag.Panic("未找到有效的骑士卡")
	}
	return item
}

// FilterEffectiveItems 按照指定条件过滤生效中的骑士卡
// 排序为按天数正序排列
func (s *planService) FilterEffectiveItems(ps ...predicate.Plan) []*ent.Plan {
	today := carbon.Now().StartOfDay().StdTime()
	plans, _ := s.orm.QueryNotDeleted().
		Where(ps...).
		Where(
			plan.Enable(true),
			plan.StartLTE(today),
			plan.EndGTE(today),
		).
		Order(ent.Asc(plan.FieldDays)).
		All(s.ctx)
	return plans
}

func (s *planService) QueryEffectiveList() ent.Plans {
	today := carbon.Now().StartOfDay().StdTime()
	items, _ := s.orm.QueryNotDeleted().
		Where(
			plan.Enable(true),
			plan.StartLTE(today),
			plan.EndGTE(today),
		).
		WithBrand().
		All(s.ctx)
	return items
}

// checkDuplicate 查询骑士卡是否冲突
func (s *planService) checkDuplicate(brandID uint64, cities []uint64, models []string, start, end time.Time, ty model.PlanType, params ...uint64) {
	q := s.orm.QueryNotDeleted().
		Where(
			plan.Enable(true),
			plan.HasCitiesWith(city.IDIn(cities...)),
			plan.StartLTE(end),
			plan.EndGTE(start),
			plan.ModelIn(models...),
		)
	if len(params) > 0 {
		parentID := params[0]
		q.Where(
			plan.IDNEQ(parentID),
			plan.ParentIDNEQ(parentID),
		)
	}
	if brandID > 0 {
		q.Where(plan.BrandID(brandID))
	} else {
		q.Where(plan.BrandIDIsNil())
	}

	if ty.Value() == model.PlanTypeEbikeRto.Value() {
		// 以租代购 类型不能重复
		q.Where(plan.Type(model.PlanTypeEbikeRto.Value()))
	} else {
		// 单电和车 不能重复
		q.Where(plan.TypeIn(model.PlanTypeBattery.Value(), model.PlanTypeEbikeWithBattery.Value()))
	}

	exists, _ := q.Exist(s.ctx)
	if exists {
		snag.Panic("骑士卡冲突")
	}
}

// Create 创建骑士卡
func (s *planService) Create(req *model.PlanCreateReq) model.PlanListRes {
	cities, _ := NewCity().QueryIDs(req.Cities)

	if len(cities) != len(req.Cities) {
		snag.Panic("城市选择错误")
	}

	// 验证车辆型号
	var brandID *uint64
	var brand *ent.EbikeBrand
	if req.BrandID > 0 {
		brand = NewEbikeBrand().QueryX(req.BrandID)
		brandID = silk.Pointer(brand.ID)
	}

	start := tools.NewTime().ParseDateStringX(req.Start)
	end := carbon.CreateFromStdTime(tools.NewTime().ParseDateStringX(req.End)).EndOfDay().StdTime()

	// 获取型号列表
	var bms []string
	mms := make(map[string]bool)
	for _, c := range req.Complexes {
		if req.Type == model.PlanTypeEbikeRto && c.RtoDays == nil {
			snag.Panic("以租代购最小使用天数必填")
		}
		if c.Model == "" {
			snag.Panic("电池型号必选")
		}
		if !mms[c.Model] {
			bms = append(bms, c.Model)
			mms[c.Model] = true
		}
	}

	NewBatteryModel().QueryModelsX(bms)

	// 查询是否重复
	s.checkDuplicate(req.BrandID, req.Cities, bms, start, end, req.Type)

	// 排序
	sort.Slice(req.Complexes, func(i, j int) bool {
		return req.Complexes[i].Days < req.Complexes[j].Days
	})

	if len(req.Notes) == 0 {
		req.Notes = make([]string, 0)
	}

	// 日租金价格
	ss := NewSetting()
	drs := ss.DailyRentItems()

	// 初始化未设定的日租金价格
	nsdrs := make(map[string]*model.PlanNotSettedDailyRent)

	// 开始创建
	var parent *ent.Plan
	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) error {
		creator := tx.Plan.Create().
			SetName(strings.TrimSpace(req.Name)).
			SetEnable(req.Enable).
			AddCityIDs(req.Cities...).
			SetStart(start).
			SetEnd(end).
			SetNotes(req.Notes).
			SetNillableBrandID(brandID).
			SetType(req.Type.Value()).
			SetIntelligent(*req.Intelligent).
			SetNillableDeposit(req.Deposit).
			SetNillableDepositAmount(req.DepositAmount).
			SetNillableDepositWechatPayscore(req.DepositWechatPayscore).
			SetNillableDepositAlipayAuthFreeze(req.DepositAlipayAuthFreeze).
			SetNillableDepositContract(req.DepositContract).
			SetNillableDepositPay(req.DepositPay).
			SetNillableAgreementID(req.AgreementID)

		for i, cl := range req.Complexes {
			c := creator.Clone().
				SetModel(strings.ToUpper(cl.Model)).
				SetPrice(cl.Price).
				SetOriginal(cl.Original).
				SetCommission(cl.Commission).
				SetDesc(cl.Desc).
				SetDays(cl.Days).
				SetDiscountNewly(cl.DiscountNewly).
				SetOverdueFee(cl.OverdueFee).
				SetDaily(cl.Daily)

			if req.Type == model.PlanTypeEbikeRto && cl.RtoDays != nil {
				c.SetRtoDays(*cl.RtoDays)
			}

			if i > 0 {
				c.SetParent(parent)
			}

			// 保存信息
			r, err := c.Save(s.ctx)
			snag.PanicIfError(err)

			if i == 0 {
				parent = r
				parent.Edges.Cities = cities
				parent.Edges.Complexes = make([]*ent.Plan, len(req.Complexes))
				parent.Edges.Complexes[i] = r
				parent.Edges.Brand = brand
			} else {
				parent.Edges.Complexes[i] = r
			}

			// 获取日租金详细
			for _, pc := range cities {
				// 若日租金是默认999 并且列表中不存在
				if key, dr := ss.DailyRent(drs, pc.ID, cl.Model, brandID); dr == model.DailyRentDefault && nsdrs[key] == nil {
					nsdrs[key] = &model.PlanNotSettedDailyRent{
						City: model.City{
							ID:   pc.ID,
							Name: pc.Name,
						},
						Model: cl.Model,
					}
					if brand != nil {
						nsdrs[key].EbikeBrand = &model.EbikeBrand{
							ID:   brand.ID,
							Name: brand.Name,
						}
					}
				}
			}
		}

		return nil
	})

	res := s.PlanWithComplexes(parent)
	for _, nsdr := range nsdrs {
		res.NotSettedDailyRent = append(res.NotSettedDailyRent, nsdr)
	}
	return res
}

// // Modify 修改骑士卡 TODO: 修改太麻烦了, 情况贼多, 暂时不做?
// func (s *planService) Modify(req *model.PlanModifyReq) model.PlanWithComplexes {
//     old, err := s.orm.QueryNotDeleted().Where(plan.ID(req.ID)).WithModels().WithCities().WithComplexes().First(s.ctx)
//     if err != nil {
//         snag.Panic("未找到骑士卡")
//     }
//     if old.ParentID != nil {
//         snag.Panic("骑士卡子项无法单独修改, 请携带原始骑士卡ID")
//     }
//
//     start := carbon.ParseByLayout(req.Start, carbon.DateLayout).StdTime()
//     end := carbon.ParseByLayout(req.End, carbon.DateLayout).StdTime()
//
//     // 查询是否重复
//     s.checkDuplicate(req.Cities, req.Models, start, end, req.ID)
//
//     // 排序
//     sort.Slice(req.Complexes, func(i, j int) bool {
//         return req.Complexes[i].Days < req.Complexes[j].Days
//     })
//
//     var parent *ent.Plan
//
//     // 判定父级骑士卡是否改变
//     first := req.Complexes[0]
//     if first.ID != old.ID {}
//
//     return model.PlanWithComplexes{}
// }

// UpdateEnable 修改骑士卡状态
func (s *planService) UpdateEnable(req *model.PlanEnableModifyReq) {
	if req.Enable == nil {
		snag.Panic("启用参数有误")
	}
	item := s.Query(req.ID)
	if item.ParentID != nil {
		snag.Panic("子项不能单独操作")
	}

	// 判断已上架的骑士卡中是否存在相同类型的套餐数据
	if *req.Enable {
		s.CheckUpdateEnable(item)
	}

	s.orm.Update().
		Where(plan.Or(
			plan.ID(req.ID),
			plan.ParentID(req.ID),
		)).
		SetEnable(*req.Enable).
		SaveX(s.ctx)
}

func (s *planService) CheckUpdateEnable(item *ent.Plan) {
	// 获取城市ids
	var cityIds []uint64
	planCities, _ := item.QueryCities().All(s.ctx)
	for _, ct := range planCities {
		cityIds = append(cityIds, ct.ID)
	}

	// 获取型号列表
	var bms []string
	mms := make(map[string]bool)
	bms = append(bms, item.Model)
	mms[item.Model] = true
	cps, _ := item.QueryComplexes().All(s.ctx)
	for _, c := range cps {
		if !mms[c.Model] {
			bms = append(bms, c.Model)
			mms[c.Model] = true
		}
	}
	NewBatteryModel().QueryModelsX(bms)

	var bId uint64
	if item.BrandID != nil {
		bId = *item.BrandID
	}

	s.checkDuplicate(bId, cityIds, bms, item.Start, item.End, model.PlanType(item.Type))
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
func (s *planService) PlanWithComplexes(item *ent.Plan) (res model.PlanListRes) {
	sort.Slice(item.Edges.Complexes, func(i, j int) bool {
		return item.Edges.Complexes[i].Days < item.Edges.Complexes[j].Days
	})

	res = model.PlanListRes{
		ID:                      item.ID,
		Type:                    model.PlanType(item.Type),
		Name:                    item.Name,
		Enable:                  item.Enable,
		Start:                   item.Start.Format(carbon.DateLayout),
		End:                     item.End.Format(carbon.DateLayout),
		Cities:                  make([]model.City, len(item.Edges.Cities)),
		Complexes:               make([]*model.PlanComplexes, 0),
		Notes:                   item.Notes,
		Intelligent:             item.Intelligent,
		Model:                   item.Model,
		DepositAlipayAuthFreeze: item.DepositAlipayAuthFreeze,
		DepositWechatPayscore:   item.DepositWechatPayscore,
		DepositPay:              item.DepositPay,
		DepositContract:         item.DepositContract,
		Deposit:                 item.Deposit,
		DepositAmount:           item.DepositAmount,
	}

	// 查询个签默认协议
	var defaultAgreement *ent.Agreement
	defaultAgreement, _ = ent.Database.Agreement.QueryNotDeleted().
		Where(
			agreement.UserType(model.AgreementUserTypePersonal.Value()),
			agreement.IsDefault(true),
		).First(s.ctx)

	if item.Edges.Agreement != nil {
		res.Agreement = &model.Agreement{
			ID:            item.Edges.Agreement.ID,
			Name:          item.Edges.Agreement.Name,
			URL:           item.Edges.Agreement.URL,
			Hash:          item.Edges.Agreement.Hash,
			ForceReadTime: item.Edges.Agreement.ForceReadTime,
		}
	} else if defaultAgreement != nil {
		// 如果没有设置协议, 则使用默认协议
		res.Agreement = &model.Agreement{
			ID:            defaultAgreement.ID,
			Name:          defaultAgreement.Name,
			URL:           defaultAgreement.URL,
			Hash:          defaultAgreement.Hash,
			ForceReadTime: defaultAgreement.ForceReadTime,
		}
	}

	// 电车型号
	eb := item.Edges.Brand
	if eb != nil {
		res.Brand = &model.EbikeBrand{
			ID:   eb.ID,
			Name: eb.Name,
		}
	}

	// 可用城市
	for i, c := range item.Edges.Cities {
		res.Cities[i] = model.City{
			ID:   c.ID,
			Name: c.Name,
		}
	}

	children := []*ent.Plan{
		item,
	}
	children = append(children, item.Edges.Complexes...)

	m := make(map[string]*model.PlanComplexes)
	for _, child := range children {
		r, ok := m[child.Model]
		if !ok {
			r = &model.PlanComplexes{}
			m[child.Model] = r
		}

		rtoDays := child.RtoDays

		pc := model.PlanComplex{
			ID:            child.ID,
			Price:         child.Price,
			Days:          child.Days,
			Original:      child.Original,
			Desc:          child.Desc,
			Commission:    child.Commission,
			Model:         child.Model,
			DiscountNewly: child.DiscountNewly,
			RtoDays:       rtoDays,
			OverdueFee:    child.OverdueFee,
			Daily:         child.Daily,
		}

		if len(child.Edges.Commissions) > 0 {
			for _, v := range child.Edges.Commissions {
				if v.Edges.PromotionCommission != nil {
					pc.CommissionName = v.Edges.PromotionCommission.Name
					pc.CommissionID = v.Edges.PromotionCommission.ID
					break
				}
			}

		}

		*r = append(*r, pc)

	}

	for _, v := range m {
		res.Complexes = append(res.Complexes, v)
	}

	return
}

// List 列举骑士卡
func (s *planService) List(req *model.PlanListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().
		Where(plan.ParentIDIsNil()).
		WithComplexes(func(pq *ent.PlanQuery) {
			pq.Where(plan.DeletedAtIsNil()).WithCommissions(func(q *ent.PromotionCommissionPlanQuery) {
				q.Where(promotioncommissionplan.DeletedAtIsNil()).WithPromotionCommission(func(q *ent.PromotionCommissionQuery) {
					q.Where(promotioncommission.TypeNEQ(promotion.CommissionCustom.Value()), promotioncommission.DeletedAtIsNil())
				})
			})
		}).
		WithCities().
		WithCommissions(func(q *ent.PromotionCommissionPlanQuery) {
			q.Where(promotioncommissionplan.DeletedAtIsNil()).WithPromotionCommission(func(q *ent.PromotionCommissionQuery) {
				q.Where(promotioncommission.TypeNEQ(promotion.CommissionCustom.Value()), promotioncommission.DeletedAtIsNil())
			})
		}).
		WithBrand().
		WithAgreement().
		Order(ent.Desc(plan.FieldStart), ent.Asc(plan.FieldEnd))

	if req.Intelligent != nil {
		q.Where(plan.Intelligent(*req.Intelligent == 1))
	}
	if req.CityID != nil {
		q.Where(plan.HasCitiesWith(city.ID(*req.CityID)))
	}
	if req.Name != nil {
		q.Where(plan.NameContainsFold(*req.Name))
	}
	if req.Enable != nil {
		q.Where(plan.Enable(*req.Enable))
	}
	if req.Model != nil {
		q.Where(plan.Model(*req.Model))
	}
	if req.Type != nil {
		q.Where(plan.Type(req.Type.Value()))
	}
	if req.BrandID != nil {
		q.Where(plan.BrandID(*req.BrandID))
	}

	if req.Deposit != nil {
		q.Where(plan.Deposit(*req.Deposit))
	}

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.Plan) model.PlanListRes {
			return s.PlanWithComplexes(item)
		},
	)
}

func (s *planService) renewalMapKey(m string, brandID *uint64) string {
	if brandID == nil {
		return m
	}
	return fmt.Sprintf("%s-%d", m, *brandID)
}

// Renewal 续签选项
func (s *planService) Renewal(req *model.PlanListRiderReq) map[string]*[]model.RiderPlanItem {
	rmap := make(map[string]*[]model.RiderPlanItem)
	today := carbon.Now().StartOfDay().StdTime()

	q := s.orm.QueryNotDeleted().
		Where(
			plan.Enable(true),
			plan.StartLTE(today),
			plan.EndGTE(today),
			plan.DaysGTE(req.Min),
			plan.HasCitiesWith(
				city.ID(req.CityID),
			),
			plan.Model(req.Model),
			plan.Intelligent(req.Intelligent),
		).
		Order(ent.Asc(plan.FieldDays))

	if req.PlanType != nil {
		q.Where(plan.Type(req.PlanType.Value()))
	}

	if req.EbikeBrandID != nil {
		q.Where(plan.BrandID(*req.EbikeBrandID))
	}

	items := q.AllX(s.ctx)

	for _, item := range items {
		key := s.renewalMapKey(item.Model, item.BrandID)
		list, ok := rmap[key]
		if !ok {
			list = new([]model.RiderPlanItem)
			rmap[key] = list
		}
		*list = append(*list, model.RiderPlanItem{
			ID:       item.ID,
			Name:     item.Name,
			Price:    item.Price,
			Days:     item.Days,
			Original: item.Original,
			Desc:     item.Desc,
		})
	}

	return rmap
}

func (s *planService) Key(model string, brandID *uint64) string {
	k := model
	if brandID != nil {
		k += fmt.Sprintf("-%d", *brandID)
	}
	return k
}

// RiderListNewly 获取新购骑士卡列表
func (s *planService) RiderListNewly(req *model.PlanListRiderReq) model.PlanNewlyRes {
	// 判断骑手是否个签
	if s.rider.EnterpriseID != nil {
		snag.Panic("仅个签骑手可购买")
	}

	// 判断骑手是否可以办理业务
	NewRider().CheckForBusiness(s.rider)

	// 判断是否有生效订阅
	_, sub := NewSubscribe().RecentDetail(s.rider.ID)
	if sub != nil && slices.Contains(model.SubscribeNotUnSubscribed(), sub.Status) {
		snag.Panic("骑手当前有其他订阅, 无法新购")
	}

	// 需缴纳押金金额
	deposit := NewRider().Deposit(s.rider.ID)

	today := carbon.Now().StartOfDay().StdTime()

	items := s.orm.QueryNotDeleted().
		Where(
			plan.Enable(true),
			plan.StartLTE(today),
			plan.EndGTE(today),
			plan.HasCitiesWith(
				city.ID(req.CityID),
			),
		).
		WithBrand().
		WithCities().
		WithAgreement().
		Order(ent.Asc(plan.FieldDays)).
		AllX(s.ctx)

	mmap := make(map[string]*model.PlanModelOption)

	bmap := make(map[uint64]*model.PlanEbikeBrandOption)

	serv := NewPlanIntroduce()
	intro := serv.QueryMap()

	// 查询个签默认协议
	var defaultAgreement *ent.Agreement
	defaultAgreement, _ = ent.Database.Agreement.QueryNotDeleted().
		Where(
			agreement.UserType(model.AgreementUserTypePersonal.Value()),
			agreement.IsDefault(true),
		).First(s.ctx)

	t, _ := NewOrder().PreconditionNewly(sub)

	for _, item := range items {
		key := s.Key(item.Model, item.BrandID)
		m, ok := mmap[key]
		if !ok {
			// 可用城市
			var cs []string
			for _, c := range item.Edges.Cities {
				cs = append(cs, c.Name)
			}
			// 封装电池型号
			m = &model.PlanModelOption{
				Children: new(model.PlanDaysPriceOptions),
				Model:    item.Model,
				Intro:    intro[serv.Key(item.Model, item.BrandID)],
				Notes:    append(item.Notes, fmt.Sprintf("仅限 %s 使用", strings.Join(cs, " / "))),
			}
			mmap[key] = m
		}

		var ramount float64
		if t == model.OrderTypeNewly && item.DiscountNewly > 0 {
			ramount = item.DiscountNewly
		}

		planDaysPriceOption := model.PlanDaysPriceOption{
			ID:                      item.ID,
			Name:                    item.Name,
			Price:                   item.Price,
			Days:                    item.Days,
			Original:                item.Original,
			DiscountNewly:           ramount,
			HasEbike:                item.BrandID != nil,
			Deposit:                 item.Deposit,
			DepositAmount:           item.DepositAmount,
			DepositWechatPayscore:   item.DepositWechatPayscore,
			DepositAlipayAuthFreeze: item.DepositAlipayAuthFreeze,
			DepositContract:         item.DepositContract,
			DepositPay:              item.DepositPay,
		}
		if item.Edges.Agreement != nil {
			planDaysPriceOption.Agreement = &model.Agreement{
				ID:            item.Edges.Agreement.ID,
				Name:          item.Edges.Agreement.Name,
				URL:           item.Edges.Agreement.URL,
				Hash:          item.Edges.Agreement.Hash,
				ForceReadTime: item.Edges.Agreement.ForceReadTime,
			}
		} else if defaultAgreement != nil {
			planDaysPriceOption.Agreement = &model.Agreement{
				ID:            defaultAgreement.ID,
				Name:          defaultAgreement.Name,
				URL:           defaultAgreement.URL,
				Hash:          defaultAgreement.Hash,
				ForceReadTime: defaultAgreement.ForceReadTime,
			}
		}

		*m.Children = append(*m.Children, planDaysPriceOption)

		if item.BrandID != nil {
			var b *model.PlanEbikeBrandOption
			bid := *item.BrandID
			b, ok = bmap[bid]
			if !ok {
				brand := item.Edges.Brand
				b = &model.PlanEbikeBrandOption{
					Children: new(model.PlanModelOptions),
					Name:     brand.Name,
					Cover:    brand.Cover,
				}
				bmap[bid] = b
			}

			var exists bool
			for _, c := range *b.Children {
				if c.Model == item.Model {
					exists = true
				}
			}
			if !exists {
				*b.Children = append(*b.Children, m)
			}
		}
	}

	res := model.PlanNewlyRes{
		Deposit:   deposit,
		Configure: NewPayment(s.rider).Configure(),
	}

	settings, _ := ent.Database.Setting.Query().Where(setting.KeyIn(model.SettingPlanBatteryDescriptionKey, model.SettingPlanEbikeDescriptionKey)).All(context.Background())
	for _, sm := range settings {
		var v model.SettingPlanDescription
		err := jsoniter.Unmarshal([]byte(sm.Content), &v)
		if err == nil {
			switch sm.Key {
			case model.SettingPlanBatteryDescriptionKey:
				res.BatteryDescription = v
			case model.SettingPlanEbikeDescriptionKey:
				res.EbikeDescription = v
			}
		}
	}

	for _, m := range mmap {
		he := false
		for _, c := range *m.Children {
			if c.HasEbike {
				he = true
			}
		}
		if !he {
			res.Models = append(res.Models, m)
		}
	}

	for _, b := range bmap {
		res.Brands = append(res.Brands, b)
	}

	return res
}

// RiderListRenewal 获取续费骑士卡列表
func (s *planService) RiderListRenewal() model.RiderPlanRenewalRes {
	sub, _ := NewSubscribe().QueryEffective(s.rider.ID)
	if sub == nil {
		snag.Panic("骑手无生效中的订阅, 无法续费")
		return model.RiderPlanRenewalRes{}
	}

	var fee float64
	var formula string
	var minDays uint

	if sub.Remaining < 0 {
		fee, formula, _ = NewSubscribe().CalculateOverdueFee(sub)
		minDays = uint(0 - sub.Remaining)
	}

	pl, _ := sub.QueryPlan().First(s.ctx)
	if pl == nil {
		snag.Panic("未找到骑士卡")
		return model.RiderPlanRenewalRes{}
	}

	planType := model.PlanType(pl.Type)
	rmap := s.Renewal(&model.PlanListRiderReq{
		CityID:       sub.CityID,
		Min:          minDays,
		Model:        sub.Model,
		EbikeBrandID: sub.BrandID,
		Intelligent:  sub.Intelligent,
		PlanType:     &planType,
	})

	items := make([]model.RiderPlanItem, 0)

	key := s.renewalMapKey(sub.Model, sub.BrandID)

	if list, ok := rmap[key]; ok {
		items = *list
	}

	return model.RiderPlanRenewalRes{
		Overdue:   sub.Remaining < 0,
		Fee:       fee,
		Formula:   formula,
		Days:      minDays,
		Items:     items,
		Configure: NewPayment(s.rider).Configure(),
	}
}

func (s *planService) NameFromID(id uint64) string {
	p, _ := ent.Database.Plan.QueryNotDeleted().Where(plan.ID(id)).First(s.ctx)
	if p == nil {
		return "-"
	}
	return p.Name
}

func (s *planService) ModifyTime(req *model.PlanModifyTimeReq) {
	s.Query(req.ID)
	s.orm.Update().
		Where(
			plan.Or(
				plan.ID(req.ID),
				plan.ParentID(req.ID),
			),
		).
		SetEnd(carbon.CreateFromStdTime(tools.NewTime().ParseDateStringX(req.End)).EndOfDay().StdTime()).
		SetStart(tools.NewTime().ParseDateStringX(req.Start)).
		ExecX(s.ctx)
}

// PlanCity 查询套餐城市
func (s *planService) PlanCity(id uint64) (res []uint64, err error) {
	pl, _ := s.orm.Query().Where(plan.ID(id)).WithCities().First(s.ctx)
	if pl == nil {
		return nil, fmt.Errorf("套餐不存在")
	}
	if pl.Edges.Cities == nil {
		return nil, fmt.Errorf("套餐城市不存在")
	}

	for _, v := range pl.Edges.Cities {
		res = append(res, v.ID)
	}
	return res, nil
}
