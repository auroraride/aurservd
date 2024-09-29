// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-10, by Jorjan

package biz

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lithammer/shortuuid/v4"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	entasset "github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assetmanager"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/assettransferdetails"
	"github.com/auroraride/aurservd/internal/ent/material"
	"github.com/auroraride/aurservd/internal/ent/warehouse"
	"github.com/auroraride/aurservd/pkg/snag"
	"github.com/auroraride/aurservd/pkg/tools"
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
		Remark:  item.Remark,
	}

	if item.Edges.City != nil {
		res.City = model.City{
			ID:   item.Edges.City.ID,
			Name: item.Edges.City.Name,
		}
		res.CityID = item.Edges.City.ID
		res.CityName = item.Edges.City.Name
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
	q := b.orm.QueryNotDeleted().
		WithCity().
		Order(ent.Desc(warehouse.FieldCreatedAt))
	b.assetsFilter(q, req)
	res = model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Warehouse) (result *definition.WareHouseAssetDetail) {
		result = &definition.WareHouseAssetDetail{
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
		return result
	})

	return res
}

// assetsFilter 条件筛选
func (b *warehouseBiz) assetsFilter(q *ent.WarehouseQuery, req *definition.WareHouseAssetListReq) {
	if req.Name != nil {
		q.Where(warehouse.NameContainsFold(*req.Name))
	}
	if req.CityID != nil {
		q.Where(warehouse.CityID(*req.CityID))
	}
	if req.ModelID != nil {
		q.Where(
			warehouse.HasAssetWith(entasset.ModelID(*req.ModelID)),
		)
	}
	if req.BrandId != nil {
		q.Where(
			warehouse.HasAssetWith(entasset.BrandID(*req.BrandId)),
		)
	}
	if req.OtherName != nil {
		q.Where(
			warehouse.HasAssetWith(entasset.HasMaterialWith(material.NameContains(*req.OtherName))),
		)
	}
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(warehouse.CreatedAtGTE(start), warehouse.CreatedAtLTE(end))
	}

}

// AssetTotal 仓库物资数据统计
func (b *warehouseBiz) AssetTotal(req *definition.WareHouseAssetListReq, wId uint64) (res definition.CommonAssetTotal) {
	// 查询仓库所属资产数据
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			entasset.LocationsType(model.AssetLocationsTypeWarehouse.Value()),
			entasset.LocationsIDIn(wId),
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

// AssetsDetail 物资详情
func (b *warehouseBiz) AssetsDetail(wId uint64) (ast *definition.CommonAssetDetail) {
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
				assettransfer.ToLocationType(model.AssetLocationsTypeWarehouse.Value()),
				assettransfer.ToLocationID(wId),
				assettransfer.DeletedAtIsNil(),
			),
		).
		WithAsset(func(query *ent.AssetQuery) {
			query.WithBrand().WithModel().WithMaterial()
		}).All(b.ctx)

	NewAssetTransferDetails().TransferInOut(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, inAts, false, wId)

	// 出库物资调拨详情
	outAts, _ := ent.Database.AssetTransferDetails.QueryNotDeleted().
		Where(
			assettransferdetails.HasTransferWith(
				assettransfer.StatusIn(model.AssetTransferStatusDelivering.Value(), model.AssetTransferStatusStock.Value()),
				assettransfer.FromLocationType(model.AssetLocationsTypeWarehouse.Value()),
				assettransfer.FromLocationID(wId),
				assettransfer.DeletedAtIsNil(),
			),
		).
		WithAsset(func(query *ent.AssetQuery) {
			query.WithBrand().WithModel().WithMaterial()
		}).All(b.ctx)

	NewAssetTransferDetails().TransferInOut(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, outAts, true, wId)

	NewAssetTransferDetails().GetCommonAssetDetail(ebikeNameMap, sBNameMap, nSbNameMap, cabAccNameMap, ebikeAccNameMap, otherAccNameMap, ast)
	return
}

// ListByCity 城市仓库列表
func (b *warehouseBiz) ListByCity() (res []*model.CascaderOptionLevel2) {
	res = make([]*model.CascaderOptionLevel2, 0)
	whList, _ := b.orm.QueryNotDeleted().WithCity().Order(ent.Asc(warehouse.FieldID)).All(b.ctx)
	cityIds := make([]uint64, 0)
	cityIdMap := make(map[uint64]*ent.City)
	cityIdListMap := make(map[uint64][]model.SelectOption)
	for _, wh := range whList {
		if wh.Edges.City != nil {
			cId := wh.Edges.City.ID
			if cityIdMap[cId] == nil {
				cityIds = append(cityIds, cId)
				cityIdMap[cId] = wh.Edges.City
			}

			cityIdListMap[cId] = append(cityIdListMap[cId], model.SelectOption{
				Label: wh.Name,
				Value: wh.ID,
			})
		}
	}

	for _, cityId := range cityIds {
		if cityIdMap[cityId] != nil && len(cityIdListMap[cityId]) != 0 {
			res = append(res, &model.CascaderOptionLevel2{
				SelectOption: model.SelectOption{
					Value: cityIdMap[cityId].ID,
					Label: cityIdMap[cityId].Name,
				},
				Children: cityIdListMap[cityId],
			})
		}
	}

	return
}

// SelectionList 仓库列表筛选项
func (b *warehouseBiz) SelectionList() (res []*model.SelectOption) {
	res = make([]*model.SelectOption, 0)
	whList, _ := b.orm.QueryNotDeleted().Order(ent.Asc(warehouse.FieldID)).All(b.ctx)

	for _, wh := range whList {
		res = append(res, &model.SelectOption{
			Value: wh.ID,
			Label: wh.Name,
		})
	}
	return
}

// ListByManager 仓管已配置仓库列表
func (b *warehouseBiz) ListByManager(am *ent.AssetManager) (res []*model.CascaderOptionLevel2) {
	res = make([]*model.CascaderOptionLevel2, 0)
	if am == nil {
		return
	}

	// 查询人员配置的仓库城市信息
	eam, err := ent.Database.AssetManager.QueryNotDeleted().
		Where(assetmanager.ID(am.ID)).
		WithBelongWarehouses(func(query *ent.WarehouseQuery) {
			query.Where(warehouse.DeletedAtIsNil()).WithCity()
		}).First(b.ctx)
	if err != nil || eam == nil {
		return
	}

	// 数据组合
	whList := eam.Edges.BelongWarehouses
	cityIds := make([]uint64, 0)
	cityIdMap := make(map[uint64]*ent.City)
	cityIdListMap := make(map[uint64][]model.SelectOption)
	for _, wh := range whList {
		if wh.Edges.City != nil {
			cId := wh.Edges.City.ID
			if cityIdMap[cId] == nil {
				cityIds = append(cityIds, cId)
				cityIdMap[cId] = wh.Edges.City
			}

			cityIdListMap[cId] = append(cityIdListMap[cId], model.SelectOption{
				Label: wh.Name,
				Value: wh.ID,
			})
		}
	}

	for _, cityId := range cityIds {
		if cityIdMap[cityId] != nil && len(cityIdListMap[cityId]) != 0 {
			res = append(res, &model.CascaderOptionLevel2{
				SelectOption: model.SelectOption{
					Value: cityIdMap[cityId].ID,
					Label: cityIdMap[cityId].Name,
				},
				Children: cityIdListMap[cityId],
			})
		}
	}

	return
}

// QuerySn 通过SN查询仓库
func (b *warehouseBiz) QuerySn(sn string) *ent.Warehouse {
	if strings.HasPrefix(sn, "WAREHOUSE:") {
		sn = strings.ReplaceAll(sn, "WAREHOUSE:", "")
	}
	item, err := b.orm.QueryNotDeleted().Where(warehouse.Sn(sn)).First(b.ctx)
	if err != nil {
		snag.Panic("未找到有效仓库")
	}
	return item
}
