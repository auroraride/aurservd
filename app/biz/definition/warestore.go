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
	ID     uint64            `json:"id"`     // ID
	Name   string            `json:"name"`   // 名称
	City   model.City        `json:"city"`   // 城市
	Detail CommonAssetDetail `json:"detail"` // 资产
}

type CommonAssetType uint8

const (
	CommonAssetTypeEbike   CommonAssetType = iota + 1 // 电车
	CommonAssetTypeBattery                            // 电池
)

func (f CommonAssetType) Value() uint8 {
	return uint8(f)
}

// WarestoreAssetsCommonReq 电池/电车资产数据请求
type WarestoreAssetsCommonReq struct {
	model.PaginationReq
	Type           CommonAssetType    `json:"type" query:"type" validate:"required" enums:"1,2"` // 资产类型 1:电车 2:电池
	WarehouseID    *uint64            `json:"warehouseID" query:"warehouseID"`                   // 仓库ID
	StoreID        *uint64            `json:"storeID" query:"storeID"`                           // 门店ID
	Status         *model.AssetStatus `json:"status" query:"status" enums:"1,2,3,4,5"`           // 资产状态 0:待入库 1:库存中 2:配送中 3:使用中 4:故障 5:报废
	ModelID        *uint64            `json:"model" query:"model"`                               // 电池型号ID
	BrandID        *uint64            `json:"brandID" query:"brandID"`                           // 电车型号ID
	BatteryKeyword *string            `json:"batteryKeyword" query:"batteryKeyword"`             // 电池编号关键字
	EbikeKeyword   *string            `json:"ebikeKeyword" query:"ebikeKeyword"`                 // 电车车牌/车架号关键字
}
