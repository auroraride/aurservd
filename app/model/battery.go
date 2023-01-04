// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-13
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type BatteryCreateReq struct {
    SN     string `json:"sn" validate:"required" trans:"电池编号"`
    CityID uint64 `json:"cityId" validate:"required" trans:"城市"`

    Enable *bool `json:"enable"` // 是否启用, 需默认为`true`
}

type BatteryModifyReq struct {
    ID     uint64  `json:"id" param:"id" validate:"required" trans:"电池ID"`
    Enable *bool   `json:"enable"` // 是否启用
    CityID *uint64 `json:"cityId"` // 城市ID
}

type BatteryFilter struct {
    SN     string `json:"sn" query:"sn"`         // 编号
    Model  string `json:"model" query:"model"`   // 型号
    CityID uint64 `json:"cityId" query:"cityId"` // 城市
    Status *int   `json:"status" query:"status"` // 状态 0:全部 1:启用(不携带默认为启用) 2:禁用
}

type BatteryListReq struct {
    PaginationReq
    BatteryFilter
}

type BatteryListRes struct {
    ID      uint64            `json:"id"`
    City    City              `json:"city"`              // 城市
    Model   string            `json:"model"`             // 型号
    Enable  bool              `json:"enable"`            // 是否启用
    SN      string            `json:"sn"`                // 编号
    Rider   *Rider            `json:"rider,omitempty"`   // 骑手
    Cabinet *CabinetBasicInfo `json:"cabinet,omitempty"` // 电柜
}

type Battery struct {
    ID    uint64 `json:"id"`
    SN    string `json:"sn"`    // 编号
    Model string `json:"model"` // 型号
}
