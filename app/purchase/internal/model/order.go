package model

type PaymentCache struct {
}

// Payway 支付方式
type Payway string

const (
	Alipay Payway = "alipay"
	Wechat Payway = "wechat"
	Cash   Payway = "cash"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待支付
	OrderStatusStaging   OrderStatus = "staging"   // 分期执行中
	OrderStatusEnded     OrderStatus = "ended"     // 已完成
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
	OrderStatusRefunded  OrderStatus = "refunded"  // 已退款
)

func (o OrderStatus) Value() string {
	return string(o)
}

func (o OrderStatus) String() string {
	switch o {
	case OrderStatusPending:
		return "待支付"
	case OrderStatusStaging:
		return "分期执行中"
	case OrderStatusEnded:
		return "已完成"
	case OrderStatusCancelled:
		return "已取消"
	case OrderStatusRefunded:
		return "已退款"
	default:
		return "未知"
	}
}

// OrderCreateReq 创建订单请求
type OrderCreateReq struct {
	GoodsID   uint64 `json:"goodsId" validate:"required"`   // 商品id
	PlanIndex *int   `json:"planIndex" validate:"required"` // 付款计划索引
}
