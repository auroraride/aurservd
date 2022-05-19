// Copyright (C) liasica. 2022-present.
//
// Created at 2022-05-19
// Based on aurservd by liasica, magicrolan@qq.com.

package model

type PlanCreateReq struct {
    Name       string   `json:"name" validate:"required" trans:"骑士卡名称"`
    Enable     bool     `json:"enable"` // 是否启用
    Cities     []uint64 `json:"cities" validate:"required,min=1" trans:"启用城市"`
    Models     []uint64 `json:"models" validate:"required,min=1" trans:"电池型号"`
    Start      string   `json:"start" validate:"required,datetime=2006-01-02" trans:"开始日期"`
    End        string   `json:"end" validate:"required,datetime=2006-01-02" trans:"结束日期"`
    Price      float64  `json:"price" validate:"required" trans:"价格"`
    Days       uint     `json:"days" validate:"required,min=1" trans:"有效天数"`
    Commission float64  `json:"commission"` // 提成
}

// PlanEnableModifyReq 骑士卡状态修改请求
type PlanEnableModifyReq struct {
    ID     uint64 `json:"id" validate:"required" param:"id"` // 骑士卡ID
    Enable *bool  `json:"enable" validate:"required"`        // 启用或禁用
}
