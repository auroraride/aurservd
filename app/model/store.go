// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-22
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    StoreStatusMaintain uint8 = iota // 维护中
    StoreStatusOpen                  // 营业中
    StoreStatusClose                 // 休息中
    StoreStatusHidden                // 隐藏
)

// StoreCreateReq 门店创建请求
type StoreCreateReq struct {
    BranchID    *uint64              `json:"branchId" validate:"required" trans:"网点"`
    Name        *string              `json:"name" validate:"required" trans:"门店名称"`
    Status      uint8                `json:"status" validate:"required" enums:"0,1,2,3"` // 门店状态 0维护 1营业 2休息 3隐藏
    Materials   []StockStoreMaterial `json:"materials"`                                  // 门店物资
    EbikeObtain bool                 `json:"ebikeObtain"`                                // 是否可以领取车辆
    EbikeRepair bool                 `json:"ebikeRepair"`                                // 是否可以维修车辆
}

// StoreModifyReq 门店修改请求
type StoreModifyReq struct {
    ID          uint64  `json:"id" validate:"required" param:"id"`
    BranchID    *uint64 `json:"branchId" trans:"网点"`
    Name        *string `json:"name" trans:"门店名称"`
    Status      *uint8  `json:"status" enums:"0,1,2,3"` // 门店状态 0维护 1营业 2休息 3隐藏
    EbikeObtain *bool   `json:"ebikeObtain"`            // 是否可以领取车辆
    EbikeRepair *bool   `json:"ebikeRepair"`            // 是否可以维修车辆
}

type StoreItem struct {
    ID          uint64    `json:"id"`
    Name        string    `json:"name"`               // 门店名称
    Status      uint8     `json:"status"`             // 状态
    City        City      `json:"city"`               // 城市
    QRCode      string    `json:"qrcode"`             // 门店二维码
    Employee    *Employee `json:"employee,omitempty"` // 店员, 有可能不存在
    BranchID    uint64    `json:"branchId"`           // 网点ID
    Branch      BranchItem
    EbikeObtain bool `json:"ebikeObtain"` // 是否可以领取车辆
    EbikeRepair bool `json:"ebikeRepair"` // 是否可以维修车辆
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
    Status uint8 `json:"status"` // 门店状态
}

type StoreListReq struct {
    PaginationReq

    CityID      *uint64 `json:"cityId" query:"cityId"`           // 城市
    Name        *string `json:"name" query:"name"`               // 门店名称
    Status      *uint8  `json:"status" query:"status"`           // 门店状态
    EbikeObtain *bool   `json:"ebikeObtain" query:"ebikeObtain"` // 是否可以领取车辆, 不携带即为查询所有
    EbikeRepair *bool   `json:"ebikeRepair" query:"ebikeRepair"` // 是否可以维修车辆, 不携带即为查询所有
}

type StoreSwtichStatusReq struct {
    Status uint8 `json:"status" validate:"required,gte=1,lte=2" enums:"1,2"` // 状态 1:营业中 2:休息中
}
