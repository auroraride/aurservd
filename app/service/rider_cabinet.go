// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-29
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/pkg/snag"
)

type riderCabinetService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
}

func NewRiderCabinet() *riderCabinetService {
    return &riderCabinetService{
        ctx: context.Background(),
    }
}

func NewRiderCabinetWithRider(rider *ent.Rider) *riderCabinetService {
    s := NewRiderCabinet()
    s.ctx = context.WithValue(s.ctx, "rider", rider)
    s.rider = rider
    return s
}

func NewRiderCabinetWithModifier(m *model.Modifier) *riderCabinetService {
    s := NewRiderCabinet()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// Process 获取待换电信息
func (s *riderCabinetService) Process(req *model.RiderCabinetOperateReq) model.RiderCabinetOperateProcess {
    // 是否有生效中套餐
    o := NewSubscribe().Recent(s.rider.ID)
    if o == nil || o.Status != model.SubscribeStatusUsing {
        snag.Panic("无生效中的骑行卡")
    }
    // 查询电柜
    cs := NewCabinet()
    cab := cs.QueryWithSerial(req.Serial)
    // 查询套餐
    if !cs.Health(cab) {
        snag.Panic("电柜目前不可用")
    }
    info := NewCabinet().Usable(cab)
    if info.EmptyBin == nil {
        snag.Panic("电柜目前不可用")
    }
    return info
}
