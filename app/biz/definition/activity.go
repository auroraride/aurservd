package definition

import "github.com/auroraride/aurservd/app/model"

// ActivityReqCommon 活动请求公共字段
type ActivityReqCommon struct {
	Name   string `json:"name" validate:"required" trans:"名称"`
	Sort   int    `json:"sort" validate:"required" trans:"排序"`
	Image  string `json:"image" validate:"required" trans:"图片"`
	Link   string `json:"link" validate:"required" trans:"链接"`
	Remark string `json:"remark,omitempty"` // 备注
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
