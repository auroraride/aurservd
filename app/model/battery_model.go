// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// BatteryModel 电池型号
type BatteryModel struct {
	ID    uint64 `json:"id,omitempty"`
	Model string `json:"model"` // 电池型号
}

// BatteryModelReq 电池型号创建请求
type BatteryModelReq struct {
	Model string `json:"model"` // 电池型号(POST), 例如60V30AH
}
