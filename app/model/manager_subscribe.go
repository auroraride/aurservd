// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-19
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type ManagerSubscribeActive struct {
    ID           uint64  `json:"id" validate:"required" trans:"订阅ID"`
    StoreID      *uint64 `json:"storeId"`                // 门店ID
    EbikeKeyword *string `json:"ebikeKeyword,omitempty"` // 车架号或车牌号
}

type ManagerSubscribeChangeEbike struct {
    ID           uint64  `json:"id" validate:"required" trans:"订阅ID"`
    StoreID      uint64  `json:"storeId" validate:"required" trans:"门店ID"` // 旧车入库至门店
    EbikeKeyword *string `json:"ebikeKeyword" validate:"required" trans:"车架号或车牌号"`
}
