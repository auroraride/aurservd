// Copyright (C) liasica. 2022-present.
//
// Created at 2022-03-01
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// 城市启用状态
const (
    CityStatusAll     = iota // 全部
    CityStatusNotOpen        // 未启用
    CityStatusOpen           // 已启用
)

// CityItem 城市
type CityItem struct {
    ID       uint64     `json:"id"`                 // 城市或省份ID
    Open     bool       `json:"open,omitempty"`     // 是否启用
    Name     string     `json:"name,omitempty"`     // 城市/省份
    Children []CityItem `json:"children,omitempty"` // 城市列表
}

// CityListReq 城市列表请求
type CityListReq struct {
    Status uint `query:"status" validate:"gte=0,lte=2"` // 启用状态 0:全部 1:未启用 2:已启用
}

// CityModifyReq 城市修改请求
type CityModifyReq struct {
    ID   uint64 `json:"id" validate:"required,number" trans:"城市"`
    Open bool   `json:"open" validate:"required" trans:"状态"`
}

// CityModifyRes 城市修改返回
type CityModifyRes struct {
    Open bool `json:"open,omitempty"`
}
