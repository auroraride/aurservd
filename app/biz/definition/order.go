package definition

// OrderDepositCreditReq  信用免押请求
type OrderDepositCreditReq struct {
	PlanID uint64 `json:"planId" validate:"required"  trans:"套餐ID"`            // 套餐ID
	Payway uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"3,4"` // 1:支付宝 2:微信 3:支付宝预授权支付 4:微信支付分支付
}

type OrderDepositCreditRes struct {
	Prepay     string `json:"prepay"`               //  预支付信息
	OutOrderNo string `json:"outOrderNo,omitempty"` //  订单号
	OutTradeNo string `json:"outTradeNo,omitempty"` //  平台订单号
}

// OrderCreateReq 订单创建请求
type OrderCreateReq struct {
	PlanID    uint64 `json:"planId" validate:"required" trans:"套餐ID"`
	Payway    uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"1,2"`              // 1支付宝 2微信 3支付宝预授权 4:微信支付分支付
	OrderType uint   `json:"orderType" validate:"required" trans:"订单类型" enums:"1,2,3,4,5,6,7"` // 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金

	CityID uint64 `json:"cityId"` // 城市ID, 新签必填

	Point   bool     `json:"point"`   // 是否使用积分
	Coupons []uint64 `json:"coupons"` // 优惠券
}

// OrderCreateRes 订单创建响应
type OrderCreateRes struct {
	OrderID string `json:"orderId"` // 订单ID
}

// OrderDepositCancelReq 取消押金订单
type OrderDepositCancelReq struct {
	OutOrderNo string `json:"outOrderNo" validate:"required"` // 商户端的唯一订单号
}

// OrderDepositFreezeToPayRes 冻结金额转支付
type OrderDepositFreezeToPayRes struct {
	OutTradeNo string `json:"outTradeNo" validate:"required"` // 订单编号
	TradeNo    string `json:"tradeNo" validate:"required"`    // 平台订单号
}

// OrderDepositUnfreezeReq 解冻金额请求
type OrderDepositUnfreezeReq struct {
	OutOrderNo string `json:"outOrderNo" validate:"required"` // 商户端的唯一订单号
}

// 解冻支付宝请求
type FandAuthUnfreezeReq struct {
	AuthNo       string  `json:"authNo" validate:"required"`       // 支付宝授权资金订单号
	OutRequestNo string  `json:"outRequestNo" validate:"required"` // 商户端的唯一订单号
	Amount       float64 `json:"amount" validate:"required"`       // 解冻金额
	Remark       string  `json:"remark" validate:"required"`       // 备注
}

// ExtraParam 预授权免押拓展参数
type ExtraParam struct {
	Category      string         `json:"category"`                // 授权业务对应的类目
	ServiceId     string         `json:"serviceId"`               // 服务ID
	CreditExtInfo *CreditExtInfo `json:"creditExtInfo,omitempty"` // 信用支付扩展参数
}

// CreditExtInfo 信用支付扩展参数
type CreditExtInfo struct {
	AssessmentAmount string `json:"assessmentAmount"` // 评估金额
}

// PostPayments 后付费项目
type PostPayments struct {
	Name   string `json:"name"`   // 项目名称
	Amount string `json:"amount"` // 项目金额
}
