// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-28, by aurb

package biz

import "github.com/auroraride/aurservd/internal/ent"

type assetMaintenanceBiz struct {
	orm *ent.AssetMaintenanceClient
}

func NewAssetMaintenance() *assetMaintenanceBiz {
	return &assetMaintenanceBiz{
		orm: ent.Database.AssetMaintenance,
	}
}
