package definition

type Payway uint8

// TODO: 禁止重复定义，统一使用 model.OrderPaywayWechatDeposit / model.OrderPaywayAlipayDeposit
const (
	PaywayAlipayDeposit Payway = iota + 1 // 支付宝芝麻免押
	PaywayWechatDeposit                   // 微信支付分支付
	PaywayAlipay                          // 预支付
)

func (p Payway) Value() uint8 {
	return uint8(p)
}

// OrderDepositCreditReq  信用免押请求
type OrderDepositCreditReq struct {
	PlanID uint64 `json:"planId" validate:"required"  trans:"套餐ID"`              // 套餐ID
	Payway Payway `json:"payway" validate:"required" trans:"支付方式" enums:"1,2,3"` //   1支付宝 2微信 3支付宝预授权
}

type OrderDepositCreditRes struct {
	Prepay string `json:"prepay"` //  预支付信息
}

// OrderCreateReq 订单创建请求
type OrderCreateReq struct {
	PlanID    uint64 `json:"planId" validate:"required" trans:"套餐ID"`
	Payway    uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"1,2,3"`            // 1支付宝 2微信 3支付宝预授权
	OrderType uint   `json:"orderType" validate:"required" trans:"订单类型" enums:"1,2,3,4,5,6,7"` // 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金

	CityID uint64 `json:"cityId"` // 城市ID, 新签必填

	Point   bool     `json:"point"`   // 是否使用积分
	Coupons []uint64 `json:"coupons"` // 优惠券
}

// OrderCreateRes 订单创建响应
type OrderCreateRes struct {
	OrderID string `json:"orderId"` // 订单ID
}
