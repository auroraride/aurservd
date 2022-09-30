// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-08
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/assistance"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/commission"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/auroraride/aurservd/pkg/utils"
    "github.com/golang-module/carbon/v2"
    "github.com/google/uuid"
    "github.com/rs/xid"
    log "github.com/sirupsen/logrus"
    "time"
)

type employeeService struct {
    cacheKeyPrefix string
    ctx            context.Context
    modifier       *model.Modifier
    rider          *ent.Rider
    employee       *ent.Employee
    orm            *ent.EmployeeClient
    employeeInfo   *model.Employee
}

func NewEmployee() *employeeService {
    return &employeeService{
        cacheKeyPrefix: "EMPLOYEE_",
        ctx:            context.Background(),
        orm:            ent.Database.Employee,
    }
}

func NewEmployeeWithRider(r *ent.Rider) *employeeService {
    s := NewEmployee()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewEmployeeWithModifier(m *model.Modifier) *employeeService {
    s := NewEmployee()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewEmployeeWithEmployee(e *ent.Employee) *employeeService {
    s := NewEmployee()
    if e != nil {
        s.employee = e
        s.employeeInfo = &model.Employee{
            ID:    e.ID,
            Name:  e.Name,
            Phone: e.Phone,
        }
        s.ctx = context.WithValue(s.ctx, "employee", s.employeeInfo)
    }
    return s
}

func (s *employeeService) Query(id uint64) (*ent.Employee, error) {
    return s.orm.QueryNotDeleted().Where(employee.ID(id)).First(s.ctx)
}

func (s *employeeService) QueryX(id uint64) *ent.Employee {
    item, _ := s.Query(id)
    if item == nil {
        snag.Panic("未找到店员")
    }
    return item
}

func (s *employeeService) GetEmployeeByID(id uint64) (*ent.Employee, error) {
    return s.orm.QueryNotDeleted().Where(employee.ID(id)).WithStore().First(s.ctx)
}

// QueryByPhone 根据phone查找店员
func (s *employeeService) QueryByPhone(phone string) *ent.Employee {
    item, _ := s.orm.QueryNotDeleted().Where(employee.Phone(phone)).First(s.ctx)
    return item
}

// Create 添加店员
func (s *employeeService) Create(req *model.EmployeeCreateReq) *ent.Employee {
    // 判断重复
    em := s.QueryByPhone(req.Phone)
    if em != nil {
        snag.Panic("店员已存在")
    }
    var err error
    em, err = s.orm.Create().
        SetPhone(req.Phone).
        SetName(req.Name).
        SetCityID(req.CityID).
        Save(s.ctx)
    if em == nil {
        log.Error(err)
        snag.Panic("店员添加失败")
    }
    return em
}

// Modify 修改店员
func (s *employeeService) Modify(req *model.EmployeeModifyReq) {
    _, err := s.orm.ModifyOne(s.QueryX(*req.ID), req).Save(s.ctx)
    if err != nil {
        snag.Panic("保存失败")
        log.Error(err)
    }
}

func (s *employeeService) List(req *model.EmployeeListReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().
        WithCity()

    if req.Keyword != nil {
        q.Where(
            employee.Or(
                employee.NameContainsFold(*req.Keyword),
                employee.PhoneContainsFold(*req.Keyword),
            ),
        )
    }

    if req.CityID != nil {
        q.Where(employee.CityID(*req.CityID))
    }

    if req.Status != 0 {
        q.Where(employee.Enable(req.Status == 1))
    }

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Employee) model.EmployeeListRes {
        res := model.EmployeeListRes{
            ID:     item.ID,
            Name:   item.Name,
            Phone:  item.Phone,
            Enable: item.Enable,
        }

        ec := item.Edges.City
        if ec != nil {
            res.City = model.City{
                ID:   ec.ID,
                Name: ec.Name,
            }
        }
        return res
    })
}

func (s *employeeService) activityListFilter(req model.EmployeeActivityListFilter, export bool) (q *ent.EmployeeQuery, info ar.Map) {
    info = make(ar.Map)
    var start, end time.Time
    if req.Start != "" {
        info["开始日期"] = req.Start
        start = tools.NewTime().ParseDateStringX(req.Start)
    }
    if req.End != "" {
        info["截止日期"] = req.Start
        end = tools.NewTime().ParseNextDateStringX(req.End)
    }

    q = s.orm.QueryNotDeleted().
        WithStore().
        WithCity().
        WithExchanges(func(query *ent.ExchangeQuery) {
            if !start.IsZero() {
                query.Where(exchange.CreatedAtGTE(start))
            }
            if !end.IsZero() {
                query.Where(exchange.CreatedAtLT(end))
            }
        }).
        WithCommissions(func(query *ent.CommissionQuery) {
            if !start.IsZero() {
                query.Where(commission.CreatedAtGTE(start))
            }
            if !end.IsZero() {
                query.Where(commission.CreatedAtLT(end))
            }
            if export {
                query.WithPlan().WithOrder().WithRider()
            }
        }).
        WithAssistances(func(query *ent.AssistanceQuery) {
            query.Where(assistance.StatusIn(model.AssistanceStatusSuccess, model.AssistanceStatusUnpaid))

            if !start.IsZero() {
                query.Where(assistance.CreatedAtGTE(start))
            }
            if !end.IsZero() {
                query.Where(assistance.CreatedAtLT(end))
            }
            if export {
                query.WithRider().WithStore()
            }
        })

    if req.Keyword != "" {
        info["关键词"] = req.Keyword
        q.Where(
            employee.Or(
                employee.NameContainsFold(req.Keyword),
                employee.PhoneContainsFold(req.Keyword),
            ),
        )
    }

    if req.StoreID != 0 {
        info["门店"] = ent.NewExportInfo(req.StoreID, store.Table)
        q.Where(employee.HasStoreWith(store.ID(req.StoreID)))
    }

    if req.CityID != 0 {
        info["城市"] = ent.NewExportInfo(req.CityID, city.Table)
        q.Where(employee.CityID(req.CityID))
    }

    return
}

// Activity 店员动态
func (s *employeeService) Activity(req *model.EmployeeActivityListReq) *model.PaginationRes {
    q, _ := s.activityListFilter(req.EmployeeActivityListFilter, false)
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Employee) model.EmployeeActivityListRes {
        res := model.EmployeeActivityListRes{
            ID:            item.ID,
            Name:          item.Name,
            Phone:         item.Phone,
            ExchangeTimes: len(item.Edges.Exchanges),
        }

        ec := item.Edges.City
        if ec != nil {
            res.City = model.City{
                ID:   ec.ID,
                Name: ec.Name,
            }
        }

        es := item.Edges.Store
        if es != nil {
            res.Store = &model.Store{
                ID:   es.ID,
                Name: es.Name,
            }
        }

        for _, cm := range item.Edges.Commissions {
            res.Amount = tools.NewDecimal().Sum(cm.Amount, res.Amount)
        }

        for _, ass := range item.Edges.Assistances {
            res.AssistanceTimes += 1
            res.AssistanceMeters += ass.Distance
        }
        return res
    })
}

func (s *employeeService) ActivityExport(req *model.EmployeeActivityExportReq) model.ExportRes {
    q, info := s.activityListFilter(req.EmployeeActivityListFilter, true)
    return NewExportWithModifier(s.modifier).Start("店员业绩", req.EmployeeActivityListFilter, info, req.Remark, func(path string) {
        items, _ := q.All(s.ctx)
        title := []any{
            "城市",          // 0
            "店员",          // 1
            "换电次数",      // 2
            "救援次数",      // 3
            "救援里程 (米)", // 4
            "里程",          // 5
            "原因",          // 6
            "骑手",          // 8
            "位置",          // 7
            "时间",          // 9
            "业绩提成",      // 10
            "骑手",          // 11
            "骑士卡",        // 12
            "订单金额",      // 13
            "提成金额",      // 14
            "时间",          // 15
        }
        rows := tools.ExcelItems{title}
        for _, item := range items {
            row := []any{
                item.Edges.City.Name,
                fmt.Sprintf("%s - %s", item.Name, item.Phone),
                len(item.Edges.Exchanges),
                len(item.Edges.Assistances),
            }
            // 救援
            asstotal := 0.0
            var assistances tools.ExcelItems
            if len(item.Edges.Assistances) == 0 {
                assistances = [][]any{{"", "", "", "", ""}}
            }
            for _, a := range item.Edges.Assistances {
                asto := "-"
                if a.Edges.Store != nil {
                    asto = a.Edges.Store.Name
                }

                asstotal += a.Distance
                assistances = append(assistances, []any{
                    a.Distance,
                    a.Breakdown,
                    fmt.Sprintf("%s - %s", a.Edges.Rider.Name, a.Edges.Rider.Phone),
                    fmt.Sprintf("[%s] %s", asto, a.Address),
                    a.CreatedAt.Format(carbon.DateTimeLayout),
                })
            }
            row = append(row, []any{asstotal, assistances}...)

            // 业绩
            comtotal := 0.0
            var coms tools.ExcelItems
            if len(item.Edges.Commissions) == 0 {
                coms = [][]any{{"", "", "", "", ""}}
            }
            for _, c := range item.Edges.Commissions {
                comtotal = tools.NewDecimal().Sum(comtotal, c.Amount)
                coms = append(coms, []any{
                    fmt.Sprintf("%s - %s", c.Edges.Rider.Name, c.Edges.Rider.Phone),
                    fmt.Sprintf("%s - %d天", c.Edges.Plan.Name, c.Edges.Plan.Days),
                    fmt.Sprintf("%.2f", c.Edges.Order.Amount),
                    fmt.Sprintf("%.2f", c.Amount),
                    c.CreatedAt.Format(carbon.DateTimeLayout),
                })
            }

            row = append(row, []any{comtotal, coms}...)
            rows = append(rows, row)
        }
        tools.NewExcel(path).AddValues(rows).Done()
    })
}

// Delete 删除骑手
func (s *employeeService) Delete(req *model.EmployeeDeleteReq) {
    item := s.QueryX(req.ID)
    _, err := s.orm.SoftDeleteOne(item).Save(s.ctx)
    if err != nil {
        log.Error(err)
        snag.Panic("店员删除失败")
    }
}

func (s *employeeService) tokenKey(id uint64) string {
    return fmt.Sprintf("%s%d", s.cacheKeyPrefix, id)
}

// ExtendTokenTime 延长登录有效期
func (s *employeeService) ExtendTokenTime(id uint64, token string) {
    ctx := context.Background()
    cache.Set(ctx, s.tokenKey(id), token, 7*24*time.Hour)
    cache.Set(ctx, token, id, 7*24*time.Hour)
}

// Profile 骑手资料获取
func (s *employeeService) Profile(e *ent.Employee) model.EmployeeProfile {
    st := e.Edges.Store
    res := model.EmployeeProfile{
        ID:     e.ID,
        Qrcode: fmt.Sprintf("EMPLOYEE:%s", e.Sn),
        Phone:  e.Phone,
        Name:   e.Name,
    }
    if st != nil {
        res.Onduty = true
        res.Store = &model.StoreWithStatus{
            Status: st.Status,
            Store: model.Store{
                ID:   st.ID,
                Name: st.Name,
            },
        }
    }
    return res
}

// Signin 店员登录
func (s *employeeService) Signin(req *model.EmployeeSignReq) model.EmployeeProfile {
    // 校验短信
    debugPhones := ar.Config.App.Debug.Phone
    if !debugPhones[req.Phone] && !NewSms().VerifyCode(req.SmsId, req.SmsCode) {
        snag.Panic("短信验证码校验失败")
    }
    e, _ := s.orm.QueryNotDeleted().Where(employee.Phone(req.Phone)).WithStore().First(s.ctx)
    if e == nil {
        snag.Panic("未找到用户")
    }

    // 被禁用
    if !e.Enable {
        snag.Panic(snag.StatusForbidden, ar.BannedMessage)
    }

    // 生成token
    token := xid.New().String() + utils.RandTokenString()
    key := s.tokenKey(e.ID)

    // 删除旧的token
    if old := cache.Get(s.ctx, key).Val(); old != "" {
        cache.Del(s.ctx, key)
        cache.Del(s.ctx, old)
    }

    // 生成UUID
    sn := uuid.New()
    e.Update().SetSn(sn).SaveX(s.ctx)

    s.ExtendTokenTime(e.ID, token)

    e.Sn = sn
    res := s.Profile(e)
    res.Token = token

    return res
}

// RefreshQrcode 重新生成二维码
func (s *employeeService) RefreshQrcode() model.EmployeeQrcodeRes {
    e := s.employee.Update().SetSn(uuid.New()).SaveX(s.ctx)
    return model.EmployeeQrcodeRes{Qrcode: fmt.Sprintf("EMPLOYEE:%s", e.Sn)}
}

func (s *employeeService) Enable(req *model.EmployeeEnableReq) {
    e := s.QueryX(req.ID)
    e.Update().SetEnable(req.Enable).SaveX(s.ctx)
}

func (s *employeeService) Signout(e *ent.Employee) {
    ctx := context.Background()
    key := s.tokenKey(e.ID)
    token := cache.Get(ctx, key).Val()
    cache.Del(ctx, key)
    cache.Del(ctx, token)
}

func (s *employeeService) NameFromID(id uint64) string {
    p, _ := ent.Database.Employee.QueryNotDeleted().Where(employee.ID(id)).First(s.ctx)
    if p == nil {
        return "-"
    }
    return p.Name
}
