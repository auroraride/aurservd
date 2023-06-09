// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-05
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/golang-module/carbon/v2"
	"github.com/jinzhu/copier"
	"github.com/r3labs/diff/v3"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/enterprisecontract"
	"github.com/auroraride/aurservd/internal/ent/enterpriseprice"
	"github.com/auroraride/aurservd/internal/ent/enterprisestatement"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/internal/ent/subscribealter"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type enterpriseService struct {
	ctx        context.Context
	modifier   *model.Modifier
	orm        *ent.EnterpriseClient
	agent      *ent.Agent
	enterprise *ent.Enterprise
}

func NewEnterprise() *enterpriseService {
	return &enterpriseService{
		ctx: context.Background(),
		orm: ent.Database.Enterprise,
	}
}

func NewEnterpriseWithModifier(m *model.Modifier) *enterpriseService {
	s := NewEnterprise()
	s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
	s.modifier = m
	return s
}

func NewEnterpriseWithAgent(ag *ent.Agent, en *ent.Enterprise) *enterpriseService {
	s := NewEnterprise()
	s.agent = ag
	s.enterprise = en
	return s
}

func (s *enterpriseService) Query(id uint64) (*ent.Enterprise, error) {
	return s.orm.QueryNotDeleted().Where(enterprise.ID(id)).First(s.ctx)
}

func (s *enterpriseService) QueryX(id uint64) *ent.Enterprise {
	e, _ := s.Query(id)
	if e == nil {
		snag.Panic("未找到有效企业")
	}
	return e
}

// Create 创建企业
func (s *enterpriseService) Create(req *model.EnterpriseDetail) uint64 {
	// 判断是否代理商并且非预付费
	if req.Agent != nil && *req.Agent && *req.Payment != model.EnterprisePaymentPrepay {
		snag.Panic("代理商模式付费方式错误")
	}

	var err error
	e := &ent.Enterprise{}
	e, err = ent.EntitySetAttributes(s.orm.Create(), e, req).Save(s.ctx)
	if err != nil {
		snag.Panic("企业创建失败")
	}
	return e.ID
}

// Modify 修改企业
// TODO 后付费转为预付费, 预付费转为后付费
func (s *enterpriseService) Modify(req *model.EnterpriseDetailWithID) {
	e := s.QueryX(req.ID)
	// 判定付费方式是否允许改变
	if *req.Payment != e.Payment {
		// set := NewEnterpriseStatementWithModifier(s.modifier).StoreCurrent(e)
		// if set.Cost > 0 {
		//     snag.Panic("企业已产生费用, 无法修改支付方式")
		// }
		snag.Panic("无法转换支付方式")
	}

	_, err := s.orm.ModifyOne(e, req.EnterpriseDetail).Save(s.ctx)
	if err != nil {
		snag.Panic("企业修改失败")
	}
}

// PriceKey 获取企业价格key (城市id-电池型号)
func (s *enterpriseService) PriceKey(cityID uint64, model string) string {
	return fmt.Sprintf("%d-%s", cityID, model)
}

// DetailQuery 企业详情查询语句
func (s *enterpriseService) DetailQuery() *ent.EnterpriseQuery {
	return s.orm.QueryNotDeleted().WithCity().
		WithPrices(func(ep *ent.EnterprisePriceQuery) {
			ep.Where(enterpriseprice.DeletedAtIsNil()).WithCity()
		}).
		WithContracts(func(ecq *ent.EnterpriseContractQuery) {
			ecq.Where(enterprisecontract.DeletedAtIsNil()).Order(ent.Desc(enterprisecontract.FieldEnd))
		}).
		WithRiders(func(erq *ent.RiderQuery) {
			erq.Where(rider.DeletedAtIsNil())
		}).
		WithStatements(func(esq *ent.EnterpriseStatementQuery) {
			esq.Where(enterprisestatement.SettledAtIsNil())
		})
}

// List 列举企业
func (s *enterpriseService) List(req *model.EnterpriseListReq) *model.PaginationRes {
	tt := tools.NewTime()

	q := s.DetailQuery()
	if req.Name != nil {
		q.Where(enterprise.NameContainsFold(*req.Name))
	}
	if req.CityID != nil {
		q.Where(enterprise.CityID(*req.CityID))
	}
	if req.Status != nil {
		q.Where(enterprise.Status(*req.Status))
	}
	if req.Payment != nil {
		q.Where(enterprise.Payment(*req.Payment))
	}
	if req.ContactKeyword != nil {
		q.Where(enterprise.Or(
			enterprise.ContactNameContainsFold(*req.ContactKeyword),
			enterprise.ContactPhoneContainsFold(*req.ContactKeyword),
			enterprise.IdcardNumberContainsFold(*req.ContactKeyword),
		))
	}
	if req.Start != nil {
		q.Where(enterprise.HasContractsWith(enterprisecontract.StartLTE(tt.ParseDateStringX(*req.Start))))
	}
	if req.End != nil {
		q.Where(enterprise.HasContractsWith(enterprisecontract.EndGTE(tt.ParseDateStringX(*req.End))))
	}
	if req.Agent != nil {
		q.Where(enterprise.Agent(*req.Agent))
	}
	return model.ParsePaginationResponse(
		q, req.PaginationReq,
		func(item *ent.Enterprise) model.EnterpriseRes {
			return s.Detail(item)
		},
	)
}

// GetDetail 获取企业详情
func (s *enterpriseService) GetDetail(req *model.IDParamReq) model.EnterpriseRes {
	item, _ := s.DetailQuery().Where(enterprise.ID(req.ID)).
		WithStatements(func(esq *ent.EnterpriseStatementQuery) {
			esq.Where(enterprisestatement.SettledAtIsNil())
		}).First(s.ctx)
	if item == nil {
		snag.Panic("未找到有效企业")
	}
	return s.Detail(item)
}

// Detail 企业详情
func (s *enterpriseService) Detail(item *ent.Enterprise) (res model.EnterpriseRes) {
	_ = copier.Copy(&res, item)
	res.Riders = len(item.Edges.Riders)
	res.City = model.City{
		ID:   item.Edges.City.ID,
		Name: item.Edges.City.Name,
	}
	contracts := item.Edges.Contracts
	if contracts != nil {
		res.Contracts = make([]model.EnterpriseContract, len(contracts))
		for i, ec := range contracts {
			res.Contracts[i] = model.EnterpriseContract{
				ID:    ec.ID,
				Start: ec.Start.Format(carbon.DateLayout),
				End:   ec.End.Format(carbon.DateLayout),
				File:  ec.File,
			}
		}
	}

	prices := item.Edges.Prices
	if prices != nil {
		res.Prices = make([]model.EnterprisePriceWithCity, len(prices))
		for i, ep := range prices {
			res.Prices[i] = model.EnterprisePriceWithCity{
				ID:          ep.ID,
				Model:       ep.Model,
				Price:       ep.Price,
				Intelligent: ep.Intelligent,
				City: model.City{
					ID:   ep.Edges.City.ID,
					Name: ep.Edges.City.Name,
				},
			}
		}
	}

	stas := item.Edges.Statements
	if item.Payment == model.EnterprisePaymentPostPay && len(stas) > 0 {
		res.Unsettlement = stas[0].Days
		res.StatementStart = stas[0].Start.Format(carbon.DateLayout)
	}

	return
}

// QueryAllCollaborated 获取合作中的企业
func (s *enterpriseService) QueryAllCollaborated() []*ent.Enterprise {
	items, _ := s.orm.QueryNotDeleted().
		Where(enterprise.Status(model.EnterpriseStatusCollaborated)).
		WithPrices(func(query *ent.EnterprisePriceQuery) {
			query.Where(enterpriseprice.DeletedAtIsNil())
		}).
		All(s.ctx)
	return items
}

// QueryAllSubscribe 获取企业所有订阅
func (s *enterpriseService) QueryAllSubscribe(enterpriseID uint64, args ...string) []*ent.Subscribe {
	q := ent.Database.Subscribe.
		QueryNotDeleted().
		WithCity().
		Where(
			// 所属企业
			subscribe.EnterpriseID(enterpriseID),
			// 已开始使用
			subscribe.StartAtNotNil(),
		)
	if len(args) > 0 {
		q.Where(subscribe.StartAtGTE(tools.NewTime().ParseDateStringX(args[0])))
	}
	if len(args) > 1 {
		q.Where(subscribe.StartAtLT(tools.NewTime().ParseNextDateStringX(args[1])))
	}
	items, _ := q.All(s.ctx)
	return items
}

// QueryAllBillingSubscribe 获取所有待结算骑手团签订阅
func (s *enterpriseService) QueryAllBillingSubscribe(enterpriseID uint64, args ...any) []*ent.Subscribe {
	q := ent.Database.Subscribe.QueryNotDeleted().
		WithCity().
		Where(
			// 所属企业
			subscribe.EnterpriseID(enterpriseID),
			// 已开始使用
			subscribe.StartAtNotNil(),
		)

	end := time.Now()
	if len(args) > 0 {
		if d, ok := args[0].(time.Time); ok {
			end = d
		}
	}
	endDate := carbon.Time2Carbon(end).StartOfDay().Carbon2Time()
	// 获取未结算订阅
	q.Where(
		subscribe.Or(
			// 或上次结算日期为空
			subscribe.LastBillDateIsNil(),
			// 或小于截止日期
			subscribe.LastBillDateLT(endDate),
		),
		subscribe.Or(
			// 或上次结算日期为空
			subscribe.LastBillDateIsNil(),
			// 或者退租日期为空的
			subscribe.EndAtIsNil(),
			// 或者终止时间晚于上次已结算时间的
			func(selector *sql.Selector) {
				selector.Where(sql.ColumnsGTE(selector.C(subscribe.FieldEndAt), fmt.Sprintf("%s + INTERVAL '1 day'", selector.C(subscribe.FieldLastBillDate))))
			},
		),
	)

	items, _ := q.All(s.ctx)
	return items
}

// DiffPrices 价格修改: 获取价格设定差异
func (s *enterpriseService) DiffPrices(item *ent.Enterprise, data []model.EnterprisePrice) (changes diff.Changelog) {
	src := s.GetPriceValues(item)
	dst := make(map[string]float64)
	for _, d := range data {
		dst[s.PriceKey(d.CityID, d.Model)] = d.Price
	}
	var err error
	changes, err = diff.Diff(src, dst, diff.Filter(func(path []string, parent reflect.Type, field reflect.StructField) bool {
		return path[1] == "Price"
	}))
	if err != nil {
		snag.Panic(fmt.Sprintf("价格对比失败: %s", err))
	}
	return
}

// GetPrices 获取企业费用表
func (s *enterpriseService) GetPrices(item *ent.Enterprise) map[string]model.EnterprisePrice {
	var items []*ent.EnterprisePrice
	if item.Edges.Prices == nil {
		items, _ = ent.Database.EnterprisePrice.QueryNotDeleted().Where(enterpriseprice.EnterpriseID(item.ID)).WithCity().All(s.ctx)
	} else {
		items = item.Edges.Prices
	}
	res := make(map[string]model.EnterprisePrice)
	for _, price := range items {
		ci := price.Edges.City
		cid := item.CityID
		cname := ""
		if ci != nil {
			cname = ci.Name
		}
		res[s.PriceKey(price.CityID, price.Model)] = model.EnterprisePrice{
			CityID:   cid,
			CityName: cname,
			Model:    price.Model,
			Price:    price.Price,
			ID:       price.ID,
		}
	}
	return res
}

// GetPriceValues 获取企业价格表
func (s *enterpriseService) GetPriceValues(item *ent.Enterprise) (res map[string]float64) {
	res = make(map[string]float64)
	for key, price := range s.GetPrices(item) {
		res[key] = price.Price
	}
	return res
}

func (s *enterpriseService) CalculateStatement(e *ent.Enterprise, end time.Time) (es *ent.EnterpriseStatement, bills []model.StatementBillData) {
	tt := tools.NewTime()
	prices := s.GetPriceValues(e)
	es = NewEnterpriseStatement().Current(e)

	// 获取所有骑手未结算订阅
	subs := s.QueryAllBillingSubscribe(e.ID, end)
	for _, sub := range subs {
		// 是否已终止并且终止时间早于结算时间
		if sub.LastBillDate != nil && sub.EndAt != nil && sub.LastBillDate.After(*sub.EndAt) {
			continue
		}

		// 使用天数
		var used int

		// 开始时间
		from := carbon.Time2Carbon(*sub.StartAt).StartOfDay().Carbon2Time()

		// 上次结算日期存在则从上次结算日次日开始计算
		if sub.LastBillDate != nil {
			from = carbon.Time2Carbon(*sub.LastBillDate).Tomorrow().AddDay().Carbon2Time()
		}

		// 结束时间
		to := end
		// 判定是否退订
		if sub.EndAt != nil && sub.EndAt.Before(end) {
			to = carbon.Time2Carbon(*sub.EndAt).StartOfDay().Carbon2Time()
		}
		if to.Before(from) {
			continue
		}

		used = tt.UsedDays(to, from)

		// 按城市/型号计算金额
		k := s.PriceKey(sub.CityID, sub.Model)
		p := prices[k]

		cost, _ := decimal.NewFromFloat(p).Mul(decimal.NewFromInt(int64(used))).Float64()

		bills = append(bills, model.StatementBillData{
			EnterpriseID: *sub.EnterpriseID,
			RiderID:      sub.RiderID,
			SubscribeID:  sub.ID,
			City: model.City{
				ID:   sub.Edges.City.ID,
				Name: sub.Edges.City.Name,
			},
			StatementID: es.ID,
			StationID:   sub.StationID,
			Days:        used,
			End:         to.Format(carbon.DateLayout),
			Start:       from.Format(carbon.DateLayout),
			Cost:        cost,
			Price:       p,
			Model:       sub.Model,
		})
	}

	return
}

func (s *enterpriseService) UpdateStatementByID(id uint64) {
	e, _ := ent.Database.Enterprise.QueryNotDeleted().Where(enterprise.ID(id)).First(s.ctx)
	if e != nil {
		s.UpdateStatement(e)
	}
}

// UpdateStatement 更新企业账单
func (s *enterpriseService) UpdateStatement(e *ent.Enterprise) {
	sta, bills := s.CalculateStatement(e, time.Now())

	// 总天数
	var days int
	// 总计费用
	var cost float64
	td := tools.NewDecimal()
	for _, bill := range bills {
		cost = td.Sum(cost, bill.Cost)
		days += bill.Days
	}

	// 统计历史轧账
	cost += NewEnterpriseStatement().HistoryCost(e.ID)

	// 企业付款方式
	var balance float64
	switch e.Payment {
	case model.EnterprisePaymentPrepay:
		// 预付费, 计算余额
		balance = tools.NewDecimal().Sub(e.PrepaymentTotal, cost)
	}

	_, err := e.Update().SetBalance(balance).Save(s.ctx)
	if err != nil {
		zap.L().Error("企业更新失败: "+strconv.FormatUint(e.ID, 10), zap.Error(err))
	}

	now := carbon.Now().StartOfDay().Carbon2Time()
	_, err = sta.Update().
		SetRiderNumber(len(bills)).
		SetDays(days).
		SetCost(cost).
		SetDate(now).
		Save(s.ctx)

	if err != nil {
		zap.L().Error("企业更新失败: "+strconv.FormatUint(e.ID, 10), zap.Error(err))
	}

	zap.L().Info("企业更新成功: " + e.Name +
		"[" + strconv.FormatUint(e.ID, 10) + "]" +
		", 总使用人数: " + strconv.Itoa(len(bills)) +
		", 账期使用总天数: " + strconv.Itoa(days) +
		", 总费用: " + strconv.FormatFloat(cost, 'f', 2, 64) +
		", 余额: " + strconv.FormatFloat(balance, 'f', 2, 64) +
		", 出账日期: " + now.Format(carbon.DateLayout))
}

// Prepayment 后台修改预付费
func (s *enterpriseService) Prepayment(req *model.EnterprisePrepaymentReq) float64 {
	balance, err := NewPrepayment(s.modifier).UpdateBalance(model.PaywayCash, &model.AgentPrepay{
		EnterpriseID: req.ID,
		Remark:       req.Remark,
		Amount:       req.Amount,
	})
	if err != nil {
		snag.Panic(err)
	}
	return balance
}

// Business 企业是否可以办理业务
func (s *enterpriseService) Business(e *ent.Enterprise) error {
	if e == nil {
		return errors.New("未找到企业信息")
	}

	if e.Status != model.EnterpriseStatusCollaborated {
		return errors.New("企业已终止合作")
	}

	if e.Payment == model.EnterprisePaymentPrepay && e.Balance < 0 {
		return errors.New("企业已欠费")
	}

	return nil
}

// QueryPrice 查找价格设置
func (s *enterpriseService) QueryPrice(id uint64) *ent.EnterprisePrice {
	p, _ := ent.Database.EnterprisePrice.QueryNotDeleted().Where(enterpriseprice.ID(id)).First(s.ctx)
	if p == nil {
		snag.Panic("未找到价格信息")
	}
	return p
}

// ModifyPrice 修改价格
func (s *enterpriseService) ModifyPrice(req *model.EnterprisePriceModifyReq) model.EnterprisePriceWithCity {
	c := NewCity().Query(req.CityID)

	var p *ent.EnterprisePrice
	var err error
	if req.ID != 0 {
		p = s.QueryPrice(req.ID)
		if p.CityID != req.CityID {
			snag.Panic("城市无法修改")
		}
		if p.Model != req.Model {
			snag.Panic("电池型号无法修改")
		}
		if p.Price != req.Price {
			// 自动轧账
			srv := NewEnterpriseStatementWithModifier(s.modifier)
			// 获取账单信息
			info := srv.GetBill(&model.StatementBillReq{
				End:   carbon.Now().Yesterday().StartOfDay().Format("Y-m-d"),
				ID:    req.EnterpriseID,
				Force: true,
			})
			// 轧账
			srv.Bill(&model.StatementClearBillReq{
				UUID:   info.UUID,
				Remark: "修改价格自动轧账",
			})
			p, err = p.Update().SetPrice(req.Price).Save(s.ctx)
			s.UpdateStatementByID(req.EnterpriseID)
		}
	} else {
		client := ent.Database.EnterprisePrice
		// 判定价格是否重复
		if exist, _ := client.QueryNotDeleted().
			Where(enterpriseprice.EnterpriseID(req.EnterpriseID), enterpriseprice.Model(req.Model), enterpriseprice.CityID(req.CityID), enterpriseprice.Intelligent(req.Intelligent)).
			Exist(s.ctx); exist {
			snag.Panic("价格设置重复")
		}
		p, err = client.
			Create().
			SetCityID(req.CityID).
			SetEnterpriseID(req.EnterpriseID).
			SetModel(req.Model).
			SetPrice(req.Price).
			SetIntelligent(req.Intelligent).
			Save(s.ctx)
	}
	if err != nil {
		snag.Panic("企业价格操作失败")
	}
	return model.EnterprisePriceWithCity{
		ID:    p.ID,
		Model: p.Model,
		Price: p.Price,
		City: model.City{
			ID:   c.ID,
			Name: c.Name,
		},
	}
}

// DeletePrice 删除价格
func (s *enterpriseService) DeletePrice(req *model.IDParamReq) {
	// 判断是否有进行中的订阅
	p := s.QueryPrice(req.ID)
	if exist, _ := ent.Database.Subscribe.QueryNotDeleted().Where(
		subscribe.EnterpriseID(p.EnterpriseID),
		subscribe.Status(model.SubscribeStatusUsing),
		subscribe.CityID(p.CityID),
		subscribe.Model(p.Model),
	).Exist(s.ctx); exist {
		snag.Panic("计费中, 无法删除")
	}
	err := ent.Database.EnterprisePrice.SoftDeleteOne(p).Exec(s.ctx)
	if err != nil {
		snag.Panic("价格删除失败")
	}
}

// QueryContract 查找合同
func (s *enterpriseService) QueryContract(contractID uint64) *ent.EnterpriseContract {
	c, _ := ent.Database.EnterpriseContract.QueryNotDeleted().Where(enterprisecontract.ID(contractID)).First(s.ctx)
	if c == nil {
		snag.Panic("未找到合同")
	}
	return c
}

// ModifyContract 编辑企业合同
func (s *enterpriseService) ModifyContract(req *model.EnterpriseContractModifyReq) {
	tt := tools.NewTime()
	start := tt.ParseDateStringX(req.Start)
	end := tt.ParseDateStringX(req.End)
	var err error
	if req.ID != 0 {
		// 查找合同
		_, err = s.QueryContract(req.ID).
			Update().
			SetStart(start).
			SetEnd(end).
			SetFile(req.File).
			Save(s.ctx)
	} else {
		client := ent.Database.EnterpriseContract
		// 判断合同日期是否重复
		if exist, _ := client.QueryNotDeleted().Where(
			enterprisecontract.EnterpriseID(req.EnterpriseID),
			enterprisecontract.StartLTE(end),
			enterprisecontract.EndGTE(start),
		).Exist(s.ctx); exist {
			snag.Panic("合同日期段重复")
		}

		err = client.Create().
			SetStart(start).
			SetEnd(end).
			SetFile(req.File).
			SetEnterpriseID(req.EnterpriseID).
			Exec(s.ctx)
	}
	if err != nil {
		snag.Panic("合同操作失败")
	}
}

// DeleteContract 删除企业合同
func (s *enterpriseService) DeleteContract(req *model.IDParamReq) {
	err := ent.Database.EnterpriseContract.SoftDeleteOne(s.QueryContract(req.ID)).Exec(s.ctx)
	if err != nil {
		snag.Panic("合同删除失败")
	}
}

func (s *enterpriseService) NameFromID(id uint64) string {
	p, _ := ent.Database.Enterprise.QueryNotDeleted().Where(enterprise.ID(id)).First(s.ctx)
	if p == nil {
		return "-"
	}
	return p.Name
}

// SubscribeApplyList 骑手订阅申请加时列表
func (s *enterpriseService) SubscribeApplyList(req *model.SubscribeAlterApplyReq, enterpriseId uint64) *model.PaginationRes {
	q := ent.Database.SubscribeAlter.QueryNotDeleted().Where(subscribealter.EnterpriseID(enterpriseId)).Where(subscribealter.HasRiderWith(rider.DeletedAtIsNil())).
		Order(ent.Desc(subscribealter.FieldCreatedAt)).WithRider().WithSubscribe()

	tt := tools.NewTime()
	if req.Start != nil {
		q.Where(subscribealter.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
	}
	if req.End != nil {
		q.Where(subscribealter.CreatedAtLT(tt.ParseNextDateStringX(*req.End)))
	}

	if req.Status != nil {
		q.Where(subscribealter.Status(*req.Status))
	}
	if req.Keyword != nil {
		q.Where(subscribealter.HasRiderWith(rider.Or(rider.NameContainsFold(*req.Keyword),
			rider.PhoneContainsFold(*req.Keyword))))
	}
	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.SubscribeAlter) model.SubscribeAlterApplyListRsp {
			rsp := model.SubscribeAlterApplyListRsp{
				ID:   item.ID,
				Days: item.Days,
				// 申请时间
				ApplyTime: item.CreatedAt.Format(carbon.DateTimeLayout),
				// 审批状态
				Status: item.Status,
			}
			if item.ExpireTime != nil {
				rsp.ReviewTime = item.ExpireTime.Format(carbon.DateTimeLayout)
			}
			if item.ReviewTime != nil {
				rsp.ReviewTime = item.ReviewTime.Format(carbon.DateTimeLayout)
			}
			if item.Edges.Rider != nil {
				// 骑手姓名
				rsp.RiderName = item.Edges.Rider.Name
				// 骑手手机号
				rsp.RiderPhone = item.Edges.Rider.Phone
			}
			return rsp
		})
}

// SubscribeApplyReviewApply 审批加时申请
func (s *enterpriseService) SubscribeApplyReviewApply(req *model.SubscribeAlterReviewReq) {
	// 查找申请记录
	q := ent.Database.SubscribeAlter.QueryNotDeleted().Where(subscribealter.IDIn(req.Ids...)).Where(subscribealter.HasRiderWith(rider.DeletedAtIsNil()))
	if req.EnterpriseID != nil {
		q.Where(subscribealter.EnterpriseID(*req.EnterpriseID))
	}
	alter, err := q.All(s.ctx)
	if err != nil || len(alter) == 0 {
		snag.Panic("申请记录不存在")
	}
	for _, v := range alter {
		if v.Status != model.SubscribeAlterUnreviewed {
			snag.Panic("申请已审批")
		}
		// 事务
		ent.WithTxPanic(s.ctx, func(tx *ent.Tx) error {
			// 查询订阅信息
			sub, _ := tx.Subscribe.Query().Where(subscribe.ID(v.SubscribeID)).First(s.ctx)
			if sub == nil || sub.Status == model.SubscribeStatusUnSubscribed || sub.Status == model.SubscribeStatusCanceled {
				zap.L().Log(zap.ErrorLevel, "订阅信息不存在")
				return errors.New("订阅信息不存在")
			}
			err = tx.SubscribeAlter.UpdateOne(v).SetStatus(req.Status).SetReviewTime(time.Now()).Exec(s.ctx)
			if err != nil {
				zap.L().Log(zap.ErrorLevel, "审批加时申请失败", zap.Error(err))
				return err
			}
			// 剩余天数
			before := tools.NewTime().LastDaysToNow(*sub.AgentEndAt)
			// 加时后的结束时间
			after := before + v.Days

			if after > 0 && sub.Status == model.SubscribeStatusOverdue {
				sub.Status = model.SubscribeStatusUsing
			}
			if after < 0 {
				sub.Status = model.SubscribeStatusOverdue
			}
			// 更新订阅时间
			if err = tx.Subscribe.UpdateOne(sub).AddAlterDays(v.Days).
				SetAgentEndAt(tools.NewTime().WillEnd(*sub.AgentEndAt, v.Days, true)).
				SetStatus(sub.Status).
				Exec(s.ctx); err != nil {
				zap.L().Log(zap.ErrorLevel, "更新订阅时间失败", zap.Error(err))
				return err
			}
			return nil
		})
	}
}

// PriceList 团签价格列表
func (s *enterpriseService) PriceList(enterpriseId uint64) (rsp []model.EnterprisePriceWithCity) {
	pr, _ := ent.Database.EnterprisePrice.QueryNotDeleted().WithCity().Where(enterpriseprice.EnterpriseID(enterpriseId)).All(s.ctx)
	if pr == nil {
		snag.Panic("团签价格不存在")
	}
	rsp = make([]model.EnterprisePriceWithCity, 0)
	for _, ep := range pr {
		data := model.EnterprisePriceWithCity{
			ID:          ep.ID,
			Model:       ep.Model,
			Price:       ep.Price,
			Intelligent: ep.Intelligent,
			City: model.City{
				ID:   ep.Edges.City.ID,
				Name: ep.Edges.City.Name,
			},
		}
		rsp = append(rsp, data)
	}
	return
}
