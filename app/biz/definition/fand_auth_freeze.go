package definition

type FandAuthFreezeStatus string

const (
	FandAuthFreezeStatusInit    FandAuthFreezeStatus = "INIT"    // 初始化
	FandAuthFreezeStatusSuccess FandAuthFreezeStatus = "SUCCESS" // 成功
	FandAuthFreezeStatusClosed  FandAuthFreezeStatus = "CLOSED"  // 关闭
)

const (
	FandAuthNotifyType   = "fund_auth_freeze"           // 资金预授权冻结成功
	FandAuthUnfreezeType = "fund_auth_unfreeze"         // 资金授权订单解冻
	FandAuthPayType      = "fund_auth_operation_cancel" // 资金预授权明细撤销
)

type FandAuthFreezeNotification struct {
	AuthNo                    string               `json:"auth_no"`                                // 支付宝资金授权订单号
	NotifyType                string               `json:"notify_type"`                            // 通知类型  fund_auth_freeze 资金预授权冻结成功 fund_auth_unfreeze 资金授权订单解冻 fund_auth_operation_cancel 资金预授权明细撤销
	FundAuthFreeze            string               `json:"fund_auth_freeze"`                       // 资金操作类型
	OutOrderNo                string               `json:"out_order_no"`                           // 商家的资金授权订单号
	OperationID               string               `json:"operation_id"`                           // 支付宝的资金操作流水号
	OutRequestNo              string               `json:"out_request_no"`                         // 商家资金操作流水号
	OperationType             string               `json:"operation_type"`                         // 资金操作类型 freeze 冻结 unfreeze 解冻
	Amount                    string               `json:"amount"`                                 // 本次操作冻结的金额
	Status                    FandAuthFreezeStatus `json:"status"`                                 // 资金预授权明细的状态
	GmtCreate                 string               `json:"gmt_create"`                             // 明细创建时间
	GmtTrans                  string               `json:"gmt_trans"`                              // 明细处理完成时间
	PayerLogonID              string               `json:"payer_logon_id"`                         // 付款方支付宝账号登录号，脱敏
	PayerUserID               string               `json:"payer_user_id"`                          // 付款方支付宝账号UID
	PayeeLogonID              string               `json:"payee_logon_id,omitempty"`               // 收款方支付宝账号，脱敏
	PayeeUserID               string               `json:"payee_user_id,omitempty"`                // 收款方支付宝账号UID
	TotalFreezeAmount         string               `json:"total_freeze_amount"`                    // 累计冻结金额
	TotalUnfreezeAmount       string               `json:"total_unfreeze_amount"`                  // 累计解冻金额
	TotalPayAmount            string               `json:"total_pay_amount"`                       // 累计支付金额
	RestAmount                string               `json:"rest_amount"`                            // 剩余冻结金额
	CreditAmount              string               `json:"credit_amount,omitempty"`                // 本次操作中信用金额
	FundAmount                string               `json:"fund_amount,omitempty"`                  // 本次操作中自有资金金额
	TotalFreezeCreditAmount   string               `json:"total_freeze_credit_amount,omitempty"`   // 累计冻结信用金额
	TotalFreezeFundAmount     string               `json:"total_freeze_fund_amount,omitempty"`     // 累计冻结自有资金金额
	TotalUnfreezeCreditAmount string               `json:"total_unfreeze_credit_amount,omitempty"` // 累计解冻信用金额
	TotalUnfreezeFundAmount   string               `json:"total_unfreeze_fund_amount,omitempty"`   // 累计解冻自有资金金额
	TotalPayCreditAmount      string               `json:"total_pay_credit_amount,omitempty"`      // 累计支付信用金额
	TotalPayFundAmount        string               `json:"total_pay_fund_amount,omitempty"`        // 累计支付自有资金金额
	RestCreditAmount          string               `json:"rest_credit_amount,omitempty"`           // 剩余冻结信用金额
	RestFundAmount            string               `json:"rest_fund_amount,omitempty"`             // 剩余冻结自有资金金额
	PreAuthType               string               `json:"pre_auth_type,omitempty"`                // 预授权类型 预授权类型，目前支持 CREDIT_AUTH(信用预授权); 商家可根据该标识来判断该笔预授权的类型，当返回值为"CREDIT_AUTH"表明该笔预授权为信用预授权，没有真实冻结资金；当返回值为空或者不为"CREDIT_AUTH"则表明该笔预授权为普通资金预授权，会冻结用户资金
	CreditMerchantExt         string               `json:"credit_merchant_ext,omitempty"`          // 芝麻透出给商家的信息
}

// FandAuthUnfreeze 解除预授权
type FandAuthUnfreeze struct {
	AuthNo       string  `json:"auth_no"`        // 支付宝资金授权订单号
	OutRequestNo string  `json:"out_request_no"` // 商户授权资金操作流水号
	Amount       float64 `json:"amount"`         // 预授权解冻金额
	Remark       string  `json:"remark"`         // 解冻说明
}

// TradePay 冻结转支付
type TradePay struct {
	AuthNo      string  `json:"auth_no"`      // 支付宝资金授权订单号
	OutTradeNo  string  `json:"out_trade_no"` // 商户订单号
	TotalAmount float64 `json:"total_amount"` // 支付金额
	Subject     string  `json:"subject"`      // 订单标题
}
