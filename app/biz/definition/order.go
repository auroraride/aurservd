package definition

// OrderDepositFreeReq  免押金支付请求
type OrderDepositFreeReq struct {
	PlanID uint64 `json:"planId" validate:"required"  trans:"套餐ID"`            // 套餐ID
	Payway uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"4,5"` // 4支付宝芝麻免押 5微信支付分支付
}

type OrderDepositFreeRes struct {
	Prepay string `json:"prepay"` //  预支付信息
}
