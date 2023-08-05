package promotion

import "github.com/auroraride/aurservd/app/model"

type WithdrawalMethod uint8

const (
	TaxExemptAmount  float64 = 800    // 免税额度
	TaxRate          float64 = 0.2    // 税率
	FransferFee      float64 = 5      // 转账手续费
	FeeExemptionCity string  = "西安市"  // 免手续费城市
	FeeExemptionBank string  = "招商银行" // 免手续费银行
	FeeRate          float64 = 0.006  // 手续费率
)

const (
	WithdrawalMethodBank WithdrawalMethod = iota + 1 // 银行卡
)

func (w WithdrawalMethod) Value() uint8 {
	return uint8(w)
}
func (w WithdrawalMethod) String() string {
	switch w {
	case WithdrawalMethodBank:
		return "银行卡"
	}
	return ""
}

type WithdrawalStatus uint8

const (
	WithdrawalStatusPending WithdrawalStatus = iota // 待审核
	WithdrawalStatusSuccess                         // 成功
	WithdrawalStatusFailed                          // 失败
)

func (w WithdrawalStatus) Value() uint8 {
	return uint8(w)
}

// WithdrawalListReq 提现列表请求
type WithdrawalListReq struct {
	model.PaginationReq
	WithdrawalFilter
	ID *uint64 `json:"id" query:"id"` // 会员id
}

// WithdrawalFilter 提现列表筛选
type WithdrawalFilter struct {
	Account *string `json:"account" query:"account"`             // 提现账户
	Status  *uint8  `json:"status" query:"status" enums:"0,1,2"` // 提现状态 0:待审核 1:成功 2:失败
	Start   *string `json:"start" query:"start" `                // 开始日期
	End     *string `json:"end" query:"end"`                     // 结束日期
	Keywork *string `json:"keywork" query:"keywork"`             // 关键字
}

// WithdrawalListRes 提现列表响应
type WithdrawalListRes struct {
	WithdrawalDetail
}

// WithdrawalDetail 提现详情
type WithdrawalDetail struct {
	ID              uint64       `json:"id"` // id
	*MemberBaseInfo              // 会员
	BankCard        *BankCardRes `json:"bankCard,omitempty"` // 银行卡
	Amount          float64      `json:"amount"`             // 提现金额
	Fee             float64      `json:"fee"`                // 提现手续费
	Status          uint8        `json:"status"`             // 状态 0:待审核 1:成功 2:失败
	CreatedAt       string       `json:"createdAt"`          // 创建时间
	Method          string       `json:"method"`             // 提现方式 1:银行卡
	Remark          string       `json:"remark"`             // 备注
	ApplyTime       string       `json:"applyTime"`          // 申请时间
	ReviewTime      string       `json:"review_time"`        // 审核时间
}

// WithdrawalAlterReq 提现申请请求
type WithdrawalAlterReq struct {
	AccountID   uint64  `json:"accountId" validate:"required"`                    // 提现账户ID
	ApplyAmount float64 `json:"applyAmount" validate:"required,min=100,max=4000"` //  提现金额
}

// WithdrawalApprovalReq 审批提现请求
type WithdrawalApprovalReq struct {
	IDs    []uint64 `json:"ids" validate:"required"`                // 提现id
	Status uint8    `json:"status" validate:"required,oneof=1 2"`   // 状态 1:成功 2:失败
	Remark string   `json:"remark" validate:"required_if=Status 2"` // 备注
}

// WithdrawalFeeRes 计算提现费用返回
type WithdrawalFeeRes struct {
	ApplyAmount    float64 `json:"applyAmount"`    // 提现金额
	AmountReceived float64 `json:"amountReceived"` // 实际到账金额
	WithdrawalFee  float64 `json:"withdrawalFee"`  // 服务费
	Taxable        float64 `json:"taxable"`        // 应缴税款
}

type WithdrawalExport struct {
	Remark string `json:"remark"`
}
