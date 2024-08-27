package definition

import "github.com/auroraride/aurservd/app/model"

// Activity 活动公共字段
type Activity struct {
	Name         string              `json:"name"`         // 名称
	Sort         int                 `json:"sort"`         // 排序
	Image        model.ActivityImage `json:"image"`        // 图片
	Link         string              `json:"link"`         // 链接
	Popup        *bool               `json:"popup"`        // 弹窗 true:是 false:否
	Home         *bool               `json:"home"`         // 首页icon true:是 false:否
	Remark       string              `json:"remark"`       // 备注
	Enable       bool                `json:"enable" `      // 是否启用 true:是 false:否
	UpdatedAt    string              `json:"updatedAt"`    // 更新时间
	Introduction string              `json:"introduction"` // 活动简介
}

// ActivityListReq 活动列表请求
type ActivityListReq struct {
	model.PaginationReq
	Keyword *string `json:"keyword" query:"keyword"` // 关键字
	Enable  *bool   `json:"enable" query:"enable"`   // 是否启用
	Start   *string `json:"start" query:"start"`     // 开始时间
	End     *string `json:"end" query:"end"`         // 结束时间
}

// ActivityDetail 活动详情
type ActivityDetail struct {
	ID uint64 `json:"id"`
	Activity
}

// ActivityCreateReq 活动保存请求
type ActivityCreateReq struct {
	Name         string              `json:"name" validate:"required" trans:"名称"`  // 名称
	Sort         int                 `json:"sort" validate:"required" trans:"排序"`  // 排序
	Image        model.ActivityImage `json:"image" validate:"required" trans:"图片"` // 图片
	Link         string              `json:"link" validate:"required" trans:"链接"`  // 链接
	Enable       *bool               `json:"enable" validate:"required"`           // 是否启用 true:是 false:否
	Introduction string              `json:"introduction" validate:"required"`     // 活动简介
	Popup        *bool               `json:"popup"`                                // 弹窗 true:是 false:否
	Home         *bool               `json:"home"`                                 // 首页icon true:是 false:否
	Remark       *string             `json:"remark"`                               // 备注
}

// ActivityModifyReq 活动修改请求
type ActivityModifyReq struct {
	model.IDParamReq
	ActivityCreateReq
}
