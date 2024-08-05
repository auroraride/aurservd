// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-02, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

// CabinetAssetListReq 电柜资产列表请求
type CabinetAssetListReq struct {
	model.PaginationReq
	CityID  *uint64 `json:"cityId" query:"cityId"`   // 城市ID
	ModelID *uint64 `json:"modelID" query:"modelID"` // 电池型号ID
	Name    *string `json:"name" query:"name"`       // 电柜名称
	Sn      *string `json:"sn" query:"sn"`           // 电柜编号
	Start   *string `json:"start" query:"start"`     // 开始时间
	End     *string `json:"end" query:"end"`         // 结束时间
}

// CabinetAssetDetail 电柜资产信息
type CabinetAssetDetail struct {
	ID           uint64       `json:"id"`             // 电柜ID
	City         model.City   `json:"city"`           // 城市
	Sn           string       `json:"sn"`             // 电柜编号
	Name         string       `json:"name"`           // 电柜名称
	CabinetAsset CabinetAsset `json:"warehouseAsset"` // 电柜资产
}

// CabinetAsset 电柜资产
type CabinetAsset struct {
	SmartBatteryTotal    int                    `json:"smartBatteryTotal"`    // 智能电池总数
	SmartBatteries       []*model.StockMaterial `json:"smartBatteries"`       // 智能电池物资详情
	NonSmartBatteryTotal int                    `json:"nonSmartBatteryTotal"` // 非智能电池总数
	NonSmartBatteries    []*model.StockMaterial `json:"nonSmartBatteries"`    // 非智能电池物资详情
}
