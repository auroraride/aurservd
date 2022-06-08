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
    BranchID *uint64 `json:"branchId" validate:"required" trans:"网点"`
    Name     *string `json:"name" validate:"required" trans:"门店名称"`
    Status   uint8   `json:"status" validate:"required" enums:"0,1,2,3"` // 门店状态 0维护 1营业 2休息 3隐藏
}

// StoreModifyReq 门店修改请求
type StoreModifyReq struct {
    ID       uint64  `json:"id" validate:"required" param:"id"`
    BranchID *uint64 `json:"branchId" trans:"网点"`
    Name     *string `json:"name" trans:"门店名称"`
    Status   *uint8  `json:"status" enums:"0,1,2,3"` // 门店状态 0维护 1营业 2休息 3隐藏
}

type StoreItem struct {
    ID       uint64    `json:"id"`
    Name     string    `json:"name"`               // 门店名称
    Status   uint8     `json:"status"`             // 状态
    City     City      `json:"city"`               // 城市
    Employee *Employee `json:"employee,omitempty"` // 店员, 有可能不存在
}

type Store struct {
    ID   uint64 `json:"id"`
    Name string `json:"name"` // 门店名称
}

type StoreListReq struct {
    PaginationReq

    CityID *uint64 `json:"cityId" query:"cityId"` // 城市
    Name   *string `json:"name" query:"name"`     // 门店名称
    Status *uint8  `json:"status" query:"status"` // 门店状态
}
