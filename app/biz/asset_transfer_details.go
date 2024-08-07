// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-07, by aurb

package biz

import (
	"github.com/auroraride/aurservd/app/biz/definition"
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
func (s *assetTransferDetailsBiz) InOutCount(items map[string]*definition.AssetMaterial, key string, outTransfer bool) {
	if _, ok := items[key]; !ok {
		items[key] = &definition.AssetMaterial{
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

}
