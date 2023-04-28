// Copyright (C) liasica. 2022-present.
//
// Created at 2022-06-11
// Based on aurservd by liasica, magicrolan@qq.com.

package model

const (
	RefundStatusPending uint8 = iota // 退款中
	RefundStatusSuccess              // 已同意
	RefundStatusRefused              // 已拒绝
	RefundStatusFail                 // 已失败
)

type Refund struct {
	Status      uint8     `json:"status" enums:"0,1,2,3"` // 退款状态 0:处理中 1:已同意 2:已拒绝 3:已失败
	Amount      float64   `json:"amount"`                 // 退款金额
	OutRefundNo string    `json:"outRefundNo"`            // 退款单号
	Reason      string    `json:"reason"`                 // 退款理由
	RefundAt    string    `json:"refundAt"`               // 退款成功时间
	CreatedAt   string    `json:"createdAt"`              // 申请退款时间
	Remark      string    `json:"remark"`                 // 备注
	Modifier    *Modifier `json:"modifier,omitempty"`     // 处理人 (可为空)
}

// RefundReq 退款申请
type RefundReq struct {
	SubscribeID *uint64 `json:"subscribeId"` // 骑士卡ID, 和deposit不能同时存在, 也不能同时为空
	Deposit     *bool   `json:"deposit"`     // 是否退押金, 押金退款条件: 1. 无最近订单; 2. 存在订单且状态为已退款或已退租
}

// RefundAuditReq 退款处理请求
type RefundAuditReq struct {
	OutRefundNo string `json:"outRefundNo" validate:"required" trans:"退款单号"`
	Remark      string `json:"remark"`                                             // 操作备注
	Status      uint8  `json:"status" enums:"1,2" validate:"required,gte=1,lte=2"` // 退款状态 1:同意 2:拒绝
}

type RefundRes struct {
	OutRefundNo string `json:"outRefundNo"` // 退款单号
}
