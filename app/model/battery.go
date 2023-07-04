// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "github.com/auroraride/adapter"

type BatteryCreateReq struct {
	SN     string `json:"sn" validate:"required" trans:"电池编号"`
	CityID uint64 `json:"cityId" validate:"required" trans:"城市"`

	Enable *bool `json:"enable"` // 是否启用, 需默认为`true`
}

type BatteryModifyReq struct {
	ID     uint64  `json:"id" param:"id" validate:"required" trans:"电池ID"`
	Enable *bool   `json:"enable"` // 是否启用
	CityID *uint64 `json:"cityId"` // 城市ID
}

type BatteryFilter struct {
	SN           string      `json:"sn" query:"sn"`                           // 编号
	Model        string      `json:"model" query:"model"`                     // 型号
	CityID       uint64      `json:"cityId" query:"cityId"`                   // 城市
	Status       *int        `json:"status" query:"status"`                   // 状态 0:全部 1:启用(不携带默认为启用) 2:禁用
	EnterpriseID *uint64     `json:"enterpriseId" query:"enterpriseId"`       // 团签ID
	StationID    *uint64     `json:"stationId" query:"stationId"`             // 站点ID
	CabinetID    *uint64     `json:"cabinetID" query:"cabinetID"`             // 电柜ID
	CabinetName  *string     `json:"cabinetName" query:"cabinetName"`         // 电柜名称
	Keyword      *string     `json:"keyword" query:"keyword"`                 // 关键词
	OwnerType    *uint8      `json:"ownerType" query:"ownerType" enums:"1,2"` // 归属类型   1:平台 2:代理商
	RiderID      *uint64     `json:"riderId" query:"riderId"`                 // 骑手ID
	Goal         BatteryGoal `json:"goal" query:"goal" enums:"0,1,2,3"`       // 查询目标, 0:不筛选 1:站点 2:电柜 3:骑手
}

type BatteryListReq struct {
	PaginationReq
	BatteryFilter
}

type BatteryStatus uint8

const (
	BatteryStatusIdle        BatteryStatus = iota // 充电中
	BatteryStatusCharging                         // 充电中
	BatteryStatusDisCharging                      // 放电中
	BatteryStatusFault                            // 异常
)

type BatteryListRes struct {
	ID      uint64               `json:"id"`
	Brand   adapter.BatteryBrand `json:"brand"`             // 品牌 TB:拓邦, XC:星创
	City    *City                `json:"city,omitempty"`    // 城市
	Model   string               `json:"model"`             // 型号
	Enable  bool                 `json:"enable"`            // 是否启用
	SN      string               `json:"sn"`                // 编号
	Rider   *Rider               `json:"rider,omitempty"`   // 骑手
	Cabinet *CabinetBasicInfo    `json:"cabinet,omitempty"` // 电柜
	*BmsBattery
	StationName    string `json:"stationName,omitempty"`    // 站点名称
	EnterpriseName string `json:"enterpriseName,omitempty"` // 团签名称
}

type Battery struct {
	ID    uint64 `json:"id"`
	SN    string `json:"sn"`    // 编号
	Model string `json:"model"` // 型号
}

type BatterySearchReq struct {
	Serial       string  `json:"serial" query:"serial" trans:"流水号" validate:"required,min=4"`
	EnterpriseID *uint64 `json:"enterpriseId" query:"enterpriseId"` // 团签ID: 0为查询非团签电池; 不携带为全部数据
	StationID    *uint64 `json:"stationId" query:"stationId"`       // 站点ID: 0为查询非站点电池; 不携带为全部数据
}

type BatteryBind struct {
	RiderID   uint64 `json:"riderId" validate:"required"`   // 骑手ID
	BatteryID uint64 `json:"batteryId" validate:"required"` // 电池ID
}

// BatteryDetail 电池信息
type BatteryDetail struct {
	ID    uint64  `json:"id"`    // 电池ID
	Model string  `json:"model"` // 电池型号
	SN    string  `json:"sn"`    // 电池编码
	Soc   float64 `json:"soc"`   // 当前电量, 暂时隐藏
}

type BatteryInCabinet struct {
	CabinetID uint64 `json:"cabinetId"` // 所在电柜ID
	Ordinal   int    `json:"ordinal"`   // 仓位序号
}

type BatteryUnbindRequest struct {
	RiderID uint64 `json:"riderId" validate:"required"` // 骑手ID
}

type BatteryBatchQueryRequest struct {
	IDs []uint64 `json:"ids" validate:"required,min=1"`
}

type BatteryQueryRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

// BatteryEnterpriseTransfer 代理商电池转移信息
type BatteryEnterpriseTransfer struct {
	Sn           string  `json:"sn"`           // 电池编码
	StationID    *uint64 `json:"stationId"`    // 站点ID
	EnterpriseID *uint64 `json:"enterpriseId"` // 团签ID
}

type BatteryGroup struct {
	Model string `json:"model"`
}
