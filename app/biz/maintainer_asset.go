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
	"github.com/auroraride/aurservd/internal/ent/maintainer"
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
	q := b.orm.Query().Where(maintainer.Enable(true))
	if req.Keyword != nil {
		q.Where(maintainer.Or(maintainer.NameContainsFold(*req.Keyword), maintainer.PhoneContainsFold(*req.Keyword)))
	}

	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Maintainer) (result *definition.MaintainerAssetDetail) {
		result = &definition.MaintainerAssetDetail{
			ID:              item.ID,
			Name:            item.Name,
			MaintainerAsset: b.assetForMaintainer(req, item.ID),
		}
		return result
	})

	return res
}

func (b *maintainerAssetBiz) assetForMaintainer(req *definition.MaintainerAssetListReq, id uint64) (asset definition.MaintainerAsset) {
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
	// 按照分类进行累计
	for _, v := range list {
		switch v.Type {
		case model.AssetTypeEbike.Value():
			asset.EbikeTotal = asset.EbikeTotal + 1
		case model.AssetTypeSmartBattery.Value():
			asset.SmartBatteryTotal = asset.SmartBatteryTotal + 1
		case model.AssetTypeNonSmartBattery.Value():
			asset.NonSmartBatteryTotal = asset.NonSmartBatteryTotal + 1
		case model.AssetTypeCabinetAccessory.Value():
			asset.CabinetAccessoryTotal = asset.CabinetAccessoryTotal + 1
		case model.AssetTypeEbikeAccessory.Value():
			asset.EbikeAccessoryTotal = asset.EbikeAccessoryTotal + 1
		case model.AssetTypeOtherAccessory.Value():
			asset.OtherAssetTotal = asset.OtherAssetTotal + 1
		}
	}

	// todo 新调拨记录计算

	return asset
}
