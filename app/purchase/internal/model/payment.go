package model

// Payway 支付方式
type Payway string

const (
	Alipay Payway = "alipay"
	Wechat Payway = "wechat"
	Cash   Payway = "cash"
)

// PaymentStatus 支付订单状态
type PaymentStatus string

const (
	PaymentStatusObligation PaymentStatus = "obligation" // 待付款
	PaymentStatusPaid       PaymentStatus = "paid"       // 已支付
	PaymentStatusCanceled   PaymentStatus = "canceled"   // 已取消
	PaymentStatusOverdue    PaymentStatus = "overdue"    // 已逾期
)

func (p PaymentStatus) Value() string {
	return string(p)
}

func (p PaymentStatus) String() string {
	switch p {
	case PaymentStatusObligation:
		return "待付款"
	case PaymentStatusPaid:
		return "已支付"
	case PaymentStatusCanceled:
		return "已取消"
	case PaymentStatusOverdue:
		return "已逾期"
	default:
		return "未知"
	}
}

type PaymentPlanCreateReq struct {
	OrderID uint64 `json:"orderId" validate:"required"` // 订单id
}

type PaymentReq struct {
	OrderID   uint64 `json:"orderId" validate:"required"`   // 订单id
	Payway    Payway `json:"payway" validate:"required"`    // 支付方式
	PlanIndex *int   `json:"planIndex" validate:"required"` // 付款计划索引
}

// PaymentDetail 分期订单详情
type PaymentDetail struct {
	ID          uint64        `json:"id"`          // 分期订单ID
	Total       float64       `json:"total"`       // 支付金额（订单金额+逾期金额）
	Amount      float64       `json:"amount"`      // 订单金额
	BillingDate string        `json:"billingDate"` // 账单日期
	OverdueDays int           `json:"overdueDays"` // 逾期天数
	Forfeit     float64       `json:"forfeit"`     // 逾期金额（滞纳金）
	PaymentDate string        `json:"paymentDate"` // 支付日期
	Payway      Payway        `json:"payway"`      // 支付方式 alipay-支付宝 wechat-微信支付 cash-现金
	OutTradeNo  string        `json:"outTradeNo"`  // 交易订单号
	Status      PaymentStatus `json:"status"`      // 支付状态 obligation:待付款, paid:已支付, canceled:已取消, overdue:已逾期
}
