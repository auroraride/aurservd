// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "fmt"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/cabinet"
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/cache"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "github.com/lithammer/shortuuid/v4"
    "math"
    "strings"
    "time"
)

type exchangeService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *ent.Employee
    orm      *ent.ExchangeClient
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
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// RiderInterval 检查用户换电间隔
func (s *exchangeService) RiderInterval(riderID uint64) {
    // 检查用户换电间隔
    iv := cache.Int(model.SettingExchangeInterval)
    if exist, _ := ent.Database.Exchange.QueryNotDeleted().Where(
        exchange.RiderID(riderID),
        exchange.Success(true),
    ).First(s.ctx); exist != nil {
        m := int(math.Ceil(time.Now().Sub(exist.FinishAt).Minutes()))
        if iv-m > 0 {
            snag.Panic(fmt.Sprintf("换电过于频繁, %d分钟可再次换电", iv-m))
        }
    }
}

// Store 扫门店二维码换电
// 换电操作有出库和入库, 所以不记录
func (s *exchangeService) Store(req *model.ExchangeStoreReq) *model.ExchangeStoreRes {
    s.RiderInterval(s.rider.ID)

    qr := strings.ReplaceAll(req.Code, "STORE:", "")
    item := NewStore().QuerySn(qr)
    // 门店状态
    if item.Status != model.StoreStatusOpen {
        snag.Panic("门店未营业")
    }

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
        SaveX(s.ctx)

    message := sub.Model
    message = strings.ReplaceAll(message, "AH", "安")
    message = strings.ReplaceAll(message, "V", "伏")
    // 发送播报信息给店员
    NewSpeech().SendSpeech(ee.ID, fmt.Sprintf("%s扫码换电%s", s.rider.Edges.Person.Name, message))

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
            store := item.Edges.Store
            if store != nil {
                res.Type = "门店"
                res.Name = store.Name
            }

            return res
        },
    )
}

// listBasicQuery 列表基础查询语句
func (s *exchangeService) listBasicQuery(req model.ExchangeListReq) *ent.ExchangeQuery {
    tt := tools.NewTime()

    q := ent.Database.Exchange.
        QueryNotDeleted().
        WithRider(func(rq *ent.RiderQuery) {
            rq.WithPerson()
        }).
        WithEnterprise()

    if req.Start != nil {
        q.Where(exchange.CreatedAtGTE(tt.ParseDateStringX(*req.Start)))
    }

    if req.End != nil {
        q.Where(exchange.CreatedAtLTE(tt.ParseNextDateStringX(*req.End)))
    }

    if req.Keyword != nil {
        q.Where(
            exchange.HasRiderWith(
                rider.Or(
                    rider.PhoneContainsFold(*req.Keyword),
                    rider.HasPersonWith(person.NameContainsFold(*req.Keyword)),
                ),
            ),
        )
    }

    switch req.Aimed {
    case model.BusinessAimedPersonal:
        q.Where(exchange.EnterpriseIDIsNil())
        break
    case model.BusinessAimedEnterprise:
        q.Where(exchange.EnterpriseIDNotNil())
        break
    }

    return q
}

func (s *exchangeService) EmployeeList(req *model.ExchangeListReq) *model.PaginationRes {
    q := s.listBasicQuery(*req).
        WithSubscribe(func(sq *ent.SubscribeQuery) {
            sq.WithPlan()
        }).
        Where(exchange.EmployeeID(s.employee.ID))

    return model.ParsePaginationResponse(
        q,
        req.PaginationReq,
        func(item *ent.Exchange) (res model.ExchangeEmployeeListRes) {
            res = model.ExchangeEmployeeListRes{
                ID:    item.ID,
                Name:  item.Edges.Rider.Edges.Person.Name,
                Phone: item.Edges.Rider.Phone,
                Time:  item.CreatedAt.Format(carbon.DateTimeLayout),
                Model: item.Model,
            }
            sub := item.Edges.Subscribe
            if sub != nil {
                p := sub.Edges.Plan
                if p != nil {
                    res.Plan = &model.Plan{
                        ID:   p.ID,
                        Name: p.Name,
                        Days: p.Days,
                    }
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
        },
    )
}

func (s *exchangeService) List(req *model.ExchangeManagerListReq) *model.PaginationRes {
    q := s.listBasicQuery(req.ExchangeListReq).
        WithCity().
        WithStore().
        WithCabinet().
        Order(ent.Desc(exchange.FieldCreatedAt))

    switch req.Target {
    case 1:
        q.Where(exchange.CabinetIDNotNil())
        break
    case 2:
        q.Where(exchange.StoreIDNotNil())
        break
    }

    if s.modifier != nil && s.modifier.Phone == "15537112255" {
        req.CityID = 410100
    }

    if req.CityID != 0 {
        q.Where(exchange.CityID(req.CityID))
    }

    if req.Employee != "" {
        q.Where(
            exchange.HasEmployeeWith(
                employee.Or(
                    employee.NameContainsFold(req.Employee),
                    employee.PhoneContainsFold(req.Employee),
                ),
            ),
        )
    }

    if req.Status != 0 {
        q.Where(
            exchange.Success(req.Status == 1),
            exchange.FinishAtNotNil(),
        )
    }

    if req.Serial != "" {
        q.Where(exchange.HasCabinetWith(cabinet.Serial(req.Serial)))
    }

    if req.Brand != "" {
        q.Where(exchange.HasCabinetWith(cabinet.Brand(req.Brand)))
    }

    // 是否备用方案 1是 2否
    if req.Alternative != 0 {
        q.Where(exchange.Alternative(req.Alternative == 2))
    }

    if req.CabinetID != 0 {
        q.Where(exchange.CabinetID(req.CabinetID))
    }

    if req.StoreID != 0 {
        q.Where(exchange.StoreID(req.StoreID))
    }

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Exchange) model.ExchangeManagerListRes {
        res := model.ExchangeManagerListRes{
            ID:          item.ID,
            Name:        item.Edges.Rider.Edges.Person.Name,
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
            res.Enterprise = &model.EnterpriseBasic{
                ID:   e.ID,
                Name: e.Name,
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
                res.Error = fmt.Sprintf("%s [%s]", item.Info.Message, ex.CurrentStep().Step)
            }
        }
        return res
    })
}
