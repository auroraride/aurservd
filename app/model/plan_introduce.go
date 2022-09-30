// Copyright (C) liasica. 2022-present.
//
// Created at 2022-09-30
// Based on aurservd by liasica, magicrolan@qq.com.

package model

// PlanIntroduceOption 未设置选项
type PlanIntroduceOption struct {
    Model  string       `json:"model"`  // 未设定的电池型号
    Brands []EbikeBrand `json:"brands"` // 该电池型号未设置的电车品牌
}

type PlanIntroduceEbike struct {
    ID    uint64 `json:"id"`
    Name  string `json:"name"`
    Model string `json:"model"`
}
type PlanIntroduceCreateReq struct {
    EbikeBrandID uint64 `json:"ebikeBrandId,omitempty"`    // 电车
    Model        string `json:"model" validate:"required"` // 电池型号
}
