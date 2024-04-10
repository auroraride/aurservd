package biz

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/golang-module/carbon/v2"
	jsoniter "github.com/json-iterator/go"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agreement"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/plan"
	"github.com/auroraride/aurservd/internal/ent/setting"
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
	}

	return res
}
