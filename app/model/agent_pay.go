// Created at 2023-05-31

package model

type AgentPrepayReq struct {
	Amount float64 `json:"amount" validate:"required"` // 充值金额
	OpenID string  `json:"openId" validate:"required"` // 微信openid
}

type AgentPrepayRes struct {
	PrepayID string `json:"prepayId"`
}

type AgentPrepay struct {
	EnterpriseID uint64  `json:"enterpriseID,omitempty"` // 企业ID
	Remark       string  `json:"remark,omitempty"`       // 充值备注
	Amount       float64 `json:"amount,omitempty"`       // 充值金额
	ID           uint64  `json:"ID,omitempty"`           // 代理用户ID
	Name         string  `json:"name,omitempty"`         // 代理姓名
	Phone        string  `json:"phone,omitempty"`        // 代理电话
}
