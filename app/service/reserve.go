// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-13
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "context"
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/business"
    "github.com/auroraride/aurservd/internal/ent/reserve"
    "github.com/auroraride/aurservd/pkg/cache"
    "time"
)

type reserveService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    orm      *ent.ReserveClient
}

func NewReserve() *reserveService {
    return &reserveService{
        ctx: context.Background(),
        orm: ent.Database.Reserve,
    }
}

func NewReserveWithRider(r *ent.Rider) *reserveService {
    s := NewReserve()
    s.ctx = context.WithValue(s.ctx, "rider", r)
    s.rider = r
    return s
}

func NewReserveWithModifier(m *model.Modifier) *reserveService {
    s := NewReserve()
    s.ctx = context.WithValue(s.ctx, "modifier", m)
    s.modifier = m
    return s
}

// CabinetCounts 获取电柜当前预约数量
// 条件: 执行中的预约或一个小时以内的未执行预约
func (s *reserveService) CabinetCounts(ids []uint64, typ business.Type) (data map[uint64]int) {
    max := time.Duration(cache.Int(model.SettingReserveDuration))
    var results []struct {
        CabinetID uint64 `json:"cabinet_id"`
        Count     int    `json:"count"`
    }
    _ = s.orm.QueryNotDeleted().
        Where(
            reserve.CabinetIDIn(ids...),
            reserve.Or(
                reserve.And(
                    reserve.Status(model.ReserveStatusPending.Value()),
                    reserve.CreatedAtGTE(time.Now().Add(-max*time.Minute)),
                    reserve.Type(typ.String()),
                ),
                reserve.Status(model.ReserveStatusProcessing.Value()),
            ),
        ).
        GroupBy(reserve.FieldCabinetID).
        Aggregate(ent.Count()).
        Scan(s.ctx, &results)
    data = make(map[uint64]int)
    for _, result := range results {
        data[result.CabinetID] = result.Count
    }
    return
}
