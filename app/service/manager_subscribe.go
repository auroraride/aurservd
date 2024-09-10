// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-19
// Based on aurservd by liasica, magicrolan@qq.com.

package service

import (
	"github.com/auroraride/aurservd/app/logging"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetattributes"
	"github.com/auroraride/aurservd/internal/ent/assetattributevalues"
	"github.com/auroraride/aurservd/internal/ent/subscribe"
	"github.com/auroraride/aurservd/pkg/silk"
	"github.com/auroraride/aurservd/pkg/snag"
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
func (s *managerSubscribeService) Active(req *model.ManagerSubscribeActive) model.AllocateCreateRes {
	return NewAllocate(s.modifier).Create(&model.AllocateCreateParams{
		SubscribeID: silk.UInt64(req.ID),
		StoreID:     req.StoreID,
		BatteryID:   req.BatteryID,
		EbikeParam: model.AllocateCreateEbikeParam{
			Keyword: req.EbikeKeyword,
		},
	})
}

// ChangeEbike 修改订阅车辆
func (s *managerSubscribeService) ChangeEbike(req *model.ManagerSubscribeChangeEbike) {
	sub, _ := ent.Database.Subscribe.QueryNotDeleted().
		Where(
			subscribe.Status(model.SubscribeStatusUsing),
			subscribe.EbikeIDNotNil(),
			subscribe.ID(req.ID),
		).
		WithBrand().
		WithEbike().
		First(s.ctx)
	if sub == nil {
		snag.Panic("未找到订阅")
	}

	if sub.StationID == nil && req.StoreID == nil {
		snag.Panic("门店为必选")
	}

	if sub.StationID != nil && req.StoreID != nil {
		snag.Panic("代理骑手无法使用门店物资")
	}

	// 获取门店和站点ID
	stationID := sub.StationID
	storeID := req.StoreID

	// 查询车辆
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			asset.Type(model.AssetTypeEbike.Value()),
			asset.Status(model.AssetStatusStock.Value()),
		)
	// 站点
	var toLocationType model.AssetLocationsType
	var toLocationID uint64
	if stationID != nil {
		q.Where(asset.LocationsType(model.AssetLocationsTypeStation.Value()), asset.LocationsID(*stationID))
		toLocationType = model.AssetLocationsTypeStation
		toLocationID = *stationID
	}

	// 门店
	if storeID != nil {
		q.Where(asset.LocationsType(model.AssetLocationsTypeStore.Value()), asset.LocationsID(*storeID))
		toLocationType = model.AssetLocationsTypeStore
		toLocationID = *storeID
	}
	if req.EbikeKeyword != nil {
		q.Where(
			asset.Or(
				asset.SnContainsFold(*req.EbikeKeyword),
			))
		attributes, _ := ent.Database.AssetAttributes.Query().Where(assetattributes.Key("plate")).First(s.ctx)
		if attributes != nil {
			q.Where(
				asset.Or(
					asset.HasValuesWith(
						assetattributevalues.AttributeID(attributes.ID),
						assetattributevalues.ValueContainsFold(*req.EbikeKeyword),
					),
				),
			)
		}
	}
	newBike, _ := q.First(s.ctx)
	if newBike == nil {
		snag.Panic("未找到可用电车")
		return
	}
	if newBike.BrandID != nil && *newBike.BrandID != *sub.BrandID {
		snag.Panic("电车型号不同")
	}

	// 旧车入库
	if sub.Edges.Ebike != nil {
		oldBike := sub.Edges.Ebike
		fromLocationType := model.AssetLocationsType(oldBike.LocationsType)
		_, failed, err := NewAssetTransfer().Transfer(s.ctx, &model.AssetTransferCreateReq{
			FromLocationType: &fromLocationType,
			FromLocationID:   &oldBike.LocationsID,
			ToLocationType:   toLocationType,
			ToLocationID:     toLocationID,
			Details: []model.AssetTransferCreateDetail{
				{
					AssetType: model.AssetTypeEbike,
					SN:        silk.String(oldBike.Sn),
				},
			},
			Reason:            "修改订阅车辆",
			AssetTransferType: model.AssetTransferTypeTransfer,
			OperatorID:        s.modifier.ID,
			OperatorType:      model.OperatorTypeManager,
			AutoIn:            true,
		}, s.modifier)
		if err != nil {
			snag.Panic(err)
			return
		}
		if len(failed) > 0 {
			snag.Panic(failed[0])
		}
	}

	// 新车出库
	fromLocationType := model.AssetLocationsType(newBike.LocationsType)
	_, failed, err := NewAssetTransfer().Transfer(s.ctx, &model.AssetTransferCreateReq{
		FromLocationType: &fromLocationType,
		FromLocationID:   &newBike.LocationsID,
		ToLocationType:   model.AssetLocationsTypeRider,
		ToLocationID:     sub.RiderID,
		Details: []model.AssetTransferCreateDetail{
			{
				AssetType: model.AssetTypeEbike,
				SN:        silk.String(newBike.Sn),
			},
		},
		Reason:            "修改订阅车辆",
		AssetTransferType: model.AssetTransferTypeTransfer,
		OperatorID:        s.modifier.ID,
		OperatorType:      model.OperatorTypeManager,
		AutoIn:            true,
	}, s.modifier)
	if err != nil {
		snag.Panic(err)
	}
	if len(failed) > 0 {
		snag.Panic(failed[0])
	}
	// 更新订阅
	err = ent.Database.Subscribe.UpdateOneID(sub.ID).SetEbikeID(newBike.ID).SetNillableBatteryID(newBike.BrandID).Exec(s.ctx)
	if err != nil {
		snag.Panic(err)
	}
}

func (s *managerSubscribeService) UnbindEbike(req *model.ManagerSubscribeUnbindEbike) {
	sub, _ := ent.Database.Subscribe.QueryNotDeleted().Where(subscribe.EbikeIDNotNil(), subscribe.ID(req.ID)).WithRider().WithEbike().WithBrand().First(s.ctx)
	if sub == nil {
		snag.Panic("未找到有效订阅")
	}
	if sub.Edges.Ebike == nil {
		snag.Panic("未绑定车辆")
	}
	oldEbike := sub.Edges.Ebike
	fromLocationType := model.AssetLocationsType(oldEbike.LocationsType)
	_, failed, err := NewAssetTransfer().Transfer(s.ctx, &model.AssetTransferCreateReq{
		FromLocationType: &fromLocationType,
		FromLocationID:   &oldEbike.LocationsID,
		ToLocationType:   model.AssetLocationsTypeStore,
		ToLocationID:     req.StoreID,
		Details: []model.AssetTransferCreateDetail{
			{
				AssetType: model.AssetTypeEbike,
				SN:        silk.String(oldEbike.Sn),
			},
		},
		Reason:            "解绑车辆",
		AssetTransferType: model.AssetTransferTypeTransfer,
		OperatorID:        s.modifier.ID,
		OperatorType:      model.OperatorTypeManager,
		AutoIn:            true,
	}, s.modifier)
	if err != nil {
		return
	}
	if len(failed) > 0 {
		snag.Panic(failed[0])
	}

	// 删除订阅标的车辆信息, 保留车辆型号以便下次绑定新车
	err = ent.Database.Subscribe.UpdateOneID(sub.ID).ClearEbikeID().Exec(s.ctx)
	if err != nil {
		snag.Panic(err)
	}

	// 记录操作日志
	go logging.NewOperateLog().
		SetOperate(model.OperateUnbindEbike).
		SetRef(sub.Edges.Rider).
		SetDiff("车辆编号: "+sub.Edges.Ebike.Sn, "无车辆").
		SetModifier(s.modifier).
		Send()
}
