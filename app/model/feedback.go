package model

// FeedbackReq 新增反馈记录请求参数
type FeedbackReq struct {
	// 反馈内容
	Content string `json:"content"`
	// 反馈图片
	Url []string `json:"url"`
	// 反馈类型
	Type uint8 `json:"type"`
	// 反馈用户团签id
	EnterpriseID int64 `json:"enterpriseID"`
}

// FeedbackListReq 反馈列表请求参数
type FeedbackListReq struct {
	PaginationReq
	// 反馈类型
	Type *uint8 `json:"type" query:"type"`
	// 关键词
	Keyword string `json:"keyword" query:"keyword"`
	// 反馈开始时间
	Start *string `json:"start" query:"start"`
	// 反馈结束时间
	End *string `json:"end" query:"end"`
	// 反馈用户团签id
	EnterpriseID *uint64 `json:"enterpriseID" query:"enterpriseId"`
}

// FeedbackDetail 反馈详情
type FeedbackDetail struct {
	// 反馈ID
	ID uint64 `json:"id"`
	// 反馈内容
	Content string `json:"content"`
	// 反馈图片
	Url []string `json:"url"`
	// 反馈类型
	Type uint8 `json:"type"`
	// 反馈用户团签id
	EnterpriseID int64 `json:"enterpriseID"`
	// 反馈用户团签名称
	EnterpriseName string `json:"enterpriseName"`
	// 反馈用户团签联系人
	EnterpriseContactName string `json:"enterpriseContactName"`
	// 反馈用户团签联系电话
	EnterpriseContactPhone string `json:"enterpriseContactPhone"`
	// 反馈时间
	CreatedAt string `json:"createdAt"`
}
