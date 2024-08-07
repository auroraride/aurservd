// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-07, by Jorjan

package definition

type AssetMaterial struct {
	// Type      model.AssetType `json:"type"`                // 物资类别
	// ID        uint64          `json:"id"`                  // 物资ID
	Name      string `json:"name"`                // 物资名称
	Outbound  int    `json:"outbound"`            // 出库数量
	Inbound   int    `json:"inbound"`             // 入库数量
	Surplus   int    `json:"surplus"`             // 剩余
	Exception int    `json:"exception,omitempty"` // 异常数量(电柜无)
}
