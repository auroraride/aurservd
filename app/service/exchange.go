// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "entgo.io/ent/dialect/sql"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/internal/ent/city"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/store"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "github.com/lithammer/shortuuid/v4"
    "golang.org/x/exp/slices"
    "math"
    "strconv"
    "strings"
    "time"
)

type exchangeService struct {
    ctx          context.Context
    modifier     *model.Modifier
    rider        *ent.Rider
    employee     *ent.Employee
    orm          *ent.ExchangeClient
    employeeInfo *model.Employee
}

func NewExchange() *exchangeService {
    return &exchangeService{
        ctx: context.Background(),
        orm: ent.Database.Exchange,
    }
}

func NewExchangeWithRider(r *ent.Rider) *exchangeService {
    s := NewExchange()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewExchangeWithModifier(m *model.Modifier) *exchangeService {
    s := NewExchange()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

func NewExchangeWithEmployee(e *ent.Employee) *exchangeService {
    s := NewExchange()
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

// RiderLimit 检查单独配置的骑手换电间隔
func (s *exchangeService) RiderLimit(r *ent.Rider, list []model.ExchangeLimit) (n int) {
    if len(list) == 0 {
        return
    }
    slices.SortFunc(list, func(a, b model.ExchangeLimit) bool {
        return a.Hours < b.Hours
    })
    // 查询骑手最大时间段内换电情况
    now := time.Now()
    max := list[len(list)-1].Hours
    items, _ := r.QueryExchanges().
        Where(
            exchange.FinishAtGTE(now.Add(-time.Duration(max)*time.Hour)),
            exchange.Success(true),
        ).
        Select(exchange.FieldFinishAt).
        Order(ent.Desc(exchange.FieldFinishAt)).
        All(s.ctx)

    // 按小时封装限制数量
    lm := make(map[int]int)

    index := 1
    for _, l := range list {
        for ; index <= l.Hours; index++ {
            lm[index] = l.Times
        }
    }

    tm := make(map[int]int)
    for _, item := range items {
        h := int(math.Ceil(now.Sub(item.FinishAt).Hours()))
        tm[h] += 1
        if t, ok := lm[h]; ok && t <= tm[h] {
            return 60 - now.Minute()
        }
    }

    return
}

// RiderInterval 检查用户换电间隔
func (s *exchangeService) RiderInterval(r *ent.Rider, cityID uint64) {
    var list []model.ExchangeLimit

    // 获取骑手换电间隔配置
    if len(r.ExchangeLimit) > 0 {
        list = r.ExchangeLimit
    } else {
        // 获取城市换电配置
        data := make(model.SettingExchangeLimits)
        _ = cache.Get(s.ctx, model.SettingExchangeLimitKey).Scan(&data)
        if cm, ok := data[cityID]; ok {
            list = cm
        }
    }

    n := s.RiderLimit(r, list)
    if n > 0 {
        snag.Panic("换电过于频繁, " + strconv.Itoa(n) + "分钟可再次换电")
    }

    if len(list) == 0 {

        // 检查用户换电间隔
        iv := cache.Int(model.SettingExchangeIntervalKey)
        if exist, _ := ent.Database.Exchange.QueryNotDeleted().Where(
            exchange.RiderID(r.ID),
            exchange.Success(true),
        ).Order(ent.Desc(exchange.FieldCreatedAt)).First(s.ctx); exist != nil {
            m := int(math.Ceil(time.Now().Sub(exist.FinishAt).Minutes()))
            n = iv - m
            if n > 0 {
                snag.Panic("换电过于频繁, " + strconv.Itoa(n) + "分钟可再次换电")
            }
        }
    }
}

// Store 扫门店二维码换电
// 换电操作有出库和入库, 所以不记录
func (s *exchangeService) Store(req *model.ExchangeStoreReq) *model.ExchangeStoreRes {
    qr := strings.ReplaceAll(req.Code, "STORE:", "")
    item := NewStore().QuerySn(qr)
    // 门店状态
    if item.Status != model.StoreStatusOpen {
        snag.Panic("门店未营业")
    }

    s.RiderInterval(s.rider, item.CityID)

    ee := item.Edges.Employee
    if ee == nil {
        snag.Panic("门店当前无工作人员")
    }

    // 获取套餐
    sub := NewSubscribe().RecentX(s.rider.ID)

    // 检查用户是否可以办理业务
    NewRiderPermissionWithRider(s.rider).BusinessX().SubscribeX(model.RiderPermissionTypeExchange, sub)

    // 存储
    uid := shortuuid.New()
    s.orm.Create().
        SetEmployee(ee).
        SetRider(s.rider).
        SetSuccess(true).
        SetStoreID(item.ID).
        SetCityID(sub.CityID).
        SetUUID(uid).
        SetModel(sub.Model).
        SetNillableEnterpriseID(sub.EnterpriseID).
        SetNillableStationID(sub.StationID).
        SetSubscribeID(sub.ID).
        SetFinishAt(time.Now()).
        SetDuration(0).
        SaveX(s.ctx)

    message := sub.Model
    message = strings.ReplaceAll(message, "AH", "安")
    message = strings.ReplaceAll(message, "V", "伏")
    // 发送播报信息给店员
    NewSpeech().SendSpeech(ee.ID, fmt.Sprintf("%s扫码换电%s", s.rider.Name, message))

    return &model.ExchangeStoreRes{
        Model:     sub.Model,
        StoreName: item.Name,
        Time:      time.Now().Unix(),
        UUID:      uid,
    }
}

// Overview 换电概览
func (s *exchangeService) Overview(riderID uint64) (res model.ExchangeOverview) {
    res.Times, _ = s.orm.QueryNotDeleted().Where(exchange.RiderID(riderID), exchange.Success(true)).Count(s.ctx)
    // 总使用天数
    items, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.RiderID(riderID)).All(s.ctx)
    for _, item := range items {
        switch item.Status {
        case model.SubscribeStatusInactive:
            break
        default:
            res.Days += item.InitialDays + item.AlterDays + item.OverdueDays + item.RenewalDays + item.PauseDays - item.Remaining + 1 // 已用天数(+1代表当前天数算作1天)
            break
        }
    }
    return
}

// RiderList 换电记录
func (s *exchangeService) RiderList(riderID uint64, req model.PaginationReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().
        Where(exchange.RiderID(riderID), exchange.Success(true)).
        WithStore().
        WithCity().
        WithCabinet().
        Order(ent.Desc(exchange.FieldCreatedAt))
    return model.ParsePaginationResponse[model.ExchangeRiderListRes, ent.Exchange](
        q,
        req,
        func(item *ent.Exchange) (res model.ExchangeRiderListRes) {
            res = model.ExchangeRiderListRes{
                ID:      item.ID,
                Time:    item.StartAt.Format(carbon.DateTimeLayout),
                Success: item.Success,
                City: model.City{
                    ID:   item.Edges.City.ID,
                    Name: item.Edges.City.Name,
                },
            }
            cab := item.Edges.Cabinet
            if cab != nil {
                res.Type = "电柜"
                res.Name = cab.Name
                if item.Info != nil && item.Info.Exchange != nil {
                    ex := item.Info.Exchange
                    res.BinInfo = model.ExchangeLogBinInfo{
                        EmptyIndex: ex.Empty.Index,
                        FullIndex:  ex.Fully.Index,
                    }
                }
            }
            st := item.Edges.Store
            if st != nil {
                res.Type = "门店"
                res.Name = st.Name
            }

            return res
        },
    )
}

// listBasicQuery 列表基础查询语句
func (s *exchangeService) listBasicQuery(req model.ExchangeListBasicFilter) (q *ent.ExchangeQuery, info ar.Map) {
    tt := tools.NewTime()

    info = make(ar.Map)

    q = ent.Database.Exchange.
        QueryNotDeleted().
        WithRider().
        WithEnterprise().
        Order(ent.Desc(exchange.FieldCreatedAt))

    if req.Start != nil {
        info["开始时间"] = *req.Start
        q.Where(exchange.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
    }

    if req.End != nil {
        info["结束时间"] = *req.End
        q.Where(exchange.CreatedAtLTE(tt.ParseNextDateStringX(*req.End)))
    }

    if req.Keyword != nil {
        info["关键词"] = *req.Keyword
        q.Where(
            exchange.HasRiderWith(
                rider.Or(
                    rider.PhoneContainsFold(*req.Keyword),
                    rider.NameContainsFold(*req.Keyword),
                ),
            ),
        )
    }

    info["对象"] = []string{"全部", "个签", "团签"}[req.Aimed]
    switch req.Aimed {
    case model.BusinessAimedPersonal:
        q.Where(exchange.EnterpriseIDIsNil())
        break
    case model.BusinessAimedEnterprise:
        q.Where(exchange.EnterpriseIDNotNil())
        break
    }

    return
}

func (s *exchangeService) EmployeeList(req *model.ExchangeEmployeeListReq) *model.PaginationRes {
    q, _ := s.listBasicQuery(req.ExchangeListBasicFilter)
    q.WithSubscribe(func(sq *ent.SubscribeQuery) {
        sq.WithPlan()
    }).
        Where(exchange.EmployeeID(s.employee.ID))

    return model.ParsePaginationResponse(
        q,
        req.PaginationReq,
        func(item *ent.Exchange) (res model.ExchangeEmployeeListRes) {
            res = model.ExchangeEmployeeListRes{
                ID:    item.ID,
                Name:  item.Edges.Rider.Name,
                Phone: item.Edges.Rider.Phone,
                Time:  item.CreatedAt.Format(carbon.DateTimeLayout),
                Model: item.Model,
            }
            sub := item.Edges.Subscribe
            if sub != nil {
                p := sub.Edges.Plan
                if p != nil {
                    res.Plan = p.BasicInfo()
                }
            }

            e := item.Edges.Enterprise
            if e != nil {
                res.Enterprise = &model.Enterprise{
                    ID:    e.ID,
                    Name:  e.Name,
                    Agent: e.Agent,
                }
            }

            return
        },
    )
}

func (s *exchangeService) listFilter(req model.ExchangeListFilter) (q *ent.ExchangeQuery, info ar.Map) {
    if s.modifier != nil && s.modifier.Phone == "15537112255" {
        req.CityID = 410100
    }

    q, info = s.listBasicQuery(req.ExchangeListBasicFilter)
    q.WithCity().
        WithStore().
        WithCabinet()

    switch req.Target {
    case 1:
        q.Where(exchange.CabinetIDNotNil())
        break
    case 2:
        q.Where(exchange.StoreIDNotNil())
        break
    }

    info["换电类别"] = []string{"全部", "电柜", "门店"}[req.Target]

    if req.CityID != 0 {
        q.Where(exchange.CityID(req.CityID))
        info["城市"] = ent.NewExportInfo(req.CityID, city.Table)
    }

    if req.Employee != "" {
        info["店员"] = req.Employee
        q.Where(
            exchange.HasEmployeeWith(
                employee.Or(
                    employee.NameContainsFold(req.Employee),
                    employee.PhoneContainsFold(req.Employee),
                ),
            ),
        )
    }

    if req.Status != nil {
        info["状态"] = []string{"进行中", "成功", "失败"}[*req.Status]
        q.Where(
            exchange.Success(*req.Status == 1),
        )
        if *req.Status != 0 {
            q.Where(
                exchange.FinishAtNotNil(),
            )
        }
    }

    if req.Serial != "" {
        info["电柜编码"] = req.Serial
        q.Where(exchange.HasCabinetWith(cabinet.Serial(req.Serial)))
    }

    if req.Brand != "" {
        info["电柜品牌"] = req.Brand
        q.Where(exchange.HasCabinetWith(cabinet.Brand(req.Brand)))
    }

    // 是否备用方案 1是 2否
    info["备选方案"] = []string{"全部", "满电", "非满电"}[req.Alternative]
    if req.Alternative != 0 {
        q.Where(exchange.Alternative(req.Alternative == 2))
    }

    if req.CabinetID != 0 {
        q.Where(exchange.CabinetID(req.CabinetID))
        info["电柜"] = ent.NewExportInfo(req.CabinetID, cabinet.Table)
    }

    if req.StoreID != 0 {
        q.Where(exchange.StoreID(req.StoreID))
        info["门店"] = ent.NewExportInfo(req.StoreID, store.Table)
    }

    if req.Model != "" {
        info["电池型号"] = req.Model
        q.Where(exchange.Model(req.Model))
    }

    if req.Times > 0 {
        info["换电次数"] = req.Times
        q.Where(func(sel *sql.Selector) {
            sel.Where(
                sql.In(
                    exchange.FieldRiderID,
                    sql.Select(exchange.FieldRiderID).
                        From(sql.Table(exchange.Table)).
                        GroupBy(exchange.FieldRiderID).
                        Having(sql.GTE("COUNT(1)", req.Times)),
                ),
            )
        })
    }

    return
}

func (s *exchangeService) listDetail(item *ent.Exchange) (res model.ExchangeManagerListRes) {
    if item.Edges.Rider == nil {
        return
    }
    res = model.ExchangeManagerListRes{
        ID:          item.ID,
        Name:        item.Edges.Rider.Name,
        Phone:       item.Edges.Rider.Phone,
        Time:        item.CreatedAt.Format(carbon.DateTimeLayout),
        Model:       item.Model,
        Alternative: item.Alternative,
    }

    if item.FinishAt.IsZero() && item.CabinetID != 0 {
        res.Status = 0
    } else {
        if item.Success {
            res.Status = 1
        } else {
            res.Status = 2
        }
    }

    e := item.Edges.Enterprise
    if e != nil {
        res.Enterprise = &model.Enterprise{
            ID:    e.ID,
            Name:  e.Name,
            Agent: e.Agent,
        }
    }

    es := item.Edges.Store
    if es != nil {
        res.Store = &model.Store{
            ID:   es.ID,
            Name: es.Name,
        }
    }

    ec := item.Edges.City
    if ec != nil {
        res.City = model.City{ID: ec.ID, Name: ec.Name}
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

    if item.Info != nil && item.Info.Exchange != nil {
        ex := item.Info.Exchange
        res.Full = fmt.Sprintf("%d号仓, %.2f%%", ex.Fully.Index+1, ex.Fully.Electricity)
        res.Empty = fmt.Sprintf("%d号仓, %.2f%%", ex.Empty.Index+1, ex.Empty.Electricity)
        if !item.Success && !item.FinishAt.IsZero() {
            if len(ex.Steps) > 0 {
                res.Error = fmt.Sprintf("%s [%s]", item.Info.Message, ex.CurrentStep().Step)
            } else {
                res.Error = "未找到换电信息"
            }
        }
    }

    if item.Remark != "" {
        res.Error = item.Remark
    }
    return res
}

func (s *exchangeService) List(req *model.ExchangeManagerListReq) *model.PaginationRes {
    q, _ := s.listFilter(req.ExchangeListFilter)

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Exchange) model.ExchangeManagerListRes {
        return s.listDetail(item)
    })
}

func (s *exchangeService) Export(req *model.ExchangeListExport) model.ExportRes {
    q, info := s.listFilter(req.ExchangeListFilter)
    return NewExportWithModifier(s.modifier).Start("换电明细", req.ExchangeListFilter, info, req.Remark, func(path string) {
        items, _ := q.All(s.ctx)
        var rows tools.ExcelItems
        title := []any{
            "城市",
            "状态",
            "姓名",
            "电话",
            "型号",
            "团签",
            "门店",
            "电柜",
            "方案",
            "满仓",
            "空仓",
            "失败原因",
            "开始时间",
            "结束时间",
            "耗时 (秒)",
        }
        rows = append(rows, title)
        for _, item := range items {
            detail := s.listDetail(item)
            row := []any{
                detail.City.Name,
                []string{"进行中", "成功", "失败"}[detail.Status],
                detail.Name,
                detail.Phone,
                detail.Model,
                "", // 团签
                "", // 门店
                "", // 电柜
                "满电",
                detail.Full,
                detail.Empty,
                detail.Error,
                detail.Time,
                "", // 结束时间
                "", // 耗时
            }
            if detail.Enterprise != nil {
                row[5] = detail.Enterprise.Name
            }
            if detail.Store != nil {
                row[6] = detail.Store.Name
            }
            if detail.Cabinet != nil {
                row[7] = fmt.Sprintf("[%s]%s - %s", detail.Cabinet.Brand, detail.Cabinet.Name, detail.Cabinet.Serial)
            }
            if !detail.Alternative {
                row[8] = "非满电"
            }
            if !item.FinishAt.IsZero() {
                row[13] = item.FinishAt.Format(carbon.DateTimeLayout)
                row[14] = item.Duration
            }
            rows = append(rows, row)
        }
    })
}
