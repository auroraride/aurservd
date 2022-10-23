// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
    "github.com/auroraride/aurservd/app/model"
    "github.com/auroraride/aurservd/internal/ent"
    "github.com/auroraride/aurservd/internal/ent/stock"
    "github.com/auroraride/aurservd/internal/ent/subscribe"
    "github.com/auroraride/aurservd/pkg/silk"
    "github.com/auroraride/aurservd/pkg/snag"
    "github.com/auroraride/aurservd/pkg/tools"
)

type managerSubscribeService struct {
    *BaseService
}

func NewManagerSubscribe(params ...any) *managerSubscribeService {
    s := &managerSubscribeService{
        BaseService: newService(params...),
    }
    if s.modifier == nil {
        snag.Panic("无权限操作")
    }
    return s
}

func (s *managerSubscribeService) Active(req *model.ManagerSubscribeActive) {
    var bikeID *uint64
    if req.EbikeKeyword != nil {
        bike := NewAllocate().UnallocatedEbikeInfo(*req.EbikeKeyword)
        bikeID = silk.UInt64(bike.ID)
    }

    NewAllocate(s.modifier).Create(&model.AllocateCreateReq{
        EbikeID:     bikeID,
        SubscribeID: silk.UInt64(req.ID),
        StoreID:     req.StoreID,
    })
}

func (s *managerSubscribeService) ChangeEBike(req *model.ManagerSubscribeChangeEbike) {
    bike := NewAllocate().UnallocatedEbikeInfo(*req.EbikeKeyword)

    sub, _ := ent.Database.Subscribe.QueryNotDeleted().
        Where(
            subscribe.Status(model.SubscribeStatusUsing),
            subscribe.EbikeIDNotNil(),
        ).
        WithBrand().
        First(s.ctx)
    if sub == nil {
        snag.Panic("未找到订阅")
    }

    if bike.Brand.ID != *sub.BrandID {
        snag.Panic("电车型号不同")
    }

    ent.WithTxPanic(s.ctx, func(tx *ent.Tx) (err error) {
        // 旧车入库
        err = tx.Stock.Create().
            SetEbikeID(*sub.EbikeID).
            SetNum(1).
            SetStoreID(req.StoreID).
            SetSn(tools.NewUnique().NewSN()).
            SetRiderID(sub.RiderID).
            SetName(sub.Edges.Brand.Name).
            SetMaterial(stock.MaterialEbike).
            SetCityID(sub.CityID).
            SetSubscribeID(sub.ID).
            SetBrandID(*sub.BrandID).
            Exec(s.ctx)
        if err != nil {
            return
        }

        // 新车出库
        err = tx.Stock.Create().
            SetEbikeID(bike.ID).
            SetNum(-1).
            SetStoreID(req.StoreID).
            SetSn(tools.NewUnique().NewSN()).
            SetRiderID(sub.RiderID).
            SetName(bike.Brand.Name).
            SetMaterial(stock.MaterialEbike).
            SetCityID(sub.CityID).
            SetSubscribeID(sub.ID).
            SetBrandID(bike.Brand.ID).
            Exec(s.ctx)
        if err != nil {
            return
        }

        // 更新新车所属
        err = tx.Ebike.UpdateOneID(bike.ID).SetRiderID(sub.RiderID).SetStatus(model.EbikeStatusUsing).Exec(s.ctx)
        if err != nil {
            return
        }

        // 删除电车所属
        err = tx.Ebike.UpdateOneID(*sub.EbikeID).ClearRiderID().SetStatus(model.EbikeStatusInStock).Exec(s.ctx)
        if err != nil {
            return
        }

        // 更新订阅
        return tx.Subscribe.UpdateOneID(sub.ID).SetEbikeID(bike.ID).SetBrandID(bike.Brand.ID).Exec(s.ctx)
    })
}
