package definition

import "github.com/auroraride/aurservd/app/model"

// ActivityReqCommon 活动请求公共字段
type ActivityReqCommon struct {
	Name   string            `json:"name" validate:"required" trans:"名称"`  // 名称
	Sort   int               `json:"sort" validate:"required" trans:"排序"`  // 排序
	Image  map[string]string `json:"image" validate:"required" trans:"图片"` // 图片  {"list":"图片地址","popup":"图片地址","index":"图片地址" }
	Link   string            `json:"link" validate:"required" trans:"链接"`  // 链接
	Popup  *bool             `json:"popup"`                                // 弹窗 true:是 false:否
	Index  *bool             `json:"index"`                                // 首页icon true:是 false:否
	Remark *string           `json:"remark"`                               // 备注
}

// ActivityDetail 活动详情
type ActivityDetail struct {
	ID uint64 `json:"id"`
	ActivityReqCommon
}

// ActivityCreateReq 活动保存请求
type ActivityCreateReq struct {
	ActivityReqCommon
}

// ActivityModifyReq 活动修改请求
type ActivityModifyReq struct {
	model.IDParamReq
	ActivityReqCommon
}
