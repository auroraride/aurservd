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
	Type *uint8 `json:"type"`
	// 关键词
	Keyword string `json:"keyword"`
	// 反馈开始时间
	StartTime *string `json:"startTime"`
	// 反馈结束时间
	EndTime *string `json:"endTime"`
	// 反馈用户团签id
	EnterpriseID *uint64 `json:"enterpriseID"`
}

// FeedbackListRsp 反馈列表
type FeedbackListRsp struct {
	Pagination Pagination        `json:"pagination"` // 分页属性
	List       []*FeedbackDetail `json:"list"`       // 列表
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
