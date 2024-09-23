// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-02, by aurb

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/internal/ent/agreement"
	entasset "github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/assettransferdetails"
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/pkg/tools"
)

type enterpriseAssetBiz struct {
	orm      *ent.EnterpriseClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewEnterpriseAsset() *enterpriseAssetBiz {
	return &enterpriseAssetBiz{
		orm: ent.Database.Enterprise,
		ctx: context.Background(),
	}
}

func NewEnterpriseAssetWithModifier(m *model.Modifier) *enterpriseAssetBiz {
	s := NewEnterpriseAsset()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// Assets 资产列表
func (b *enterpriseAssetBiz) Assets(req *definition.EnterpriseAssetListReq) (res *model.PaginationRes) {
	// 查询分页数据
	q := ent.Database.EnterpriseStation.QueryNotDeleted().WithCity().WithEnterprise().Order(ent.Desc(agreement.FieldCreatedAt))
	b.assetsFilter(q, req)
	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.EnterpriseStation) (result *definition.EnterpriseAssetDetail) {
		result = &definition.EnterpriseAssetDetail{
			ID:    item.ID,
			Name:  item.Name,
			Total: b.AssetTotal(req, item.ID),
		}
		if item.Edges.City != nil {
			result.City = model.City{
				ID:   item.Edges.City.ID,
				Name: item.Edges.City.Name,
			}
		}

		if item.Edges.Enterprise != nil {
			result.Enterprise = definition.EnterpriseDetail{
				ID:   item.Edges.Enterprise.ID,
				Name: item.Edges.Enterprise.Name,
			}
		}

		return result
	})

	return res
}

// assetsFilter 条件筛选
func (b *enterpriseAssetBiz) assetsFilter(q *ent.EnterpriseStationQuery, req *definition.EnterpriseAssetListReq) {
	if req.EnterpriseID != nil {
		q.Where(enterprisestation.HasEnterpriseWith(enterprise.ID(*req.EnterpriseID)))
	}
	if req.StationID != nil {
		q.Where(enterprisestation.ID(*req.StationID))
	}
	if req.CityID != nil {
		q.Where(enterprisestation.CityID(*req.CityID))
	}
	if req.ModelID != nil {
		// 查询型号资产
		ids := make([]uint64, 0)
		list, _ := ent.Database.Asset.QueryNotDeleted().WithStation().Where(
			entasset.ModelID(*req.ModelID),
			entasset.LocationsType(model.AssetLocationsTypeStation.Value()),
			entasset.Status(model.AssetStatusStock.Value()),
		).All(b.ctx)
		for _, v := range list {
			if v.Edges.Station != nil {
				ids = append(ids, v.Edges.Station.ID)
			}
		}

		q.Where(
			enterprisestation.IDIn(ids...),
		)
	}

	if req.OtherName != nil {
		// 查询其他物资资产
		ids := make([]uint64, 0)
		list, _ := ent.Database.Asset.QueryNotDeleted().WithStation().Where(
			entasset.LocationsType(model.AssetLocationsTypeStation.Value()),
			entasset.Status(model.AssetStatusStock.Value()),
			entasset.HasMaterialWith(material.NameContainsFold(*req.OtherName)),
		).All(b.ctx)
		for _, v := range list {
			if v.Edges.Station != nil {
				ids = append(ids, v.Edges.Station.ID)
			}
		}

		q.Where(
			enterprisestation.IDIn(ids...),
		)
	}
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(enterprisestation.CreatedAtGTE(start), enterprisestation.CreatedAtLTE(end))
	}

}

// AssetTotal 物资数据统计
func (b *enterpriseAssetBiz) AssetTotal(req *definition.EnterpriseAssetListReq, id uint64) (res definition.CommonAssetTotal) {
	// 查询所属资产数据
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			entasset.LocationsType(model.AssetLocationsTypeStation.Value()),
			entasset.LocationsIDIn(id),
			entasset.Status(model.AssetStatusStock.Value()),
		)

	if req.ModelID != nil {
		q.Where(
			entasset.ModelID(*req.ModelID),
			entasset.TypeIn(model.AssetTypeSmartBattery.Value(), model.AssetTypeNonSmartBattery.Value()),
		)
	}
	if req.BrandId != nil {
		q.Where(
			entasset.BrandID(*req.BrandId),
			entasset.Type(model.AssetTypeEbike.Value()),
		)
	}
	if req.OtherName != nil {
		q.Where(
			entasset.NameContains(*req.OtherName),
			entasset.TypeIn(model.AssetTypeCabinetAccessory.Value(), model.AssetTypeEbikeAccessory.Value(), model.AssetTypeOtherAccessory.Value()),
		)
	}
	list, _ := q.All(b.ctx)
	for _, v := range list {
		switch v.Type {
		case model.AssetTypeEbike.Value():
			res.EbikeTotal += 1
		case model.AssetTypeSmartBattery.Value():
			res.SmartBatteryTotal += 1
		case model.AssetTypeNonSmartBattery.Value():
			res.NonSmartBatteryTotal += 1
		case model.AssetTypeEbikeAccessory.Value():
			res.EbikeAccessoryTotal += 1
		case model.AssetTypeCabinetAccessory.Value():
			res.CabinetAccessoryTotal += 1
		case model.AssetTypeOtherAccessory.Value():
			res.OtherAssetTotal += 1
		}
	}
	return
}

// AssetDetail 物资详情
func (b *enterpriseAssetBiz) AssetDetail(id uint64) (ast *definition.CommonAssetDetail) {
	ast = &definition.CommonAssetDetail{
		Ebikes:             make([]*definition.AssetMaterial, 0),
		SmartBatteries:     make([]*definition.AssetMaterial, 0),
		NonSmartBatteries:  make([]*definition.AssetMaterial, 0),
		CabinetAccessories: make([]*definition.AssetMaterial, 0),
		EbikeAccessories:   make([]*definition.AssetMaterial, 0),
		OtherAssets:        make([]*definition.AssetMaterial, 0),
	}

	ebikeNameMap := make(map[string]*definition.AssetMaterial)
	sBNameMap := make(map[string]*definition.AssetMaterial)
	nSbNameMap := make(map[string]*definition.AssetMaterial)
	cabAccNameMap := make(map[string]*definition.AssetMaterial)
	ebikeAccNameMap := make(map[string]*definition.AssetMaterial)
	otherAccNameMap := make(map[string]*definition.AssetMaterial)

	// 入库物资调拨详情
	inAts, _ := ent.Database.AssetTransferDetails.QueryNotDeleted().
		Where(
			assettransferdetails.IsIn(true),
			assettransferdetails.HasTransferWith(
				assettransfer.Status(model.AssetTransferStatusStock.Value()),
				assettransfer.ToLocationType(model.AssetLocationsTypeStation.Value()),
				assettransfer.ToLocationID(id),
				assettransfer.DeletedAtIsNil(),
			),
		).
		WithAsset(func(query *ent.AssetQuery) {
			query.WithBrand().WithModel().WithMaterial()
		}).All(b.ctx)

	NewAssetTransferDetails().TransferInOut(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, inAts, false)

	// 出库物资调拨详情
	outAts, _ := ent.Database.AssetTransferDetails.QueryNotDeleted().
		Where(
			assettransferdetails.HasTransferWith(
				assettransfer.StatusIn(model.AssetTransferStatusDelivering.Value(), model.AssetTransferStatusStock.Value()),
				assettransfer.FromLocationType(model.AssetLocationsTypeStation.Value()),
				assettransfer.FromLocationID(id),
				assettransfer.DeletedAtIsNil(),
			),
		).
		WithAsset(func(query *ent.AssetQuery) {
			query.WithBrand().WithModel().WithMaterial()
		}).All(b.ctx)

	NewAssetTransferDetails().TransferInOut(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, outAts, true)

	NewAssetTransferDetails().GetCommonAssetDetail(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, ast)

	return
}
