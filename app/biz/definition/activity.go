package definition

import "github.com/auroraride/aurservd/app/model"

// ActivityReqCommon 活动请求公共字段
type ActivityReqCommon struct {
	Name   string `json:"name" validate:"required" trans:"名称"`  // 名称
	Sort   int    `json:"sort" validate:"required" trans:"排序"`  // 排序
	Image  string `json:"image" validate:"required" trans:"图片"` // 图片
	Link   string `json:"link" validate:"required" trans:"连接"`  // 连接
	Remark string `json:"remark,omitempty" trans:"备注"`          // 备注
}

// ActivityDetail 活动详情
type ActivityDetail struct {
	ID     uint64 `json:"id" trans:"活动ID"`             // ID
	Name   string `json:"name" trans:"名称"`             // 名称
	Sort   int    `json:"sort" trans:"排序"`             // 排序
	Image  string `json:"image" trans:"图片"`            // 图片
	Link   string `json:"link" trans:"连接"`             // 连接
	Remark string `json:"remark,omitempty" trans:"备注"` // 备注
}

// ActivitySaveReq 活动保存请求
type ActivitySaveReq struct {
	ActivityReqCommon
}

// ActivityModifyReq 活动修改请求
type ActivityModifyReq struct {
	model.IDParamReq
	ActivityReqCommon
}

// ActivityListReq 活动列表请求
type ActivityListReq struct {
	model.PaginationReq
}
