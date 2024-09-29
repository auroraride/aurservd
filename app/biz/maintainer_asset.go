// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-02, by aurb

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	entasset "github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/assettransferdetails"
	"github.com/auroraride/aurservd/internal/ent/maintainer"
	"github.com/auroraride/aurservd/internal/ent/material"
)

type maintainerAssetBiz struct {
	orm      *ent.MaintainerClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewMaintainerAsset() *maintainerAssetBiz {
	return &maintainerAssetBiz{
		orm: ent.Database.Maintainer,
		ctx: context.Background(),
	}
}

func NewMaintainerAssetWithModifier(m *model.Modifier) *maintainerAssetBiz {
	s := NewMaintainerAsset()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// Assets 资产列表
func (b *maintainerAssetBiz) Assets(req *definition.MaintainerAssetListReq) (res *model.PaginationRes) {
	// 查询分页数据
	q := b.orm.Query().
		Where(maintainer.Enable(true)).
		Order(ent.Desc(maintainer.FieldID))
	b.assetsFilter(q, req)
	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Maintainer) (result *definition.MaintainerAssetDetail) {
		result = &definition.MaintainerAssetDetail{
			ID:    item.ID,
			Name:  item.Name,
			Phone: item.Phone,
			Total: b.AssetTotal(req, item.ID),
		}
		return result
	})

	return res
}

// assetsFilter 条件筛选
func (b *maintainerAssetBiz) assetsFilter(q *ent.MaintainerQuery, req *definition.MaintainerAssetListReq) {
	if req.Keyword != nil {
		q.Where(maintainer.Or(maintainer.NameContainsFold(*req.Keyword), maintainer.PhoneContainsFold(*req.Keyword)))
	}

	if req.ModelID != nil {
		q.Where(
			maintainer.HasAssetWith(entasset.ModelID(*req.ModelID)),
		)
	}
	if req.BrandID != nil {
		q.Where(
			maintainer.HasAssetWith(entasset.BrandID(*req.BrandID)),
		)
	}
	if req.OtherName != nil {
		q.Where(
			maintainer.HasAssetWith(entasset.HasMaterialWith(material.NameContains(*req.OtherName))),
		)
	}
}

// AssetTotal 仓库物资数据统计
func (b *maintainerAssetBiz) AssetTotal(req *definition.MaintainerAssetListReq, id uint64) (res definition.CommonAssetTotal) {
	// 查询所属资产数据
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			entasset.LocationsType(model.AssetLocationsTypeOperation.Value()),
			entasset.LocationsIDIn(id),
			entasset.Status(model.AssetStatusStock.Value()),
		)
	if req.ModelID != nil {
		q.Where(
			entasset.ModelID(*req.ModelID),
		)
	}
	if req.BrandID != nil {
		q.Where(
			entasset.BrandID(*req.BrandID),
		)
	}
	if req.OtherName != nil {
		q.Where(
			entasset.HasMaterialWith(material.NameContains(*req.OtherName)),
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
func (b *maintainerAssetBiz) AssetDetail(id uint64) (ast *definition.CommonAssetDetail) {
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
				assettransfer.ToLocationType(model.AssetLocationsTypeOperation.Value()),
				assettransfer.ToLocationID(id),
				assettransfer.DeletedAtIsNil(),
			),
		).
		WithAsset(func(query *ent.AssetQuery) {
			query.WithBrand().WithModel().WithMaterial()
		}).All(b.ctx)

	NewAssetTransferDetails().TransferInOut(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, inAts, false, id)

	// 出库物资调拨详情
	outAts, _ := ent.Database.AssetTransferDetails.QueryNotDeleted().
		Where(
			assettransferdetails.HasTransferWith(
				assettransfer.StatusIn(model.AssetTransferStatusDelivering.Value(), model.AssetTransferStatusStock.Value()),
				assettransfer.FromLocationType(model.AssetLocationsTypeOperation.Value()),
				assettransfer.FromLocationID(id),
				assettransfer.DeletedAtIsNil(),
			),
		).
		WithAsset(func(query *ent.AssetQuery) {
			query.WithBrand().WithModel().WithMaterial()
		}).All(b.ctx)

	NewAssetTransferDetails().TransferInOut(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, outAts, true, id)

	NewAssetTransferDetails().GetCommonAssetDetail(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, ast)

	return
}
