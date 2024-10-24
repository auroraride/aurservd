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
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	entasset "github.com/auroraride/aurservd/internal/ent/asset"
	"github.com/auroraride/aurservd/internal/ent/assettransfer"
	"github.com/auroraride/aurservd/internal/ent/assettransferdetails"
	"github.com/auroraride/aurservd/internal/ent/cabinet"
	"github.com/auroraride/aurservd/pkg/tools"
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
	q := b.orm.QueryNotDeleted().
		WithCity().
		Order(ent.Desc(cabinet.FieldCreatedAt))
	b.assetsFilter(q, req)
	return model.ParsePaginationResponse(q, req.PaginationReq, func(item *ent.Cabinet) (result *definition.CabinetAssetDetail) {
		result = &definition.CabinetAssetDetail{
			ID:    item.ID,
			Name:  item.Name,
			Sn:    item.Serial,
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
}

// assetsFilter 条件筛选
func (b *cabinetAssetBiz) assetsFilter(q *ent.CabinetQuery, req *definition.CabinetAssetListReq) {
	if req.CityID != nil {
		q.Where(cabinet.CityID(*req.CityID))
	}
	if req.ModelID != nil {
		q.Where(
			cabinet.HasAssetWith(entasset.ModelID(*req.ModelID)),
		)
	}
	if req.Name != nil {
		q.Where(cabinet.NameContainsFold(*req.Name))
	}
	if req.Sn != nil {
		q.Where(cabinet.SerialContainsFold(*req.Sn))
	}
	if req.Start != nil && req.End != nil {
		start := tools.NewTime().ParseDateStringX(*req.Start)
		end := tools.NewTime().ParseNextDateStringX(*req.End)
		q.Where(cabinet.CreatedAtGTE(start), cabinet.CreatedAtLTE(end))
	}

}

// AssetTotal 物资数据统计
func (b *cabinetAssetBiz) AssetTotal(req *definition.CabinetAssetListReq, id uint64) (res definition.CabinetTotal) {
	// 查询所属资产数据
	q := ent.Database.Asset.QueryNotDeleted().
		Where(
			entasset.LocationsType(model.AssetLocationsTypeCabinet.Value()),
			entasset.LocationsIDIn(id),
			entasset.Status(model.AssetStatusStock.Value()),
		)

	if req.ModelID != nil {
		q.Where(
			entasset.ModelID(*req.ModelID),
		)
	}

	list, _ := q.All(b.ctx)
	for _, v := range list {
		switch v.Type {
		case model.AssetTypeSmartBattery.Value():
			res.SmartBatteryTotal += 1
		case model.AssetTypeNonSmartBattery.Value():
			res.NonSmartBatteryTotal += 1
		}
	}
	return
}

// AssetDetail 物资详情
func (b *cabinetAssetBiz) AssetDetail(id uint64) (ast *definition.CabinetTotalDetail) {
	ast = &definition.CabinetTotalDetail{
		SmartBatteries:    make([]*definition.AssetMaterial, 0),
		NonSmartBatteries: make([]*definition.AssetMaterial, 0),
	}

	sBNameMap := make(map[string]*definition.AssetMaterial)
	nSbNameMap := make(map[string]*definition.AssetMaterial)

	// 入库物资调拨详情
	inAts, _ := ent.Database.AssetTransferDetails.QueryNotDeleted().
		Where(
			assettransferdetails.IsIn(true),
			assettransferdetails.HasTransferWith(
				assettransfer.Status(model.AssetTransferStatusStock.Value()),
				assettransfer.ToLocationType(model.AssetLocationsTypeCabinet.Value()),
				assettransfer.ToLocationID(id),
				assettransfer.DeletedAtIsNil(),
			),
		).
		WithAsset(func(query *ent.AssetQuery) {
			query.WithBrand().WithModel().WithMaterial()
		}).All(b.ctx)

	NewAssetTransferDetails().TransferInOut(nil, sBNameMap, nSbNameMap, nil, nil, nil, inAts, false, id)

	// 出库物资调拨详情
	outAts, _ := ent.Database.AssetTransferDetails.QueryNotDeleted().
		Where(
			assettransferdetails.HasTransferWith(
				assettransfer.StatusIn(model.AssetTransferStatusDelivering.Value(), model.AssetTransferStatusStock.Value()),
				assettransfer.FromLocationType(model.AssetLocationsTypeCabinet.Value()),
				assettransfer.FromLocationID(id),
				assettransfer.DeletedAtIsNil(),
			),
		).
		WithAsset(func(query *ent.AssetQuery) {
			query.WithBrand().WithModel().WithMaterial()
		}).All(b.ctx)

	NewAssetTransferDetails().TransferInOut(nil, sBNameMap, nSbNameMap, nil, nil, nil, outAts, true, id)

	// 组装出入库数据

	for _, v := range sBNameMap {
		ast.SmartBatteries = append(ast.SmartBatteries, v)
	}
	for _, v := range nSbNameMap {
		ast.NonSmartBatteries = append(ast.NonSmartBatteries, v)
	}

	// 排序
	sort.Slice(ast.SmartBatteries, func(i, j int) bool {
		return strings.Compare(ast.SmartBatteries[i].Name, ast.SmartBatteries[j].Name) < 0
	})
	sort.Slice(ast.NonSmartBatteries, func(i, j int) bool {
		return strings.Compare(ast.NonSmartBatteries[i].Name, ast.NonSmartBatteries[j].Name) < 0
	})
	return
}

// AssetsExport 资产列表导出
func (b *cabinetAssetBiz) AssetsExport(req *definition.CabinetAssetListReq) model.AssetExportRes {
	q := b.orm.QueryNotDeleted().
		WithCity().
		Order(ent.Desc(cabinet.FieldCreatedAt))
	b.assetsFilter(q, req)

	return service.NewAssetExportWithModifier(b.modifier).Start("电柜物资", req.CabinetAssetListFilter, nil, "", func(path string) {
		items, _ := q.All(b.ctx)
		var rows tools.ExcelItems
		title := []any{
			"城市",
			"电柜编号",
			"电柜名称",
			"智能电池",
			"非智能电池",
		}
		rows = append(rows, title)
		for _, item := range items {
			assetTotal := b.AssetTotal(req, item.ID)
			var cityName string
			if item.Edges.City != nil {
				cityName = item.Edges.City.Name
			}
			row := []any{
				cityName,
				item.Serial,
				item.Name,
				assetTotal.SmartBatteryTotal,
				assetTotal.NonSmartBatteryTotal,
			}
			rows = append(rows, row)
		}
		tools.NewExcel(path).AddValues(rows).Done()
	})
}
