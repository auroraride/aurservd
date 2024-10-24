// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type StoreStatus uint8

const (
	StoreStatusMaintain StoreStatus = iota // 维护中
	StoreStatusOpen                        // 营业中
	StoreStatusClose                       // 休息中
	StoreStatusHidden                      // 隐藏
)

func (t StoreStatus) Value() uint8 {
	return uint8(t)
}

// StoreCreateReq 门店创建请求
type StoreCreateReq struct {
	BranchID      *uint64              `json:"branchId" validate:"required" trans:"网点"`
	Name          *string              `json:"name" validate:"required" trans:"门店名称"`
	Status        *StoreStatus         `json:"status" validate:"required" enums:"0,1,2,3"`    // 门店状态 0维护 1营业 2休息 3隐藏
	Materials     []StockStoreMaterial `json:"materials"`                                     // 门店物资
	EbikeObtain   bool                 `json:"ebikeObtain"`                                   // 是否可以领取车辆
	EbikeRepair   bool                 `json:"ebikeRepair"`                                   // 是否可以维修车辆
	EbikeSale     bool                 `json:"ebikeSale"`                                     // 是否可以买车
	BusinessHours string               `json:"businessHours"`                                 // 营业时间
	Rest          bool                 `json:"rest"`                                          // 是否驿站
	Photos        []string             `json:"photos" validate:"required" trans:"门店照片"`       // 门店照片
	Phone         string               `json:"phone" validate:"required,max=20" trans:"门店电话"` // 门店电话
	HeadPic       string               `json:"headPic" validate:"required" trans:"门店头图"`      // 门店头图
	GroupID       *uint64              `json:"groupID"`                                       // 门店集合ID
}

// StoreModifyReq 门店修改请求
type StoreModifyReq struct {
	ID            uint64       `json:"id" validate:"required" param:"id"`
	BranchID      *uint64      `json:"branchId" trans:"网点"`
	Name          *string      `json:"name" trans:"门店名称"`
	Status        *StoreStatus `json:"status" enums:"0,1,2,3"`  // 门店状态 0维护 1营业 2休息 3隐藏
	EbikeObtain   *bool        `json:"ebikeObtain"`             // 是否可以领取车辆
	EbikeRepair   *bool        `json:"ebikeRepair"`             // 是否可以维修车辆
	EbikeSale     *bool        `json:"ebikeSale"`               // 是否可以买车
	BusinessHours *string      `json:"businessHours"`           // 营业时间
	Rest          *bool        `json:"rest"`                    // 是否驿站
	Photos        *[]string    `json:"photos"`                  // 门店照片
	Phone         *string      `json:"phone" validate:"max=20"` // 门店电话
	HeadPic       *string      `json:"headPic"`                 // 门店头图
	GroupID       *uint64      `json:"groupID"`                 // 门店集合ID
}

type StoreItem struct {
	ID            uint64      `json:"id"`
	Name          string      `json:"name"`               // 门店名称
	Status        StoreStatus `json:"status"`             // 状态
	City          City        `json:"city"`               // 城市
	QRCode        string      `json:"qrcode"`             // 门店二维码
	Employee      *Employee   `json:"employee,omitempty"` // 店员, 有可能不存在
	BranchID      uint64      `json:"branchId"`           // 网点ID
	Branch        BranchItem
	EbikeObtain   bool        `json:"ebikeObtain"`     // 是否可以领取车辆
	EbikeRepair   bool        `json:"ebikeRepair"`     // 是否可以维修车辆
	EbikeSale     bool        `json:"ebikeSale"`       // 是否可以买车
	BusinessHours string      `json:"businessHours"`   // 营业时间
	Rest          bool        `json:"rest"`            // 是否驿站
	Photos        []string    `json:"photos"`          // 照片
	Phone         string      `json:"phone"`           // 门店电话
	HeadPic       string      `json:"headPic"`         // 门店头图
	Group         *StoreGroup `json:"group,omitempty"` // 门店集合
}

type StoreGroup struct {
	ID   uint64 `json:"id"`   // 集合id
	Name string `json:"name"` // 集合名称
}

type Store struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"` // 门店名称
}

type StoreLngLat struct {
	Store
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

type StoreWithStatus struct {
	Store
	Status StoreStatus `json:"status"` // 门店状态
}

type StoreBusinessType uint8

const (
	StoreBusinessTypeObtain StoreBusinessType = iota + 1 // 领取车辆(租车)
	StoreBusinessTypeRepair                              // 维修
	StoreBusinessTypeSale                                // 购买
	StoreBusinessTypeRest                                // 驿站
)

func (t StoreBusinessType) Value() uint8 {
	return uint8(t)
}

type StoreListReq struct {
	PaginationReq

	CityID *uint64      `json:"cityId" query:"cityId"` // 城市
	Name   *string      `json:"name" query:"name"`     // 门店名称
	Status *StoreStatus `json:"status" query:"status"` // 门店状态

	BusinessType *StoreBusinessType `json:"businessType" query:"businessType"` // 业务类型 1:领取车辆(租车) 2:维修车辆 3:买车 4:驿站
	GroupID      *uint64            `json:"groupID" query:"groupID"`           // 门店集合ID
}

type StoreSwtichStatusReq struct {
	Status StoreStatus `json:"status" validate:"required,gte=1,lte=2" enums:"1,2"` // 状态 1:营业中 2:休息中
}
