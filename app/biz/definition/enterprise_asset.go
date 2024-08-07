// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-02, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

// EnterpriseAssetListReq 团签资产列表请求
type EnterpriseAssetListReq struct {
	model.PaginationReq
	CityID    *uint64 `json:"cityId" query:"cityId"`       // 城市ID
	StationID *uint64 `json:"stationId" query:"stationId"` // 站点ID
	ModelID   *uint64 `json:"modelID" query:"modelID"`     // 电池型号ID
	BrandId   *uint64 `json:"brandId" query:"brandId"`     // 电车型号ID
	OtherName *string `json:"otherName" query:"otherName"` // 其他物资名称
	Start     *string `json:"start" query:"start"`         // 开始时间
	End       *string `json:"end" query:"end"`             // 结束时间
}

// EnterpriseAssetDetail 团签资产信息
type EnterpriseAssetDetail struct {
	ID              uint64               `json:"id"`             // 团签ID
	Name            string               `json:"name"`           // 团签名称
	Stations        []*EnterpriseStation `json:"stations"`       // 团签站点
	City            model.City           `json:"city"`           // 城市
	EnterpriseAsset EnterpriseAsset      `json:"warehouseAsset"` // 团签资产
}

type EnterpriseStation struct {
	ID   uint64 `json:"id"`   // 站点ID
	Name string `json:"name"` // 站点名称
}

// EnterpriseAsset 团签资产
type EnterpriseAsset struct {
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
