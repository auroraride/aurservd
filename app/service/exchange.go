// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "entgo.io/ent/dialect/sql"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/exchange"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "github.com/google/uuid"
    "time"
)

type exchangeService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *model.Employee
    orm      *ent.ExchangeClient
}

func NewExchange() *exchangeService {
    return &exchangeService{
        ctx: context.Background(),
        orm: ar.Ent.Exchange,
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

func NewExchangeWithEmployee(e *model.Employee) *exchangeService {
    s := NewExchange()
    s.ctx = context.WithValue(s.ctx, "employee", e)
    s.employee = e
    return s
}

// Store 扫门店二维码换电
func (s *exchangeService) Store(req *model.ExchangeStoreReq) *model.ExchangeStoreRes {
    item := NewStore().QuerySn(req.Code[6:])
    // 门店状态
    if item.Status != model.StoreStatusOpen {
        snag.Panic("门店未营业")
    }

    ee := item.Edges.Employee
    if ee == nil {
        snag.Panic("门店当前无工作人员")
    }

    // 获取套餐
    o := NewSubscribe().Recent(s.rider.ID)

    // TODO 判定门店物资是否匹配电压型号
    if o.Status != model.SubscribeStatusUsing {
        snag.Panic("骑士卡状态异常")
    }

    // 存储
    uid := uuid.New().String()
    s.orm.Create().
        SetEmployee(ee).
        SetRider(s.rider).
        SetSuccess(true).
        SetStore(item).
        SetCityID(o.City.ID).
        SetUUID(uid).
        SaveX(s.ctx)

    return &model.ExchangeStoreRes{
        Voltage:   o.Voltage,
        StoreName: item.Name,
        Time:      time.Now().Unix(),
        UUID:      uid,
    }
}

// Overview 换电概览
func (s *exchangeService) Overview(riderID uint64) (res model.ExchangeOverview) {
    var result []struct {
        Count int    `json:"count"`
        Date  string `json:"date"`
    }
    s.orm.QueryNotDeleted().Where(exchange.RiderID(riderID), exchange.Success(true)).Modify(func(s *sql.Selector) {
        s.Select(
            sql.As(sql.Count(exchange.FieldID), "count"),
        ).AppendSelectExprAs(sql.Raw(`created_at::DATE`), "date").GroupBy("date")
    }).ScanX(s.ctx, &result)

    res.Days = len(result)
    for _, r := range result {
        res.Times += r.Count
    }
    return
}

// Log 换电记录
func (s *exchangeService) Log(riderID uint64, req *model.PaginationReq) *model.PaginationRes {
    q := s.orm.QueryNotDeleted().Where(exchange.RiderID(riderID), exchange.Success(true)).WithStore().WithCity().WithCabinet()
    return model.ParsePaginationResponse[model.ExchangeLogRes, ent.Exchange](
        q,
        *req,
        func(item *ent.Exchange) (res model.ExchangeLogRes) {
            res = model.ExchangeLogRes{
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
