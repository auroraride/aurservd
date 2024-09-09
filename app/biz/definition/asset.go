// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-20, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

type ReqType uint8

const (
	ReqTypeBattery ReqType = iota + 1 // 电池
	ReqTypeEbike                      // 电车
)

// AssetListReq 资产列表请求
type AssetListReq struct {
	model.PaginationReq
	WareHousID *uint64            `json:"wareHousID" query:"wareHousID"`           // 仓库ID
	StoreID    *uint64            `json:"storeID" query:"storeID"`                 // 门店ID
	Status     *model.AssetStatus `json:"status" query:"status" enums:"1,2,3,4,5"` // 资产状态 0:待入库 1:库存中 2:配送中 3:使用中 4:故障 5:报废
	ModelID    *uint64            `json:"modelID" query:"modelID"`                 // 电池型号ID
	BrandID    *uint64            `json:"brandId" query:"brandId"`                 // 电车品牌ID
	SN         *string            `json:"sn" query:"sn"`                           // 编号(电池编号或车架号或车牌号)
	Type       ReqType            `json:"type" query:"type" enums:"1,2"`           // 请求资产类型 1:电池 2:电车
}

// CommonAssetTotal 资产总计通用
type CommonAssetTotal struct {
	EbikeTotal            int `json:"ebikeTotal"`            // 电车总数
	SmartBatteryTotal     int `json:"smartBatteryTotal"`     // 智能电池总数
	NonSmartBatteryTotal  int `json:"nonSmartBatteryTotal"`  // 非智能电池总数
	CabinetAccessoryTotal int `json:"cabinetAccessoryTotal"` // 电柜配件总数
	EbikeAccessoryTotal   int `json:"ebikeAccessoryTotal"`   // 电车配件总数
	OtherAssetTotal       int `json:"otherAssetTotal"`       // 其他物资总数
}

// CommonAssetDetail 资产详情通用
type CommonAssetDetail struct {
	Ebikes             []*AssetMaterial `json:"ebikes"`             // 电车物资详情
	SmartBatteries     []*AssetMaterial `json:"smartBatteries"`     // 智能电池物资详情
	NonSmartBatteries  []*AssetMaterial `json:"nonSmartBatteries"`  // 非智能电池物资详情
	CabinetAccessories []*AssetMaterial `json:"cabinetAccessories"` // 电柜配件物资详情
	EbikeAccessories   []*AssetMaterial `json:"ebikeAccessories"`   // 电车配件物资详情
	OtherAssets        []*AssetMaterial `json:"otherAssets"`        // 其他物资详情
}

// WarestoreAssetDetail 小程序资产统计数据
type WarestoreAssetDetail struct {
	EbikeTotal            int                  `json:"ebikeTotal"`            // 电车总数
	Ebikes                []*WarestoreMaterial `json:"ebikes"`                // 电车物资详情
	SmartBatteryTotal     int                  `json:"smartBatteryTotal"`     // 智能电池总数
	SmartBatteries        []*WarestoreMaterial `json:"smartBatteries"`        // 智能电池物资详情
	NonSmartBatteryTotal  int                  `json:"nonSmartBatteryTotal"`  // 非智能电池总数
	NonSmartBatteries     []*WarestoreMaterial `json:"nonSmartBatteries"`     // 非智能电池物资详情
	CabinetAccessoryTotal int                  `json:"cabinetAccessoryTotal"` // 电柜配件总数
	CabinetAccessories    []*WarestoreMaterial `json:"cabinetAccessories"`    // 电柜配件物资详情
	EbikeAccessoryTotal   int                  `json:"ebikeAccessoryTotal"`   // 电车配件总数
	EbikeAccessories      []*WarestoreMaterial `json:"ebikeAccessories"`      // 电车配件物资详情
	OtherAssetTotal       int                  `json:"otherAssetTotal"`       // 其他物资总数
	OtherAssets           []*WarestoreMaterial `json:"otherAssets"`           // 其他物资详情
}

// WarestoreMaterial 小程序资产明细数据
type WarestoreMaterial struct {
	ID   uint64 `json:"id"`   // 电池类型ID、电车品牌ID、其他物资ID
	Name string `json:"name"` // 物资名称
	Num  int    `json:"num"`  // 物资数量
}

// AssetSnReq 通过sn查询资产请求
type AssetSnReq struct {
	SN string `json:"sn" query:"sn" param:"sn"` // 资产编号
}
