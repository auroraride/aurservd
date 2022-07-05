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
    "github.com/auroraride/aurservd/internal/ent/employee"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    "github.com/auroraride/aurservd/internal/ent/person"
    "github.com/auroraride/aurservd/internal/ent/rider"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
    "github.com/golang-module/carbon/v2"
    "github.com/lithammer/shortuuid/v4"
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

// Store 扫门店二维码换电
// 换电操作有出库和入库, 所以不记录
func (s *exchangeService) Store(req *model.ExchangeStoreReq) *model.ExchangeStoreRes {
    // 检查用户是否可以办理业务
    NewRiderPermissionWithRider(s.rider).BusinessX()

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
    subd, sub := NewSubscribe().RecentDetail(s.rider.ID)

    if subd == nil {
        snag.Panic("未找到有效订阅")
    }

    if subd.Status != model.SubscribeStatusUsing {
        snag.Panic("骑士卡状态异常")
    }

    // 存储
    uid := shortuuid.New()
    s.orm.Create().
        SetEmployee(ee).
        SetRider(s.rider).
        SetSuccess(true).
        SetStoreID(item.ID).
        SetCityID(subd.City.ID).
        SetUUID(uid).
        SetModel(subd.Model).
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
        Model:     subd.Model,
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
                Time:    item.CreatedAt.Format(carbon.DateTimeLayout),
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
                info := item.Detail.Info
                if info != nil {
                    res.BinInfo = model.ExchangeLogBinInfo{
                        EmptyIndex: info.EmptyIndex,
                        FullIndex:  info.FullIndex,
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
        q.Where(exchange.CreatedAtLTE(tt.ParseDateStringX(*req.End)))
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
        WithCabinet()

    switch req.Target {
    case 1:
        q.Where(exchange.CabinetIDNotNil())
        break
    case 2:
        q.Where(exchange.StoreIDNotNil())
        break
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

    tt := tools.NewTime()
    if req.Start != "" {
        q.Where(exchange.CreatedAtGTE(tt.ParseDateStringX(req.Start)))
    }
    if req.End != "" {
        q.Where(exchange.CreatedAtLT(tt.ParseNextDateStringX(req.End)))
    }

    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Exchange) model.ExchangeManagerListRes {
        res := model.ExchangeManagerListRes{
            ID:    item.ID,
            Name:  item.Edges.Rider.Edges.Person.Name,
            Phone: item.Edges.Rider.Phone,
            Time:  item.CreatedAt.Format(carbon.DateTimeLayout),
            Model: item.Model,
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
        return res
    })
}
