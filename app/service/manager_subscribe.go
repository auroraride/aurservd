// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"github.com/auroraride/aurservd/app/logging"
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

// Active 激活订阅
// 团签无需签约
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
		BatteryID:   req.BatteryID,
	})
}

// ChangeEbike 修改订阅车辆
func (s *managerSubscribeService) ChangeEbike(req *model.ManagerSubscribeChangeEbike) {
	bike := NewAllocate().UnallocatedEbikeInfo(*req.EbikeKeyword)

	sub, _ := ent.Database.Subscribe.QueryNotDeleted().
		Where(
			subscribe.Status(model.SubscribeStatusUsing),
			subscribe.EbikeIDNotNil(),
			subscribe.ID(req.ID),
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

func (s *managerSubscribeService) UnbindEbike(req *model.ManagerSubscribeUnbindEbike) {
	sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.EbikeIDNotNil(), subscribe.ID(req.ID)).WithRider().WithEbike().WithBrand().First(s.ctx)
	if sub == nil {
		snag.Panic("未找到有效订阅")
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

		// 删除电车所属
		err = tx.Ebike.UpdateOneID(*sub.EbikeID).ClearRiderID().SetStatus(model.EbikeStatusInStock).Exec(s.ctx)
		if err != nil {
			return
		}

		// 删除订阅标的车辆信息, 保留车辆型号以便下次绑定新车
		return tx.Subscribe.UpdateOneID(sub.ID).ClearEbikeID().Exec(s.ctx)
	})

	// 记录操作日志
	go logging.NewOperateLog().
		SetOperate(model.OperateUnbindEbike).
		SetRef(sub.Edges.Rider).
		SetDiff("车辆编号: "+sub.Edges.Ebike.Sn, "无车辆").
		SetModifier(s.modifier).
		Send()
}
