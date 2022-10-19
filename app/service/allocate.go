// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/app/socket"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/allocate"
    "github.com/auroraride/aurservd/internal/ent/ebike"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "time"
)

type allocateService struct {
    *BaseService
    orm *ent.AllocateClient
}

func NewAllocate(params ...any) *allocateService {
    return &allocateService{
        BaseService: newService(params...),
        orm:         ent.Database.Allocate,
    }
}

func (s *allocateService) QueryID(id uint64) (*ent.Allocate, error) {
    return s.orm.Query().Where(allocate.ID(id)).First(s.ctx)
}

func (s *allocateService) QueryIDX(id uint64) *ent.Allocate {
    ea, _ := s.QueryID(id)
    if ea == nil {
        snag.Panic("未找到信息")
    }
    return ea
}

// QueryEffectiveSubscribeID 查询生效中的分配信息
func (s *allocateService) QueryEffectiveSubscribeID(subscribeID uint64) (*ent.Allocate, error) {
    return s.orm.Query().
        Where(
            allocate.SubscribeID(subscribeID),
            allocate.TimeGTE(carbon.Now().SubSeconds(model.AllocateExpiration).Carbon2Time()),
        ).
        First(s.ctx)
}

func (s *allocateService) QueryEffectiveSubscribeIDX(subscribeID uint64) *ent.Allocate {
    ea, _ := s.QueryEffectiveSubscribeID(subscribeID)
    if ea == nil {
        snag.Panic("未找到有效分配信息")
    }
    return ea
}

// UnallocatedEbikeInfo 获取未分配车辆信息
func (s *allocateService) UnallocatedEbikeInfo(keyword string) model.EbikeInfo {
    if s.entStore == nil {
        snag.Panic("店员未上班")
    }

    bike, _ := NewEbike().AllocatableBaseFilter().Where(
        ebike.Or(
            ebike.Sn(keyword),
            ebike.Plate(keyword),
        ),
    ).First(s.ctx)
    if bike == nil {
        snag.Panic("未找到可分配电车")
    }

    // 查询是否已被分配
    NewEbike().IsAllocatedX(bike.ID)

    return model.EbikeInfo{
        ID:        bike.ID,
        SN:        bike.Sn,
        ExFactory: bike.ExFactory,
        Plate:     bike.Plate,
        Color:     bike.Color,
    }
}

// AllocateEbike 分配车辆
func (s *allocateService) AllocateEbike(req *model.AllocateCreateReq) model.IDPostReq {
    if req.EbikeID == nil {
        snag.Panic("需要分配电车")
    }

    if s.entStore == nil {
        snag.Panic("未找到门店信息")
    }

    // 查找订阅
    sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(
        subscribe.ID(req.SubscribeID),
        subscribe.Status(model.SubscribeStatusInactive),
        subscribe.NeedContract(true),
        subscribe.EbikeIDIsNil(),
        subscribe.BrandIDNotNil(),
    ).WithRider().First(s.ctx)

    if sub == nil {
        snag.Panic("未找到订阅信息")
    }

    if sub.CityID != s.entStore.CityID {
        snag.Panic("无法跨城市操作")
    }

    // 查找骑手
    r := sub.Edges.Rider
    if r == nil {
        snag.Panic("骑手查询失败")
    }

    // 是否被分配过
    ea, _ := s.orm.Query().Where(allocate.SubscribeID(req.SubscribeID)).First(s.ctx)
    if ea != nil {
        if ea.Time.After(carbon.Now().SubSeconds(model.AllocateExpiration).Carbon2Time()) {
            snag.Panic("已被分配过")
        }
    }

    // 查找电车
    bike := NewEbike().QueryAllocatableX(*req.EbikeID, s.entStore.ID)

    // 比对型号
    if bike.BrandID != *sub.BrandID {
        snag.Panic("待分配车辆型号错误")
    }

    // 查找型号信息
    brand := bike.Edges.Brand
    if brand == nil {
        snag.Panic("电车型号查询失败")
    }

    // 判定电池库存
    if NewStockBatchable().Fetch(model.StockTargetStore, s.entStore.ID, sub.Model) < 1 {
        snag.Panic("电池库存不足")
    }

    info := &model.Allocate{
        Rider: model.Rider{
            ID:    r.ID,
            Phone: r.Phone,
            Name:  r.Name,
        },
        Ebike: NewEbike().Detail(bike, brand),
    }

    // 存储分配信息
    id, err := s.orm.Create().
        SetType(allocate.TypeEbike).
        SetEmployee(s.entEmployee).
        SetStore(s.entStore).
        SetEbike(bike).
        SetBrand(brand).
        SetSubscribe(sub).
        SetRider(r).
        SetStatus(model.AllocateStatusPending.Value()).
        SetInfo(info).
        SetTime(time.Now()).
        SetModel(sub.Model).
        OnConflictColumns(allocate.FieldSubscribeID).
        UpdateNewValues().
        ID(s.ctx)

    if err != nil {
        log.Errorf("分配失败: %v", err)
        snag.Panic("分配失败")
    }

    // 推送签约消息
    socket.SendMessage(NewRiderSocket(), r.ID, &model.RiderSocketMessage{ContractSign: &model.ContractSignReq{
        SubscribeID: sub.ID,
        StoreID:     silk.Pointer(s.entStore.ID),
    }})

    return model.IDPostReq{
        ID: id,
    }
}

func (s *allocateService) detail(ea *ent.Allocate) model.AllocateDetail {
    return model.AllocateDetail{
        Status:   model.AllocateStatus(ea.Status),
        Time:     ea.Time.Format(carbon.DateTimeLayout),
        ID:       ea.ID,
        Allocate: ea.Info,
        Type:     ea.Type.String(),
    }
}

// Info 分配信息
func (s *allocateService) Info(req *model.IDParamReq) model.AllocateDetail {
    ea := s.QueryIDX(req.ID)
    return s.detail(ea)
}

// EmployeeList 电车分配店员列表
func (s *allocateService) EmployeeList(req *model.AllocateEmployeeListReq) *model.PaginationRes {
    q := s.orm.Query().Where(allocate.EmployeeID(s.employee.ID)).Order(ent.Desc(allocate.FieldTime))
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Allocate) model.AllocateDetail {
        return s.detail(item)
    })
}
