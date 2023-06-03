package model

const (
	SubscribeAlterUnreviewed = iota // 未审核
	SubscribeAlterNormal            // 已通过
	SubscribeAlterUnpass            // 未通过
)

// SubscribeAlterApplyReq 申请请求
type SubscribeAlterApplyReq struct {
	PaginationReq
	Start   *string `json:"start" query:"start"`     // 开始时间
	End     *string `json:"end" query:"end"`         // 结束时间
	Keyword *string `json:"keyword" query:"keyword"` // 关键词
	Status  *int    `json:"status" query:"status"`   // 审核状态
	ID      uint64  `json:"id" query:"id"`           // 团签id
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
