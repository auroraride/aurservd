// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-01, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

// StoreAssetListReq 门店资产列表请求
type StoreAssetListReq struct {
	model.PaginationReq
	CityID    *uint64 `json:"cityId" query:"cityId"`       // 城市ID
	GroupID   *uint64 `json:"groupId" query:"groupId"`     // 门店集合ID
	StoreID   *uint64 `json:"storeId" query:"storeId"`     // 门店ID
	ModelID   *uint64 `json:"modelID" query:"modelID"`     // 电池型号ID
	BrandId   *uint64 `json:"brandId" query:"brandId"`     // 电车型号ID
	OtherName *string `json:"otherName" query:"otherName"` // 其他物资名称
	Start     *string `json:"start" query:"start"`         // 开始时间
	End       *string `json:"end" query:"end"`             // 结束时间
}

// StoreAssetDetail 门店资产信息
type StoreAssetDetail struct {
	ID         uint64     `json:"id"`             // 门店ID
	Name       string     `json:"name"`           // 门店名称
	GroupName  string     `json:"groupName"`      // 门店集合名称
	City       model.City `json:"city"`           // 城市
	Lng        float64    `json:"lng"`            // 经度
	Lat        float64    `json:"lat"`            // 纬度
	StoreAsset StoreAsset `json:"warehouseAsset"` // 门店资产
}

// StoreAsset 门店资产
type StoreAsset struct {
	EbikeTotal            int                    `json:"ebikeTotal"`            // 电车总数
	Ebikes                []*model.StockMaterial `json:"ebikes"`                // 电车物资详情
	SmartBatteryTotal     int                    `json:"smartBatteryTotal"`     // 智能电池总数
	SmartBatteries        []*model.StockMaterial `json:"smartBatteries"`        // 智能电池物资详情
	NonSmartBatteryTotal  int                    `json:"nonSmartBatteryTotal"`  // 非智能电池总数
	NonSmartBatteries     []*model.StockMaterial `json:"nonSmartBatteries"`     // 非智能电池物资详情
	CabinetAccessoryTotal int                    `json:"cabinetAccessoryTotal"` // 电柜配件总数
	CabinetAccessories    []*model.StockMaterial `json:"cabinetAccessories"`    // 电柜配件物资详情
	EbikeAccessoryTotal   int                    `json:"ebikeAccessoryTotal"`   // 电车配件总数
	EbikeAccessories      []*model.StockMaterial `json:"ebikeAccessories"`      // 电车配件物资详情
	OtherAssetTotal       int                    `json:"otherAssetTotal"`       // 其他物资总数
	OtherAssets           []*model.StockMaterial `json:"otherAssets"`           // 其他物资详情
}
