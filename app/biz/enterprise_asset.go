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
	"github.com/auroraride/aurservd/internal/ent/enterprise"
	"github.com/auroraride/aurservd/internal/ent/enterprisestation"
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
	// 查询分页的门店数据
	q := b.orm.QueryNotDeleted().WithCity().WithStations()
	if req.CityID != nil {
		q.Where(enterprise.CityID(*req.CityID))
	}
	if req.StationID != nil {
		q.Where(enterprise.HasStationsWith(enterprisestation.ID(*req.StationID)))
	}
	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Enterprise) (result *definition.EnterpriseAssetDetail) {
		result = &definition.EnterpriseAssetDetail{
			ID:              item.ID,
			Name:            item.Name,
			City:            model.City{},
			EnterpriseAsset: b.assetForEnterprise(req, item.ID),
		}
		if item.Edges.City != nil {
			result.City = model.City{
				ID:   item.Edges.City.ID,
				Name: item.Edges.City.Name,
			}
		}
		stations := make([]*definition.EnterpriseStation, 0)
		for _, s := range item.Edges.Stations {
			stations = append(stations, &definition.EnterpriseStation{
				ID:   s.ID,
				Name: s.Name,
			})
		}

		result.Stations = stations
		return result
	})

	return res
}

func (b *enterpriseAssetBiz) assetForEnterprise(req *definition.EnterpriseAssetListReq, id uint64) (asset definition.EnterpriseAsset) {
	// 查询所属资产数据
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			entasset.LocationsType(model.AssetLocationsTypeStation.Value()),
			entasset.LocationsIDIn(id),
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
