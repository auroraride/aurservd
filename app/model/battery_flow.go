// Copyright (C) liasica. 2023-present.
//
// Created at 2023-02-07
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type BatteryFlowCreateReq struct {
    SN          string  `json:"sn" validate:"required"`
    BatteryID   uint64  `json:"batteryId" validate:"required"`
    RiderID     *uint64 `json:"riderId"`
    CabinetID   *uint64 `json:"cabinetId"`
    Serial      *string `json:"serial"`
    Ordinal     *int    `json:"ordinal"`
    SubscribeID *uint64 `json:"subscribeId"`
}