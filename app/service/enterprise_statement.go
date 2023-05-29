// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/google/uuid"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/enterprisebill"
	"github.com/auroraride/aurservd/internal/ent/enterprisestatement"
	"github.com/auroraride/aurservd/internal/ent/predicate"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/pkg/cache"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

type enterpriseStatementService struct {
	ctx      context.Context
	modifier *model.Modifier
	rider    *ent.Rider
	employee *ent.Employee
	orm      *ent.EnterpriseStatementClient
}

func NewEnterpriseStatement() *enterpriseStatementService {
	return &enterpriseStatementService{
		ctx: context.Background(),
		orm: ent.Database.EnterpriseStatement,
	}
}

func NewEnterpriseStatementWithModifier(m *model.Modifier) *enterpriseStatementService {
	s := NewEnterpriseStatement()
	s.ctx = context.WithValue(s.ctx, "modifier", m)
	s.modifier = m
	return s
}

func (s *enterpriseStatementService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// Current 获取企业当前账单, 若无则新增
func (s *enterpriseStatementService) Current(e *ent.Enterprise) *ent.EnterpriseStatement {
	res, _ := s.orm.QueryNotDeleted().Where(
		enterprisestatement.EnterpriseID(e.ID),
		enterprisestatement.SettledAtIsNil(),
	).First(s.ctx)

	if res == nil {
		res, _ = s.orm.Create().
			SetEnterpriseID(e.ID).
			SetStart(carbon.Time2Carbon(e.CreatedAt).StartOfDay().Carbon2Time()).
			Save(s.ctx)
	}
	return res
}

func (s *enterpriseStatementService) GetBill(req *model.StatementBillReq) *model.StatementBillRes {
	// 判定结算日是否早于当天
	end := tools.NewTime().ParseDateStringX(req.End)
	if !end.Before(carbon.Now().StartOfDay().Carbon2Time()) {
		snag.Panic("结算日必须早于当天")
	}

	// 判定企业
	e, _ := ent.Database.Enterprise.QueryNotDeleted().Where(enterprise.ID(req.ID)).WithCity().First(s.ctx)
	if e == nil {
		snag.Panic("未找到企业")
	}

	if e.Payment != model.EnterprisePaymentPostPay && !req.Force {
		snag.Panic("只有后付费企业才可请求")
	}

	// 查询未结账账单
	es, bills := NewEnterprise().CalculateStatement(e, end)
	if es.Start.Sub(end).Seconds() > 0 {
		msg := fmt.Sprintf("无账单信息, 账单开始日期: %s, 所选截止日期: %s", es.Start.Format(carbon.DateLayout), end.Format(carbon.DateLayout))
		snag.Panic(msg)
	}

	uid := uuid.New().String()

	res := &model.StatementBillRes{
		ID:   e.ID,
		UUID: uid,
		City: model.City{
			ID:   e.CityID,
			Name: e.Edges.City.Name,
		},
		ContactName:  e.ContactName,
		ContactPhone: e.ContactPhone,
		Start:        es.Start.Format(carbon.DateLayout),
		End:          req.End,
		Bills:        bills,
		StatementID:  es.ID,
	}

	srv := NewEnterprise()

	overview := make(map[string]*model.BillOverview)
	var cost float64
	td := tools.NewDecimal()
	for _, bill := range bills {
		cost = td.Sum(cost, bill.Cost)
		res.Days += bill.Days

		key := srv.PriceKey(bill.City.ID, bill.Model)
		ov, ok := overview[key]
		if !ok {
			ov = &model.BillOverview{
				Model:  bill.Model,
				Number: 0,
				Price:  bill.Price,
				Days:   0,
				Cost:   0,
				City:   bill.City,
			}
			overview[key] = ov
		}
		ov.Number += 1
		ov.Days += bill.Days
		ov.Cost = td.Sum(ov.Cost, bill.Cost)
	}

	for _, ov := range overview {
		res.Overview = append(res.Overview, ov)
	}

	res.Cost = math.Round(cost*100) / 100

	// 缓存账单
	cache.Set(s.ctx, uid, res, 600*time.Second)

	return res
}

// Bill 后付费企业结账
func (s *enterpriseStatementService) Bill(req *model.StatementClearBillReq) {
	// 查找账单
	br := new(model.StatementBillRes)
	err := cache.Get(s.ctx, req.UUID).Scan(br)
	if err != nil {
		snag.Panic("未找到账单信息")
		return
	}

	start := tools.NewTime().ParseDateStringX(br.Start)
	end := tools.NewTime().ParseDateStringX(br.End)

	if start.After(end) {
		snag.Panic("账单信息错误")
	}

	es, _ := s.orm.QueryNotDeleted().
		Where(
			enterprisestatement.ID(br.StatementID),
			enterprisestatement.SettledAtIsNil(),
		).
		First(s.ctx)
	if es == nil {
		snag.Panic("未找到账单信息")
	}

	// 下个账单开始日
	next := carbon.Time2Carbon(end).StartOfDay().AddDay().Carbon2Time()

	ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
		_, err = tx.EnterpriseStatement.
			UpdateOneID(br.StatementID).
			SetSettledAt(time.Now()).
			SetEnd(end).
			SetRiderNumber(len(br.Bills)).
			SetDays(br.Days).
			SetCost(br.Cost).
			SetRemark(req.Remark).
			Save(s.ctx)
		if err != nil {
			return
		}

		// 更新所有订阅
		riders := 0
		for _, bill := range br.Bills {
			var sub *ent.Subscribe
			sub, err = tx.Subscribe.UpdateOneID(bill.SubscribeID).SetLastBillDate(end).Save(s.ctx)

			if err != nil {
				return
			}

			if sub.EndAt == nil || !sub.EndAt.Before(next) {
				riders += 1
			}

			// 保存账单
			_, err = tx.EnterpriseBill.Create().
				SetEnterpriseID(es.EnterpriseID).
				SetStatementID(es.ID).
				SetStart(tools.NewTime().ParseDateStringX(bill.Start)).
				SetEnd(tools.NewTime().ParseDateStringX(bill.End)).
				SetDays(bill.Days).
				SetPrice(bill.Price).
				SetCost(bill.Cost).
				SetCityID(bill.City.ID).
				SetNillableStationID(bill.StationID).
				SetRiderID(bill.RiderID).
				SetSubscribeID(bill.SubscribeID).
				SetModel(bill.Model).
				Save(s.ctx)

			if err != nil {
				return
			}
		}

		// 创建新账单
		_, err = tx.EnterpriseStatement.Create().
			SetEnterpriseID(es.EnterpriseID).
			SetStart(next).
			SetCost(tools.NewDecimal().Sub(es.Cost, br.Cost)).
			SetDays(es.Days - br.Days).
			SetRiderNumber(riders).
			Save(s.ctx)
		return
	})

	cache.Del(s.ctx, req.UUID)
}

// Historical 获取企业结账历史
func (s *enterpriseStatementService) Historical(req *model.StatementBillHistoricalListReq) *model.PaginationRes {
	q := s.orm.QueryNotDeleted().
		Where(
			enterprisestatement.EnterpriseID(req.EnterpriseID),
			enterprisestatement.SettledAtNotNil(),
		).Order(ent.Desc(enterprisestatement.FieldSettledAt))

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.EnterpriseStatement) model.StatementBillHistoricalListRes {
			return model.StatementBillHistoricalListRes{
				ID:        item.ID,
				Cost:      item.Cost,
				Remark:    item.Remark,
				Creator:   item.Creator,
				Days:      item.Days,
				Start:     item.Start.Format(carbon.DateLayout),
				End:       item.End.Format(carbon.DateLayout),
				SettledAt: item.SettledAt.Format(carbon.DateTimeLayout),
			}
		},
	)
}

func (s *enterpriseStatementService) info(id uint64) (*ent.EnterpriseStatement, *ent.Enterprise) {
	es, _ := ent.Database.EnterpriseStatement.QueryNotDeleted().
		Where(
			enterprisestatement.ID(id),
			enterprisestatement.SettledAtNotNil(),
		).
		WithEnterprise().
		First(s.ctx)
	if es == nil {
		snag.Panic("未找到账单")
	}
	return es, es.Edges.Enterprise
}

// Statement 账单明细
func (s *enterpriseStatementService) Statement(id uint64) []model.StatementDetail {
	s.info(id)
	return s.detail(id)
}

// detail 账单明细
func (s *enterpriseStatementService) detail(id uint64) []model.StatementDetail {
	items, _ := ent.Database.EnterpriseBill.
		QueryNotDeleted().
		WithRider().
		WithCity().
		WithStation().
		Where(enterprisebill.StatementID(id)).
		All(s.ctx)

	res := make([]model.StatementDetail, len(items))
	for i, item := range items {
		r := item.Edges.Rider
		c := item.Edges.City
		t := item.Edges.Station
		res[i] = model.StatementDetail{
			Rider: model.Rider{
				ID:    r.ID,
				Phone: r.Phone,
				Name:  r.Name,
			},
			City: model.City{
				ID:   c.ID,
				Name: c.Name,
			},
			Start: item.Start.Format(carbon.DateLayout),
			End:   item.End.Format(carbon.DateLayout),
			Days:  item.Days,
			Model: item.Model,
			Price: item.Price,
			Cost:  item.Cost,
		}
		if t != nil {
			res[i].Station = &model.EnterpriseStation{
				ID:   t.ID,
				Name: t.Name,
			}
		}
	}

	return res
}

// DetailExport 账单明细导出
func (s *enterpriseStatementService) DetailExport(req *model.StatementBillDetailExport) model.ExportRes {
	es, e := s.info(req.ID)
	info := ar.Map{"企业": e.Name, "开始": es.Start.Format(carbon.ShortDateLayout), "结束": es.End.Format(carbon.ShortDateLayout)}
	return NewExportWithModifier(s.modifier).Start(e.Name+"账单明细", req.ID, info, req.Remark, func(path string) {
		items := s.Statement(req.ID)
		sheet := fmt.Sprintf("%s-%s", es.Start.Format(carbon.ShortDateLayout), es.End.Format(carbon.ShortDateLayout))
		ex := tools.NewExcel(path, sheet)
		// 设置数据
		var rows tools.ExcelItems
		rows = append(rows, []any{"姓名", "电话", "城市", "站点", "开始日期", "结束日期", "使用天数", "电池型号", "日单价", "费用"})
		for _, x := range items {
			so := ""
			if x.Station != nil {
				so = x.Station.Name
			}
			rows = append(rows, []any{
				x.Rider.Name,
				x.Rider.Phone,
				x.City.Name,
				so,
				x.Start,
				x.End,
				x.Days,
				x.Model,
				fmt.Sprintf("%.2f", x.Price),
				fmt.Sprintf("%.2f", x.Cost),
			})
		}

		ex.AddValues(rows).Done()
	})
}

func (s *enterpriseStatementService) usageFilter(e *ent.Enterprise, req model.StatementUsageFilter) (q *ent.SubscribeQuery, start time.Time, end time.Time) {
	q = ent.Database.Subscribe.QueryNotDeleted().
		WithRider().
		WithCity().
		WithStation().
		Where(
			subscribe.EnterpriseID(req.ID),
			subscribe.StartAtNotNil(),
		)

	var bw []predicate.EnterpriseBill

	if req.Start == "" {
		start = carbon.Time2Carbon(e.CreatedAt).StartOfDay().Carbon2Time()
	} else {
		start = tools.NewTime().ParseDateStringX(req.Start)
	}

	today := carbon.Now().StartOfDay().Carbon2Time()

	q.Where(
		subscribe.Or(
			subscribe.EndAtIsNil(),
			subscribe.EndAtGTE(start),
		),
	)
	bw = append(bw, enterprisebill.EndGTE(start))

	if req.End == "" {
		end = carbon.Now().StartOfDay().Carbon2Time()
	} else {
		end = tools.NewTime().ParseDateStringX(req.End)
	}

	if end.After(today) {
		end = today
	}

	next := end.AddDate(0, 0, 1)

	// 开始时间早于结束时间
	q.Where(subscribe.StartAtLT(next))
	bw = append(bw, enterprisebill.StartLT(next))

	q.WithBills(func(bq *ent.EnterpriseBillQuery) {
		if len(bw) > 0 {
			bq.Where(bw...)
		}
		bq.Order(ent.Asc(enterprisebill.FieldEnd))
	})

	if req.CityID != 0 {
		q.Where(subscribe.CityID(req.CityID))
	}

	if req.StationID != 0 {
		q.Where(subscribe.StationID(req.StationID))
	}

	if req.Model != "" {
		q.Where(subscribe.Model(req.Model))
	}

	if req.Keyword != "" {
		q.Where(
			subscribe.HasRiderWith(rider.Or(
				rider.NameContainsFold(req.Keyword),
				rider.PhoneContainsFold(req.Keyword),
			)),
		)
	}
	return
}

// Usage 使用明细
func (s *enterpriseStatementService) Usage(req *model.StatementUsageReq) *model.PaginationRes {
	e := NewEnterprise().QueryX(req.ID)
	prices := NewEnterprise().GetPriceValues(e)
	q, start, end := s.usageFilter(e, req.StatementUsageFilter)
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Subscribe) model.StatementUsageRes {
		return s.usageDetail(e, item, start, end, prices)
	})
}

func (s *enterpriseStatementService) UsageExport(req *model.StatementUsageExport) model.ExportRes {
	e := NewEnterprise().QueryX(req.ID)
	prices := NewEnterprise().GetPriceValues(e)
	filter := model.StatementUsageFilter{
		ID:    req.ID,
		Start: req.Start,
		End:   req.End,
	}
	q, start, end := s.usageFilter(e, filter)
	taxonomy := fmt.Sprintf("%s使用明细", e.Name)

	info := map[string]interface{}{
		"企业":   e.Name,
		"开始时间": start.Format(carbon.DateLayout),
		"结束时间": end.Format(carbon.DateLayout),
	}

	return NewExportWithModifier(s.modifier).Start(taxonomy, filter, info, req.Remark, func(path string) {
		sheet := fmt.Sprintf("%s-%s", start.Format(carbon.ShortDateLayout), end.Format(carbon.ShortDateLayout))
		ex := tools.NewExcel(path, sheet)

		items, _ := q.All(s.ctx)

		var rows tools.ExcelItems
		rows = append(rows, []any{"城市", "姓名", "电话", "站点", "型号", "状态", "删除时间", "开始日期", "结束日期", "使用天数", "日单价", "费用"})
		for _, item := range items {
			detail := s.usageDetail(e, item, start, end, prices)
			sta := ""
			if detail.Station != nil {
				sta = detail.Station.Name
			}
			row := []any{
				detail.City.Name,
				detail.Rider.Name,
				detail.Rider.Phone,
				sta,
				detail.Model,
				detail.Status,
				detail.DeletedAt,
			}
			var sub tools.ExcelItems
			for _, ui := range detail.Items {
				sub = append(sub, []any{
					ui.Start,
					ui.End,
					ui.Days,
					ui.Price,
					ui.Cost,
				})
			}
			row = append(row, sub)

			rows = append(rows, row)
		}

		ex.AddValues(rows).Done()
	})
}

func (s *enterpriseStatementService) usageDetail(en *ent.Enterprise, item *ent.Subscribe, start, end time.Time, prices map[string]float64) model.StatementUsageRes {
	r := item.Edges.Rider
	c := item.Edges.City
	del := ""
	if r.DeletedAt != nil {
		del = r.DeletedAt.Format(carbon.DateTimeLayout)
	}
	status := model.SubscribeStatusText(item.Status)
	today := carbon.Now().StartOfDay().Carbon2Time()
	if en.Agent && item.AgentEndAt != nil && item.AgentEndAt.Before(today) {
		status = "已超期"
	}
	res := model.StatementUsageRes{
		Model: item.Model,
		City: model.City{
			ID:   c.ID,
			Name: c.Name,
		},
		Rider: model.Rider{
			ID:    r.ID,
			Phone: r.Phone,
			Name:  r.Name,
		},
		Status:    status,
		DeletedAt: del,
		Items:     s.usageItems(item, start, end, prices),
	}
	a := item.Edges.Station
	if a != nil {
		res.Station = &model.EnterpriseStation{
			ID:   a.ID,
			Name: a.Name,
		}
	}
	return res
}

// usageItems 计算使用项
func (s *enterpriseStatementService) usageItems(sub *ent.Subscribe, start time.Time, end time.Time, prices map[string]float64) (items []*model.StatementUsageItem) {
	if sub.Edges.Bills != nil {
		for _, bill := range sub.Edges.Bills {
			if data := s.usageItemCalculate(start, end, bill.Start, bill.End, bill.Price); data != nil {
				items = append(items, data)
			}
		}
	}
	subStart := *sub.StartAt
	// 订阅 上次结账日如果晚于结束日期第二天 直接跳过
	next := end.AddDate(0, 0, 1)
	if sub.LastBillDate != nil {
		if sub.LastBillDate.After(next) {
			return
		}
		subStart = *sub.LastBillDate
	}
	var endAt time.Time
	if sub.EndAt != nil {
		endAt = *sub.EndAt
	}
	if data := s.usageItemCalculate(start, end, subStart, endAt, prices[NewEnterprise().PriceKey(sub.CityID, sub.Model)]); data != nil {
		items = append(items, data)
	}
	return
}

// usageItemCalculate 计算日期时间和价格
func (s *enterpriseStatementService) usageItemCalculate(srcStart time.Time, srcEnd time.Time, destStart time.Time, destEnd time.Time, price float64) *model.StatementUsageItem {
	from := destStart
	to := srcEnd

	if !srcStart.IsZero() {
		if srcStart.After(destStart) {
			from = srcStart
		}
	}
	if !destEnd.IsZero() {
		if destEnd.Before(srcEnd) {
			to = destEnd
		}
	}

	days := tools.NewTime().UsedDays(to, from)
	if from.After(to) {
		return nil
	}
	return &model.StatementUsageItem{
		Start: from.Format(carbon.DateLayout),
		End:   to.Format(carbon.DateLayout),
		Days:  days,
		Price: price,
		Cost:  tools.NewDecimal().Mul(float64(days), price),
	}
}

func (s *enterpriseStatementService) HistoryCost(enterpriseID uint64) float64 {
	var result []struct {
		ID  uint64  `json:"id"`
		Sum float64 `json:"sum"`
	}

	_ = s.orm.Query().Modify().
		Where(enterprisestatement.EnterpriseID(enterpriseID), enterprisestatement.SettledAtNotNil()).
		GroupBy(enterprisestatement.FieldID).
		Aggregate(ent.Sum(enterprisestatement.FieldCost)).
		Scan(s.ctx, &result)

	if result == nil {
		return 0
	}

	return result[0].Sum
}
