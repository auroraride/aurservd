// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-06
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import jsoniter "github.com/json-iterator/go"

type StatementBillReq struct {
    End string `json:"end" validate:"required,datetime=2006-01-02" query:"end" trans:"账单截止日期"`
    ID  uint64 `json:"id" validate:"required" query:"id" trans:"企业ID"`
}

type BillOverview struct {
    Model  string  `json:"model"`  // 电池型号
    Number int     `json:"number"` // 使用骑手数量
    Price  float64 `json:"price"`  // 单价
    Days   int     `json:"days"`   // 天数
    Cost   float64 `json:"cost"`   // 账单金额
    City   City    `json:"city"`   // 城市
}

type StatementBillRes struct {
    ID           uint64              `json:"id"`           // 企业ID
    StatementID  uint64              `json:"statementId"`  // 账单ID
    UUID         string              `json:"uuid"`         // 账单编码, 结账时使用
    City         City                `json:"city"`         // 企业城市
    ContactName  string              `json:"contactName"`  // 联系人
    ContactPhone string              `json:"contactPhone"` // 联系电话
    Start        string              `json:"start"`        // 账单周期开始日期
    End          string              `json:"end"`          // 账单周期结束日期
    Cost         float64             `json:"cost"`         // 账单总额
    Days         int                 `json:"days"`         // 总使用天数
    Overview     []*BillOverview     `json:"overview"`     // 账单概览
    Bills        []StatementBillData `json:"bills"`        // 详情
}

type StatementBillData struct {
    EnterpriseID uint64  `json:"enterpriseId"` // 企业ID
    RiderID      uint64  `json:"riderId"`      // 骑手ID
    SubscribeID  uint64  `json:"subscribeId"`  // 订阅ID
    StatementID  uint64  `json:"statementId"`  // 账单ID
    Start        string  `json:"start"`        // 开始日期
    End          string  `json:"end"`          // 结束日期
    Days         int     `json:"days"`         // 天数
    Price        float64 `json:"price"`        // 单价
    Cost         float64 `json:"cost"`         // 金额小计
    Model        string  `json:"model"`        // 电池型号
    City         City    `json:"city"`         // 城市
}

func (d *StatementBillRes) MarshalBinary() ([]byte, error) {
    return jsoniter.Marshal(d)
}

func (d *StatementBillRes) UnmarshalBinary(data []byte) error {
    return jsoniter.Unmarshal(data, d)
}

type StatementClearBillReq struct {
    UUID   string  `json:"uuid"`   // 账单编码
    Remark *string `json:"remark"` // 备注信息
}
