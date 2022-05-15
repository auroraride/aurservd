// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// BatteryModel 电池型号
type BatteryModel struct {
    ID       uint64 `json:"id"`
    Voltage  string `json:"voltage"`
    Capacity string `json:"capacity"`
}

// BatteryModelCreateReq 电池型号创建请求
type BatteryModelCreateReq struct {
    Voltage  string `json:"voltage" validate:"required" trans:"电压"`
    Capacity string `json:"capacity" validate:"required" trans:"容量"`
}
