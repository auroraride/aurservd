package definition

import "github.com/auroraride/aurservd/app/model"

const (
	SourceRider = uint8(iota) + 1 // 骑手
	SourceAgent                   // 代理
)

// FeedbackReq 反馈请求参数
type FeedbackReq struct {
	Content string   `json:"content" validate:"required"` // 反馈内容
	Url     []string `json:"url"`                         // 反馈图片
	Type    uint8    `json:"type" validate:"required"`    // 反馈类型
}

// FeedbackListReq 反馈列表请求参数
type FeedbackListReq struct {
	model.PaginationReq
	Type         *uint8  `json:"type" query:"type"`                       // 反馈类型
	Source       *uint8  `json:"source" query:"source"`                   // 反馈来源
	Keyword      string  `json:"keyword" query:"keyword"`                 // 关键词
	Start        *string `json:"start" query:"start"`                     // 反馈开始时间
	End          *string `json:"end" query:"end"`                         // 反馈结束时间
	Enterprise   *uint8  `json:"enterprise,omitempty" query:"enterprise"` // 是否团签, 0:全部 1:团签 2:个签
	EnterpriseID *uint64 `json:"enterpriseID" query:"enterpriseId"`       // 团签企业ID, `enterprise = 1`时才会生效
}

// FeedbackDetail 反馈详情
type FeedbackDetail struct {
	ID                     uint64   `json:"id"`                     // 反馈ID
	Content                string   `json:"content"`                // 反馈内容
	Url                    []string `json:"url"`                    // 反馈图片
	Type                   uint8    `json:"type"`                   // 反馈类型
	Source                 uint8    `json:"source"`                 // 反馈来源
	EnterpriseID           *uint64  `json:"enterpriseID"`           // 反馈用户团签id
	EnterpriseName         string   `json:"enterpriseName"`         // 反馈用户团签名称
	EnterpriseContactName  string   `json:"enterpriseContactName"`  // 反馈用户团签联系人
	EnterpriseContactPhone string   `json:"enterpriseContactPhone"` // 反馈用户团签联系电话
	CreatedAt              string   `json:"createdAt"`              // 反馈时间
}
