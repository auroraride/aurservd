// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-04
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ar"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/google/uuid"
    "time"
)

type exchangeService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    employee *model.Employee
}

func NewExchange() *exchangeService {
    return &exchangeService{
        ctx: context.Background(),
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
    ar.Ent.Exchange.Create().
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
