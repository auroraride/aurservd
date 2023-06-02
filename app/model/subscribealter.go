package model

const (
	// SubscribealterUnreviewed 未审核
	SubscribeAlterUnreviewed = iota
	// SubscribeAlterNormal 已通过
	SubscribeAlterNormal
	// SubscribeAlterUnpass 未通过
	SubscribeAlterUnpass
)

// SubscribeAlterApplyReq 申请请求
type SubscribeAlterApplyReq struct {
	PaginationReq
	// 开始时间
	StartTime *string `json:"start"`
	// 结束时间
	EndTime *string `json:"end"`
	// 关键词
	Keyword *string `json:"keyword"`
	// 审核状态
	Status *int `json:"status"`
	// 团签id
	ID uint64 `json:"id" param:"id"`
}

// SubscribeAlterApplyListRsp 申请列表返回信息
type SubscribeAlterApplyListRsp struct {
	// id
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
	ExpireTime int `json:"expireTime"`
	// 骑手姓名
	RiderName string `json:"riderName"`
	// 骑手手机号
	RiderPhone string `json:"riderPhone"`
}

// SubscribeAlterReviewReq 审批订阅申请
type SubscribeAlterReviewReq struct {
	Ids []uint64 `json:"ids"`
	// 审批状态
	Status int `json:"status"`
}
