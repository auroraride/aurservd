// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-21, by aurb

package biz

import (
	"context"

	"github.com/auroraride/aurservd/app/biz/definition"
	"github.com/auroraride/aurservd/app/model"
	"github.com/auroraride/aurservd/app/service"
	"github.com/auroraride/aurservd/internal/ent"
	"github.com/auroraride/aurservd/pkg/silk"
)

type assetCheckBiz struct {
	orm *ent.AssetCheckClient
	ctx context.Context
}

func NewAssetCheck() *assetCheckBiz {
	return &assetCheckBiz{
		orm: ent.Database.AssetCheck,
		ctx: context.Background(),
	}
}

// GetAssetBySN 通过sn查询资产
func (b *assetCheckBiz) GetAssetBySN(am *ent.AssetManager, ep *ent.Employee, req *definition.AssetCheckByAssetSnReq) (res *model.AssetCheckByAssetSnRes, err error) {
	newReq := model.AssetCheckByAssetSnReq{
		SN: req.SN,
	}

	if am != nil {
		wType := model.AssetOperateRoleTypeManager
		newReq.OpratorType = wType
		newReq.OpratorID = am.ID
	}
	if ep != nil {
		sType := model.AssetOperateRoleTypeStore
		newReq.OpratorType = sType
		newReq.OpratorID = ep.ID
	}

	return service.NewAssetCheck().GetAssetBySN(b.ctx, &newReq)
}

// Create 创建资产盘点
func (b *assetCheckBiz) Create(am *ent.AssetManager, ep *ent.Employee, req *definition.AssetCheckCreateReq) (res *definition.AssetCheckCreateRes, err error) {
	var md model.Modifier

	newReq := model.AssetCheckCreateReq{
		AssetCheckCreateDetail: req.AssetCheckCreateDetail,
		StartAt:                req.StartAt,
		EndAt:                  req.EndAt,
	}

	if am != nil {
		wType := model.AssetLocationsTypeWarehouse
		newReq.LocationsType = wType
		// todo 查询上班所属位置
		// newReq.LocationsID = 1
		newReq.OpratorID = am.ID
		newReq.OpratorType = model.AssetOperateRoleTypeManager

		md = model.Modifier{
			ID:    am.ID,
			Name:  am.Name,
			Phone: am.Phone,
		}
	}
	if ep != nil {
		sType := model.AssetLocationsTypeStore
		newReq.LocationsType = sType
		// todo 查询上班所属位置
		// newReq.LocationsID = 1
		newReq.OpratorID = ep.ID
		newReq.OpratorType = model.AssetOperateRoleTypeStore
		md = model.Modifier{
			ID:    ep.ID,
			Name:  ep.Name,
			Phone: ep.Phone,
		}
	}

	cId, err := service.NewAssetCheck().CreateAssetCheck(b.ctx, &newReq, &md)
	if err != nil {
		return nil, err
	}

	return &definition.AssetCheckCreateRes{ID: cId}, nil
}

// List 盘点记录
func (b *assetCheckBiz) List(am *ent.AssetManager, ep *ent.Employee, req *definition.AssetCheckListReq) (res *model.PaginationRes, err error) {
	newReq := model.AssetCheckListReq{
		PaginationReq: req.PaginationReq,
		Keyword:       req.Keyword,
		StartAt:       req.StartAt,
		EndAt:         req.EndAt,
		CheckResult:   req.CheckResult,
	}

	if am != nil {
		wType := model.AssetLocationsTypeWarehouse
		newReq.LocationsType = &wType
		newReq.LocationsID = silk.UInt64(am.ID)
	}
	if ep != nil {
		sType := model.AssetLocationsTypeStore
		newReq.LocationsType = &sType
		newReq.LocationsID = silk.UInt64(ep.ID)
	}

	return service.NewAssetCheck().List(b.ctx, &newReq)
}
