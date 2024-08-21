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
)

type assetTransferBiz struct {
	orm *ent.AssetTransferClient
	ctx context.Context
}

func NewAssetTransfer() *assetTransferBiz {
	return &assetTransferBiz{
		orm: ent.Database.AssetTransfer,
		ctx: context.Background(),
	}
}

// Transfer 创建资产调拨
func (b *assetTransferBiz) Transfer(am *ent.AssetManager, ep *ent.Employee, req *definition.AssetTransferCreateReq) (failed []string, err error) {

	var md model.Modifier

	aType := model.AssetTransferTypeTransfer
	newReq := model.AssetTransferCreateReq{
		ToLocationType:    req.ToLocationType,
		ToLocationID:      req.ToLocationID,
		Details:           req.Details,
		Reason:            req.Reason,
		AssetTransferType: &aType,
	}

	if am != nil {
		wType := model.AssetLocationsTypeWarehouse
		newReq.FromLocationType = &wType
		newReq.FromLocationID = req.FromLocationID
		md = model.Modifier{
			ID:    am.ID,
			Name:  am.Name,
			Phone: am.Phone,
		}
	}
	if ep != nil {
		sType := model.AssetLocationsTypeStore
		newReq.FromLocationType = &sType
		newReq.FromLocationID = req.FromLocationID
		md = model.Modifier{
			ID:    ep.ID,
			Name:  ep.Name,
			Phone: ep.Phone,
		}
	}

	return service.NewAssetTransfer().Transfer(b.ctx, &newReq, &md)
}

// Flow 资产流转明细
func (b *assetTransferBiz) Flow(req *model.AssetTransferFlowReq) []*model.AssetTransferFlow {
	return service.NewAssetTransfer().Flow(b.ctx, req)
}

// GetTransferBySn 扫码入库根据Sn获取调拨信息
func (b *assetTransferBiz) GetTransferBySn(req *model.GetTransferBySNReq) (res *model.AssetTransferListRes, err error) {
	return service.NewAssetTransfer().GetTransferBySN(b.ctx, req)
}
