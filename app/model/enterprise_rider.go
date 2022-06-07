// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-07
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// EnterpriseRiderCreateReq 企业骑手创建请求
type EnterpriseRiderCreateReq struct {
    EnterpriseID uint64 `json:"enterpriseId" validate:"required" trans:"企业ID"`
    StationID    uint64 `json:"stationId" validate:"required" trans:"站点ID"`
    Name         string `json:"name" validate:"required" trans:"姓名"`
    Phone        string `json:"phone" validate:"required" trans:"电话号"`
}
