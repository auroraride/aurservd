// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-03
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type BusinessSubscribeReq struct {
    ID        uint64  `json:"id" validate:"required" trans:"订阅ID"`
    StoreID   *uint64 `json:"storeId" trans:"门店ID"`
    CabinetID *uint64 `json:"cabinetId" trans:"电柜ID"`
}
