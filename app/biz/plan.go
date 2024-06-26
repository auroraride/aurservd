package biz

import (
	"context"
	"fmt"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agreement"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/ebikebrand"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/setting"
	"github.com/auroraride/aurservd/internal/ent/stock"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/pkg/utils"
)

type planBiz struct {
	orm *ent.PlanClient
	ctx context.Context
}

func NewPlanBiz() *planBiz {
	return &planBiz{
		orm: ent.Database.Plan,
		ctx: context.Background(),
	}
}

// RiderListNewly 套餐列表
func (s *planBiz) RiderListNewly(r *ent.Rider, req *model.PlanListRiderReq) *definition.PlanNewlyRes {
	var state uint

	today := carbon.Now().StartOfDay().StdTime()

	q := s.orm.QueryNotDeleted().
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
		Order(ent.Asc(plan.FieldDays))

	if req.StoreId != nil {
		// 查询门店库存电车所属brandId
		var brandIds []uint64
		storeItem, _ := ent.Database.Store.QueryNotDeleted().
			WithStocks().
			Where(store.ID(*req.StoreId)).
			First(s.ctx)
		if storeItem.Edges.Stocks != nil {
			for _, st := range storeItem.Edges.Stocks {
				if st.BrandID != nil {
					brandIds = append(brandIds, *st.BrandID)
				}
			}
		}
		if len(brandIds) > 0 {
			q.Where(
				plan.HasBrandWith(ebikebrand.IDIn(brandIds...)),
			)
		}
	}

	items := q.AllX(s.ctx)
	mmap := make(map[string]*model.PlanModelOption)

	bmap := make(map[uint64]*model.PlanEbikeBrandOption)

	rbmap := make(map[uint64]*model.PlanEbikeBrandOption)

	serv := service.NewPlanIntroduce()
	intro := serv.QueryMap()

	// 查询个签默认协议
	var defaultAgreement *ent.Agreement
	defaultAgreement, _ = ent.Database.Agreement.QueryNotDeleted().
		Where(
			agreement.UserType(model.AgreementUserTypePersonal.Value()),
			agreement.IsDefault(true),
		).First(s.ctx)

	for _, item := range items {
		key := serv.Key(item.Model, item.BrandID)
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
		if r != nil {
			// 判断是否有生效订阅
			_, sub := service.NewSubscribe().RecentDetail(r.ID)
			if sub != nil && slices.Contains(model.SubscribeNotUnSubscribed(), sub.Status) {
				ramount = 0
			} else {
				state, _ = service.NewOrder().PreconditionNewly(sub)
				if state == model.OrderTypeNewly && item.DiscountNewly > 0 {
					ramount = item.DiscountNewly
				}
			}
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
			RtoDays:                 item.RtoDays,
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

		SortIDOptions(*m.Children)

		if item.BrandID != nil {
			switch item.Type {
			case model.PlanTypeEbikeWithBattery.Value():
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
			case model.PlanTypeEbikeRto.Value():
				var b *model.PlanEbikeBrandOption
				bid := *item.BrandID
				b, ok = rbmap[bid]
				if !ok {
					brand := item.Edges.Brand
					b = &model.PlanEbikeBrandOption{
						Children: new(model.PlanModelOptions),
						Name:     brand.Name,
						Cover:    brand.Cover,
					}
					rbmap[bid] = b
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
			default:
			}

		}
	}

	res := &definition.PlanNewlyRes{}

	if r != nil {
		res.Configure = service.NewPayment(r).Configure()
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
		SortPlanEbikeModelByName(*b.Children)
	}

	for _, rb := range rbmap {
		res.RtoBrands = append(res.RtoBrands, rb)
		SortPlanEbikeModelByName(*rb.Children)
	}

	SortPlanEbikeBrandByName(res.Brands)
	SortPlanModelByName(res.Models)
	SortPlanEbikeBrandByName(res.RtoBrands)

	return res
}

func SortPlanEbikeBrandByName(options []*model.PlanEbikeBrandOption) {
	sort.Slice(options, func(i, j int) bool {
		return options[i].Name < options[j].Name
	})
}

func SortPlanModelByName(options []*model.PlanModelOption) {
	sort.Slice(options, func(i, j int) bool {
		return options[i].Model < options[j].Model
	})
}

func SortPlanEbikeModelByName(options model.PlanModelOptions) {
	sort.Slice(options, func(i, j int) bool {
		return options[i].Model < options[j].Model
	})
}

func SortIDOptions(options model.PlanDaysPriceOptions) {
	sort.Slice(options, func(i, j int) bool {
		numI, _ := strconv.Atoi(strconv.FormatUint(options[i].ID, 10))
		numJ, _ := strconv.Atoi(strconv.FormatUint(options[j].ID, 10))
		return numI < numJ
	})
}

// EbikeList 车电套餐列表
func (s *planBiz) EbikeList(brandIds []uint64) (res definition.PlanNewlyRes) {

	today := carbon.Now().StartOfDay().StdTime()

	items := s.orm.QueryNotDeleted().
		Where(
			plan.Enable(true),
			plan.StartLTE(today),
			plan.EndGTE(today),
			plan.BrandIDIn(brandIds...),
		).
		WithBrand().
		WithCities().
		WithAgreement().
		Order(ent.Asc(plan.FieldDays)).
		AllX(s.ctx)

	bmap := make(map[uint64]*model.PlanEbikeBrandOption)

	rbmap := make(map[uint64]*model.PlanEbikeBrandOption)

	serv := service.NewPlanIntroduce()
	intro := serv.QueryMap()

	for _, item := range items {
		// 可用城市
		var cs []string
		for _, c := range item.Edges.Cities {
			cs = append(cs, c.Name)
		}
		// 封装电池型号
		m := &model.PlanModelOption{
			Children: new(model.PlanDaysPriceOptions),
			Model:    item.Model,
			Intro:    intro[serv.Key(item.Model, item.BrandID)],
			Notes:    append(item.Notes, fmt.Sprintf("仅限 %s 使用", strings.Join(cs, " / "))),
		}
		switch item.Type {
		case model.PlanTypeEbikeWithBattery.Value():
			var b *model.PlanEbikeBrandOption
			bid := *item.BrandID
			b, ok := bmap[bid]
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
		case model.PlanTypeEbikeRto.Value():
			var b *model.PlanEbikeBrandOption
			bid := *item.BrandID
			b, ok := rbmap[bid]
			if !ok {
				brand := item.Edges.Brand
				b = &model.PlanEbikeBrandOption{
					Children: new(model.PlanModelOptions),
					Name:     brand.Name,
					Cover:    brand.Cover,
				}
				rbmap[bid] = b
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
		default:
		}
	}

	settings, _ := ent.Database.Setting.Query().Where(setting.KeyIn(model.SettingPlanEbikeDescriptionKey)).All(context.Background())
	for _, sm := range settings {
		var v model.SettingPlanDescription
		err := jsoniter.Unmarshal([]byte(sm.Content), &v)
		if err == nil {
			switch sm.Key {
			case model.SettingPlanEbikeDescriptionKey:
				res.EbikeDescription = v
			}
		}
	}

	for _, b := range bmap {
		res.Brands = append(res.Brands, b)
		SortPlanEbikeModelByName(*b.Children)
	}

	for _, rb := range rbmap {
		res.RtoBrands = append(res.RtoBrands, rb)
		SortPlanEbikeModelByName(*rb.Children)
	}

	SortPlanEbikeBrandByName(res.Brands)
	SortPlanEbikeBrandByName(res.RtoBrands)

	return
}

// Detail 套餐详情
func (s *planBiz) Detail(req *definition.PlanDetailReq) (*definition.PlanDetailRes, error) {
	d, _ := s.orm.QueryNotDeleted().
		Where(plan.ID(req.ID)).
		WithBrand().
		First(s.ctx)
	if d == nil {
		return nil, nil
	}
	res := &definition.PlanDetailRes{
		Plan: model.Plan{
			ID:          d.ID,
			Name:        d.Name,
			Price:       d.Price,
			Days:        d.Days,
			Intelligent: d.Intelligent,
			Type:        model.PlanType(d.Type),
		},
		Notes: d.Notes,
	}

	return res, nil
}

// ListByStore 基于门店的套餐列表
func (s *planBiz) ListByStore(req *definition.StorePlanReq) []*definition.StoreEbikePlan {
	// 查询附近门店数据
	maxDistance := definition.DefaultMaxDistance
	if req.Distance != nil && *req.Distance < maxDistance {
		maxDistance = *req.Distance
	}
	stq := ent.Database.Store.QueryNotDeleted().
		Where(
			store.CityID(req.CityId),
			store.HasStocksWith(stock.BrandIDNotNil()),
		).
		WithCity().
		WithStocks().
		Modify(func(sel *sql.Selector) {
			sel.
				AppendSelectExprAs(sql.Raw(fmt.Sprintf(`ST_Distance(ST_GeographyFromText('SRID=4326;POINT(' || "store"."lng" || ' ' || "store"."lat" || ')'),ST_GeographyFromText('SRID=4326;POINT(%f  %f)'))`, req.Lng, req.Lat)), "distance").
				OrderBy(sql.Asc("distance")).
				Where(sql.P(func(b *sql.Builder) {
					b.WriteString(fmt.Sprintf(`ST_DWithin(ST_GeographyFromText('SRID=4326;POINT(' || "store"."lng" || ' ' || "store"."lat" || ')'),ST_GeographyFromText('SRID=4326;POINT(%f  %f)'), %f)`, req.Lng, req.Lat, maxDistance))
				}))
		})

	storelist, _ := stq.All(s.ctx)
	if len(storelist) == 0 {
		return []*definition.StoreEbikePlan{}
	}

	storeMap := make(map[uint64]*ent.Store)
	for _, v := range storelist {
		storeMap[v.ID] = v
	}

	// 门店电车库存品牌ID
	var brandIds []uint64
	// 门店电车品牌集合
	cityStoreBrandExist := make(map[string]bool)
	cityBrand2StoresMap := make(map[string][]uint64)
	for _, st := range storelist {
		// 计算该门店电车品牌库存
		brandStockMap := s.CalStoreEbikeStock(st.Edges.Stocks)
		for _, stc := range st.Edges.Stocks {
			if stc.BrandID != nil {
				existKey := fmt.Sprintf("%d-%d-%d", req.CityId, st.ID, *stc.BrandID)
				if !cityStoreBrandExist[existKey] && brandStockMap[*stc.BrandID] > 0 {
					cityStoreBrandExist[existKey] = true
					brandIds = append(brandIds, *stc.BrandID)
					// 保存存在库存的电车门店对应门店ID
					cbKey := fmt.Sprintf("%d-%d", req.CityId, *stc.BrandID)
					cityBrand2StoresMap[cbKey] = append(cityBrand2StoresMap[cbKey], st.ID)
				}
			}
		}
	}

	// 查询车电套餐
	today := carbon.Now().StartOfDay().StdTime()
	q := s.orm.QueryNotDeleted().
		Where(
			plan.TypeIn(model.PlanTypeEbikeWithBattery.Value(), model.PlanTypeEbikeRto.Value()),
			plan.Enable(true),
			plan.StartLTE(today),
			plan.EndGTE(today),
			plan.HasCitiesWith(city.ID(req.CityId)),
		).
		WithBrand().
		WithCities().
		WithAgreement().
		Order(ent.Asc(plan.FieldDays))

	if req.BrandId != nil {
		q.Where(plan.HasBrandWith(ebikebrand.ID(*req.BrandId)))
	} else {
		q.Where(plan.HasBrandWith(ebikebrand.IDIn(brandIds...)))
	}

	items := q.AllX(s.ctx)

	// 骑士卡筛选
	items = s.FilterPlanForStore(items)

	storeEbikePlansMap := make(map[uint64][]*definition.StoreEbikePlan)

	// 所有门店骑士卡车电套餐
	for _, item := range items {
		// 查找骑士卡所属门店
		storeCheckMap := make(map[uint64]bool)
		if item.BrandID == nil {
			continue
		}
		storeIds := cityBrand2StoresMap[fmt.Sprintf("%d-%d", req.CityId, *item.BrandID)]
		for _, storeId := range storeIds {
			// 判断门店信息是否有效
			if storeMap[storeId] == nil {
				continue
			}

			// 门店电车品牌库存验证
			if !cityStoreBrandExist[fmt.Sprintf("%d-%d-%d", req.CityId, storeId, *item.BrandID)] {
				continue
			}

			// 门店查重
			if storeCheckMap[storeId] {
				continue
			}
			storeCheckMap[storeId] = true

			// 赋值
			sep := &definition.StoreEbikePlan{
				StoreId:    storeId,
				StoreName:  storeMap[storeId].Name,
				PlanId:     item.ID,
				PlanName:   item.Name,
				Rto:        item.Type == model.PlanTypeEbikeRto.Value(),
				Daily:      item.Daily,
				DailyPrice: utils.NewNumber().Decimal(item.Price / float64(item.Days)),
				MonthPrice: utils.NewNumber().Decimal(30 * item.Price / float64(item.Days)),
			}

			brand := item.Edges.Brand
			if brand != nil {
				sep.BrandId = brand.ID
				sep.BrandName = brand.Name
				sep.Cover = brand.Cover
			}

			distance, err := storeMap[storeId].Value("distance")
			if distance != nil || err == nil {
				distanceFloat, ok := distance.(float64)
				if ok {
					sep.Distance = distanceFloat
				}
			}
			storeEbikePlansMap[storeId] = append(storeEbikePlansMap[storeId], sep)
		}
	}

	var allPlans []*definition.StoreEbikePlan

	// 按照门店分组排序
	for _, v := range storelist {
		seps := storeEbikePlansMap[v.ID]
		if len(seps) > 0 {
			// 是否需要价格排序
			if req.SortType != nil && *req.SortType == definition.StorePlanSortTypePrice {
				SortPlanEbikeModelByDailyPrice(seps)
			}
			allPlans = append(allPlans, seps...)
		}
	}

	// 再次判断排序方式，重新排序数据
	if req.SortType != nil && *req.SortType == definition.StorePlanSortTypePrice {
		SortPlanEbikeModelByDailyPrice(allPlans)
	}

	return allPlans
}

// StorePlanDetail 基于门店的套餐列表详情
func (s *planBiz) StorePlanDetail(r *ent.Rider, req *definition.StorePlanDetailReq) *definition.StorePlanDetail {
	var state uint

	q := s.orm.QueryNotDeleted().
		Where(
			plan.Or(
				plan.ParentID(req.PlanId),
				plan.ID(req.PlanId),
			),
		).
		WithBrand().
		WithCities().
		WithAgreement().
		Order(ent.Asc(plan.FieldDays))

	// 查询门店库存电车所属brandId
	var brandIds []uint64
	storeItem, _ := ent.Database.Store.QueryNotDeleted().
		WithStocks().
		Where(store.ID(req.StoreId)).
		First(s.ctx)
	if storeItem == nil {
		return nil
	}

	if storeItem.Edges.Stocks != nil {
		for _, st := range storeItem.Edges.Stocks {
			if st.BrandID != nil {
				brandIds = append(brandIds, *st.BrandID)
			}
		}
	}
	if len(brandIds) > 0 {
		q.Where(
			plan.HasBrandWith(ebikebrand.IDIn(brandIds...)),
		)
	}

	item := q.FirstX(s.ctx)
	mmap := make(map[string]*model.PlanModelOption)

	var res definition.StorePlanDetail

	serv := service.NewPlanIntroduce()
	intro := serv.QueryMap()

	// 查询个签默认协议
	var defaultAgreement *ent.Agreement
	defaultAgreement, _ = ent.Database.Agreement.QueryNotDeleted().
		Where(
			agreement.UserType(model.AgreementUserTypePersonal.Value()),
			agreement.IsDefault(true),
		).First(s.ctx)

	key := serv.Key(item.Model, item.BrandID)
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
	if r != nil {
		// 判断是否有生效订阅
		_, sub := service.NewSubscribe().RecentDetail(r.ID)
		if sub != nil && slices.Contains(model.SubscribeNotUnSubscribed(), sub.Status) {
			ramount = 0
		} else {
			state, _ = service.NewOrder().PreconditionNewly(sub)
			if state == model.OrderTypeNewly && item.DiscountNewly > 0 {
				ramount = item.DiscountNewly
			}
		}
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
		RtoDays:                 item.RtoDays,
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

	SortIDOptions(*m.Children)

	if item.BrandID != nil {
		brand := item.Edges.Brand
		res = definition.StorePlanDetail{
			Children: new(model.PlanModelOptions),
			Name:     brand.Name,
			Cover:    brand.Cover,
		}

		var exists bool
		for _, c := range *res.Children {
			if c.Model == item.Model {
				exists = true
			}
		}
		if !exists {
			*res.Children = append(*res.Children, m)
		}

	}

	if r != nil {
		res.Configure = service.NewPayment(r).Configure()
	}

	settings, _ := ent.Database.Setting.Query().Where(setting.Key(model.SettingPlanEbikeDescriptionKey)).All(context.Background())
	for _, sm := range settings {
		var v model.SettingPlanDescription
		err := jsoniter.Unmarshal([]byte(sm.Content), &v)
		if err == nil {
			res.EbikeDescription = v
		}
	}

	SortPlanEbikeModelByName(*res.Children)

	return &res
}

func SortPlanEbikeModelByDailyPrice(options []*definition.StoreEbikePlan) {
	sort.Slice(options, func(i, j int) bool {
		return options[i].DailyPrice < options[j].DailyPrice
	})
}

// ListByStoreById 通过门店ID查询门店的套餐列表
func (s *planBiz) ListByStoreById(storeId uint64) []*definition.StoreEbikePlan {
	stq := ent.Database.Store.QueryNotDeleted().
		Where(
			store.ID(storeId),
			store.HasStocksWith(stock.BrandIDNotNil()),
		).
		WithCity().
		WithEmployee().
		WithStocks()

	str, _ := stq.First(s.ctx)
	if str == nil {
		return []*definition.StoreEbikePlan{}
	}

	storeMap := make(map[uint64]*ent.Store)
	storeMap[storeId] = str

	// 门店电车库存品牌ID
	var brandIds []uint64
	// 门店电车品牌集合
	cityStoreBrandExist := make(map[string]bool)
	cityBrand2StoresMap := make(map[string][]uint64)
	if str.Edges.Stocks != nil {
		// 计算该门店电车品牌库存
		brandStockMap := s.CalStoreEbikeStock(str.Edges.Stocks)
		for _, stc := range str.Edges.Stocks {
			if stc.BrandID != nil {
				existKey := fmt.Sprintf("%d-%d-%d", str.CityID, str.ID, *stc.BrandID)
				if !cityStoreBrandExist[existKey] && brandStockMap[*stc.BrandID] > 0 {
					cityStoreBrandExist[existKey] = true
					brandIds = append(brandIds, *stc.BrandID)
					// 保存存在库存的电车门店对应门店ID
					cbKey := fmt.Sprintf("%d-%d", str.CityID, *stc.BrandID)
					cityBrand2StoresMap[cbKey] = append(cityBrand2StoresMap[cbKey], str.ID)
				}
			}
		}
	}

	// 查询车电套餐
	today := carbon.Now().StartOfDay().StdTime()
	q := s.orm.QueryNotDeleted().
		Where(
			plan.TypeIn(model.PlanTypeEbikeWithBattery.Value(), model.PlanTypeEbikeRto.Value()),
			plan.Enable(true),
			plan.StartLTE(today),
			plan.EndGTE(today),
			plan.HasBrandWith(ebikebrand.IDIn(brandIds...)),
			plan.HasCitiesWith(city.ID(str.CityID)),
		).
		WithBrand().
		WithCities().
		WithAgreement().
		Order(ent.Asc(plan.FieldDays))

	items := q.AllX(s.ctx)

	// 骑士卡筛选
	items = s.FilterPlanForStore(items)

	storeEbikePlansMap := make(map[uint64][]*definition.StoreEbikePlan)

	// 所有门店骑士卡车电套餐
	for _, item := range items {
		// 查找骑士卡所属门店
		storeCheckMap := make(map[uint64]bool)
		storeIds := cityBrand2StoresMap[fmt.Sprintf("%d-%d", str.CityID, *item.BrandID)]
		if len(storeIds) != 0 {
			for _, stId := range storeIds {
				// 判断门店信息是否有效
				if storeMap[storeId] == nil {
					continue
				}

				// 门店电车品牌库存验证
				if !cityStoreBrandExist[fmt.Sprintf("%d-%d-%d", str.CityID, storeId, *item.BrandID)] {
					continue
				}

				// 门店查重
				if storeCheckMap[storeId] {
					continue
				}

				// 赋值
				sep := &definition.StoreEbikePlan{
					StoreId:    stId,
					StoreName:  storeMap[stId].Name,
					PlanId:     item.ID,
					PlanName:   item.Name,
					Rto:        item.Type == model.PlanTypeEbikeRto.Value(),
					Daily:      item.Daily,
					DailyPrice: utils.NewNumber().Decimal(item.Price / float64(item.Days)),
					MonthPrice: utils.NewNumber().Decimal(30 * item.Price / float64(item.Days)),
				}

				brand := item.Edges.Brand
				if brand != nil {
					sep.BrandId = brand.ID
					sep.BrandName = brand.Name
					sep.Cover = brand.Cover
				}

				distance, err := storeMap[stId].Value("distance")
				if distance != nil || err == nil {
					distanceFloat, ok := distance.(float64)
					if ok {
						sep.Distance = distanceFloat
					}
				}
				storeEbikePlansMap[stId] = append(storeEbikePlansMap[stId], sep)
			}
		}
	}

	var allPlans []*definition.StoreEbikePlan
	allPlans = append(allPlans, storeEbikePlansMap[str.ID]...)

	return allPlans
}

// CalStoreEbikeStock 计算库存数
func (s *planBiz) CalStoreEbikeStock(stocks []*ent.Stock) map[uint64]int {
	resMap := make(map[uint64]int)
	for _, stc := range stocks {
		if stc.BrandID != nil {
			resMap[*stc.BrandID] += stc.Num
		}
	}
	return resMap
}

func (s *planBiz) FilterPlanForStore(plans []*ent.Plan) []*ent.Plan {
	result := make([]*ent.Plan, 0)

	// 首先按照父子级整理数据
	parentPlanIds := make([]uint64, 0)
	planMap := make(map[uint64][]*ent.Plan)
	for _, pl := range plans {
		if pl.ParentID == nil {
			parentPlanIds = append(parentPlanIds, pl.ID)
			planMap[pl.ID] = append(planMap[pl.ID], pl)
		} else {
			planMap[*pl.ParentID] = append(planMap[*pl.ParentID], pl)
		}
	}

	// 每一个骑士卡集合开始筛选符合要求的骑士卡（日租套餐都需保留，月租套餐保留最接近30天的数据）
	for _, pId := range parentPlanIds {
		// 轮询同一骑士卡的所有套餐数据
		var minDiffDay float64
		var nearMonthPlan ent.Plan
		for _, pl := range planMap[pId] {
			// 日租直接放进结果集
			if pl.Daily {
				result = append(result, pl)
				continue
			}

			absDays := math.Abs(float64(pl.Days - 30))

			// 月套餐初始化
			if nearMonthPlan.ID == 0 {
				minDiffDay = absDays
				nearMonthPlan = *pl
				continue
			}

			// 找出与30天最接近的套餐
			if absDays < minDiffDay {
				minDiffDay = math.Abs(float64(pl.Days - 30))
				nearMonthPlan = *pl
			}
		}
		// 结果集加入筛选后的月套餐
		result = append(result, &nearMonthPlan)
	}

	return result
}
