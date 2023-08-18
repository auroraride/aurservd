package promotion

import "github.com/auroraride/aurservd/app/model"

type EarningsStatus uint8

// 收益状态
const (
	EarningsStatusUnsettled EarningsStatus = iota // 未结算
	EarningsStatusSettled                         // 已结算
	EarningsStatusCanceled                        // 已取消
)

func (a EarningsStatus) Value() uint8 {
	return uint8(a)
}

func (a EarningsStatus) String() string {
	switch a {
	case EarningsStatusUnsettled:
		return "未结算"
	case EarningsStatusSettled:
		return "已结算"
	case EarningsStatusCanceled:
		return "已取消"
	}
	return ""
}

// EarningsReq 请求参数
type EarningsReq struct {
	model.PaginationReq
	EarningsFilter
	ID *uint64 `json:"id"  param:"id"` // id
}

// EarningsFilter 收益筛选Filter
type EarningsFilter struct {
	Keyword           *string `json:"keyword"  query:"keyword"`                     // 关键词 手机号/姓名
	Status            *uint8  `json:"status" enums:"0,1,2"  query:"status"`         // 状态 0:未结算 1:已结算 2:已取消
	CommissionRuleKey *string `json:"commissionRuleKey"  query:"commissionRuleKey"` // 任务类型
	Start             *string `json:"start"  query:"start"`                         // 开始日期
	End               *string `json:"end" query:"end"`                              // 结束日期
}

// EarningsRes 收益列表信息
type EarningsRes struct {
	EarningsDetail
	Phone string `json:"phone" ` // 手机号
	Name  string `json:"name" `  // 姓名

}

// EarningsDetail 收益明细
type EarningsDetail struct {
	ID                 uint64  `json:"id" `                 // id
	MemberID           uint64  `json:"memberId" `           // 会员id
	CommissionID       uint64  `json:"commissionId" `       // 返佣方案id
	CommissionRuleName string  `json:"commissionRuleName" ` // 返佣方案类型
	Amount             float64 `json:"amount" `             // 金额
	Status             uint8   `json:"status" `             // 状态 0:未结算 1:已结算 2:已取消
	CreateTime         string  `json:"createTime" `         // 返佣时间
	Remark             string  `json:"remark" `             // 备注
}

// EarningsCreateReq 创建收益
type EarningsCreateReq struct {
	MemberID          uint64            `json:"memberId" validate:"required"`          // 会员id
	RiderID           uint64            `json:"riderId" validate:"required"`           // 骑手id
	CommissionID      uint64            `json:"commissionId" validate:"required"`      // 返佣方案id
	CommissionRuleKey CommissionRuleKey `json:"commissionRuleKey" validate:"required"` // 返佣方案类型
	Amount            float64           `json:"amount" validate:"required"`            // 金额
	OrderID           uint64            `json:"orderId" validate:"required"`           // 订单id
	PlanID            uint64            `json:"planId" validate:"required"`            // 骑士卡id
}

// EarningsCancelReq 取消收益
type EarningsCancelReq struct {
	ID     uint64  `json:"id" validate:"required"`     // id
	Remark *string `json:"remark" validate:"required"` // 备注
}
