// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-14
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/ebike"
    "github.com/auroraride/aurservd/internal/ent/ebikeallocate"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    log "github.com/sirupsen/logrus"
    "time"
)

type ebikeAllocateService struct {
    *BaseService
    orm *ent.EbikeAllocateClient
}

func NewEbikeAllocate(params ...any) *ebikeAllocateService {
    return &ebikeAllocateService{
        BaseService: newService(params...),
        orm:         ent.Database.EbikeAllocate,
    }
}

func (s *ebikeAllocateService) QueryID(id uint64) (*ent.EbikeAllocate, error) {
    return s.orm.Query().Where(ebikeallocate.ID(id)).First(s.ctx)
}

func (s *ebikeAllocateService) QueryIDX(id uint64) *ent.EbikeAllocate {
    ea, _ := s.QueryID(id)
    if ea == nil {
        snag.Panic("未找到信息")
    }
    return ea
}

// QueryEffectiveSubscribeID 查询生效中的分配信息
func (s *ebikeAllocateService) QueryEffectiveSubscribeID(subscribeID uint64) (*ent.EbikeAllocate, error) {
    return s.orm.Query().
        Where(
            ebikeallocate.SubscribeID(subscribeID),
            ebikeallocate.TimeGTE(carbon.Now().SubSeconds(model.EbikeAllocateExpiration).Carbon2Time()),
        ).
        First(s.ctx)
}

func (s *ebikeAllocateService) QueryEffectiveSubscribeIDX(subscribeID uint64) *ent.EbikeAllocate {
    ea, _ := s.QueryEffectiveSubscribeID(subscribeID)
    if ea == nil {
        snag.Panic("未找到有效分配信息")
    }
    return ea
}

// UnallocatedInfo 获取未分配车辆信息
func (s *ebikeAllocateService) UnallocatedInfo(keyword string) model.EbikeInfo {
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

// Allocate 分配车辆
func (s *ebikeAllocateService) Allocate(req *model.EbikeAllocateReq) model.IDPostReq {
    if s.entEmployee == nil {
        snag.Panic("需要店员操作")
    }

    if s.entStore == nil {
        snag.Panic("店员未上班")
    }

    // 查找订阅
    sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(
        subscribe.ID(req.SubscribeID),
        subscribe.Status(model.SubscribeStatusInactive),
        subscribe.NeedContract(true),
        subscribe.EbikeIDIsNil(),
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
    ea, _ := s.orm.Query().Where(ebikeallocate.SubscribeID(req.SubscribeID)).First(s.ctx)
    if ea != nil {
        if ea.Time.After(carbon.Now().SubSeconds(model.EbikeAllocateExpiration).Carbon2Time()) {
            snag.Panic("已被分配过")
        }
    }

    // 查找电车
    bike := NewEbike().QueryAllocatableX(req.EbikeID, s.entStore.ID)
    brand := bike.Edges.Brand
    if brand == nil {
        snag.Panic("电车型号查询失败")
    }

    // 判定库存
    if NewStockBatchable().Fetch(model.StockTargetStore, s.entStore.ID, sub.Model) < 1 {
        snag.Panic("电池库存不足")
    }

    info := &model.EbikeAllocate{
        Rider: model.Rider{
            ID:    r.ID,
            Phone: r.Phone,
            Name:  r.Name,
        },
        Ebike: model.Ebike{
            EbikeInfo: model.EbikeInfo{
                ID:        bike.ID,
                SN:        bike.Sn,
                ExFactory: bike.ExFactory,
                Plate:     bike.Plate,
                Color:     bike.Color,
            },
            Brand: model.EbikeBrand{
                ID:   brand.ID,
                Name: brand.Name,
            },
        },
        Model: sub.Model,
    }

    // 存储分配信息
    id, err := s.orm.Create().
        SetEmployee(s.entEmployee).
        SetStore(s.entStore).
        SetEbike(bike).
        SetBrand(brand).
        SetSubscribe(sub).
        SetRider(r).
        SetStatus(model.EbikeAllocateStatusPending.Value()).
        SetInfo(info).
        SetTime(time.Now()).
        OnConflictColumns(ebikeallocate.FieldSubscribeID).
        UpdateNewValues().
        ID(s.ctx)

    if err != nil {
        log.Errorf("分配失败: %v", err)
        snag.Panic("分配失败")
    }

    return model.IDPostReq{
        ID: id,
    }
}

func (s *ebikeAllocateService) detail(ea *ent.EbikeAllocate) model.EbikeAllocateDetail {
    return model.EbikeAllocateDetail{
        Status:        model.EbikeAllocateStatus(ea.Status),
        Time:          ea.Time.Format(carbon.DateTimeLayout),
        ID:            ea.ID,
        EbikeAllocate: ea.Info,
    }
}

// Info 电车分配信息
func (s *ebikeAllocateService) Info(req *model.IDParamReq) model.EbikeAllocateDetail {
    ea := s.QueryIDX(req.ID)
    return s.detail(ea)
}

// EmployeeList 电车分配店员列表
func (s *ebikeAllocateService) EmployeeList(req *model.EbikeAllocateEmployeeListReq) *model.PaginationRes {
    q := s.orm.Query().Where(ebikeallocate.EmployeeID(s.employee.ID)).Order(ent.Desc(ebikeallocate.FieldTime))
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.EbikeAllocate) model.EbikeAllocateDetail {
        return s.detail(item)
    })
}
