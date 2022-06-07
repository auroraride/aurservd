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
    Phone        string `json:"phone" validate:"required,phone" trans:"电话号"`
}

type EnterpriseRider struct {
    ID        uint64            `json:"id"`
    Name      string            `json:"name"`              // 姓名
    Phone     string            `json:"phone"`             // 电话
    Days      int               `json:"days"`              // 总天数
    Unsettled int               `json:"unsettled"`         // 未结算天数
    Voltage   float64           `json:"voltage,omitempty"` // 可用电压型号, 骑手未开通订阅则此字段不存在
    CreatedAt string            `json:"createdAt"`         // 添加时间
    Station   EnterpriseStation `json:"station"`           // 站点
}

type EnterpriseRiderListReq struct {
    EnterpriseID uint64  `json:"enterpriseId" validate:"required" query:"enterpriseId" trans:"企业ID"`
    Keyword      *string `json:"keyword"` // 搜索关键词
    Start        *string `json:"start"`   // 使用开始时间
    End          *string `json:"end"`     // 使用结束时间
    PaginationReq
}
