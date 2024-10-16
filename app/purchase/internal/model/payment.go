package model

// Payway 支付方式
type Payway string

const (
	Alipay Payway = "alipay"
	Wechat Payway = "wechat"
	Cash   Payway = "cash"
)

type PaymentPlanCreateReq struct {
	OrderID uint64 `json:"orderId" validate:"required"` // 订单id
}

type PaymentReq struct {
	OrderID   uint64 `json:"orderId" validate:"required"`   // 订单id
	Payway    Payway `json:"payway" validate:"required"`    // 支付方式
	PlanIndex *int   `json:"planIndex" validate:"required"` // 付款计划索引
}
