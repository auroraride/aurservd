// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// BatteryModel 电池型号
type BatteryModel struct {
    ID    uint64 `json:"id"`
    Model string `json:"model"` // 电池型号
}

// BatteryModelCreateReq 电池型号创建请求
type BatteryModelCreateReq struct {
    Model string `json:"model"` // 电池型号, 例如60V30AH
}
