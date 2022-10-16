// Copyright (C) liasica. 2022-present.
//
// Created at 2022-10-14
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import (
    "github.com/qiniu/qmgo/field"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

const (
    EbikeAllocateExpiration = 600 // 分配过期时间(s), 超过一定时间不签约后分配失效, 需要重新分配
)

type EbikeAllocateStatus uint8

const (
    EbikeAllocateStatusPending EbikeAllocateStatus = iota + 1 // 未激活
    EbikeAllocateStatusSigned                                 // 已签约
    EbikeAllocateStatusVoid                                   // 已作废
)

type EbikeAllocate struct {
    field.DefaultField `bson:",inline"`

    SubscribeID uint64              `json:"subscribeId" bson:"subscribeId"`
    Status      EbikeAllocateStatus `json:"status" bson:"status"`
    Rider       Rider               `json:"rider" bson:"rider"`
    Ebike       Ebike               `json:"ebike" bson:"ebike"`
    Model       string              `json:"model" bson:"model"`
    StoreID     uint64              `json:"storeId" bson:"storeId"`
    EmployeeID  uint64              `json:"employeeId" bson:"employeeId"`
    SN          string              `json:"sn" bson:"sn"`
}

type EbikeAllocateReq struct {
    EbikeID     uint64 `json:"ebikeId" validate:"required" trans:"电车ID"`
    SubscribeID uint64 `json:"subscribeId" validate:"required" trans:"订阅ID"`
}

type EbikeAllocateRes struct {
    AllocateID string `json:"allocateId"` // 分配ID
}

type EbikeAllocateIDQueryReq struct {
    AllocateID primitive.ObjectID `json:"allocateId" validate:"required" query:"allocateId" trans:"分配ID"`
}

type EbikeAllocateInfo struct {
    Status EbikeAllocateStatus `json:"status" bson:"status"` // 签约状态 1:未签约 2:已签约 3:已作废
    Rider  Rider               `json:"rider" bson:"rider"`   // 骑手信息
    Ebike  Ebike               `json:"ebike" bson:"ebike"`   // 电车信息
    Model  string              `json:"model" bson:"model"`   // 电池型号
}

type EbikeAllocateEmployeeListReq struct {
    PaginationReq
    Status EbikeAllocateStatus `json:"status" query:"status"` // 签约状态 1:未签约(默认) 2:已签约
}
