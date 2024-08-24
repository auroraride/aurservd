// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-02, by aurb

package biz

import (
	"context"
	"sort"
	"strings"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	entasset "github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/assettransferdetails"
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
			ID:    item.ID,
			Name:  item.Name,
			Phone: item.Phone,
			Total: b.AssetTotal(req, item.ID),
		}
		return result
	})

	return res
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
			entasset.TypeIn(model.AssetTypeSmartBattery.Value(), model.AssetTypeNonSmartBattery.Value()),
		)
	}
	if req.BrandID != nil {
		q.Where(
			entasset.BrandID(*req.BrandID),
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

	b.transferInOut(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, inAts, false)

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

	b.transferInOut(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, outAts, true)

	// 组装出入库数据
	for _, v := range ebikeNameMap {
		ast.Ebikes = append(ast.Ebikes, v)
	}
	for _, v := range sBNameMap {
		ast.SmartBatteries = append(ast.SmartBatteries, v)
	}
	for _, v := range nSbNameMap {
		ast.NonSmartBatteries = append(ast.NonSmartBatteries, v)
	}
	for _, v := range cabAccNameMap {
		ast.CabinetAccessories = append(ast.CabinetAccessories, v)
	}
	for _, v := range ebikeAccNameMap {
		ast.EbikeAccessories = append(ast.EbikeAccessories, v)
	}
	for _, v := range otherAccNameMap {
		ast.OtherAssets = append(ast.OtherAssets, v)
	}

	// 排序
	sort.Slice(ast.Ebikes, func(i, j int) bool {
		return strings.Compare(ast.Ebikes[i].Name, ast.Ebikes[j].Name) < 0
	})
	sort.Slice(ast.SmartBatteries, func(i, j int) bool {
		return strings.Compare(ast.SmartBatteries[i].Name, ast.SmartBatteries[j].Name) < 0
	})
	sort.Slice(ast.NonSmartBatteries, func(i, j int) bool {
		return strings.Compare(ast.NonSmartBatteries[i].Name, ast.NonSmartBatteries[j].Name) < 0
	})
	sort.Slice(ast.CabinetAccessories, func(i, j int) bool {
		return strings.Compare(ast.CabinetAccessories[i].Name, ast.CabinetAccessories[j].Name) < 0
	})
	sort.Slice(ast.EbikeAccessories, func(i, j int) bool {
		return strings.Compare(ast.EbikeAccessories[i].Name, ast.EbikeAccessories[j].Name) < 0
	})
	sort.Slice(ast.OtherAssets, func(i, j int) bool {
		return strings.Compare(ast.OtherAssets[i].Name, ast.OtherAssets[j].Name) < 0
	})
	return
}

// transferInOut 物资出入库统计
func (b *maintainerAssetBiz) transferInOut(ebikeNameMap, sBNameMap, nSbNameMap,
	cabAccNameMap, ebikeAccNameMap, otherAccNameMap map[string]*definition.AssetMaterial,
	ats []*ent.AssetTransferDetails, outTrans bool) {
	for _, inAt := range ats {
		ws := inAt.Edges.Asset
		if ws != nil {
			switch ws.Type {
			case model.AssetTypeEbike.Value():
				if ws.Edges.Brand != nil {
					brandName := ws.Edges.Brand.Name
					NewAssetTransferDetails().InOutCount(ebikeNameMap, brandName, outTrans)
				}
			case model.AssetTypeSmartBattery.Value():
				if ws.Edges.Model != nil {
					modelName := ws.Edges.Model.Model
					NewAssetTransferDetails().InOutCount(sBNameMap, modelName, outTrans)
				}
			case model.AssetTypeNonSmartBattery.Value():
				if ws.Edges.Model != nil {
					modelName := ws.Edges.Model.Model
					NewAssetTransferDetails().InOutCount(nSbNameMap, modelName, outTrans)
				}
			case model.AssetTypeCabinetAccessory.Value():
				if ws.Edges.Material != nil {
					materialName := ws.Edges.Material.Name
					NewAssetTransferDetails().InOutCount(cabAccNameMap, materialName, outTrans)
				}
			case model.AssetTypeEbikeAccessory.Value():
				if ws.Edges.Material != nil {
					materialName := ws.Edges.Material.Name
					NewAssetTransferDetails().InOutCount(ebikeAccNameMap, materialName, outTrans)
				}
			case model.AssetTypeOtherAccessory.Value():
				if ws.Edges.Material != nil {
					materialName := ws.Edges.Material.Name
					NewAssetTransferDetails().InOutCount(otherAccNameMap, materialName, outTrans)
				}
			}
		}
	}
}
