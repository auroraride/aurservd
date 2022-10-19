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
    return s.orm.Query().Where(allocate.ID(id)).WithBrand().WithRider().First(s.ctx)
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

func (s *allocateService) Create(req *model.AllocateCreateReq) model.IDPostReq {

    if s.entStore == nil {
        snag.Panic("未找到门店信息")
    }

    // 查找订阅
    _, sub := NewBusinessRider(nil).Inactive(req.SubscribeID)

    if sub == nil {
        snag.Panic("未找到订阅信息")
    }

    if sub.CityID != s.entStore.CityID {
        snag.Panic("无法跨城市操作")
    }

    if sub.BrandID != nil && req.EbikeID == nil {
        snag.Panic("需要分配电车")
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
    var bikeID, brandID *uint64
    if req.EbikeID != nil {
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

        bikeID = silk.UInt64(bike.ID)
        brandID = silk.UInt64(brand.ID)
    }

    // 判定电池库存
    if NewStockBatchable().Fetch(model.StockTargetStore, s.entStore.ID, sub.Model) < 1 {
        snag.Panic("电池库存不足")
    }

    // 存储分配信息
    id, err := s.orm.Create().
        SetType(allocate.TypeEbike).
        SetEmployee(s.entEmployee).
        SetStore(s.entStore).
        SetNillableEbikeID(bikeID).
        SetNillableBrandID(brandID).
        SetSubscribe(sub).
        SetRider(r).
        SetStatus(model.AllocateStatusPending.Value()).
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
    r := ea.Edges.Rider
    res := model.AllocateDetail{
        ID:     ea.ID,
        Type:   ea.Type.String(),
        Status: model.AllocateStatus(ea.Status),
        Time:   ea.Time.Format(carbon.DateTimeLayout),
        Model:  ea.Model,
        Rider: model.Rider{
            ID:    r.ID,
            Phone: r.Phone,
            Name:  r.Name,
        },
    }
    bike := ea.Edges.Ebike
    brand := ea.Edges.Brand
    if bike != nil && brand != nil {
        res.Ebike = &model.Ebike{
            EbikeInfo: model.EbikeInfo{
                ID:        bike.ID,
                SN:        bike.Sn,
                ExFactory: bike.ExFactory,
                Plate:     bike.Plate,
                Color:     bike.Color,
            },
            Brand: &model.EbikeBrand{
                ID:    brand.ID,
                Name:  brand.Name,
                Cover: brand.Cover,
            },
        }
    }
    return res
}

// Info 分配信息
func (s *allocateService) Info(req *model.IDParamReq) model.AllocateDetail {
    ea := s.QueryIDX(req.ID)
    return s.detail(ea)
}

// EmployeeList 电车分配店员列表
func (s *allocateService) EmployeeList(req *model.AllocateEmployeeListReq) *model.PaginationRes {
    q := s.orm.Query().
        WithRider().
        WithEbike().
        WithBrand().
        Where(allocate.EmployeeID(s.employee.ID)).
        Order(ent.Desc(allocate.FieldTime))
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Allocate) model.AllocateDetail {
        return s.detail(item)
    })
}
