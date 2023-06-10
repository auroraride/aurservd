package model

const (
	SubscribeAlterStatusPending = iota // 未审核
	SubscribeAlterStatusAgree          // 已通过
	SubscribeAlterStatusRefuse         // 拒绝申请
)

// SubscribeAlterApplyReq 申请请求
type SubscribeAlterApplyReq struct {
	PaginationReq
	Start   *string `json:"start" query:"start"`     // 开始时间
	End     *string `json:"end" query:"end"`         // 结束时间
	Keyword *string `json:"keyword" query:"keyword"` // 关键词
	Status  *int    `json:"status" query:"status"`   // 审核状态
}

type SubscribeAlterApplyManagerReq struct {
	SubscribeAlterApplyReq
	EnterpriseID uint64 `json:"enterpriseId" query:"enterpriseId" param:"enterpriseId"` // 团签ID
}

// SubscribeAlterApplyListRsp 申请列表返回信息
type SubscribeAlterApplyListRsp struct {
	ID uint64 `json:"id"`
	// 天数
	Days int `json:"days"`
	// 申请时间
	ApplyTime string `json:"applyTime"`
	// 审核时间
	ReviewTime string `json:"reviewTime"`
	// 审核状态
	Status int `json:"status"`
	// 到期时间
	ExpireTime string `json:"expireTime,omitempty"`
	// 骑手姓名
	RiderName string `json:"riderName,omitempty"`
	// 骑手手机号
	RiderPhone string `json:"riderPhone,omitempty"`
}

// SubscribeAlterReviewReq 审批订阅申请
type SubscribeAlterReviewReq struct {
	Ids    []uint64 `json:"ids"`
	Status int      `json:"status" validate:"required" enums:"1,2"` // 审批状态 1:通过 2:拒绝
}
