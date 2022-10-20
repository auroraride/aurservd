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
    al, _ := s.QueryID(id)
    if al == nil {
        snag.Panic("未找到信息")
    }
    return al
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
    al, _ := s.QueryEffectiveSubscribeID(subscribeID)
    if al == nil {
        snag.Panic("未找到有效分配信息")
    }
    return al
}

// UnallocatedEbikeInfo 获取未分配车辆信息
func (s *allocateService) UnallocatedEbikeInfo(keyword string) model.Ebike {
    bike, _ := NewEbike().AllocatableBaseFilter().Where(
        ebike.Or(
            ebike.Sn(keyword),
            ebike.Plate(keyword),
        ),
        ebike.StoreIDNotNil(),
    ).WithBrand().First(s.ctx)
    if bike == nil {
        snag.Panic("未找到可分配电车")
    }

    // 查询是否已被分配
    NewEbike().IsAllocatedX(bike.ID)

    brand := bike.Edges.Brand

    return model.Ebike{
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

func (s *allocateService) Create(req *model.AllocateCreateReq) model.IDPostReq {
    if req.StoreID == nil && req.CabinetID == nil {
        snag.Panic("必须由门店或电柜参与")
    }

    if req.StoreID != nil && req.CabinetID != nil {
        snag.Panic("门店和电柜不能同时存在")
    }

    // 查找订阅
    _, sub := NewBusinessRider(nil).Inactive(req.SubscribeID)

    if sub == nil {
        snag.Panic("未找到订阅信息")
    }

    var (
        cityID     uint64
        entStore   *ent.Store
        entCabinet *ent.Cabinet
    )

    if req.StoreID != nil {
        entStore = NewStore().Query(*req.StoreID)
        cityID = entStore.CityID
    }

    if req.CabinetID != nil {
        entCabinet = NewCabinet().QueryOne(*req.CabinetID)
        if entCabinet.CityID != nil {
            cityID = *entCabinet.CityID
        }
    }

    if sub.CityID != cityID {
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
    exists, _ := s.orm.Query().Where(allocate.SubscribeID(req.SubscribeID)).First(s.ctx)
    if exists != nil {
        if exists.Time.After(carbon.Now().SubSeconds(model.AllocateExpiration).Carbon2Time()) {
            snag.Panic("已被分配过")
        }
    }

    // 查找电车
    var bikeID, brandID *uint64
    typ := allocate.TypeBattery
    if req.EbikeID != nil {
        typ = allocate.TypeEbike
        if req.StoreID == nil {
            snag.Panic("车电必须由门店参与")
        }
        bike := NewEbike().QueryAllocatableX(*req.EbikeID, *req.StoreID)

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
    if req.StoreID != nil {
        if !NewStock().CheckStore(s.entStore.ID, sub.Model, 1) {
            snag.Panic("电池库存不足")
        }
    }

    // 存储分配信息
    id, err := s.orm.Create().
        SetType(typ).
        SetNillableEmployeeID(req.EmployeeID).
        SetNillableStoreID(req.StoreID).
        SetNillableCabinetID(req.CabinetID).
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
    }})

    return model.IDPostReq{
        ID: id,
    }
}

func (s *allocateService) detail(al *ent.Allocate) model.AllocateDetail {
    r := al.Edges.Rider
    res := model.AllocateDetail{
        ID:     al.ID,
        Type:   al.Type.String(),
        Status: model.AllocateStatus(al.Status),
        Time:   al.Time.Format(carbon.DateTimeLayout),
        Model:  al.Model,
        Rider: model.Rider{
            ID:    r.ID,
            Phone: r.Phone,
            Name:  r.Name,
        },
    }

    bike := al.Edges.Ebike
    if bike != nil {
        res.Ebike = &model.EbikeInfo{
            ID:        bike.ID,
            SN:        bike.Sn,
            ExFactory: bike.ExFactory,
            Plate:     bike.Plate,
            Color:     bike.Color,
        }
    }

    brand := al.Edges.Brand
    if brand != nil {
        res.EbikeBrand = &model.EbikeBrand{
            ID:    brand.ID,
            Name:  brand.Name,
            Cover: brand.Cover,
        }
    }
    return res
}

// Info 分配信息
func (s *allocateService) Info(req *model.IDParamReq) model.AllocateDetail {
    al := s.QueryIDX(req.ID)
    return s.detail(al)
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
