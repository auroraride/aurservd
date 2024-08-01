// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-01, by aurb

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	entasset "github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/store"
)

type storeAssetBiz struct {
	orm      *ent.StoreClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewStoreAsset() *storeAssetBiz {
	return &storeAssetBiz{
		orm: ent.Database.Store,
		ctx: context.Background(),
	}
}

func NewStoreAssetWithModifier(m *model.Modifier) *storeAssetBiz {
	s := NewStoreAsset()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// Assets 资产列表
func (b *storeAssetBiz) Assets(req *definition.StoreAssetListReq) (res *model.PaginationRes) {
	// 查询分页的门店数据
	q := b.orm.QueryNotDeleted().WithCity()
	// if req.GroupID != nil {
	// 	// q.Where(store.GroupID(*req.GroupID))
	// }
	if req.StoreID != nil {
		q.Where(store.ID(*req.StoreID))
	}
	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Store) (result *definition.StoreAssetDetail) {
		result = &definition.StoreAssetDetail{
			ID:         item.ID,
			Name:       item.Name,
			Lng:        item.Lng,
			Lat:        item.Lat,
			StoreAsset: b.assetForStore(req, item.ID),
		}
		if item.Edges.City != nil {
			result.City = model.City{
				ID:   item.Edges.City.ID,
				Name: item.Edges.City.Name,
			}
		}
		return result
	})

	return res
}

func (b *storeAssetBiz) assetForStore(req *definition.StoreAssetListReq, wId uint64) (asset definition.StoreAsset) {
	// 查询所属资产数据
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			entasset.LocationsType(model.AssetLocationsTypeStore.Value()),
			entasset.LocationsIDIn(wId),
			entasset.Status(model.AssetStatusStock.Value()),
		)
	if req.CityID != nil {
		q.Where(entasset.CityID(*req.CityID))
	}

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
