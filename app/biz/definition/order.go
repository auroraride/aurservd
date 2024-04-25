package definition

import "github.com/auroraride/aurservd/app/model"

// OrderDepositCreditReq  信用免押请求
type OrderDepositCreditReq struct {
	PlanID uint64 `json:"planId" validate:"required"  trans:"套餐ID"`                // 套餐ID
	Payway uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"1,2,3,4"` //  1支付宝 2微信 3支付宝信用免押 4 微信支付分
}

type OrderDepositCreditRes struct {
	Prepay     string `json:"prepay"`               //  支付信息
	OutOrderNo string `json:"outOrderNo,omitempty"` //  支付宝授权资金订单号
	OutTradeNo string `json:"outTradeNo,omitempty"` //  平台订单号
}

// OrderCreateReq 订单创建请求
type OrderCreateReq struct {
	PlanID    uint64 `json:"planId" trans:"套餐ID"`
	Payway    uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"1,2,3"`            // 1支付宝 2微信 3支付宝预授权(仅限V2使用)
	OrderType uint   `json:"orderType" validate:"required" trans:"订单类型" enums:"1,2,3,4,5,6,7"` // 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金

	CityID uint64 `json:"cityId"` // 城市ID, 新签必填

	Point   bool     `json:"point"`   // 是否使用积分
	Coupons []uint64 `json:"coupons"` // 优惠券

	// 以下字段仅限V2使用
	DepositAlipayAuthFreeze bool               `json:"depositAlipayAuthFreeze"` // 是否使用支付宝预授权信用分支付押金
	DepositOrderNo          *string            `json:"depositOrderNo"`          // 押金订单编号 (如果分开支付的押金此参必填 例如 选择了信用免押,支付为支付宝,则此参数必填)
	StoreID                 *uint64            `json:"storeId"`                 // 门店ID
	AgreementHash           *string            `json:"agreementHash"`           // 协议hash
	DepositType             *model.DepositType `json:"depositType"`             // 押金类型 1:芝麻免押 2:微信支付分免押 3:合同押金 4:支付押金
	PointNum                *int64             `json:"pointNum"`                // 积分数量
}

// OrderCreateRes 订单创建响应
type OrderCreateRes struct {
	OrderID string `json:"orderId"` // 订单ID
}

// OrderDepositCancelReq 取消押金订单
type OrderDepositCancelReq struct {
	OutOrderNo string `json:"outOrderNo" validate:"required"` // 支付宝授权资金订单号
}

// OrderDepositFreezeToPayRes 冻结金额转支付
type OrderDepositFreezeToPayRes struct {
	OutTradeNo string `json:"outTradeNo" validate:"required"` // 订单编号
	TradeNo    string `json:"tradeNo" validate:"required"`    // 平台订单号
}

// OrderDepositUnfreezeReq 解冻金额请求
type OrderDepositUnfreezeReq struct {
	OutOrderNo string `json:"outOrderNo" validate:"required"` // 支付宝授权资金订单号
}

// FandAuthUnfreezeReq 解冻支付宝请求
type FandAuthUnfreezeReq struct {
	AuthNo       string  `json:"authNo" validate:"required"`       // 支付宝授权资金订单号
	OutRequestNo string  `json:"outRequestNo" validate:"required"` // 商户授权资金操作流水号
	Amount       float64 `json:"amount" validate:"required"`       // 解冻金额
	Remark       string  `json:"remark" validate:"required"`       // 备注
	IsDeposit    bool    `json:"isDeposit"`                        // 是否是押金
}

type FandAuthUnfreezeRes struct {
	OutOrderNo string `json:"outOrderNo,omitempty"` //  授权资金订单号
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
	Name        string `json:"name"`        // 项目名称
	Amount      string `json:"amount"`      // 项目金额
	Description string `json:"description"` // 项目描述
}

// FundAuthOperationDetailReq 查询资金授权冻结订单
type FundAuthOperationDetailReq struct {
	AuthNo       string `json:"authNo"`       // 支付宝资金授权订单号
	OutRequestNo string `json:"outRequestNo"` // 商户授权资金操作流水号
	OutOrderNo   string `json:"out_order_no"` // 商户的资金授权订单号
}

type FreezeToPay struct {
	OutOrderNo string `json:"outOrderNo" validate:"required"` // 授权资金订单号
}
