// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-01
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
	"database/sql/driver"

	"github.com/auroraride/aurservd/app/model/types"
)

const (
	EbikeColorDefault = "橘黄"
)

type EbikeStatus uint8

const (
	EbikeStatusInStock     EbikeStatus = iota // 库存中
	EbikeStatusUsing                          // 使用中
	EbikeStatusMaintenance                    // 维修中
	EbikeStatusScrapped                       // 已报废
)

func (p *EbikeStatus) Scan(t interface{}) (err error) {
	var v uint8
	v, err = types.Uint8(t)
	if err == nil {
		*p = EbikeStatus(v)
	}

	return
}

func (p EbikeStatus) Value() (driver.Value, error) {
	return uint8(p), nil
}

func (p EbikeStatus) RawValue() any {
	return uint8(p)
}

func (p EbikeStatus) String() string {
	switch p {
	case EbikeStatusInStock:
		return "库存中"
	case EbikeStatusUsing:
		return "使用中"
	case EbikeStatusMaintenance:
		return "维修中"
	case EbikeStatusScrapped:
		return "已报废"
	}
	return " - "
}

type EbikeListFilter struct {
	RiderID      uint64       `json:"riderId" query:"riderId"`                 // 骑手ID
	StoreID      uint64       `json:"storeId" query:"storeId"`                 // 门店ID
	BrandID      uint64       `json:"brandId" query:"brandId"`                 // 品牌ID
	Enable       *bool        `json:"enable" query:"enable"`                   // 是否启用, 默认`true`, 不携带为获取全部
	Status       *EbikeStatus `json:"status" query:"status" enums:"0,1,2,3"`   // 状态, 0:库存中 1:使用中 2:维修中 3:已报废, 不携带为获取全部
	ExFactory    string       `json:"exFactory" query:"exFactory"`             // 生产批次
	Keyword      string       `json:"keyword" query:"keyword"`                 // 搜索关键词<骑手:电话/姓名, 车辆:车架号/车牌号/终端编号/SIM卡号>
	OwnerType    *uint8       `json:"ownerType" query:"ownerType" enums:"1,2"` // 归属类型   1:平台 2:代理商
	EnterpriseID *uint64      `json:"enterpriseId" query:"enterpriseId"`       // 团签ID
	StationID    *uint64      `json:"stationId" query:"stationId"`             // 站点ID
	Goal         EbikeGoal    `json:"goal" query:"goal" enums:"0,1,2"`         // 查询目标, 0:不筛选 1:站点 2:骑手
}

type EbikeAttributes struct {
	Enable  *bool   `json:"enable,omitempty"`  // 是否启用, 默认要启用
	Plate   *string `json:"plate,omitempty"`   // 车牌号
	Machine *string `json:"machine,omitempty"` // 终端编号
	Sim     *string `json:"sim,omitempty"`     // SIM卡号
	Color   *string `json:"color,omitempty"`   // 颜色, 默认`橘黄`, 创建或编辑时用选择列表, 选项为: `橘黄` / `红` / `白` / `黑`
}

type EbikeListReq struct {
	PaginationReq
	EbikeListFilter
}

type EbikeListRes struct {
	ID uint64 `json:"id"`
	EbikeAttributes

	ExFactory      string  `json:"exFactory"`                // 生产批次
	Rider          string  `json:"rider,omitempty"`          // 骑手
	Store          string  `json:"store,omitempty"`          // 门店
	Brand          string  `json:"brand,omitempty"`          // 品牌
	SN             string  `json:"sn"`                       // 车架号
	BrandID        uint64  `json:"brandId"`                  // 品牌ID
	Status         string  `json:"status"`                   // 状态
	StationName    string  `json:"stationName,omitempty"`    // 所属站点名称
	StationID      *uint64 `json:"stationId,omitempty"`      // 站点id
	EnterpriseName string  `json:"enterpriseName,omitempty"` // 团签名称
	EnterpriseID   *uint64 `json:"enterpriseId,omitempty"`   // 团签ID

}

type EbikeCreateReq struct {
	SN        string `json:"sn" validate:"required" trans:"车架号"`
	ExFactory string `json:"exFactory" validate:"required" trans:"生产批次"`
	BrandID   uint64 `json:"brandId" validate:"required" trans:"型号"` // 关联: `MB015 车辆型号列表`
	EbikeAttributes
}

type EbikeModifyReq struct {
	ID        uint64  `json:"id" param:"id" validate:"required"`
	ExFactory *string `json:"exFactory"` // 生产批次
	EbikeAttributes
}

type EbikeInfo struct {
	ID           uint64  `json:"id,omitempty"`
	SN           string  `json:"sn,omitempty"`           // 车架号
	ExFactory    string  `json:"exFactory,omitempty"`    // 生产批次
	Plate        *string `json:"plate,omitempty"`        // 车牌号
	Color        string  `json:"color,omitempty"`        // 颜色
	EnterpriseID *uint64 `json:"enterpriseId,omitempty"` // 团签ID
	StationID    *uint64 `json:"stationId,omitempty"`    // 站点ID
}

type Ebike struct {
	EbikeInfo
	Brand *EbikeBrand `json:"brand,omitempty" bson:"brand,omitempty"` // 品牌信息
}

type EbikeBusinessInfo struct {
	ID        uint64 `json:"id"`
	BrandID   uint64 `json:"brandId"`
	BrandName string `json:"brandName"`
}

type EbikeUnallocatedParams struct {
	ID        *uint64 `json:"id"`        // 电车ID
	StoreID   *uint64 `json:"storeId"`   // 门店ID
	StationID *uint64 `json:"stationId"` // 站点ID
	Keyword   *string `json:"keyword"`   // 车辆关键词 (编码或牌照)
}

type EbikeAgentSearchReq struct {
	StationID *uint64 `json:"stationId" query:"stationId" validate:"required" trans:"站点ID (一般是subscribe.stationId)"`
	Keyword   *string `json:"keyword" query:"keyword" validate:"required" trans:"车辆关键词 (编码或牌照)"`
}
