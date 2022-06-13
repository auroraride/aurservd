// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type AttendancePrecheck struct {
    Duty    bool     `json:"duty"` // 上下班 `true`上班 `false`下班
    SN      *string  `json:"sn" validate:"required" trans:"门店编号"`
    Lat     *float64 `json:"lat" validate:"required" trans:"纬度"`
    Lng     *float64 `json:"lng" validate:"required" trans:"经度"`
    Address *string  `json:"address" validate:"required" trans:"详细地址"`
}

type AttendanceCreateReq struct {
    Photo     *string        `json:"photo"`                                              // 上班照片
    Inventory map[string]int `json:"inventory" validate:"required" trans:"物资盘点清单"` // 格式为 [名称]:数量
    *AttendancePrecheck
}
