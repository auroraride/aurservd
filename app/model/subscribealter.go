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
}

// ApplyListRsp 申请列表返回信息
type ApplyListRsp struct {
	// 天数
	Days int `json:"days"`
	// 申请时间
	ApplyTime string `json:"apply_time"`
	// 审核时间
	ReviewTime string `json:"review_time"`
	// 审核状态
	Status int `json:"status"`
}
