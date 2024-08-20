// Copyright (C) aurservd. 2024-present.
//
// Created at 2024-08-19, by aurb

package definition

import "github.com/auroraride/aurservd/app/model"

const (
	SignTokenWarehouse = "WAREHOUSE" // 仓库平台
	SignTokenStore     = "STORE"     // 门店平台
)

type PlatType uint8

const (
	PlatTypeWarehouse PlatType = iota + 1 // 仓库平台登录
	PlatTypeStore                         // 门店平台登录
)

func (f PlatType) Value() uint8 {
	return uint8(f)
}

type WarestorePeopleSigninReq struct {
	Phone    string   `json:"phone" validate:"required" trans:"电话"`
	Password string   `json:"password" validate:"required" trans:"密码"`
	PlatType PlatType `json:"platType" validate:"required,oneof=1 2" trans:"登录平台类型"`
}

type WarestorePeopleSigninRes struct {
	Profile WarestorePeopleProfile `json:"profile"`
	Token   string                 `json:"token"`
}

type OpenidReq struct {
	Code string `json:"code" query:"code"`
}

type OpenidRes struct {
	Openid string `json:"openid"`
}

type WarestorePeopleProfile struct {
	ID           uint64 `json:"id"`
	Phone        string `json:"phone"`        // 手机号
	Name         string `json:"name"`         // 姓名
	RoleName     string `json:"roleName"`     // 角色名称
	Duty         bool   `json:"duty"`         // 上下班 `true`上班 `false`下班
	DutyLocation string `json:"dutyLocation"` // 上班位置
}

// TransferListReq 调拨记录筛选条件
type TransferListReq struct {
	model.PaginationReq
	model.AssetTransferFilter
}

// TransferDetailRes 调拨记录详情信息
type TransferDetailRes struct {
	model.AssetTransferListRes
	Detail []*model.AssetTransferDetail `json:"detail"` // 调拨详情数据
}

type AssetCountRes struct {
	ReceivingCount       int              `json:"receivingCount"`       // 待接收数量
	DeliveringCount      int              `json:"deliveringCount"`      // 配送中数量
	ExceptionCount       int              `json:"exceptionCount"`       // 异常告警数量
	EbikeAsset           AssetCountDetail `json:"ebikeAsset"`           // 电车资产统计
	SmartBatteryAsset    AssetCountDetail `json:"smartBatteryAsset"`    // 智能电池资产统计
	NonSmartBatteryAsset AssetCountDetail `json:"nonSmartBatteryAsset"` // 非智能电池资产统计
	OtherAsset           AssetCountDetail `json:"otherAsset"`           // 其他资产统计
}

type AssetCountDetail struct {
	StockCount int `json:"stockCount"` // 库存中数量
	FaultCount int `json:"faultCount"` // 故障数量
	TotalCount int `json:"totalCount"` // 合计数量
}

// AssetTransferReceiveBatchReq 批量接收资产
type AssetTransferReceiveBatchReq struct {
	AssetTransferReceive []model.AssetTransferReceiveReq `json:"assetTransferReceive" validate:"required,dive,required"`
}

// WarestoreAssetsReq 资产数据请求
type WarestoreAssetsReq struct {
	CityID      *uint64 `json:"cityId" query:"cityId"`           // 城市ID
	WarehouseID *uint64 `json:"warehouseID" query:"warehouseID"` // 仓库ID
	StoreID     *uint64 `json:"storeID" query:"storeID"`         // 门店ID
}

// WarestoreAssetRes 资产信息
type WarestoreAssetRes struct {
	ID     uint64               `json:"id"`     // ID
	Name   string               `json:"name"`   // 名称
	City   model.City           `json:"city"`   // 城市
	Detail WarestoreAssetDetail `json:"detail"` // 资产
}

// WarestoreAssetDetail 资产详情
type WarestoreAssetDetail struct {
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
