// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-01, by aurb

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	entasset "github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/assettransferdetails"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/internal/ent/store"
	"github.com/auroraride/aurservd/pkg/tools"
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
	q := b.orm.QueryNotDeleted().
		WithCity().
		WithGroup().
		Order(ent.Desc(store.FieldCreatedAt))
	b.assetsFilter(q, req)
	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Store) (result *definition.StoreAssetDetail) {
		result = &definition.StoreAssetDetail{
			ID:    item.ID,
			Name:  item.Name,
			Lng:   item.Lng,
			Lat:   item.Lat,
			Total: b.AssetTotal(req, item.ID),
		}
		if item.Edges.City != nil {
			result.City = model.City{
				ID:   item.Edges.City.ID,
				Name: item.Edges.City.Name,
			}
		}
		if item.Edges.Group != nil {
			result.GroupName = item.Edges.Group.Name
		}
		return result
	})

	return res
}

// assetsFilter 条件筛选
func (b *storeAssetBiz) assetsFilter(q *ent.StoreQuery, req *definition.StoreAssetListReq) {
	if req.CityID != nil {
		q.Where(store.CityID(*req.CityID))
	}
	if req.GroupID != nil {
		q.Where(store.GroupID(*req.GroupID))
	}
	if req.StoreID != nil {
		q.Where(store.ID(*req.StoreID))
	}
	if req.ModelID != nil {
		q.Where(
			store.HasAssetWith(entasset.ModelID(*req.ModelID)),
		)
	}
	if req.BrandId != nil {
		q.Where(
			store.HasAssetWith(entasset.BrandID(*req.BrandId)),
		)
	}
	if req.OtherName != nil {
		q.Where(
			store.HasAssetWith(entasset.HasMaterialWith(material.NameContains(*req.OtherName))),
		)
	}
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(store.CreatedAtGTE(start), store.CreatedAtLTE(end))
	}

}

// AssetTotal 物资数据统计
func (b *storeAssetBiz) AssetTotal(req *definition.StoreAssetListReq, sId uint64) (res definition.CommonAssetTotal) {
	// 查询所属资产数据
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			entasset.LocationsType(model.AssetLocationsTypeStore.Value()),
			entasset.LocationsIDIn(sId),
			entasset.Status(model.AssetStatusStock.Value()),
		)
	if req.ModelID != nil {
		q.Where(
			entasset.ModelID(*req.ModelID),
		)
	}
	if req.BrandId != nil {
		q.Where(
			entasset.BrandID(*req.BrandId),
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
func (b *storeAssetBiz) AssetDetail(id uint64) (ast *definition.CommonAssetDetail) {
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
				assettransfer.ToLocationType(model.AssetLocationsTypeStore.Value()),
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
				assettransfer.FromLocationType(model.AssetLocationsTypeStore.Value()),
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

// AssetsExport 资产列表导出
func (b *storeAssetBiz) AssetsExport(req *definition.StoreAssetListReq) model.AssetExportRes {
	// 查询分页的门店数据
	q := b.orm.QueryNotDeleted().
		WithCity().
		WithGroup().
		Order(ent.Desc(store.FieldCreatedAt))
	b.assetsFilter(q, req)

	return service.NewAssetExportWithModifier(b.modifier).Start("门店物资", req.StoreAssetListFilter, nil, "", func(path string) {
		items, _ := q.All(b.ctx)
		var rows tools.ExcelItems
		title := []any{
			"门店",
			"门店集合",
			"城市",
			"智能电池",
			"非智能电池",
			"电车",
			"电车配件",
			"电柜配件",
			"其他物资",
		}
		rows = append(rows, title)
		for _, item := range items {
			assetTotal := b.AssetTotal(req, item.ID)

			var cityName, groupName string
			if item.Edges.City != nil {
				cityName = item.Edges.City.Name
			}
			if item.Edges.Group != nil {
				groupName = item.Edges.Group.Name
			}

			row := []any{
				item.Name,
				groupName,
				cityName,
				assetTotal.SmartBatteryTotal,
				assetTotal.NonSmartBatteryTotal,
				assetTotal.EbikeTotal,
				assetTotal.EbikeAccessoryTotal,
				assetTotal.CabinetAccessoryTotal,
				assetTotal.OtherAssetTotal,
			}
			rows = append(rows, row)
		}
		tools.NewExcel(path).AddValues(rows).Done()
	})
}
