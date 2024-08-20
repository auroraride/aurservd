// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-20, by aurb

package biz

import (
	"context"

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

// Flow 资产流转明细
func (b *assetTransferBiz) Flow(req *model.AssetTransferFlowReq) []*model.AssetTransferFlow {
	return service.NewAssetTransfer().Flow(b.ctx, req)
}
