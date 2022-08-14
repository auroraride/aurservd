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
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "time"
)

type reserveService struct {
    ctx      context.Context
    modifier *model.Modifier
    rider    *ent.Rider
    orm      *ent.ReserveClient
    max      time.Duration
}

func NewReserve() *reserveService {
    max := time.Duration(cache.Int(model.SettingReserveDuration))
    return &reserveService{
        ctx: context.Background(),
        orm: ent.Database.Reserve,
        max: max,
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

func (s *reserveService) Query(id uint64) (*ent.Reserve, error) {
    return s.orm.QueryNotDeleted().Where(reserve.ID(id)).First(s.ctx)
}

func (s *reserveService) QueryX(id uint64) *ent.Reserve {
    rev, _ := s.Query(id)
    if rev == nil {
        snag.Panic("未找到预约")
    }
    return rev
}

// RiderUnfinished 查询用户未完成的预约
func (s *reserveService) RiderUnfinished(riderID uint64) *ent.Reserve {
    rev, _ := s.orm.QueryNotDeleted().
        Where(
            reserve.StatusIn(model.ReserveStatusPending.Value(), model.ReserveStatusProcessing.Value()),
            reserve.CreatedAtGTE(time.Now().Add(-s.max*time.Minute)),
            reserve.RiderID(riderID),
        ).
        First(s.ctx)
    return rev
}

func (s *reserveService) RiderUnfinishedDetail(riderID uint64) *model.RiderUnfinishedRes {
    rev := s.RiderUnfinished(riderID)
    if rev == nil {
        return nil
    }
    return s.Detail(rev)
}

func (s *reserveService) Detail(rev *ent.Reserve) *model.RiderUnfinishedRes {
    return &model.RiderUnfinishedRes{
        ID:        rev.ID,
        CabinetID: rev.CabinetID,
        Business:  rev.Type,
        Time:      rev.CreatedAt.Format(carbon.DateTimeLayout),
        Status:    model.ReserveStatus(rev.Status),
        Fid:       NewBranch().EncodeCabinetID(rev.CabinetID),
    }
}

// CabinetCounts 获取电柜当前预约数量
// 条件: 执行中的预约或一个小时以内的未执行预约
func (s *reserveService) CabinetCounts(ids []uint64, typ business.Type) (data map[uint64]int) {
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
                    reserve.CreatedAtGTE(time.Now().Add(-s.max*time.Minute)),
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

// Create 创建预约
func (s *reserveService) Create(req *model.ReserveCreateReq) *model.RiderUnfinishedRes {
    typ := business.Type(req.Business)

    // 检查订阅状态
    sub := NewSubscribeWithRider(s.rider).RecentX(s.rider.ID)

    // 检查骑手权限以及是否可办理业务
    NewRiderPermissionWithRider(s.rider).BusinessX().SubscribeX(model.RiderPermissionTypeBusiness, sub)

    // 骑士卡状态
    if !NewRiderBusiness(s.rider).Executable(sub, typ) {
        snag.Panic("骑士卡状态错误")
    }

    // 判断电柜是否可预约
    cab := NewCabinet().QueryOne(req.CabinetID)
    m := s.CabinetCounts([]uint64{cab.ID}, typ)
    if !cab.ReserveAble(typ, m[cab.ID]) {
        snag.Panic("电柜无法预约")
    }

    // 判断骑手是否可预约
    if s.RiderUnfinished(s.rider.ID) != nil {
        snag.Panic("当前已有其他预约")
    }

    // 创建预约
    rev, err := s.orm.Create().
        SetRiderID(s.rider.ID).
        SetCabinetID(cab.ID).
        SetCityID(*cab.CityID).
        SetType(req.Business).
        Save(s.ctx)
    if err != nil {
        snag.Panic("预约失败")
    }

    return s.Detail(rev)
}

// Cancel 取消预约
func (s *reserveService) Cancel(req *model.IDParamReq) {
    _, _ = s.QueryX(req.ID).Update().SetStatus(model.ReserveStatusCancel.Value()).Save(s.ctx)
}
