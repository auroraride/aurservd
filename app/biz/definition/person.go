// Copyright (C) liasica. 2024-present.
//
// Created at 2024-01-11
// Based on aurservd by liasica, magicrolan@qq.com.

package definition

// PersonCertificationOcrRes 实名认证Ocr参数
type PersonCertificationOcrRes struct {
	AppID   string `json:"appId"`   // WBAppid
	UserId  string `json:"userId"`  // 用户唯一标识
	OrderNo string `json:"orderNo"` // 订单号
	Ticket  string `json:"ticket"`  // NONCE ticket
}
