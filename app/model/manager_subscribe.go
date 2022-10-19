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
