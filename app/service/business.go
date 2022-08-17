// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "entgo.io/ent/dialect/sql"
    "entgo.io/ent/dialect/sql/sqljson"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/business"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/internal/ent/subscribepause"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "strings"
)

// 门店业务处理专用
type businessService struct {
    ctx      context.Context
    employee *ent.Employee
    modifer  *model.Modifier
}

func NewBusiness() *businessService {
    return &businessService{
        ctx: context.Background(),
    }
}

func NewBusinessWithEmployee(e *ent.Employee) *businessService {
    s := NewBusiness()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

func NewBusinessWithModifier(m *model.Modifier) *businessService {
    s := NewBusiness()
    if m != nil {
        s.ctx = context.WithValue(s.ctx, "modifier", m)
        s.modifer = m
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

    ic := r.Edges.Person.IDCardNumber
    res = model.SubscribeBusiness{
        ID:           r.ID,
        Status:       subd.Status,
        Name:         r.Edges.Person.Name,
        Phone:        r.Phone,
        IDCardNumber: ic[len(ic)-4:],
        Model:        subd.Model,
        SubscribeID:  subd.ID,
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

// listBasicQuery 列表基础查询语句
func (s *businessService) listBasicQuery(req *model.BusinessListReq) *ent.BusinessQuery {
    tt := tools.NewTime()

    q := ent.Database.Business.
        QueryNotDeleted().
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).
        WithEnterprise().
        WithPlan().
        Order(ent.Desc(business.FieldCreatedAt))

    if req.Type != nil {
        q.Where(business.TypeEQ(business.Type(*req.Type)))
    }

    if req.Start != nil {
        q.Where(business.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
    }

    if req.End != nil {
        q.Where(business.CreatedAtLT(tt.ParseNextDateStringX(*req.End)))
    }

    if req.EmployeeID != 0 {
        q.Where(business.EmployeeID(req.EmployeeID))
    }

    if req.EnterpriseID != 0 {
        q.Where(business.EnterpriseID(req.EnterpriseID))
    }

    if req.Keyword != nil {
        q.Where(business.HasRiderWith(
            rider.Or(
                rider.PhoneContainsFold(*req.Keyword),
                rider.HasPersonWith(person.NameContainsFold(*req.Keyword)),
            ),
        ))
    }

    switch req.Aimed {
    case model.BusinessAimedPersonal:
        q.Where(business.EnterpriseIDIsNil())
        break
    case model.BusinessAimedEnterprise:
        q.Where(business.EnterpriseIDNotNil())
        break
    }

    return q
}

func (s *businessService) basicDetail(item *ent.Business) (res model.BusinessEmployeeListRes) {
    res = model.BusinessEmployeeListRes{
        ID:    item.ID,
        Name:  item.Edges.Rider.Edges.Person.Name,
        Phone: item.Edges.Rider.Phone,
        Type:  item.Type.String(),
        Time:  item.CreatedAt.Format(carbon.DateTimeLayout),
    }
    p := item.Edges.Plan
    if p != nil {
        res.Plan = &model.Plan{
            ID:   p.ID,
            Name: p.Name,
            Days: p.Days,
        }
    }

    e := item.Edges.Enterprise
    if e != nil {
        res.Enterprise = &model.EnterpriseBasic{
            ID:   e.ID,
            Name: e.Name,
        }
    }

    return
}

// ListEmployee 业务列表 - 门店
func (s *businessService) ListEmployee(req *model.BusinessListReq) *model.PaginationRes {
    req.EmployeeID = s.employee.ID
    q := s.listBasicQuery(req)

    return model.ParsePaginationResponse(
        q,
        req.PaginationReq,
        func(item *ent.Business) (res model.BusinessEmployeeListRes) {
            return s.basicDetail(item)
        },
    )
}

// ListManager 业务列表 - 后台
func (s *businessService) ListManager(req *model.BusinessListReq) *model.PaginationRes {
    q := s.listBasicQuery(req).WithEmployee().WithCabinet().WithStore()

    return model.ParsePaginationResponse(
        q,
        req.PaginationReq,
        func(item *ent.Business) (res model.BusinessListRes) {
            res = model.BusinessListRes{
                BusinessEmployeeListRes: s.basicDetail(item),
                Employee:                nil,
            }
            emp := item.Edges.Employee
            if emp != nil {
                res.Employee = &model.Employee{
                    ID:    emp.ID,
                    Name:  emp.Name,
                    Phone: emp.Phone,
                }
            }

            st := item.Edges.Store
            if st != nil {
                res.Store = &model.Store{
                    ID:   st.ID,
                    Name: st.Name,
                }
            }

            cab := item.Edges.Cabinet
            if cab != nil {
                res.Cabinet = &model.CabinetBasicInfo{
                    ID:     cab.ID,
                    Brand:  model.CabinetBrand(cab.Brand),
                    Serial: cab.Serial,
                    Name:   cab.Name,
                }
            }
            return
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
        WithRider(func(query *ent.RiderQuery) {
            query.WithPerson()
        }).
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
        break
    case 2:
        q.Where(subscribepause.EndAtNotNil())
        break
    }

    // 是否逾期
    if req.Overdue {
        q.Where(subscribepause.OverdueDaysGT(0))
    }

    switch req.StartAscription {
    case 1:
        q.Where(subscribepause.StoreIDNotNil())
        break
    case 2:
        q.Where(subscribepause.CabinetIDNotNil())
        break
    }

    switch req.EndAscription {
    case 1:
        q.Where(subscribepause.EndStoreIDNotNil())
        break
    case 2:
        q.Where(subscribepause.EndCabinetIDNotNil())
        break
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
    ep := sub.Edges.Plan
    res = model.BusinessPauseListRes{
        City:            item.Edges.City.Name,
        Name:            item.Edges.Rider.Edges.Person.Name,
        Phone:           item.Edges.Rider.Phone,
        Plan:            fmt.Sprintf("%s - %d天", ep.Name, ep.Days),
        Start:           item.StartAt.Format(carbon.DateLayout),
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
        res.End = item.EndAt.Format(carbon.DateLayout)
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
