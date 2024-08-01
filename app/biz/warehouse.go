// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-07-10, by Jorjan

package biz

import (
	"context"
	"errors"
	"fmt"

	"github.com/lithammer/shortuuid/v4"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/internal/ent"
	entasset "github.com/auroraride/aurservd/internal/ent/asset"
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
	s := NewWarehouse()
	if m != nil {
		s.ctx = context.WithValue(s.ctx, model.CtxModifierKey{}, m)
		s.modifier = m
	}
	return s
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
			WarehouseAsset: b.assetForWarehouse(req, item.ID),
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

func (b *warehouseBiz) assetForWarehouse(req *definition.WareHouseAssetListReq, wIds uint64) (asset definition.WarehouseAsset) {
	// 查询仓库所属资产数据
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			entasset.LocationsType(model.AssetLocationsTypeWarehouse.Value()),
			entasset.LocationsIDIn(wIds),
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
