// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-10, by Jorjan

package biz

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/lithammer/shortuuid/v4"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	entasset "github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/assettransferdetails"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
)

type warehouseBiz struct {
	orm      *ent.WarehouseClient
	ctx      context.Context
	modifier *model.Modifier
}

func NewWarehouse() *warehouseBiz {
	return &warehouseBiz{
		orm: ent.Database.Warehouse,
		ctx: context.Background(),
	}
}

func NewWarehouseWithModifier(m *model.Modifier) *warehouseBiz {
	b := NewWarehouse()
	if m != nil {
		b.ctx = context.WithValue(b.ctx, model.CtxModifierKey{}, m)
		b.modifier = m
	}
	return b
}

// List 仓库列表
func (b *warehouseBiz) List(req *definition.WareHouseListReq) (res *model.PaginationRes) {
	q := b.orm.QueryNotDeleted().WithCity()

	if req.CityID != nil {
		q.Where(warehouse.CityID(*req.CityID))
	}
	if req.Keyword != nil {
		q.Where(warehouse.NameContains(*req.Keyword))
	}
	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Warehouse) (result *definition.WarehouseDetail) {
		return b.detail(item)
	})
	return
}

// detail 详情数据
func (b *warehouseBiz) detail(item *ent.Warehouse) (res *definition.WarehouseDetail) {
	res = &definition.WarehouseDetail{
		ID:      item.ID,
		Name:    item.Name,
		Lng:     item.Lng,
		Lat:     item.Lat,
		Address: item.Address,
		QRCode:  fmt.Sprintf("WAREHOUSE:%s", item.Sn),
	}

	if item.Edges.City != nil {
		res.City = model.City{
			ID:   item.Edges.City.ID,
			Name: item.Edges.City.Name,
		}
	}
	return
}

// Detail 查询仓库详情
func (b *warehouseBiz) Detail(id uint64) (*definition.WarehouseDetail, error) {
	g, _ := b.queryById(id)
	if g == nil {
		return nil, errors.New("仓库不存在")
	}
	return b.detail(g), nil
}

// Create 创建仓库
func (b *warehouseBiz) Create(req *definition.WarehouseCreateReq) (err error) {
	_, err = b.orm.Create().
		SetName(req.Name).
		SetCityID(req.CityID).
		SetAddress(req.Address).
		SetRemark(req.Remark).
		SetLat(req.Lat).
		SetLng(req.Lng).
		SetSn(shortuuid.New()).
		Save(b.ctx)
	if err != nil {
		return err
	}
	return
}

// queryById 通过ID查询仓库
func (b *warehouseBiz) queryById(id uint64) (item *ent.Warehouse, err error) {
	return b.orm.QueryNotDeleted().Where(warehouse.ID(id)).First(b.ctx)
}

// Modify 编辑仓库
func (b *warehouseBiz) Modify(req *definition.WarehouseModifyReq) (err error) {
	g, _ := b.queryById(req.ID)
	if g == nil {
		return errors.New("仓库不存在")
	}

	_, err = b.orm.UpdateOneID(req.ID).
		SetName(req.Name).
		SetName(req.Name).
		SetCityID(req.CityID).
		SetAddress(req.Address).
		SetRemark(req.Remark).
		SetLat(req.Lat).
		SetLng(req.Lng).
		Save(b.ctx)
	if err != nil {
		return err
	}

	return
}

// Delete 删除仓库
func (b *warehouseBiz) Delete(id uint64) (err error) {
	g, _ := b.queryById(id)
	if g == nil {
		return errors.New("仓库不存在")
	}
	_, err = b.orm.SoftDeleteOne(g).Save(b.ctx)
	if err != nil {
		return err
	}
	return
}

// Assets 仓库资产列表
func (b *warehouseBiz) Assets(req *definition.WareHouseAssetListReq) (res *model.PaginationRes) {
	// 查询分页的仓库数据
	q := b.orm.QueryNotDeleted().WithCity()
	if req.Name != nil {
		q.Where(warehouse.NameContains(*req.Name))
	}
	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Warehouse) (result *definition.WareHouseAssetDetail) {
		result = &definition.WareHouseAssetDetail{
			ID:             item.ID,
			Name:           item.Name,
			Lng:            item.Lng,
			Lat:            item.Lat,
			WarehouseAsset: b.AssetForWarehouse(req, item.ID),
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

// AssetForWarehouse 仓库物资数据
func (b *warehouseBiz) AssetForWarehouse(req *definition.WareHouseAssetListReq, wId uint64) definition.WarehouseAsset {
	astRes := definition.WarehouseAsset{
		Ebikes:             make([]*definition.AssetMaterial, 0),
		SmartBatteries:     make([]*definition.AssetMaterial, 0),
		NonSmartBatteries:  make([]*definition.AssetMaterial, 0),
		CabinetAccessories: make([]*definition.AssetMaterial, 0),
		EbikeAccessories:   make([]*definition.AssetMaterial, 0),
		OtherAssets:        make([]*definition.AssetMaterial, 0),
	}
	// 查询仓库所属资产数据
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			entasset.LocationsType(model.AssetLocationsTypeWarehouse.Value()),
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

	astIds := make([]uint64, 0)
	for _, v := range list {
		astIds = append(astIds, v.ID)
	}

	// 物资调拨详情
	b.assetTransferDetail(wId, astIds, &astRes)

	return astRes
}

// assetTransferDetail 物资调拨详情
func (b *warehouseBiz) assetTransferDetail(wId uint64, astIds []uint64, ast *definition.WarehouseAsset) {
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
				assettransfer.ToLocationType(model.AssetLocationsTypeWarehouse.Value()),
				assettransfer.ToLocationID(wId),
				assettransfer.DeletedAtIsNil(),
			),
			assettransferdetails.HasAssetWith(
				entasset.IDIn(astIds...),
				entasset.DeletedAtIsNil(),
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
				assettransfer.FromLocationType(model.AssetLocationsTypeWarehouse.Value()),
				assettransfer.FromLocationID(wId),
				assettransfer.DeletedAtIsNil(),
			),
			assettransferdetails.HasAssetWith(
				entasset.IDIn(astIds...),
				entasset.DeletedAtIsNil(),
			),
		).
		WithAsset(func(query *ent.AssetQuery) {
			query.WithBrand().WithModel().WithMaterial()
		}).All(b.ctx)

	b.transferInOut(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, outAts, true)

	// 组装出入库数据
	for _, v := range ebikeNameMap {
		ast.Ebikes = append(ast.Ebikes, v)
		ast.EbikeTotal += v.Surplus
	}
	for _, v := range sBNameMap {
		ast.SmartBatteries = append(ast.SmartBatteries, v)
		ast.SmartBatteryTotal += v.Surplus
	}
	for _, v := range nSbNameMap {
		ast.NonSmartBatteries = append(ast.NonSmartBatteries, v)
		ast.NonSmartBatteryTotal += v.Surplus
	}
	for _, v := range cabAccNameMap {
		ast.CabinetAccessories = append(ast.CabinetAccessories, v)
		ast.CabinetAccessoryTotal += v.Surplus
	}
	for _, v := range ebikeAccNameMap {
		ast.EbikeAccessories = append(ast.EbikeAccessories, v)
		ast.EbikeAccessoryTotal += v.Surplus
	}
	for _, v := range otherAccNameMap {
		ast.OtherAssets = append(ast.OtherAssets, v)
		ast.OtherAssetTotal += v.Surplus
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

// transferInOut 物资出入库统计
func (b *warehouseBiz) transferInOut(ebikeNameMap, sBNameMap, nSbNameMap,
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

// ListByCity 城市仓库列表
func (b *warehouseBiz) ListByCity() (res []*definition.WarehouseByCityRes) {
	whList, _ := b.orm.QueryNotDeleted().WithCity().Order(ent.Asc(warehouse.FieldID)).All(b.ctx)
	cityIds := make([]uint64, 0)
	cityIdMap := make(map[uint64]*ent.City)
	cityIdListMap := make(map[uint64][]*definition.WarehouseByCityDetail)
	for _, wh := range whList {
		if wh.Edges.City != nil {
			cityIds = append(cityIds, wh.CityID)
			cityIdMap[wh.CityID] = wh.Edges.City
			cityIdListMap[wh.CityID] = append(cityIdListMap[wh.CityID], &definition.WarehouseByCityDetail{
				ID:   wh.ID,
				Name: wh.Name,
			})
		}
	}

	for _, cityId := range cityIds {
		if cityIdMap[cityId] != nil && len(cityIdListMap[cityId]) != 0 {
			res = append(res, &definition.WarehouseByCityRes{
				City: model.City{
					ID:   cityIdMap[cityId].ID,
					Name: cityIdMap[cityId].Name,
				},
				WarehouseList: cityIdListMap[cityId],
			})
		}
	}

	return
}
