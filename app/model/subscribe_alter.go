package model

const (
	SubscribeAlterStatusPending = iota // 未审核
	SubscribeAlterStatusAgree          // 已通过
	SubscribeAlterStatusRefuse         // 拒绝申请
)

// SubscribeAlterRiderReq 骑手加时申请
type SubscribeAlterRiderReq struct {
	Days int `json:"days" validate:"required" trans:"天数"`
}

// SubscribeAlterFilter 加时申请列表筛选
type SubscribeAlterFilter struct {
	PaginationReq
	Start   *string `json:"start" query:"start"`     // 开始时间
	End     *string `json:"end" query:"end"`         // 结束时间
	Keyword *string `json:"keyword" query:"keyword"` // 骑手关键词 (姓名或手机号)
	Status  *int    `json:"status" query:"status"`   // 审核状态
	RiderID *uint64 `json:"riderId" query:"riderId"` // 骑手ID
}

type SubscribeAlterListReq struct {
	SubscribeAlterFilter
	EnterpriseID uint64 `json:"enterpriseId" query:"enterpriseId" param:"enterpriseId"` // 团签ID
}

type SubscribeListRiderReq struct {
	PaginationReq
	Start  *string `json:"start" query:"start"`   // 开始时间
	End    *string `json:"end" query:"end"`       // 结束时间
	Status *int    `json:"status" query:"status"` // 审核状态
}

// SubscribeAlterApplyListRes 申请列表返回信息
type SubscribeAlterApplyListRes struct {
	ID             uint64 `json:"id"`
	Days           int    `json:"days"`                     // 天数
	ApplyTime      string `json:"applyTime"`                // 申请时间
	ReviewTime     string `json:"reviewTime"`               // 审核时间
	Status         int    `json:"status"`                   // 审核状态
	SubscribeEndAt string `json:"subscribeEndAt,omitempty"` // 团签到期时间
	Rider          *Rider `json:"rider,omitempty"`          // 骑手
}

// SubscribeAlterReviewReq 审批订阅申请
type SubscribeAlterReviewReq struct {
	Ids    []uint64 `json:"ids"`
	Status int      `json:"status" validate:"required" enums:"1,2" trans:"审批状态"` // 1:通过 2:拒绝
}

// SubscribeAlter 订阅天数调整
type SubscribeAlter struct {
	ID     uint64 `json:"id" validate:"required"`     // 订阅ID
	Days   int    `json:"days" validate:"required"`   // 调整天数, 正加负减
	Reason string `json:"reason" validate:"required"` // 调整理由
}

// SubscribeAlterReq 订阅天数调整请求
type SubscribeAlterReq struct {
	SubscribeAlter
	EnterpriseID uint64 `json:"enterpriseId"` // 团签id
	AgentID      uint64 `json:"agentId"`      // 代理商操作人id
}
