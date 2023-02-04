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
    "github.com/auroraride/aurservd/internal/ent/contract"
    "github.com/auroraride/aurservd/internal/ent/ebike"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/golang-module/carbon/v2"
    "go.uber.org/zap"
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
    q := NewEbike().AllocatableBaseFilter().Where(
        ebike.Or(
            ebike.Sn(keyword),
            ebike.Plate(keyword),
        ),
        ebike.StoreIDNotNil(),
    )
    if s.entStore != nil {
        q.Where(ebike.StoreID(s.entStore.ID))
    }
    bike, _ := q.WithBrand().First(s.ctx)
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

// Create 订单激活分配
func (s *allocateService) Create(req *model.AllocateCreateReq) model.AllocateCreateRes {
    if req.StoreID == nil && req.CabinetID == nil {
        snag.Panic("必须由门店或电柜参与")
    }

    if req.StoreID != nil && req.CabinetID != nil {
        snag.Panic("门店和电柜不能同时存在")
    }

    // 查找订阅
    _, sub := NewBusinessRider(nil).Inactive(*req.SubscribeID)

    if sub == nil {
        snag.Panic("未找到订阅信息")
    }

    // 查询是否已签约
    if exists, _ := ent.Database.Contract.QueryNotDeleted().Where(
        contract.SubscribeID(sub.ID),
        contract.Status(model.ContractStatusSuccess.Value()),
        contract.Effective(true),
    ).Exist(s.ctx); exists {
        snag.Panic("已签约, 无法重新分配")
    }

    e := sub.Edges.Enterprise
    if e != nil {
        // if e.Agent && !e.UseStore && s.employee != nil {
        // 就算使用门店也禁止门店端激活
        if e.Agent && s.employee != nil {
            snag.Panic("代理骑手无法分配")
        }

        if e.Payment == model.EnterprisePaymentPrepay && e.Balance < 0 {
            snag.Panic("企业已欠费")
        }
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
    exists, _ := s.orm.Query().Where(allocate.SubscribeID(*req.SubscribeID)).First(s.ctx)
    if exists != nil {
        if exists.Time.After(carbon.Now().SubSeconds(model.AllocateExpiration).Carbon2Time()) {
            snag.Panic("已被分配过")
        }
    }

    // 查找电车
    var bikeID, brandID *uint64
    var bikeInfo *model.EbikeBusinessInfo
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
        bikeInfo = &model.EbikeBusinessInfo{
            ID:        bike.ID,
            BrandID:   brand.ID,
            BrandName: brand.Name,
        }
    }

    // 判定电池库存
    if req.StoreID != nil {
        if !NewStock().CheckStore(*req.StoreID, sub.Model, 1) {
            snag.Panic("电池库存不足")
        }
    }

    // 存储分配信息
    status := model.AllocateStatusPending.Value()
    if !sub.NeedContract {
        status = model.AllocateStatusSigned.Value()
    }

    // 强制删除原有分配信息
    s.SubscribeDeleteIfExists(sub.ID)

    // 保存分配信息
    allo, err := s.orm.Create().
        SetType(typ).
        SetNillableEmployeeID(req.EmployeeID).
        SetNillableStoreID(req.StoreID).
        SetNillableCabinetID(req.CabinetID).
        SetNillableEbikeID(bikeID).
        SetNillableBrandID(brandID).
        SetSubscribe(sub).
        SetRider(r).
        SetStatus(status).
        SetTime(time.Now()).
        SetModel(sub.Model).
        Save(s.ctx)

    if err != nil {
        zap.L().Error("分配失败", zap.Error(err))
        snag.Panic("分配失败")
    }

    if sub.NeedContract {
        // 需要签约, 推送签约消息
        socket.SendMessage(NewRiderSocket(), r.ID, &model.RiderSocketMessage{ContractSign: &model.ContractSignReq{
            SubscribeID: sub.ID,
        }})
    } else {
        // 无须签约, 直接激活
        var srv *businessRiderService
        // 直接激活
        if s.modifier != nil {
            srv = NewBusinessRiderWithModifier(s.modifier)
        } else {
            srv = NewBusinessRider(r)
        }
        srv.SetStoreID(allo.StoreID).
            SetCabinetID(allo.CabinetID).
            SetEmployeeID(allo.EmployeeID)
        // TODO 电车直接激活?
        if bikeID != nil || brandID != nil {
            snag.Panic("暂不支持此业务")
            srv.SetEbike(bikeInfo)
        }
        srv.Active(sub, allo)
    }

    return model.AllocateCreateRes{
        ID:           allo.ID,
        NeedContract: sub.NeedContract,
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
        Ebike: NewEbike().Detail(al.Edges.Ebike, al.Edges.Brand),
    }

    if time.Now().Sub(al.Time).Seconds() > model.AllocateExpiration && res.Status == model.AllocateStatusPending {
        res.Status = model.AllocateStatusVoid
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
    status := req.Status
    if req.Status != 2 {
        status = model.AllocateStatusPending
    }
    q := s.orm.Query().
        WithRider().
        WithEbike().
        WithBrand().
        Where(allocate.EmployeeID(s.employee.ID), allocate.Status(status.Value())).
        Order(ent.Desc(allocate.FieldTime))
    return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Allocate) model.AllocateDetail {
        return s.detail(item)
    })
}

// LoopStatus 长连接查询是否已分配
func (s *allocateService) LoopStatus(riderID, subscribeID uint64) (res model.AllocateRiderRes) {
    ticker := time.NewTicker(3 * time.Second)
    defer ticker.Stop()

    start := time.Now()
    for ; true; <-ticker.C {

        allo, _ := s.orm.Query().Where(
            allocate.RiderID(riderID),
            allocate.SubscribeID(subscribeID),
            allocate.TimeGT(carbon.Now().SubSeconds(model.AllocateExpiration).Carbon2Time()),
        ).First(s.ctx)

        // 如果有分配信息 并且 状态为待签约 并且 非电柜扫码
        if allo != nil && allo.Status == model.AllocateStatusPending.Value() && allo.CabinetID == nil {
            res.Allocated = true
            return
        }

        if time.Now().Sub(start).Seconds() > 50 {
            return
        }
    }

    return
}

// SubscribeDeleteIfExists 根据subscribeID强制删除分配信息
func (s *allocateService) SubscribeDeleteIfExists(subscribeID uint64) {
    _, err := s.orm.Delete().Where(allocate.SubscribeID(subscribeID)).Exec(s.ctx)
    if err != nil {
        zap.L().Error("分配信息强制删除失败", zap.Error(err))
    }
}
