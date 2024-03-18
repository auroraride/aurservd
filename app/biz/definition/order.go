package definition

// OrderDepositCreditReq  信用免押请求
type OrderDepositCreditReq struct {
	PlanID uint64 `json:"planId" validate:"required"  trans:"套餐ID"`            // 套餐ID
	Payway uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"3,4"` //  3:支付宝预授权支付 4:微信支付分支付
}

type OrderDepositCreditRes struct {
	Prepay     string `json:"prepay"`               //  支付信息
	OutOrderNo string `json:"outOrderNo,omitempty"` //  预授权订单号
	OutTradeNo string `json:"outTradeNo,omitempty"` //  平台订单号
}

// OrderCreateReq 订单创建请求
type OrderCreateReq struct {
	PlanID    uint64 `json:"planId" validate:"required" trans:"套餐ID"`
	Payway    uint8  `json:"payway" validate:"required" trans:"支付方式" enums:"1,2,3"`            // 1支付宝 2微信 3支付宝预授权(仅限V2使用)
	OrderType uint   `json:"orderType" validate:"required" trans:"订单类型" enums:"1,2,3,4,5,6,7"` // 1新签 2续签 3重签 4更改电池 5救援 6滞纳金 7押金

	CityID uint64 `json:"cityId"` // 城市ID, 新签必填

	Point   bool     `json:"point"`   // 是否使用积分
	Coupons []uint64 `json:"coupons"` // 优惠券

	// 以下字段仅限V2使用
	DepositAlipayAuthFreeze bool    `json:"depositAlipayAuthFreeze"` // 是否使用支付宝预授权信用分支付押金
	NeedContract            *bool   `json:"needContract"`            // 是否需要签约
	DepositOrderNo          *string `json:"depositOrderNo"`          // 押金订单编号 (如果分开支付的押金此参必填 例如 选择了信用免押,支付为支付宝支付,则此参数必填)
}

// OrderCreateRes 订单创建响应
type OrderCreateRes struct {
	OrderID string `json:"orderId"` // 订单ID
}

// ILLEGAL_ARGUMENT	参数异常或参数缺失	 请求参数有错，重新检查请求后，再调用资金预授权明细撤销操作
// ACCESS_FORBIDDEN	无权限使用接口	 未签约条码支付或者合同已到期
// PAYER_USER_STATUS_LIMIT	付款方状态受限	 买家支付宝账户受限，请登录支付宝认证升级，详情咨询4007585858。
// ORDER_ALREADY_FINISH	授权订单已经完结，无法再进行资金操作	 本笔授权订单已经完结，不允许进行资金撤销操作，确认请求资金解冻的资金授权订单号是否正确
// CANCEL_OPERATION_TIME_OUT	撤销操作超过允许时间范围	 撤销操作需要在预授权冻结操作之后的24小时之内发起
// REQUEST_AMOUNT_EXCEED	撤销触发的解冻金额超过冻结金额 请商户确认该笔预授权单是否已发生过解冻操作，如若发生过，则无法再发起撤销操作
// SYSTEM_ERROR	系统错误 请使用相同的参数再次调用
// PAYER_NOT_EXIST	买家账号不存在 确认买家是否已经注销账号，如有问题可联系支付宝小二处理（联系电话：4007585858）

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

// FandAuthUnfreezeReq 解冻支付宝请求
type FandAuthUnfreezeReq struct {
	AuthNo       string  `json:"authNo" validate:"required"`       // 支付宝授权资金订单号
	OutRequestNo string  `json:"outRequestNo" validate:"required"` // 商户端的唯一订单号
	Amount       float64 `json:"amount" validate:"required"`       // 解冻金额
	Remark       string  `json:"remark" validate:"required"`       // 备注
}

// 完结异常码
// ILLEGAL_ARGUMENT	参数异常或参数缺失	 请求参数有错，重新检查请求后，再重启发起资金解冻操作
// UNIQUE_VIOLATION	解冻信息被篡改	 更换商户请求流水号后，重新发起请求
// SYSTEM_ERROR	系统错误	 请使用相同的参数再次调用
// PAYER_USER_STATUS_LIMIT	付款方状态受限	 买家支付宝账户受限，请登录支付宝认证升级，详情咨询 4007585858
// AUTH_ORDER_NOT_EXIST	授权订单不存在 本笔授权订单不存在，确认请求资金解冻的资金授权订单号是否正确
// REQUEST_AMOUNT_EXCEED	请求解冻金额超限 更改解冻金额，重新发起请求
// ILLEGAL_STATUS	订单状态非法	 查询该笔授权操作信息，确认用户资金授权冻结成功
// ORDER_ALREADY_FINISH	授权订单已经完结，无法再进行资金操作 本笔授权订单已经完结，不允许进行资金解冻操作，确认请求资金解冻的资金授权订单号是否正确
// ORDER_ALREADY_CLOSED	授权订单已经关闭，无法再进行资金操作 本笔授权订单已经关闭，不允许进行资金解冻操作，确认请求资金解冻的资金授权订单号是否正确，该笔授权订单号是否已经发起过解冻
// PAYER_NOT_EXIST	买家不存在	 买家信息不存在，请联系支付宝小二确认买家是否销户。
// BIZ_ERROR	业务异常，	 商户自行确认该笔预授权订单是否被用于其他业务，或者联系支付宝客服
// ACCESS_FORBIDDEN	授权失败，本商户没有权限使用该产品，建议顾客使用其他方式付款 未签约合同或者合同已到期

type FandAuthUnfreezeRes struct {
	OutOrderNo string `json:"outOrderNo,omitempty"` // 商户端的唯一订单号
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

// 查询资金授权冻结订单异常码
// ILLEGAL_ARGUMENT	参数异常或参数缺失	 检查请求参数，修改后重新发起请求
// AUTH_ORDER_NOT_EXIST	支付宝资金授权订单不存在 检查传入参数中的支付宝资金授权订单号或商户授权订单号，修改后重新发起请求
// AUTH_OPERATION_NOT_EXIST	支付宝资金操作流水不存在 检查传入参数中的支付宝的授权资金操作流水号或商户的授权资金操作流水号，修改后重新发起请求
// ACCESS_FORBIDDEN	无权限使用该产品	未签约或签约已到期，请检查合约
// SYSTEM_ERROR	系统错误	 请使用相同的参数再次调用
// ENTERPRISE_PAY_BIZ_ERROR	因公付业务异常 请使用相同参数再次调用
// HAS_NO_PRIVILEGE	商户无权限查看该笔订单信息 请检查该商户是否创建了该笔授权订单

// FundAuthOperationDetailReq 查询资金授权冻结订单
type FundAuthOperationDetailReq struct {
	AuthNo       string `json:"authNo" validate:"required"`       // 支付宝资金授权订单号
	OutRequestNo string `json:"outRequestNo" validate:"required"` // 商户授权资金操作流水号
}
