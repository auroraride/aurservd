// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-06
// Based on aurservd by liasica, magicrolan@qq.com.

package model

import jsoniter "github.com/json-iterator/go"

type StatementBillReq struct {
	End string `json:"end" validate:"required,datetime=2006-01-02" query:"end" trans:"账单截止日期"`
	ID  uint64 `json:"id" validate:"required" query:"id" trans:"企业ID"`

	Force bool `json:"-" swaggerignore:"true"`
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
	EnterpriseID uint64  `json:"enterpriseId"`        // 企业ID
	RiderID      uint64  `json:"riderId"`             // 骑手ID
	SubscribeID  uint64  `json:"subscribeId"`         // 订阅ID
	StatementID  uint64  `json:"statementId"`         // 账单ID
	Start        string  `json:"start"`               // 开始日期
	End          string  `json:"end"`                 // 结束日期
	Days         int     `json:"days"`                // 天数
	Price        float64 `json:"price"`               // 单价
	Cost         float64 `json:"cost"`                // 金额小计
	Model        string  `json:"model"`               // 电池型号
	City         City    `json:"city"`                // 城市
	StationID    *uint64 `json:"stationId,omitempty"` // 站点ID
}

func (d *StatementBillRes) MarshalBinary() ([]byte, error) {
	return jsoniter.Marshal(d)
}

func (d *StatementBillRes) UnmarshalBinary(data []byte) error {
	return jsoniter.Unmarshal(data, d)
}

type StatementClearBillReq struct {
	UUID   string `json:"uuid" validate:"required" trans:"账单编码"`
	Remark string `json:"remark" validate:"required" trans:"备注信息"`
}

type StatementBillHistoricalListReq struct {
	PaginationReq
	Start        *string `json:"start" query:"start"` // 账单开始日期
	End          *string `json:"end" query:"end"`     // 账单结束日期
	EnterpriseID uint64  `json:"enterpriseId" query:"enterpriseId" trans:"企业ID"`
}

type StatementBillHistoricalListRes struct {
	ID        uint64    `json:"id"`
	Cost      float64   `json:"cost"`      // 账单费用
	Remark    string    `json:"remark"`    // 结账备注
	Creator   *Modifier `json:"creator"`   // 结账人
	Days      int       `json:"days"`      // 使用日期
	Start     string    `json:"start"`     // 账单开始日期
	End       string    `json:"end"`       // 账单结束日期
	SettledAt string    `json:"settledAt"` // 结账时间
}

type StatementBillDetailReq struct {
	ID uint64 `json:"id" query:"id" validate:"required" trans:"账单ID"`
}

type StatementBillDetailExport struct {
	ID     uint64 `json:"id" validate:"required" trans:"账单ID"`
	Remark string `json:"remark" validate:"required" trans:"备注"`
}

type StatementDetail struct {
	Rider   Rider              `json:"rider"`             // 骑手信息
	City    City               `json:"city"`              // 城市
	Start   string             `json:"start"`             // 开始日期
	End     string             `json:"end"`               // 结束日期
	Days    int                `json:"days"`              // 使用天数
	Model   string             `json:"model"`             // 电池型号
	Price   float64            `json:"price"`             // 日单价
	Cost    float64            `json:"cost"`              // 费用
	Station *EnterpriseStation `json:"station,omitempty"` // 站点, 可能为空
}

type StatementUsageFilter struct {
	ID        uint64 `json:"id" query:"id"`               // 企业ID
	Start     string `json:"start" query:"start"`         // 开始时间
	End       string `json:"end" query:"end"`             // 结束时间
	CityID    uint64 `json:"cityId" query:"cityId"`       // 城市ID
	StationID uint64 `json:"stationId" query:"stationId"` // 站点ID
	Model     string `json:"model" query:"model"`         // 型号
	Keyword   string `json:"keyword" query:"keyword"`     // 骑手关键词 (手机号或姓名)
}

type StatementUsageReq struct {
	PaginationReq
	StatementUsageFilter
}

type StatementUsageExport struct {
	ID     uint64 `json:"id" validate:"required" trans:"企业ID"`
	Remark string `json:"remark" validate:"required" trans:"备注"`
	Start  string `json:"start"` // 开始时间
	End    string `json:"end"`   // 结束时间
}

type StatementUsageItem struct {
	Start string  `json:"start"` // 开始日期
	End   string  `json:"end"`   // 结束日期
	Days  int     `json:"days"`  // 使用天数
	Price float64 `json:"price"` // 日单价
	Cost  float64 `json:"cost"`  // 费用
}

type StatementUsageRes struct {
	Model     string                `json:"model"`               // 电池型号
	Rider     Rider                 `json:"rider"`               // 骑手信息
	City      City                  `json:"city"`                // 城市
	Station   *EnterpriseStation    `json:"station,omitempty"`   // 站点, 可能为空
	Status    string                `json:"status"`              // 骑手状态 计费中,已退租
	DeletedAt string                `json:"deletedAt,omitempty"` // 删除时间
	Items     []*StatementUsageItem `json:"items"`               // 使用详情
}
