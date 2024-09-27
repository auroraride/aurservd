// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-07, by aurb

package biz

import (
	"sort"
	"strings"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
)

type assetTransferDetailsBiz struct {
	orm *ent.AssetTransferDetailsClient
}

func NewAssetTransferDetails() *assetTransferDetailsBiz {
	return &assetTransferDetailsBiz{
		orm: ent.Database.AssetTransferDetails,
	}
}

// InOutCount 调拨详情出入库统计
func (s *assetTransferDetailsBiz) InOutCount(items map[string]*definition.AssetMaterial, key string, outTransfer bool, id uint64, scrap bool) {
	if _, ok := items[key]; !ok {
		items[key] = &definition.AssetMaterial{
			ID:       id,
			Name:     key,
			Outbound: 0,
			Inbound:  0,
			Surplus:  0,
		}
	}

	// 判断出库还是入库
	if outTransfer {
		items[key].Outbound += 1
		items[key].Surplus -= 1

	} else {
		items[key].Inbound += 1
		items[key].Surplus += 1
	}

	if scrap {
		items[key].Scrap += 1
		items[key].Surplus -= 1
	}

}

// TransferInOut 物资出入库统计
func (b *assetTransferDetailsBiz) TransferInOut(ebikeNameMap, sBNameMap, nSbNameMap,
	cabAccNameMap, ebikeAccNameMap, otherAccNameMap map[string]*definition.AssetMaterial,
	ats []*ent.AssetTransferDetails, outTrans bool) {
	for _, inAt := range ats {
		ws := inAt.Edges.Asset
		if ws != nil {
			scrap := ws.Status == model.AssetStatusScrap.Value()

			switch ws.Type {
			case model.AssetTypeEbike.Value():
				if ws.Edges.Brand != nil {
					brandName := ws.Edges.Brand.Name
					b.InOutCount(ebikeNameMap, brandName, outTrans, ws.Edges.Brand.ID, scrap)
				}
			case model.AssetTypeSmartBattery.Value():
				if ws.Edges.Model != nil {
					modelName := ws.Edges.Model.Model
					b.InOutCount(sBNameMap, modelName, outTrans, ws.Edges.Model.ID, scrap)
				}
			case model.AssetTypeNonSmartBattery.Value():
				if ws.Edges.Model != nil {
					modelName := ws.Edges.Model.Model
					b.InOutCount(nSbNameMap, modelName, outTrans, ws.Edges.Model.ID, scrap)
				}
			case model.AssetTypeCabinetAccessory.Value():
				if ws.Edges.Material != nil {
					materialName := ws.Edges.Material.Name
					b.InOutCount(cabAccNameMap, materialName, outTrans, ws.Edges.Material.ID, scrap)
				}
			case model.AssetTypeEbikeAccessory.Value():
				if ws.Edges.Material != nil {
					materialName := ws.Edges.Material.Name
					b.InOutCount(ebikeAccNameMap, materialName, outTrans, ws.Edges.Material.ID, scrap)
				}
			case model.AssetTypeOtherAccessory.Value():
				if ws.Edges.Material != nil {
					materialName := ws.Edges.Material.Name
					b.InOutCount(otherAccNameMap, materialName, outTrans, ws.Edges.Material.ID, scrap)
				}
			}
		}
	}
}

// GetCommonAssetDetail 组装物资出入详情
func (b *assetTransferDetailsBiz) GetCommonAssetDetail(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap,
	otherAccNameMap map[string]*definition.AssetMaterial, ast *definition.CommonAssetDetail) {
	if ast == nil {
		return
	}
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
}
