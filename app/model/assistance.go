// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-23
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
    AssistanceStatusPending  uint8 = iota // 等待接单
    AssistanceStatusApproved              // 已接单
    AssistanceStatusRefused               // 已拒绝
    AssistanceStatusFailed                // 救援失败
    AssistanceStatusSuccess               // 救援成功
)

type AssistanceCreateReq struct {
    Lng             float64  `json:"lng" validate:"required" trans:"经度"`
    Lat             float64  `json:"lat" validate:"required" trans:"纬度"`
    Address         string   `json:"address" validate:"required" trans:"详细地址"`
    Breakdown       string   `json:"breakdown" validate:"required" trans:"故障原因"`
    BreakdownDesc   string   `json:"breakdownDesc" validate:"required" trans:"故障描述"`
    BreakdownPhotos []string `json:"breakdownPhotos" validate:"required,min=1,max=3" trans:"故障图片"`
}

type AssistanceCreateRes struct {
    OutTradeNo string `json:"outTradeNo"` // 救援订单编码
}
