// Created at 2023-05-31

package model

type AgentPrepayReq struct {
	Amount float64 `json:"amount" validate:"required"` // 充值金额
	OpenID string  `json:"openId" validate:"required"` // 微信openid
}

type AgentPrepayRes struct {
	PrepayId  *string `json:"prepay_id"` // 预支付交易会话标识
	Appid     *string `json:"appId"`     // 应用ID
	TimeStamp *string `json:"timeStamp"` // 时间戳
	NonceStr  *string `json:"nonceStr"`  // 随机字符串
	Package   *string `json:"package"`   // 订单详情扩展字符串
	SignType  *string `json:"signType"`  // 签名方式
	PaySign   *string `json:"paySign"`   // 签名
}

type AgentPrepay struct {
	EnterpriseID uint64  `json:"enterpriseID,omitempty"` // 企业ID
	Remark       string  `json:"remark,omitempty"`       // 充值备注
	Amount       float64 `json:"amount,omitempty"`       // 充值金额

	ID    uint64 `json:"ID,omitempty"`    // 代理用户ID
	Name  string `json:"name,omitempty"`  // 代理姓名
	Phone string `json:"phone,omitempty"` // 代理电话
}
