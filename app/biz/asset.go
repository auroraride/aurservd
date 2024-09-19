// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-20, by aurb

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/silk"
)

type assetBiz struct {
	orm *ent.AssetClient
	ctx context.Context
}

func NewAsset() *assetBiz {
	return &assetBiz{
		orm: ent.Database.Asset,
		ctx: context.Background(),
	}
}

// List 资产列表
func (b *assetBiz) List(am *ent.AssetManager, ep *ent.Employee, req *definition.AssetListReq) (res *model.PaginationRes) {
	newReq := &model.AssetListReq{
		PaginationReq: req.PaginationReq,
		AssetFilter: model.AssetFilter{
			SN:      req.SN,
			ModelID: req.ModelID,
			Status:  req.Status,
			BrandID: req.BrandID,
		},
	}
	switch {
	case am != nil && ep == nil:
		// 确认为仓库管理员
		newReq.AssetFilter.AssetManagerID = am.ID
		lType := model.AssetLocationsTypeWarehouse
		newReq.AssetFilter.LocationsType = &lType

		switch req.Type {
		case definition.ReqTypeBattery:
			newReq.Battery = silk.Bool(true)
		case definition.ReqTypeEbike:
			aType := model.AssetTypeEbike
			newReq.AssetType = &aType
		}
		return service.NewAsset().List(b.ctx, newReq)
	case am == nil && ep != nil:
		// 确认为门店管理员
		newReq.AssetFilter.AssetManagerID = ep.ID
		lType := model.AssetLocationsTypeStore
		newReq.AssetFilter.LocationsType = &lType
		switch req.Type {
		case definition.ReqTypeBattery:
			newReq.Battery = silk.Bool(true)
		case definition.ReqTypeEbike:
			aType := model.AssetTypeEbike
			newReq.AssetType = &aType
		}
		return service.NewAsset().List(b.ctx, newReq)
	default:
		return
	}
}
