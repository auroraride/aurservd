// Copyright (C) liasica. 2022-present.
//
// Created at 2022-08-03
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type BusinessSubscribeReq struct {
    ID        uint64  `json:"id" validate:"required" trans:"订阅ID"`
    StoreID   *uint64 `json:"storeId" trans:"门店ID"`
    CabinetID *uint64 `json:"cabinetId" trans:"电柜ID"`
}

type BusinessCabinetReq struct {
    ID     uint64 `json:"id" validate:"required" trans:"订阅ID"`
    Serial string `json:"serial" validate:"required" trans:"电柜编码"`
}

type BusinessCabinetStatus struct {
    UUID  string `json:"uuid"`  // 操作ID, 使用此参数轮询获取状态
    Index int    `json:"index"` // 仓位Index, +1是仓位号
}

type BusinessCabinetStatusReq struct {
    UUID primitive.ObjectID `json:"uuid" validate:"required" query:"uuid" trans:"操作ID"`
}

type BusinessCabinetStatusRes struct {
    Success bool   `json:"success"` // 是否成功
    Stop    bool   `json:"stop"`    // 是否终止
    Message string `json:"message"` // 失败消息
}