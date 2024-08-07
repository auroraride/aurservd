// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-02, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

// MaintainerAssetListReq 运维资产列表请求
type MaintainerAssetListReq struct {
	model.PaginationReq
	Keyword   *string `json:"keyword" query:"keyword"`     // 关键字 姓名，手机号
	ModelID   *uint64 `json:"modelID" query:"modelID"`     // 电池型号ID
	BrandId   *uint64 `json:"brandId" query:"brandId"`     // 电车型号ID
	OtherName *string `json:"otherName" query:"otherName"` // 其他物资名称
}

// MaintainerAssetDetail 运维资产信息
type MaintainerAssetDetail struct {
	ID              uint64          `json:"id"`             // 运维ID
	Name            string          `json:"name"`           // 运维名称
	MaintainerAsset MaintainerAsset `json:"warehouseAsset"` // 运维资产
}

// MaintainerAsset 运维资产
type MaintainerAsset struct {
	EbikeTotal            int              `json:"ebikeTotal"`            // 电车总数
	Ebikes                []*AssetMaterial `json:"ebikes"`                // 电车物资详情
	SmartBatteryTotal     int              `json:"smartBatteryTotal"`     // 智能电池总数
	SmartBatteries        []*AssetMaterial `json:"smartBatteries"`        // 智能电池物资详情
	NonSmartBatteryTotal  int              `json:"nonSmartBatteryTotal"`  // 非智能电池总数
	NonSmartBatteries     []*AssetMaterial `json:"nonSmartBatteries"`     // 非智能电池物资详情
	CabinetAccessoryTotal int              `json:"cabinetAccessoryTotal"` // 电柜配件总数
	CabinetAccessories    []*AssetMaterial `json:"cabinetAccessories"`    // 电柜配件物资详情
	EbikeAccessoryTotal   int              `json:"ebikeAccessoryTotal"`   // 电车配件总数
	EbikeAccessories      []*AssetMaterial `json:"ebikeAccessories"`      // 电车配件物资详情
	OtherAssetTotal       int              `json:"otherAssetTotal"`       // 其他物资总数
	OtherAssets           []*AssetMaterial `json:"otherAssets"`           // 其他物资详情
}
