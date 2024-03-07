package model

// FeedbackReq 新增反馈记录请求参数
type FeedbackReq struct {
	// 反馈内容
	Content string `json:"content"`
	// 反馈图片
	Url []string `json:"url"`
	// 反馈类型
	Type uint8 `json:"type"`
}

// FeedbackListReq 反馈列表请求参数
type FeedbackListReq struct {
	PaginationReq
	// 反馈类型
	Type *uint8 `json:"type" query:"type"`
	// 反馈来源
	Source *uint8 `json:"source" query:"source"`
	// 关键词
	Keyword string `json:"keyword" query:"keyword"`
	// 反馈开始时间
	Start *string `json:"start" query:"start"`
	// 反馈结束时间
	End *string `json:"end" query:"end"`
	// 是否团签, 0:全部 1:团签 2:个签
	Enterprise *uint8 `json:"enterprise,omitempty" query:"enterprise"`
	// 团签企业ID, `enterprise = 1`时才会生效
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
	// 反馈来源
	Source uint8 `json:"source"`
	// 反馈用户团签id
	EnterpriseID *uint64 `json:"enterpriseID"`
	// 反馈用户团签名称
	EnterpriseName string `json:"enterpriseName"`
	// 反馈用户团签联系人
	EnterpriseContactName string `json:"enterpriseContactName"`
	// 反馈用户团签联系电话
	EnterpriseContactPhone string `json:"enterpriseContactPhone"`
	// 反馈时间
	CreatedAt string `json:"createdAt"`
}
