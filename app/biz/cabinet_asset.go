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
	"github.com/auroraride/aurservd/internal/ent/cabinet"
)

type cabinetAssetBiz struct {
	orm      *ent.CabinetClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewCabinetAsset() *cabinetAssetBiz {
	return &cabinetAssetBiz{
		orm: ent.Database.Cabinet,
		ctx: context.Background(),
	}
}

func NewCabinetAssetWithModifier(m *model.Modifier) *cabinetAssetBiz {
	s := NewCabinetAsset()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
}

// Assets 资产列表
func (b *cabinetAssetBiz) Assets(req *definition.CabinetAssetListReq) (res *model.PaginationRes) {
	// 查询分页的门店数据
	q := b.orm.QueryNotDeleted().WithCity()
	if req.CityID != nil {
		q.Where(cabinet.CityID(*req.CityID))
	}
	if req.Name != nil {
		q.Where(cabinet.NameContains(*req.Name))
	}
	if req.Sn != nil {
		q.Where(cabinet.SnContains(*req.Sn))
	}
	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Cabinet) (result *definition.CabinetAssetDetail) {
		result = &definition.CabinetAssetDetail{
			ID:           item.ID,
			Name:         item.Name,
			Sn:           item.Sn,
			CabinetAsset: b.assetForCabinet(req, item.ID),
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

func (b *cabinetAssetBiz) assetForCabinet(req *definition.CabinetAssetListReq, id uint64) (asset definition.CabinetAsset) {
	// 查询所属资产数据
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			entasset.LocationsType(model.AssetLocationsTypeCabinet.Value()),
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

	list, _ := q.All(b.ctx)
	// 按照分类进行累计
	for _, v := range list {
		switch v.Type {
		case model.AssetTypeSmartBattery.Value():
			asset.SmartBatteryTotal = asset.SmartBatteryTotal + 1
		case model.AssetTypeNonSmartBattery.Value():
			asset.NonSmartBatteryTotal = asset.NonSmartBatteryTotal + 1
		}
	}

	// todo 新调拨记录计算

	return asset
}
