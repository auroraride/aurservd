package model

const (
	SourceRider = iota + 1 // 骑手
	SourceAgent            // 代理
)

// FeedbackReq 反馈请求参数
type FeedbackReq struct {
	Content     string       `json:"content" validate:"required"` // 反馈内容
	Url         []string     `json:"url"`                         // 反馈图片
	Type        uint8        `json:"type" validate:"required"`    // 反馈类型
	CityID      *uint64      `json:"cityId" `                     // 城市ID
	VersionInfo *VersionInfo `json:"versionInfo"`                 // 版本信息
}

// FeedbackListReq 反馈列表请求参数
type FeedbackListReq struct {
	PaginationReq
	Type         *uint8  `json:"type" query:"type"`                 // 反馈类型
	Source       *uint8  `json:"source" query:"source"`             // 反馈来源 1:骑手 2:代理
	Keyword      string  `json:"keyword" query:"keyword"`           // 关键词
	Start        *string `json:"start" query:"start"`               // 反馈开始时间
	End          *string `json:"end" query:"end"`                   // 反馈结束时间
	Enterprise   *uint8  `json:"enterprise" query:"enterprise"`     // 是否团签, 0:全部 1:团签 2:个签
	EnterpriseID *uint64 `json:"enterpriseID" query:"enterpriseId"` // 团签企业ID, `enterprise = 1`时才会生效
}

// FeedbackDetail 反馈详情
type FeedbackDetail struct {
	ID             uint64      `json:"id"`                       // 反馈ID
	Content        string      `json:"content"`                  // 反馈内容
	Url            []string    `json:"url"`                      // 反馈图片
	Type           uint8       `json:"type"`                     // 反馈类型
	Source         uint8       `json:"source,omitempty"`         // 反馈来源
	EnterpriseID   uint64      `json:"enterpriseId,omitempty"`   // 反馈用户团签id
	EnterpriseName string      `json:"enterpriseName,omitempty"` // 反馈用户团签名称
	Name           string      `json:"name"`                     // 反馈用户名称
	Phone          string      `json:"phone"`                    // 反馈用户电话
	CreatedAt      string      `json:"createdAt"`                // 反馈时间
	CityName       string      `json:"cityName"`                 // 城市名称
	VersionInfo    VersionInfo `json:"versionInfo"`              // 版本信息
}

// VersionInfo 版本信息
type VersionInfo struct {
	Version   string `json:"version"`   // 版本号
	CommId    string `json:"commitId"`  // 提交ID
	BuildTime string `json:"buildTime"` // 编译时间
	CiJobId   string `json:"ciJobId"`   // CI任务ID
}
