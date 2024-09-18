// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"context"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/auroraride/adapter"
	"github.com/golang-module/carbon/v2"

	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ar"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/business"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/internal/ent/city"
	"github.com/auroraride/aurservd/internal/ent/employee"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/rider"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/internal/ent/subscribepause"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
)

// 门店业务处理专用
type businessService struct {
	ctx          context.Context
	employee     *ent.Employee
	modifier     *model.Modifier
	employeeInfo *model.Employee
}

func NewBusiness() *businessService {
	return &businessService{
		ctx: context.Background(),
	}
}

func (s *businessService) Text(str interface{}) string {
	var t model.BusinessType
	switch x := str.(type) {
	case model.BusinessType:
		t = x
	case string:
		t = model.BusinessType(x)
	case *string:
		t = model.BusinessType(*x)
	}
	m := map[model.BusinessType]string{
		model.BusinessTypeActive:      "激活",
		model.BusinessTypeUnsubscribe: "退租",
		model.BusinessTypePause:       "寄存",
		model.BusinessTypeContinue:    "结束寄存",
	}
	return m[t]
}

func NewBusinessWithEmployee(e *ent.Employee) *businessService {
	s := NewBusiness()
	if e != nil {
		s.employee = e
		s.employeeInfo = &model.Employee{
			ID:    e.ID,
			Name:  e.Name,
			Phone: e.Phone,
		}
		s.ctx = context.WithValue(s.ctx, model.CtxEmployeeKey{}, s.employeeInfo)
	}
	return s
}

func NewBusinessWithModifier(m *model.Modifier) *businessService {
	s := NewBusiness()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// CheckCity 检查城市
func (s *businessService) CheckCity(cityID uint64, sto *ent.Store) {
	if s.employee != nil && sto != nil && sto.CityID != cityID {
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

	s.CheckCity(subd.City.ID, s.employee.Edges.Store)

	var ic string
	if len(r.IDCardNumber) > 0 {
		ic = r.IDCardNumber
		ic = ic[len(ic)-4:]
	}

	res = model.SubscribeBusiness{
		ID:           r.ID,
		Status:       subd.Status,
		Name:         r.Name,
		Phone:        r.Phone,
		IDCardNumber: ic,
		Model:        subd.Model,
		SubscribeID:  subd.ID,
		Ebike:        subd.Ebike,
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

	s.CheckCity(sub.CityID, s.employee.Edges.Store)
	NewRider().CheckForBusiness(sub.Edges.Rider)

	if sub.Status != model.SubscribeStatusUsing {
		snag.Panic("当前为非使用中")
	}

	// 获取全部的电压列表
}

// listFilter 列表基础查询语句
func (s *businessService) listFilter(req model.BusinessFilter) (q *ent.BusinessQuery, info ar.Map) {
	info = make(ar.Map)
	tt := tools.NewTime()

	q = ent.Database.Business.
		QueryNotDeleted().
		WithRider().
		WithEnterprise().
		WithStation().
		WithPlan().
		WithCity().
		Order(ent.Desc(business.FieldCreatedAt))

	if req.Type != nil {
		info["业务类型"] = s.Text(req.Type)
		q.Where(business.TypeEQ(model.BusinessType(*req.Type)))
	}

	if req.Start != nil {
		info["开始时间"] = *req.Start
		q.Where(business.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
	}

	if req.End != nil {
		info["结束时间"] = *req.End
		q.Where(business.CreatedAtLT(tt.ParseNextDateStringX(*req.End)))
	}

	if req.EmployeeID != 0 {
		info["店员"] = ent.NewExportInfo(req.EmployeeID, employee.Table)
		q.Where(business.EmployeeID(req.EmployeeID))
	}

	if req.EnterpriseID != 0 {
		info["团签"] = ent.NewExportInfo(req.EnterpriseID, enterprise.Table)
		q.Where(business.EnterpriseID(req.EnterpriseID))
	}

	if req.Goal != model.StockGoalAll {
		info["查询目标"] = req.Goal.String()
		switch req.Goal {
		case model.StockGoalStore:
			q.Where(business.StoreIDNotNil())
		case model.StockGoalCabinet:
			q.Where(business.CabinetIDNotNil())
		}
	}

	if req.Keyword != nil {
		info["关键词"] = *req.Keyword
		q.Where(business.HasRiderWith(
			rider.Or(
				rider.PhoneContainsFold(*req.Keyword),
				rider.NameContainsFold(*req.Keyword),
			),
		))
	}

	if req.StoreID != 0 {
		info["门店"] = ent.NewExportInfo(req.StoreID, store.Table)
		q.Where(business.StoreID(req.StoreID))
	}

	if req.CabinetID != 0 {
		info["电柜"] = ent.NewExportInfo(req.CabinetID, cabinet.Table)
		q.Where(business.CabinetID(req.CabinetID))
	}

	switch req.Aimed {
	case model.BusinessAimedPersonal:
		info["业务对象"] = "个签"
		q.Where(business.EnterpriseIDIsNil())
	case model.BusinessAimedEnterprise:
		info["业务对象"] = "团签"
		q.Where(business.EnterpriseIDNotNil())
	}

	if req.CityID != 0 {
		info["城市"] = ent.NewExportInfo(req.CityID, city.Table)
		q.Where(business.CityID(req.CityID))
	}
	return
}

func (s *businessService) basicDetail(item *ent.Business) (res model.BusinessEmployeeListRes) {
	res = model.BusinessEmployeeListRes{
		ID:   item.ID,
		Type: s.Text(item.Type),
		Time: item.CreatedAt.Format(carbon.DateTimeLayout),
	}

	r := item.Edges.Rider
	if r != nil {
		res.Name = r.Name
		res.Phone = r.Phone
	}

	cit := item.Edges.City
	if cit != nil {
		res.City = cit.Name
	}

	p := item.Edges.Plan
	if p != nil {
		res.Plan = p.BasicInfo()
	}

	e := item.Edges.Enterprise
	if e != nil {
		res.Enterprise = &model.Enterprise{
			ID:    e.ID,
			Name:  e.Name,
			Agent: e.Agent,
		}
	}

	es := item.Edges.Station
	if es != nil {
		res.EnterpriseStation = &model.EnterpriseStation{
			ID:   es.ID,
			Name: es.Name,
		}
	}

	rtoBike := item.Edges.RtoEbike
	if rtoBike != nil {
		res.RtoEbikeSn = &rtoBike.Sn
	}

	if item.Remark != "" {
		res.Remark = &item.Remark
	}

	return
}

// ListEmployee 业务列表 - 门店
func (s *businessService) ListEmployee(req *model.BusinessListReq) *model.PaginationRes {
	req.EmployeeID = s.employee.ID
	q, _ := s.listFilter(req.BusinessFilter)

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.Business) (res model.BusinessEmployeeListRes) {
			return s.basicDetail(item)
		},
	)
}

// ListEnterprise 业务列表 - 团签
func (s *businessService) ListEnterprise(req *model.BusinessListReq) *model.PaginationRes {
	q, _ := s.listFilter(req.BusinessFilter)
	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.Business) (res model.BusinessListRes) {
			return s.detailInfo(item)
		},
	)
}

func (s *businessService) detailInfo(item *ent.Business) model.BusinessListRes {
	detail := model.BusinessListRes{
		BusinessEmployeeListRes: s.basicDetail(item),
		Employee:                nil,
	}
	emp := item.Edges.Employee
	if emp != nil {
		detail.Employee = &model.Employee{
			ID:    emp.ID,
			Name:  emp.Name,
			Phone: emp.Phone,
		}
	}

	st := item.Edges.Store
	if st != nil {
		detail.Store = &model.Store{
			ID:   st.ID,
			Name: st.Name,
		}
	}

	cab := item.Edges.Cabinet
	if cab != nil {
		detail.Cabinet = &model.CabinetBasicInfo{
			ID:     cab.ID,
			Brand:  cab.Brand,
			Serial: cab.Serial,
			Name:   cab.Name,
		}
	}
	// 操作人
	var operator string
	switch {
	case item.EmployeeID == nil && item.CabinetID == nil && item.AgentID == nil && item.Creator != nil:
		// 操作人是平台
		operator = "平台-" + item.Creator.Name
	case item.CabinetID != nil:
		// 操作人是骑手
		if item.Edges.Rider != nil {
			operator = "骑手-" + item.Edges.Rider.Name
		}
	case item.EmployeeID != nil:
		// 操作人是店员
		if item.Edges.Employee != nil {
			operator = "店员-" + item.Edges.Employee.Name
		}
	case item.AgentID != nil:
		// 操作人是代理
		if item.Edges.Agent != nil {
			operator = "代理-" + item.Edges.Agent.Name
		}
	}
	detail.Operator = operator
	return detail
}

// ListManager 业务列表 - 后台
func (s *businessService) ListManager(req *model.BusinessListReq) *model.PaginationRes {
	q, _ := s.listFilter(req.BusinessFilter)
	q.WithEmployee().WithCabinet().WithStore().WithAgent().WithPlan().WithRtoEbike()

	return model.ParsePaginationResponse(
		q,
		req.PaginationReq,
		func(item *ent.Business) (res model.BusinessListRes) {
			return s.detailInfo(item)
		},
	)
}

func (s *businessService) ListPause(req *model.BusinessPauseList) *model.PaginationRes {
	q := ent.Database.SubscribePause.
		QueryNotDeleted().
		WithCity().
		WithEmployee().
		WithEndEmployee().
		WithCabinet().
		WithEndCabinet().
		WithStore().
		WithEndStore().
		WithSubscribe(func(query *ent.SubscribeQuery) {
			query.WithPlan()
		}).
		WithRider().
		Order(ent.Desc(subscribepause.FieldCreatedAt))

	// 筛选城市
	if req.CityID != 0 {
		q.Where(subscribepause.CityID(req.CityID))
	}

	// 筛选骑手
	if req.RiderID != 0 {
		q.Where(subscribepause.RiderID(req.RiderID))
	}

	// 状态筛选
	switch req.Status {
	case 1:
		q.Where(subscribepause.EndAtIsNil())
	case 2:
		q.Where(subscribepause.EndAtNotNil())
	}

	// 是否逾期
	if req.Overdue {
		q.Where(subscribepause.OverdueDaysGT(0))
	}

	switch req.StartAscription {
	case 1:
		q.Where(subscribepause.StoreIDNotNil())
	case 2:
		q.Where(subscribepause.CabinetIDNotNil())
	}

	switch req.EndAscription {
	case 1:
		q.Where(subscribepause.EndStoreIDNotNil())
	case 2:
		q.Where(subscribepause.EndCabinetIDNotNil())
	}

	if req.StartDate != "" {
		start := strings.Split(strings.ReplaceAll(req.StartDate, " ", ""), ",")
		if len(start) != 2 {
			snag.Panic("寄存时间段参数错误")
		}
		q.Where(
			subscribepause.StartAtGTE(tools.NewTime().ParseDateStringX(start[0])),
			subscribepause.StartAtLT(tools.NewTime().ParseNextDateStringX(start[1])),
		)
	}

	if req.EndDate != "" {
		end := strings.Split(strings.ReplaceAll(req.EndDate, " ", ""), ",")
		if len(end) != 2 {
			snag.Panic("结束寄存时间段参数错误")
		}
		q.Where(
			subscribepause.EndAtGTE(tools.NewTime().ParseDateStringX(end[0])),
			subscribepause.EndAtLT(tools.NewTime().ParseNextDateStringX(end[1])),
		)
	}

	if req.StartBy != "" {
		q.Where(
			subscribepause.Or(
				func(sel *sql.Selector) {
					sel.Where(sqljson.StringContains(sel.C(subscribepause.FieldCreator), req.StartBy, sqljson.Path("name")))
				},
				subscribepause.HasEmployeeWith(employee.NameContainsFold(req.StartBy)),
			),
		)
	}

	if req.EndBy != "" {
		q.Where(
			subscribepause.Or(
				func(sel *sql.Selector) {
					sel.Where(sqljson.StringContains(sel.C(subscribepause.FieldEndModifier), req.EndBy, sqljson.Path("name")))
				},
				subscribepause.HasEndEmployeeWith(employee.NameContainsFold(req.EndBy)),
			),
		)
	}

	if req.StartTarget != "" {
		q.Where(
			subscribepause.Or(
				subscribepause.HasStoreWith(store.NameContainsFold(req.StartTarget)),
				subscribepause.HasCabinetWith(
					cabinet.Or(
						cabinet.NameContainsFold(req.StartTarget),
						cabinet.SerialContainsFold(req.StartTarget),
					),
				),
			),
		)
	}

	if req.EndTarget != "" {
		q.Where(
			subscribepause.Or(
				subscribepause.HasEndStoreWith(store.NameContainsFold(req.EndTarget)),
				subscribepause.HasEndCabinetWith(
					cabinet.Or(
						cabinet.NameContainsFold(req.EndTarget),
						cabinet.SerialContainsFold(req.EndTarget),
					),
				),
			),
		)
	}

	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.SubscribePause) (res model.BusinessPauseListRes) {
		return s.pauseDetail(item)
	})
}

func (s *businessService) pauseDetail(item *ent.SubscribePause) (res model.BusinessPauseListRes) {
	sub := item.Edges.Subscribe
	if sub == nil {
		return
	}
	ep := sub.Edges.Plan
	res = model.BusinessPauseListRes{
		City:            item.Edges.City.Name,
		Name:            item.Edges.Rider.Name,
		Phone:           item.Edges.Rider.Phone,
		Plan:            fmt.Sprintf("%s - %d天", ep.Name, ep.Days),
		Start:           item.StartAt.Format(carbon.DateTimeLayout),
		StartTarget:     s.pauseTarget(item.Edges.Store, item.Edges.Cabinet),
		StartAscription: s.pauseAscription(item.Edges.Store, item.Edges.Cabinet),
		StartBy:         s.pauseBy(item.Creator, item.Edges.Employee, item.Edges.Cabinet),
		EndTarget:       s.pauseTarget(item.Edges.EndStore, item.Edges.EndCabinet),
		EndAscription:   s.pauseAscription(item.Edges.EndStore, item.Edges.EndCabinet),
		EndBy:           s.pauseBy(item.EndModifier, item.Edges.EndEmployee, item.Edges.EndCabinet),
		Days:            item.Days,
		OverdueDays:     item.OverdueDays,
		Remaining:       sub.Remaining,
		SuspendDays:     item.SuspendDays,
	}

	if item.EndAt.IsZero() {
		res.Status = "寄存中"
	} else {
		res.End = item.EndAt.Format(carbon.DateTimeLayout)
		res.Status = "已结束"
	}

	if item.PauseOverdue {
		res.Status = "超期退租"
	}

	return
}

func (s *businessService) pauseTarget(st *ent.Store, cab *ent.Cabinet) string {
	if st != nil {
		return fmt.Sprintf("[门店] %s", st.Name)
	}
	if cab != nil {
		return fmt.Sprintf("[电柜] %s - %s", cab.Name, cab.Serial)
	}
	return ""
}

func (s *businessService) pauseAscription(st *ent.Store, cab *ent.Cabinet) string {
	if st != nil {
		return "门店"
	}
	if cab != nil {
		return "电柜"
	}
	return ""
}

func (s *businessService) pauseBy(m *model.Modifier, e *ent.Employee, cab *ent.Cabinet) string {
	if cab != nil {
		return cab.Serial
	}
	if e != nil {
		return e.Name
	}
	if m != nil {
		return m.Name
	}
	return ""
}

func (s *businessService) Export(req *model.BusinessExportReq) model.ExportRes {
	q, info := s.listFilter(req.BusinessFilter)
	q.WithEmployee().WithCabinet().WithStore()
	return NewExportWithModifier(s.modifier).Start("业务记录", req.BusinessFilter, info, req.Remark, func(path string) {
		items, _ := q.All(s.ctx)
		var rows tools.ExcelItems
		title := []any{
			"类型",
			"城市",
			"电柜",
			"门店",
			"店员",
			"骑手",
			"电话",
			"骑士卡",
			"骑士卡类型",
			"骑士卡天数",
			"团签",
			"时间",
		}
		rows = append(rows, title)
		for _, item := range items {
			detail := s.detailInfo(item)
			var cab, sto, emp, pla, en, palnType string
			var days uint
			if detail.Cabinet != nil {
				cab = detail.Cabinet.Serial
			}
			if detail.Store != nil {
				sto = detail.Store.Name
			}
			if detail.Plan != nil {
				pla = detail.Plan.Name
				days = detail.Plan.Days
				palnType = detail.Plan.Type.String()
			}
			if detail.Enterprise != nil {
				en = detail.Enterprise.Name
			}
			if detail.Employee != nil {
				emp = detail.Employee.Name
			}
			rows = append(rows, []any{
				detail.Type,
				detail.City,
				cab,
				sto,
				emp,
				detail.Name,
				detail.Phone,
				pla,
				palnType,
				days,
				en,
				detail.Time,
			})
		}
		tools.NewExcel(path).AddValues(rows).Done()
	})
}

func (s *businessService) Convert(bus model.BusinessType) (adapter.Business, error) {
	switch bus {
	default:
		return "", adapter.ErrorBusiness
	case model.BusinessTypeActive:
		return adapter.BusinessActive, nil
	case model.BusinessTypePause:
		return adapter.BusinessPause, nil
	case model.BusinessTypeContinue:
		return adapter.BusinessContinue, nil
	case model.BusinessTypeUnsubscribe:
		return adapter.BusinessUnsubscribe, nil
	}
}
