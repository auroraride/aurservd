package model

const (
	// Unreviewed 未审核
	Unreviewed = iota
	// Normal 已通过
	Normal
	// Unpass 未通过
	Unpass
)

// ApplyReq 申请请求
type ApplyReq struct {
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

// ApplyListRsp 申请列表返回信息
type ApplyListRsp struct {
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

// ReviewReq 审批订阅申请
type ReviewReq struct {
	Ids []uint64 `json:"ids"`
	// 审批状态
	Status int `json:"status"`
}
